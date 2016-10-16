package client

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/howeyc/gopass"
)

// Configure logger
func configureLogger() {
	var err error
	// Open the log file so that we can create the logger object
	logFile, err = os.OpenFile(logLocation, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Failed to open log file")
	}
	// Make sure the output is redirected to file instead of STDOUT
	//	log.SetOutput(logFile)
	logger = log.New(logFile, "logger: ", log.Ldate|log.Ltime|log.Llongfile)
}

// Currently reads both keys at once, might want to seperate if we're going to implement a login page
func initPubKey(publicKeyLocation string) {

	publicKeyByte, err := ioutil.ReadFile(publicKeyLocation)
	if err != nil {
		fmt.Println("Public key wasn't found")
		return err
	}
	publicKey = string(publicKeyByte)
}

func initPrivKey(privateKeyLocation string) {
	// Get the password for the private key/ id
	// Configure the package-wide variable 'decryptor'
	privateKeyByte, err := ioutil.ReadFile(privateKeyLocation)
	if err != nil {
		fmt.Println("Private Key wasn't found")
		return err
	}

	decryptor.PrivateKey = string(privateKeyByte)

	// List identities for user to choose
	if err := decryptor.GetEntities(); err != nil {
		fmt.Println("Internal error. Check key location")
		return err
	}
	// Choose which Identity to use
	scanner := bufio.NewReader(os.Stdin)
	fmt.Printf("Choose number: ")
	entityIndexStr, err := scanner.ReadString('\n')
	if err != nil {
		fmt.Println("Error getting input")
		return err
	}
	entityIndexStr = strings.TrimSpace(entityIndexStr)
	entityIndex, err := strconv.Atoi(entityIndexStr)
	if err != nil {
		fmt.Println("Internal error")
		return err
	}
	decryptor.IdentityIndex = entityIndex
	// Private keys usually have a password associated with them
	fmt.Printf("Private Key Password: ")
	privateKeyPass, err = gopass.GetPasswd()
	if err != nil {
		fmt.Println("Error getting password for private key")
		return err
	}
	// FUTURE: Private key is stored in memory, is this OK?
	// Pretty sure this is fine, openssh client (temporarily) stores ssh private key in a file at a user's homedir
	decryptor.Password = string(privateKeyPass)
	currentUser, _, _, err = decryptor.GetEntity()
	if err != nil {
		fmt.Println("Error getting entities")
		return err
	}
	// Make sure that password given by private key is correct
	if err := testPassword(); err != nil {
		fmt.Println("The password is incorrect")
		logger.Fatal("Password Failed:", err)
	}
	return nil
}
