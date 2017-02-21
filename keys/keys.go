package keys

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/user"

	"github.com/robfig/config"
)

// GenerateKeyPair uses the gpg command to create a gpg key
func GenerateKeyPair() error {
	cmd := exec.Command("gpg", "--gen-key")
	err := cmd.Run()
	if err != nil {
		fmt.Println("There was an error generating the key")
	}
	return nil

}

// ListKeys uses the gpg command to list the keys you own
func ListKeys() {
	cmd := exec.Command("gpg", "--list-keys")
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		fmt.Println("There was an error listing the key(s)")
	}

}

// ChangeDefaultKey changes the default key that is used by mwnn
func ChangeDefaultKey(defKey string) error {
	usr, err := user.Current()
	if err != nil {
		fmt.Println("Error getting home directory")
		return err
	}
	home := usr.HomeDir

	c, err := config.ReadDefault(home + "/.mwnn/config")
	if err != nil {
		return err
	}

	pub, err := c.String("Keys", "Public-Key-Location")
	if err != nil {
		return err
	}
	pubFile, err := os.Create(pub)
	if err != nil {
		return err
	}
	pubWriter := bufio.NewWriter(pubFile)
	cmd := exec.Command("gpg", "--armor", "--export", defKey)
	cmd.Stdout = pubWriter
	err = cmd.Run()
	if err != nil {
		fmt.Println("There was an error changing the default public key")
		return err
	}
	if err := pubWriter.Flush(); err != nil {
		fmt.Println("There was an error changing the default public key")
		return err
	}
	if err := pubFile.Close(); err != nil {
		fmt.Println("There was an error changing the default public key")
		return err
	}

	priv, err := c.String("Keys", "Private-Key-Location")
	if err != nil {
		return err
	}
	privFile, err := os.Create(priv)
	if err != nil {
		return err
	}
	privWriter := bufio.NewWriter(privFile)
	cmd = exec.Command("gpg", "--armor", "--export-secret-key", defKey)
	cmd.Stdout = privWriter
	err = cmd.Run()
	if err != nil {
		fmt.Println("There was an error changing the default private key")
		return err
	}
	if err := privWriter.Flush(); err != nil {
		fmt.Println("There was an error changing the default private key")
		return err
	}
	if err := privFile.Close(); err != nil {
		fmt.Println("There was an error changing the default private key")
		return err
	}

	return nil

}
