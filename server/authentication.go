package server

import (
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/t94j0/mwnn/gpg"
)

func loginUser(connections map[string]User, sender, password, publicKey string, c net.Conn) error {
	//TODO: Check here for the same user being logged in at the same time

	// TODO: If someone has pwned the server, would they be able to steal the password from
	// memory if we are putting it in clear text?

	if connections[sender].TestPassword != password {
		if err := sendMessage(connections, 1, "server", sender, "error"); err != nil {
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
	if err := sendMessage(connections, 1, sender, "server", userObject); err != nil {
		delete(connections, sender)
		return err
	}

	//TODO: Name can't have a "|" in it because of the split later
	// Send message to every connection with username and public key of user on appropriate
	// channel. The sender and public key are delimited by a "|"
	if err := sendToEveryone(connections, 2, sender+"|"+publicKey); err != nil {
		delete(connections, sender)
		return err
	}

	return nil
}

// Sets up user object to be authenticated
// Make skeleton user with a random password to set up connection for login process
func preLoginUser(connections map[string]User, id, sender, publicKey string, c net.Conn) User {
	newUser := User{c, sender, publicKey, ""}

	// TODO: Once I get wifi, use the randomstring module to generate a random password
	// TODO: Secure the transfer of the password, becuase right now it's being recieved in cleartext.
	newUser.TestPassword = "abc123"

	connections[sender] = newUser
	delete(connections, id)
	encryptedPassword, err := gpg.Encrypt(newUser.TestPassword, publicKey)
	if err != nil {
		fmt.Println(err)
		if err := sendMessage(connections, 6, sender, "server", "error"); err != nil {
			fmt.Println(err)
		}
	}
	if err := sendMessage(connections, 6, sender, "server", encryptedPassword); err != nil {
		fmt.Println(err)
	}
	// TODO return tuple with error and results for error handling
	return newUser
}
