package keys

import (

	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"github.com/robfig/config"

)

// Generates a key pair given parameters
func GenerateKeyPair() error {

	cmd := exec.Command("gpg", "--gen-key")
	err := cmd.Run()
	if err != nil{
		fmt.Println("There was an error generating the key")
	}
	return nil

}

func ListKeys() {

	cmd := exec.Command("gpg", "--list-keys")
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil{
		fmt.Println("There was an error listing the key(s)")
	}

}

func ChangeDefaultKey(defKey string) {

	usr, err := user.Current()
    if err != nil {
        fmt.Println("Error getting home directory")
		return
    }
	home := usr.HomeDir

	c, _ := config.ReadDefault(home + "/.mwnn/config")

	pub, _ := c.String("Keys", "Public-Key-Location")
	pubFile, _ := os.Create(pub)
	pubWriter := bufio.NewWriter(pubFile)
	cmd := exec.Command("gpg", "--armor", "--export", defKey)
	cmd.Stdout = pubWriter
	err = cmd.Run()
	pubWriter.Flush()
	pubFile.Close()
	if err != nil{
		fmt.Println("There was an error changing the default public key")
		return
	}

	priv, _ := c.String("Keys", "Private-Key-Location")
	privFile, _ := os.Create(priv)
	privWriter := bufio.NewWriter(privFile)
	cmd = exec.Command("gpg", "--armor", "--export-secret-key", defKey)
	cmd.Stdout = privWriter
	err = cmd.Run()
	privWriter.Flush()
	privFile.Close()
	if err != nil{
		fmt.Println("There was an error changing the default private key")
		fmt.Println(err)
		return
	}

}
