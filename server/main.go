package server

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/dchest/uniuri"
	"github.com/golang/protobuf/proto"
	"github.com/t94j0/mwnn/client/gpg"
	pb "github.com/t94j0/mwnn/client/message"
)

// Create connection stack for managing all user connections
var connections = make(map[string]User, 0)

// The user type gives the server minimal information while still making it easy to access all of
// the user information
type User struct {
	Connection   net.Conn
	Username     string
	PublicKey    string
	TestPassword string
}

// MessageType codes:
// 0: Messaging. Any messages that should be displayed should have a MessageType of 0
// 1: Pre-Auth. Send to server to start log-in process
// 2: Log in. Display to group that someone is logging in and has the public key in it
// 3: Log out. Display to group that someone is logging out. It tells the user to remove connection
// 4: Public Key. Change of public key between group
// 5: Public Key Sync. Gets all users in the current group.
// 6: Auth. Sends a test password to crack with private key
// 7: Private message. Sends a private message to a user

//////////////////////
// Helper functions //
//////////////////////

// Sends message to user by giving an unparsed proto and the raw message
func sendMessage(messageType int, recipient, sender, message string) error {
	newMessage := &pb.Message{
		MessageType: proto.Int(messageType),
		Recipient:   proto.String(recipient),
		Sender:      proto.String(sender),
		Message:     proto.String(message),
	}
	newMessageByte, err := proto.Marshal(newMessage)
	if err != nil {
		return err
	}
	// 0 is used as a delimiting byte. The client will trim \x00 from the right and then split
	// by \x00 to make sure that each message is delimited properly. This has to be done because
	// if two messages reach the client too quicky, then the user only parses one of the
	// messages
	newMessageByte = append(newMessageByte, 0)

	if _, isUser := connections[recipient]; !isUser {
		return errors.New("That is not a user")
		delete(connections, recipient)
	}

	if _, err := connections[recipient].Connection.Write(newMessageByte); err != nil {
		return err
	}
	return nil
}

// Write a message to every connection. The message shows it's from the server
func sendToEveryone(messageType int, message string) error {
	// loop over connections to reach everybody
	for _, user := range connections {
		if err := sendMessage(messageType, user.Username, "server", message); err != nil {
			return err
		}
	}
	return nil
}

// Sets up user object to be authenticated
// Make skeleton user with a random password to set up connection for login process
func preLoginUser(id, sender, publicKey string, c net.Conn) User {
	newUser := User{c, sender, publicKey, ""}

	// TODO: Once I get wifi, use the randomstring module to generate a random password
	// TODO: Secure the transfer of the password, becuase right now it's being recieved in cleartext.
	newUser.TestPassword = "abc123"

	connections[sender] = newUser
	delete(connections, id)
	encryptedPassword, err := gpg.Encrypt(newUser.TestPassword, publicKey)
	if err != nil {
		fmt.Println(err)
		if err := sendMessage(6, sender, "server", "error"); err != nil {
			fmt.Println(err)
		}
	}
	if err := sendMessage(6, sender, "server", encryptedPassword); err != nil {
		fmt.Println(err)
	}
	// TODO return tuple with error and results for error handling
	return newUser
}

func loginUser(sender, password, publicKey string, c net.Conn) error {
	//TODO: Check here for the same user being logged in at the same time

	// TODO: If someone has pwned the server, would they be able to steal the password from
	// memory if we are putting it in clear text?

	if connections[sender].TestPassword != password {
		if err := sendMessage(1, "server", sender, "error"); err != nil {
			return err
		}
		delete(connections, sender)
		return errors.New("User (" + sender + ") got the password wrong")
	}

	// Send public keys to users as the login message
	// usernames and public keys will be stored by uname|pk,uname|pk
	var userObject string
	for username, object := range connections {
		userObject += fmt.Sprintf("%s|%s,", username, object.PublicKey)
	}
	userObject = strings.TrimRight(userObject, ",")

	// Make sure that the data syncs with the user
	if err := sendMessage(1, sender, "server", userObject); err != nil {
		delete(connections, sender)
		return err
	}

	//TODO: Name can't have a "|" in it because of the split later
	// Send message to every connection with username and public key of user on appropriate
	// channel. The sender and public key are delimited by a "|"
	if err := sendToEveryone(2, sender+"|"+publicKey); err != nil {
		delete(connections, sender)
		return err
	}

	return nil
}

