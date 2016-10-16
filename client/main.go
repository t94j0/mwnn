package client

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/t94j0/gocui"
	"github.com/t94j0/mwnn/gpg"
	pb "github.com/t94j0/mwnn/message"
)

/*
TODO:
  1. Make the message_box scrollable
  2. Change keybinding issue with seekeys
  3. Fix error that is cluttering log files
*/

////////////////////////////
// Configuration variables /
////////////////////////////

var (
	serviceHost        string
	servicePort        string
	publicKeyLocation  string
	publicKey          string
	privateKeyPass     []byte
	privateKeyLocation string
	logLocation        string
)

///////////
// Types //
///////////

// User holds all the user information that we need to know as the client
type User struct {
	Username  string
	PublicKey []byte
}

//////////////////////
// Global Variables //
//////////////////////

var decryptor gpg.Decryptor
var connectedUsers = make(map[string]User, 0)
var currentUser string
var logger *log.Logger
var logFile *os.File
var conn net.Conn

////////////////////
// Loop Functions //
////////////////////

func checkMessageLoop(view *gocui.View) {
	for {
		// Create buffer to put recieved data in
		buf := make([]byte, 10000)
		// Read the data and throw error if it reached EOF. That means the server is down
		if _, err := conn.Read(buf); err == io.EOF {
			logger.Panic("Connection to server broken")
		} else if err != nil {
			logger.Println(err)
		}

		// The buffer is really big, so delete the whitespace
		buf = bytes.TrimRight(buf, "\x00")
		// Just continue if there is nothing but whitespace
		if len(buf) == 0 {
			continue
		}

		messages := bytes.Split(buf, []byte{0})
		for _, message := range messages {
			parseMessage(message, view)
		}
	}
}

//////////////////////
// Helper Functions //
//////////////////////

// Refrenced in main.go under messageHandler
// Clear screen and set cursor before all of the spaces that will be created
// beacuse of gocui.
func clearInput(g *gocui.Gui) {
	g.Execute(func(g *gocui.Gui) error {
		v, err := g.View("input_box")
		if err != nil {
			return err
		}
		v.Clear()

		if err := v.SetCursor(0, 0); err != nil {
			return err
		}

		return nil
	})
}

/*

Idea for commandHandler:
	* Make this a function a wrapper that calls a function that takes a commandType string and returns a function.
	  This allows us to reap a few benefits:
		- We only have to write the default commands
		- Modularizes commands
		- Makes it easier for third parties to add command plugins
		- Basically makes the way for handling commands a library

*/
// Refrenced in main.go under messageHandler
func commandHandler(command string, g *gocui.Gui) error {
	command = strings.TrimLeft(command, "/")
	commandArr := strings.Split(command, " ")
	g.Execute(func(g *gocui.Gui) error {
		v, err := g.View("messages_box")
		if err != nil {
			return err
		}

		switch commandArr[0] {
		case "exit":
			return gocui.ErrQuit
		case "quit":
			return gocui.ErrQuit
		case "message":
			messageUsage := "usage: /message [username] [text]"
			if len(commandArr) < 3 {
				fmt.Fprintln(v, messageUsage)
				return nil
			}

			// Put the message back together
			message := strings.Join(commandArr[2:], " ")
			user := commandArr[1]

			if _, ok := connectedUsers[user]; !ok {
				fmt.Fprintln(v, "That user does not exist")
				return nil // Possibly return custom "Error: User not found" instead of printing?
			}

			messageEnc, err := gpg.Encrypt(message, string(connectedUsers[user].PublicKey))
			if err != nil {
				fmt.Fprintln(v, "There was a problem encrypting message. Aborting")
				return err
			}
			fmt.Fprintln(v, "<- ", user+":", message) // Maybe "user <- You: message"?
			if err := sendMessage(7, user, messageEnc); err != nil {
				fmt.Fprintln(v, "There was a problem sending the message")
				return err
			}
		case "seekeys":
			// Create a view that block the whole screen to display all other user's public keys. This makes sure that the server can't send fake public keys because everyone should be able to identify eachother's public keys.
			keysView, err := createViewKeysWindow(g)
			if err != nil {
				fmt.Fprintln(v, "Problem opening keys window")
				return err
			}
			// Prints out the username and public key of all connected users
			for _, connection := range connectedUsers {
				fmt.Fprintln(keysView, blue(connection.Username), "\n", string(connection.PublicKey), "\n\n")
			}

		default:
			fmt.Fprintln(v, "This is not a command")
		}
		return nil
	})
	return nil
}

