package server

import (
	"errors"

	"github.com/golang/protobuf/proto"
	pb "github.com/t94j0/mwnn/message"
)

// Sends message to user by giving an unparsed proto and the raw message
func sendMessage(connections map[string]User, messageType int, recipient, sender, message string) error {
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
func sendToEveryone(connections map[string]User, messageType int, message string) error {
	// loop over connections to reach everybody
	for _, user := range connections {
		if err := sendMessage(connections, messageType, user.Username, "server", message); err != nil {
			return err
		}
	}
	return nil
}