//////////////
// Handlers //
//////////////

func handleListener(conn net.Conn) {
	defer conn.Close()

	// Create temporary ID for new user since we haven't gotten login message yet
	id := uniuri.New()
	connections[id] = User{conn, id, "", ""}

	var newUser User

	// Listener loop
	for {
		// Create variable for the message to read into
		buf := make([]byte, 10000)
		// Read from connection into the buffer
		if _, err := conn.Read(buf); err != nil {
			// If the message is EOF, it means that they have disconnected, so send every user
			// a message about them logging out
			if err.Error() == "EOF" {
				fmt.Println(newUser.Username + " logged out")
				if err := sendToEveryone(3, newUser.Username); err != nil {
					fmt.Println(err)
					break
				}
			} else {
				fmt.Println(err)
				break
			}
		}

		// Create variable for unmarshalling buffer
		incomingMessage := &pb.Message{}

		// Remove whitespace from buffer
		buf = bytes.Trim(buf, "\x00")
		// Discard empty buffer
		if len(buf) == 0 {
			continue
		}

		// Unmarshal message into `incomingMessage`
		if err := proto.Unmarshal(buf, incomingMessage); err != nil {
			fmt.Println("Error unmarshaling message", err)
		}

		fmt.Println(incomingMessage.GetMessage())
		switch incomingMessage.GetMessageType() {
		// Message
		case 0:
			// We have to wrap `GetMessageType` in an int because it's an int32, not an int.
			// Kinda dumb, but whatever
			if err := sendMessage(int(incomingMessage.GetMessageType()), incomingMessage.GetRecipient(), incomingMessage.GetSender(), incomingMessage.GetMessage()); err != nil {
				fmt.Println("an error:", err)
			}
		// Login
		case 1:
			newUser = preLoginUser(id, incomingMessage.GetSender(), incomingMessage.GetMessage(), conn)
			fmt.Println(newUser.Username, "initiated login")
		// Server should be sending status 2, 3, 4, but not getting them
		case 2:
			fmt.Println("Server shouldn't be getting status 2")
		case 3:
			fmt.Println("Server shouldn't be getting status 3")
		case 4:
			fmt.Println("Server shouldn't be getting status 4")
		case 5:
			fmt.Println("Server shouldn't get getting status 5")
		case 6:
			if err := loginUser(incomingMessage.GetSender(), incomingMessage.GetMessage(), newUser.PublicKey, conn); err != nil {
				fmt.Println(err)
			}
			fmt.Println(newUser.Username, "logged in")
		case 7:
			if err := sendMessage(int(incomingMessage.GetMessageType()), incomingMessage.GetRecipient(), incomingMessage.GetSender(), incomingMessage.GetMessage()); err != nil {
				fmt.Println("an error:", err)
			}

		}
	}

	// Remove connection from connection stack just in case it didn't get cleaned up
	delete(connections, id)
	delete(connections, newUser.Username)
}

//////////////////////
// System Functions //
//////////////////////

func StartServer(port string) error {
	// FOR THE GLORY OF SATAN OF COURSE! -Max
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}
	fmt.Println("Started listener on port", port)
	for {
		// We accept all connections, but if we want to block users from certian ip addresses,
		// then we can do so here
		conn, err := listener.Accept()
		if err != nil {
			// We don't want to return this error because the server can still function without this
			fmt.Println(err)
		}
		fmt.Println("A user has connected")

		go handleListener(conn)
	}
}