// Refrenced in gui.go under keybindings
// It is used to take the message from the text box and send it to all other users as a text message
func messageHandler(g *gocui.Gui, v *gocui.View) error {
	// Get message from the view buffer
	messageStr := v.ViewBuffer()
	// TODO: Fix problem where when multiple things are said, spaces start
	// getting added to the beginning.
	// Solution: This is happening because gocui's clearRunes function replaces
	// all characters with spaces. A temporary fix is to set the cursor to 0,0
	// when clearing.
	// TODO: Fix newline being added to end of every message
	messageParsed := strings.Trim(messageStr, "\n")
	messageParsed = strings.TrimLeft(messageParsed, " ")

	// If the first letter is a forward slash, then we know that there is a command coming. If
	// it isn't, then it is a message that should be shown to everyone connected.
	if len(messageParsed) > 1 && string([]byte(messageParsed)[0]) == "/" {
		if err := commandHandler(messageParsed, g); err != nil {
			return err
		}
	} else {
		for _, user := range connectedUsers {
			messageEncrypt, err := gpg.Encrypt(messageParsed, string(user.PublicKey))
			if err != nil {
				logger.Println(err)
				fmt.Println(err)
			}
			if err := sendMessage(0, user.Username, messageEncrypt); err != nil {
				logger.Println(err)
			}
		}
	}

	clearInput(g)

	return nil
}

// Refrenced in main.go in checkMessageLoop
// Used to take a protobuf byte array and figure out what to do with the data based on the message
// type
func parseMessage(buf []byte, view *gocui.View) {
	// Create message variable for the incoming proto
	incomingMessage := &pb.Message{}

	// Unmarshal the proto into incomingMessage
	if err := proto.Unmarshal(buf, incomingMessage); err != nil {
		logger.Println(err)
	}

	// This just makes sure that all of the data is legit because you can't make a gnupg key
	// that has a username less than 4.
	if len(incomingMessage.GetSender()) < 5 {
		return
	}

	switch incomingMessage.GetMessageType() {
	// Type = 0; Message is a text message
	case 0:
		decryptedMessage, err := decryptor.Decrypt(string(incomingMessage.GetMessage()))
		if err != nil && err != io.EOF {
			logger.Println(err)
		}
		fmt.Fprintln(view, incomingMessage.GetSender()+":", decryptedMessage)
	// Type = 1; Message is a login message
	case 1:
		if incomingMessage.GetMessage() != "error" {
			// Print all synced users
			fmt.Fprintln(view, blue("\nOnline Users:"))
			for _, userPair := range strings.Split(incomingMessage.GetMessage(), ",") {
				userArr := strings.Split(userPair, "|")
				username := userArr[0]
				usersPublicKey := userArr[1]
				connectedUsers[username] = User{username, []byte(usersPublicKey)}
				fmt.Fprintf(view, "%s ", blue(username))
			}
			fmt.Fprintln(view, "\n")
		} else {
			logger.Println("There was an error with the loging message")
		}
	case 2:
		// Create empty user with no public key. If they do not share their public
		// key, they can't send messages to anybody
		loginArr := strings.Split(incomingMessage.GetMessage(), "|")
		publicKey := string(loginArr[1])
		connectedUsers[loginArr[0]] = User{loginArr[0], []byte(publicKey)}
		fmt.Fprintln(view, green(loginArr[0]+" logged in"))
	case 3:
		fmt.Fprintln(view, green(incomingMessage.GetMessage()+"logged out"))
		delete(connectedUsers, incomingMessage.GetMessage())
	case 4:
		logger.Println("it's a 4")
	case 6:
		// First, check if they passed or failed the decrypt test and then finish the login
		// process.
		if incomingMessage.GetMessage() == "error" {
			fmt.Fprintln(view, red("There was a problem logging in, please try again. Check'"+logLocation+"'for more information."))
			logger.Println("Got an error message when trying to log in.")
			return
		}
		if err := logIn(incomingMessage); err != nil {
			logger.Println(err)
		}
	case 7:
		decryptedMessage, err := decryptor.Decrypt(string(incomingMessage.GetMessage()))
		if err != nil && err != io.EOF {
			logger.Println(err)
		}
		fmt.Fprintln(view, purple("->"+" "+incomingMessage.GetSender()+":"+decryptedMessage))

	}
}

