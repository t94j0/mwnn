package server

import (
	"fmt"
	"net"
)

// The user type gives the server minimal information while still making it easy to access all of
// the user information
type User struct {
	Connection   net.Conn
	Username     string
	PublicKey    string
	TestPassword string
}

var debugLevel int

// MessageType codes:
// 0: Messaging. Any messages that should be displayed should have a MessageType of 0
// 1: Pre-Auth. Send to server to start log-in process
// 2: Log in. Display to group that someone is logging in and has the public key in it
// 3: Log out. Display to group that someone is logging out. It tells the user to remove connection
// 4: Public Key. Change of public key between group
// 5: Public Key Sync. Gets all users in the current group.
// 6: Auth. Sends a test password to crack with private key
// 7: Private message. Sends a private message to a user

func StartServer(port string, debugLvl int) error {
	debugLevel = debugLvl

	// Create connection stack for managing all user connections
	var connections = make(map[string]User, 0)

	// FOR THE GLORY OF SATAN OF COURSE! -Max
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	handleDebug(1, "Server started on port "+port)

	for {
		// We accept all connections, but if we want to block users from certian ip addresses,
		// then we can do so here
		conn, err := listener.Accept()
		if err != nil {
			// We don't want to return this error because the server can still function without this
			fmt.Println(err)
		}
		handleDebug(1, "A user has connected")

		go handleListener(connections, conn)
	}
}

func handleDebug(level int, debug string) {
	if level <= debugLevel {
		fmt.Println(debug)
	}
}
