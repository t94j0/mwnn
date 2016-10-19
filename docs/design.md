# MWNN design docs
* Startup
	1. "Load configuration/parse parameters."
	2. "Enforce key existence"
	3. "Pass onto Login Process"

* Login Process
	1. "Pre-Login (server) / Login (client)"
		+ Client:
			1. The client gets their keys ready for encryption/decryption
			2. Upon successfully opening keys, the client sends a message with a type of 1 to the server
		+ Server:
			3. The server gets the message from the client and creates a local user object.
			4. The server attaches a password to the local user object and encodes it with the user's public key
			5. The server sends the encrypted password to the client as a message type of 6
	2. "Login (server) / Authenticating (client)"
		+ Client:
			1. The client decrypts the encrypted password given by the server and sends it to the server as a type of 6
		+ Server:
			2. The server validates the given password and adds the user to the list of active users.
			3. The server also sends everybody a message that a new user has connected, as well as letting the client know they're authenticated.

* Spawn Client Interface/Handlers
	1. "Client opens messenger interface"
		+ Client:
			1. Draws initial screen.
			2. Spawn goroutine to manage writing channel/user board. (ChanBox)
			3. Spawn goroutine to manage writing recieved messages area. (Inbox)
			4. Spawn goroutine to manage reading outgoing messages area. (Outbox)
	2. "Client Spawns Handlers"	
		+ Client:
			5. Spawn goroutine to manage routing incoming messages/commands/log(ons|outs). (InRouter)
			6. Spawn goroutine to manage routing outgoing messages/commands/logout. (OutRouter)

* Spawn Server Handlers
	1. "Server Spawns message Handler"
		+ Server:
			1. Spawn route handler. (ServRouter)
			2. Spawn one to all message handler. (OTAHandle)
			3. Spawn one to one message handler. (OTOHandle)
			4. Spawn serverside command handler. (CommHandle)

* Handler Specifications
	1. "ChanBox"
		+ Input(s):
			1. Input Channel for InRouter (Parameter)
			2. A Cui box for printing channel/user list in (Parameter)
			3. Login Messages (Runtime)
			4. Logout Messages (Runtime)
			5. Channel Move Messages (Runtime)
			6. Channel Creation/Deletion Messages (Runtime)
		+ Output(s):
			1. Text written to screen

	2. "Inbox"
		+ Input(s):
			1. Input Channel for InRouter (Parameter)
			2. A Cui box for printing messages in (Parameter)
			3. Chat Messages (Runtime)
			4. Command Messages (Runtime)
		+ Output(s):
			1. Text written to screen
			2. Command output (Possibly not text?)

	3. "Outbox"
		+ Input(s):
			1. Output Channel for OutRouter (Parameter)
			2. A Cui box for reading messages from (Parameter)
			3. Text from input device (Runtime)
		+ Output(s):
			1. Strings to OutRouter through Output Channel

	4. "InRouter"
		+ Input(s):
			1. Output Channel for ChanBox (Parameter)
			2. Output Channel for Inbox (Parameter)
			3. Input recieved from Server (Runtime)
		+ Output(s):
			1. Messages for ChanBox through Output Channel
			2. Messages for Inbox through Output Channel

	5. "OutRouter"
		+ Input(s):
			1. Input Channel for Outbox (Parameter)
			2. Messages to send to Server (Runtime)
		+ Output(s):
			1. Messages to Server

	6. "ServRouter"
		+ Input(s):
			1. Output Channel for OTAHandle (Parameter)
			2. Output Channel for OTOHandle (Parameter)
			3. Output Channel for CommHandle (Parameter)
			4. Messages from Client (Runtime)
		+ Output(s):
			1. Output for OTAHandle
			2. Output for OTOHandle
			3. Output for CommHandle
			4. Logs?

	7. "OTAHandle"
		+ Input(s):
			1. Input Channel for ServRouter (Parameter)
			2. Messages to push (Runtime)
		+ Output(s):
			1. Messages sent to all Clients

	8. "OTOHandle"
		+ Input(s):
			1. Input Channel for ServRouter (Parameter)
			2. Messages to push (Runtime)
		+ Output(s):
			1. Messages sent to one Client

	9. "CommHandle"
		+ Input(s):
			1. Input Channel for ServRouter (Parameter)
			2. Commands (Runtime)
		+ Output(s):
			1. Logs?
			2. Command dependant output

* Possible Future Features
	1. Plugin Folder
	2. Client-side customization (Colors, Message Pings, etc.)
	3. Admin Features (Mute, Ban, Kick, Channel Modification, Hide)
	4. PGP server integration
	5. Web interface? (rather ambitious)
