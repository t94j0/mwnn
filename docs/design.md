# MWNN design docs
* Login Process
	1. "Pre-Login"
		1. The client sends a message with a type of 1 to the server
		2. The server gets the message from the client and creates a local user object.
		3. The server attaches a password to the local user object and encodes it with the user's public key
		4. The server sends the encrypted password to the client as a message type of 6
	2. "Login"
		1. The client decrypts the encrypted password given by the server and sends it to the server as a type of 6
		2. The server validates the given password and adds the user to the list of active users.
		3. The server also sends everybody a message that a new user has connected