// Refrenced in main.go under many places.
// Takes the variables for a message to the server and construcs the protobuf and sends it through
// the connection.
func sendMessage(messageType int32, recipient, message string) error {
	newMessage := &pb.Message{
		MessageType: proto.Int32(messageType),
		Sender:      proto.String(currentUser),
		Recipient:   proto.String(recipient),
		Message:     proto.String(message),
	}

	messageByte, err := proto.Marshal(newMessage)
	if err != nil {
		return err
	}
	if _, err := conn.Write(messageByte); err != nil {
		return err
	}
	return nil
}

// Refrenced in main.go under main
// Run before checkMessageLoop to make sure that the user can log into the server
func preLogIn() error {
	if err := sendMessage(1, "server", publicKey); err != nil {
		return err
	}
	return nil
}

// Refrenced in main.go under parseMessage
func logIn(incomingMessage *pb.Message) error {
	decryptedMessage, err := decryptor.Decrypt(incomingMessage.GetMessage())
	if err != nil && err != io.EOF {
		return err
	}
	if err := sendMessage(6, "server", decryptedMessage); err != nil {
		return err
	}
	return nil
}

// Refrenced in main.go under init
// Use the decryptor to figure out if password used is correct
func testPassword() error {
	message, err := gpg.Encrypt("test", publicKey)
	if err != nil {
		return err
	}

	decryptedMessage, err := decryptor.Decrypt(message)
	if err != nil {
		return err
	}
	if decryptedMessage != "test" {
		return errors.New("The password is bad")
	}
	return nil
}

func StartClient(host, port, pubKeyLoc, prvKeyLoc, logLoc string) error {
	// Set function argument variables to public variables
	serviceHost = host
	servicePort = port

	if err := initPubKey(pubKeyLoc); err != nil {
		fmt.Println("Error loading public key")
		return err
	}

	if err := initPrivKey(prvKeyLoc); err != nil {
		fmt.Println("Error loading private key")
		return err
	}

	logLocation = logLoc

	if err := configureLogger(); err != nil {
		fmt.Println("Internal Error")
		return err
	}

	defer logFile.Close()
	// Dial server to get a net.Conn object and to make sure that the host is up

	var err error
	conn, err = net.Dial("tcp", serviceHost+":"+servicePort)
	if err != nil {
		fmt.Println("Server is down")
		return err
	}
	defer conn.Close()

	g := gocui.NewGui()
	if err := g.Init(); err != nil {
		return err
	}
	defer g.Close()
	g.SetLayout(gocuiLayout)
	if err := keybindings(g, conn); err != nil {
		return err
	}

	g.Execute(func(g *gocui.Gui) error {
		v, err := g.View("messages_box")
		if err != nil {
			return err
		}
		// Log in user, the login process is explained better in the function
		if err := preLogIn(); err != nil {
			logger.Println(err)
		}
		// Process check message loop in background (goroutine)
		go checkMessageLoop(v)

		return nil
	})

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		return err
	}

	return nil
}
