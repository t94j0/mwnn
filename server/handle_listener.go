package server

import (
	"bytes"
	"fmt"
	"net"

	"github.com/dchest/uniuri"
	"github.com/golang/protobuf/proto"
	pb "github.com/t94j0/mwnn/message"
)

func handleListener(connections map[string]User, conn net.Conn) {
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
			if err.Error() != "EOF" {
				handleDebug(2, err.Error())
			}
		}

		// Create variable for unmarshalling buffer
		incomingMessage := &pb.Message{}

		// Remove whitespace from buffer
		buf = bytes.Trim(buf, "\x00")
		// Discard empty buffer
		if buf == nil {
			handleDebug(1, newUser.Username+" logged out")
			if err := sendToEveryone(connections, 3, newUser.Username); err != nil {
				handleDebug(2, err.Error())
			}
			break
		}

		// Unmarshal message into `incomingMessage`
		if err := proto.Unmarshal(buf, incomingMessage); err != nil {
			handleDebug(2, err.Error())
		}
		switch incomingMessage.GetMessageType() {
		// Message
		case 0:
			// We have to wrap `GetMessageType` in an int because it's an int32, not an int.
			// Kinda dumb, but whatever
			if err := sendMessage(connections, int(incomingMessage.GetMessageType()), incomingMessage.GetRecipient(), incomingMessage.GetSender(), incomingMessage.GetMessage()); err != nil {
				handleDebug(2, err.Error())
			}
		// Login
		case 1:
			newUser = preLoginUser(connections, id, incomingMessage.GetSender(), incomingMessage.GetMessage(), conn)
			handleDebug(1, newUser.Username+" initiated login")
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
			if err := loginUser(connections, incomingMessage.GetSender(), incomingMessage.GetMessage(), newUser.PublicKey, conn); err != nil {
				handleDebug(2, err.Error())
			}
			handleDebug(1, newUser.Username+" logged in")
		case 7:
			if err := sendMessage(connections, int(incomingMessage.GetMessageType()), incomingMessage.GetRecipient(), incomingMessage.GetSender(), incomingMessage.GetMessage()); err != nil {
				handleDebug(2, err.Error())
			}

		}
	}

	// Remove connection from connection stack just in case it didn't get cleaned up
	delete(connections, id)
	delete(connections, newUser.Username)
}
