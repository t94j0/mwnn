package gpg

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
)

/*
  This is https://gist.github.com/stuart-warren/93750a142d3de4e8fdd2, but repurposed to fit this project better
*/

//////////////////////
// Global Variables //
//////////////////////

// Initializers
var pathToPriv string
var password string

///////////
// Types //
///////////

type Decryptor struct {
	PrivateKey    string
	Password      string
	IdentityIndex int
}

func (d *Decryptor) GetEntities() error {
	var err error

	privateKey := bytes.NewBufferString(d.PrivateKey)

	privateKeyAsc, err := armor.Decode(privateKey)
	if err != nil {
		return err
	}

	entityList, err := openpgp.ReadKeyRing(privateKeyAsc.Body)
	if err != nil {
		return err
	}

	if len(entityList) == 0 {
		return errors.New("There are no keys!")
	}

	for i, entity := range entityList {
		for k, _ := range entity.Identities {
			//TODO: Just put this in an array and send it to main.go
			fmt.Println(i, k)
		}
	}
	return nil
}

func (d *Decryptor) GetEntity() (name string, entity *openpgp.Entity, entityList openpgp.EntityList, err error) {

	//TODO: This code is duplicated, but I can't find a way to fix this yet

	privateKey := bytes.NewBufferString(d.PrivateKey)

	privateKeyAsc, err := armor.Decode(privateKey)
	if err != nil {
		return "", &openpgp.Entity{}, openpgp.EntityList{}, err
	}

	// Open the private key file
	entityList, err = openpgp.ReadKeyRing(privateKeyAsc.Body)
	if err != nil {
		return "", &openpgp.Entity{}, openpgp.EntityList{}, err
	}

	entity = entityList[d.IdentityIndex]

	for entityName, _ := range entityList[d.IdentityIndex].Identities {
		name = entityName
	}
	re, err := regexp.Compile("(.+?)[<(]")
	if err != nil {
		fmt.Println(err)
	}
	//TODO: The ID must have an email or comment associated with it
	name = re.FindStringSubmatch(name)[1]
	name = strings.TrimSpace(name)
	return
}

func (d *Decryptor) Decrypt(message string) (string, error) {
	_, entity, entityList, err := d.GetEntity()
	if err != nil {
		fmt.Println(err)
	}

	// Get the passphrase and read the private key.
	passphraseByte := []byte(d.Password)
	entity.PrivateKey.Decrypt(passphraseByte)
	for _, subkey := range entity.Subkeys {
		subkey.PrivateKey.Decrypt(passphraseByte)
	}

	// Decode the base64 string
	dec, err := base64.StdEncoding.DecodeString(message)
	if err != nil {
		return "", err
	}

	// Decrypt it with the contents of the private key
	md, err := openpgp.ReadMessage(bytes.NewBuffer(dec), entityList, nil, nil)
	if err != nil {
		return "", err
	}
	bytes, err := ioutil.ReadAll(md.UnverifiedBody)
	if err != nil {
		return "", err
	}
	decStr := string(bytes)

	return decStr, nil
}

//////////////////////
// Public Functions //
//////////////////////

func Encrypt(secretString, publicKeyStr string) (string, error) {
	publicKey := bytes.NewBufferString(publicKeyStr)
	ascBlock, err := armor.Decode(publicKey)
	if err != nil {
		return "", err
	}
	entityList, err := openpgp.ReadKeyRing(ascBlock.Body)
	if err != nil {
		return "", err
	}

	// encrypt string
	buf := new(bytes.Buffer)
	w, err := openpgp.Encrypt(buf, entityList, nil, nil, nil)
	if err != nil {
		return "", err
	}

	_, err = w.Write([]byte(secretString))
	if err != nil {
		return "", err
	}
	err = w.Close()
	if err != nil {
		return "", err
	}

	bytes, err := ioutil.ReadAll(buf)
	if err != nil {
		return "", err
	}
	encString := base64.StdEncoding.EncodeToString(bytes)

	return encString, nil
}
