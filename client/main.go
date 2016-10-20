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
  1. Fix error with double printing names (blue/red) pretty sure it's server-side
  2. Make the message_box scrollable
  3. Change keybinding issue with seekeys
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

func inRouter(out []chan string) {
	for {
		// Create buffer to put recieved data in
		buf := make([]byte, 10000)
		// Read the data and throw error if it reached EOF. That means the server is down
		if _, err := conn.Read(buf); err != io.EOF && err != nil {
			logger.Println(err)
		}

		// The buffer is really big, so delete the whitespace
		buf = bytes.TrimRight(buf, "\x00")
		// Just continue if there is nothing but whitespace
		if buf == nil {
			logger.Panic("Connection to server broken")
			continue
		}

		messages := bytes.Split(buf, []byte{0})
		for _, message := range messages {
			handleMessage(message, out)
		}
	}
}

func handleMessage(buf []byte, out []chan string) {
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
			// Message gets sent to the Inbox
			out[0] <- incomingMessage.GetSender()+":"+ decryptedMessage
			break

		// Type = 1; Message is a login message
		case 1:
			if incomingMessage.GetMessage() != "error" {
				for _, userPair := range strings.Split(incomingMessage.GetMessage(), ",") {
					userArr := strings.Split(userPair, "|")
					username := userArr[0]
					usersPublicKey := userArr[1]
					connectedUsers[username] = User{username, []byte(usersPublicKey)}
					// Login message gets sent to Chanbox
					out[1] <- blue(username)+"\n"
				}
			} else {
				logger.Println("There was an error with the loging message")
			}
			break

		// Create empty user with no public key. If they do not share their public
		// key, they can't send messages to anybody
		case 2:
			loginArr := strings.Split(incomingMessage.GetMessage(), "|")
			publicKey := string(loginArr[1])
			connectedUsers[loginArr[0]] = User{loginArr[0], []byte(publicKey)}

			// Send messages to individual boxes
			out[0]<- green(loginArr[0]+" logged in.")
			out[1]<- red(loginArr[0])
			break

		// Someone logged out
		case 3:
			// Send appropriate logout messages
			out[0]<- green(incomingMessage.GetMessage()+"logged out")
			out[1]<- "Logout^:^"+incomingMessage.GetMessage()
			delete(connectedUsers, incomingMessage.GetMessage())
			break

		// 4's are one of the squarest numbers there are!
		case 4:
			logger.Println("It's a 4")
			break

		// Some Error occured
		case 6:
			// First, check if they passed or failed the decrypt test and then finish the login
			// process.
			if incomingMessage.GetMessage() == "error" {
				out[0]<- red("There was a problem logging in, please try again. Check'"+logLocation+"'for more information.")
				logger.Println("Got an error message when trying to log in.")
				return
			}
			if err := logIn(incomingMessage); err != nil {
				logger.Println(err)
			}
			break

		// Personal Message
		case 7:
			decryptedMessage, err := decryptor.Decrypt(string(incomingMessage.GetMessage()))
			if err != nil && err != io.EOF {
				logger.Println(err)
			}
			// Send PM to inbox
			out[0] <- purple("->"+" "+incomingMessage.GetSender()+":"+decryptedMessage)
			break
	}
}

func chanBox(view *gocui.View ,in <-chan string) {
	fmt.Fprintln(view, blue("\nOnline Users:"))
	for newM := range in {
		logout := strings.HasPrefix(newM,"Logout^:^")
		switch logout {
			// User is logging in, handle appropriatly
			case true:
				// Placeholder
				break

			// User is logging in
			case false:
				fmt.Fprintf(view, newM)
				break
		}
	}
}

func inBox(view *gocui.View ,in <-chan string) {
	for newM :=  range in {
	newM = strings.Trim(newM, "\n")
	newM = strings.Trim(newM, " ")
	fmt.Fprintln(view, newM)
	}
}

func outRouter(messageType int32, recipient, message string) error {
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

// Starts server authentication process
func preLogIn() error {
	if err := outRouter(1, "server", publicKey); err != nil {
		return err
	}
	return nil
}

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
func commandHandler(command string, g *gocui.Gui) error {
	command = strings.Trim(command, "/")
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
			if err := outRouter(7, user, messageEnc); err != nil {
				fmt.Fprintln(v, "There was a problem sending the message")
				return err
			}
			break
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
			break

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
	message := v.ViewBuffer()
	message = strings.Trim(message, "\n")
	message = strings.Trim(message, " ")

	// If the first letter is a forward slash, then we know that there is a command coming. If
	// it isn't, then it is a message that should be shown to everyone connected.
	if len(message) > 1 && string([]byte(message)[0]) == "/" {
		if err := commandHandler(message, g); err != nil {
			return err
		}
	} else {
		for _, user := range connectedUsers {
			messageEncrypt, err := gpg.Encrypt(message, string(user.PublicKey))
			if err != nil {
				logger.Println(err)
				fmt.Println(err)
			}
			if err := outRouter(0, user.Username, messageEncrypt); err != nil {
				logger.Println(err)
			}
		}
	}
	clearInput(g)
	return nil
}

// Handles any failed attempts to authenticate with the server
func logIn(incomingMessage *pb.Message) error {
	decryptedMessage, err := decryptor.Decrypt(incomingMessage.GetMessage())
	if err != nil && err != io.EOF {
		return err
	}
	if err := outRouter(6, "server", decryptedMessage); err != nil {
		return err
	}
	return nil
}

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

// Kicks off the client
func StartClient(host, port, pubKeyLoc, prvKeyLoc, logLoc string) error {
	// Set function argument variables to public variables
	serviceHost = host
	servicePort = port

	// Init the pubkey
	if err := initPubKey(pubKeyLoc); err != nil {
		fmt.Println("Error loading public key")
		return err
	}

	// Init the privkey
	if err := initPrivKey(prvKeyLoc); err != nil {
		fmt.Println("Error loading private key")
		return err
	}

	// Init the logger
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

	// Init the gui
	g := gocui.NewGui()
	if err := g.Init(); err != nil {
		return err
	}
	defer g.Close()
	g.SetLayout(gocuiLayout)
	if err := keybindings(g, conn); err != nil {
		return err
	}

	// Open communication channels between channels
	var commuChans = []chan string {
		// Inbox channel: index 0
		make(chan string),
		// Chanbox channel: index 1
		make(chan string),
	}

	// Start all gui functions
	g.Execute(func(g *gocui.Gui) error {

		// Grab the guis that will be written to
		inbox, err := g.View("messages_box")
		chanbox, err2 := g.View("channel_box")
		if err != nil || err2 != nil {
			return err
		}
		// Log in user, the login process is explained better in the function
		if err := preLogIn(); err != nil {
			logger.Println(err)
		}

		// Start the inRouter
		go inRouter(commuChans)
		// Start Inbox
		go inBox(inbox, commuChans[0])
		// Start Chanbox
		go chanBox(chanbox, commuChans[1])
		// Start Command Handler, it needs entire gui and its own channel
		// If we are on any view and the enter button is pressed, submit whats in the editbox buffer
		// to the server.
		// messageHandler is in main.go
		if err := g.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, messageHandler); err != nil {
			return err
		}
		// Start outRouter

		// Start Outbox

		return nil
	})

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		return err
	}

	return nil
}
