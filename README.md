# Messenger With No Name [MWNN (Pronounced Minn)]

![](http://i.imgur.com/232MLcS.gif)

# What is MWNN?
IRC has never been a good tool for securely talking to other people because the data is sent in clear text over the internet. (Unless you are using SSL, but there are inherent (vulnerabilies with it)[http://www.howtogeek.com/182425/5-serious-problems-with-https-and-ssl-security-on-the-web/]). With MWNN, I am trying to create an irc-like service that uses GPG encryption to send text data across the internet safely so that there is no possibility for a MITM attack.

For more information on GPG, go [here](https://www.gnupg.org/)

# Installing

## Binary
**Not implimented yet**

## From Source

1. Download Go and set GOPATH
2. `go get github.com/t94j0/mwnn` 
3. `go get github.com/t94j0/gocui`

### Running the Server
1. Run `./server/build`
2. Run `./server/mwnnserver`

### Running the Client

#### Generating private key
1. gpg --gen-key
2. gpg --armor --export maxh@maxh.io > [location for public key]
3. gpg --armor --export-secret-key maxh@maxh.io > [location for private key]

#### Run mwnn
1. Use the `./client/build.sh` to make `./mwnn`
2. Run with `./mwnn -c [ip of server] -p [port] -puk [location of exported public key] -prk [location of exported private key]`

# Messenger
**Sending PM's**
* /message [target username] [message to send]


**Quitting**
* /quit
* Ctrl+c

**Viewing Public Keys**
/seekeys - This will display all connected user's keys. Always be sure that the keys you see here match what the other user's see. To exit, press Ctrl-q
