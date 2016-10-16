package server

import (
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/jmcvetta/randutil"
	"github.com/t94j0/mwnn/gpg"
)

// To understand how logging in is implimented, go to `docs/design.md` and look for the login process

// Sets up user object to be authenticated
// Make skeleton user with a random password to set up connection for login process
func preLoginUser(connections map[string]User, id, sender, publicKey string, c net.Conn) User {
	// Create user object to build user from
	newUser := User{c, sender, publicKey, ""}

	// Create password for the user
	newPassword, err := randutil.String(10, randutil.Alphanumeric)
	if err != nil {
		handleDebug(2, err.Error())
	}
	newUser.TestPassword = newPassword

	// Upgrade user from id to username
	connections[sender] = newUser
	delete(connections, id)
	// Encrypt password and send to specified user
	encryptedPassword, err := gpg.Encrypt(newUser.TestPassword, publicKey)
	if err != nil {
		fmt.Println(err)
		if err := sendMessage(connections, 6, sender, "server", "error"); err != nil {
			handleDebug(2, err.Error())
		}
	}
	if err := sendMessage(connections, 6, sender, "server", encryptedPassword); err != nil {
		handleDebug(2, err.Error())
	}

	// TODO return tuple with error and results for error handling
	return newUser
}

func loginUser(connections map[string]User, sender, password, publicKey string, c net.Conn) error {
	//TODO: Check here for the same user being logged in at the same time

	// TODO: If someone has pwned the server, would they be able to steal the password from
	// memory if we are putting it in clear text?
	// Yes, it can happen, but this can happen to any other desktop application. The solution
	// is to not get a memory dump of your machine.

	// Check if the password is the same as the one that was specified in `preLogin`
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
