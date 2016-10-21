package keys

import (
	"errors"
	"fmt"
	"os"
	"bytes"
	"strings"

	"github.com/howeyc/gopass"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
)

var RecievedBadName = errors.New("Names cannot contain \"()<>|\x00\"")
var RecievedBadEmail = errors.New("Email cannot contain \"()<>|\x00\"")
var PasswordMismatch = errors.New("Input passwords do not match")

//var RecievedBadComment = errors.New("Comment cannot contain \"()<>|\x00\"")

// Generates a key pair given parameters
func GenerateKeyPair(pubKeyLoc, privKeyLoc, name, email string) error {

	if strings.ContainsAny(name, "()<>|\x00") {
		return RecievedBadName
	}
	if strings.ContainsAny(email, "()<>|\x00") {
		return RecievedBadEmail
	}
	/*
		if strings.ContainsAny(comment, "()<>|\x00") {
			return RecievedBadComment
		}
	*/
	fmt.Printf("Private Key Password: ")
	privateKeyPass, err := gopass.GetPasswd()
	if err != nil {
		return err
	}
	fmt.Printf("Re-enter Private Key Password: ")
	privateKeyPass1, err := gopass.GetPasswd()
	if err != nil {
		return err
	}
	if string(privateKeyPass) != string(privateKeyPass1) {
		return PasswordMismatch
	}
	newPair, err := openpgp.NewEntity(name, " ", email, nil)
	if err != nil {
		return err
	}

	for _, identity := range newPair.Identities {
		if err := newPair.SignIdentity(identity.UserId.Id, newPair, nil); err != nil {
			return err
		}
	}

	// Create private key file
	privateKeyFile, err := os.Create(privKeyLoc)
	if err != nil {
		return err
	}
	armoredBuff := bytes.NewBuffer(nil)
	privateKeyBuff := bytes.NewBuffer(nil)
	newPair.SerializePrivate(privateKeyBuff, nil)
	w, err := armor.Encode(armoredBuff, openpgp.PrivateKeyType, nil)
	if err != nil {
		fmt.Printf("%v", err)
		return nil
	}
	plain, err := openpgp.SymmetricallyEncrypt(w, privateKeyPass, nil, nil)
	if err != nil{
		fmt.Printf(err.Error())
		return err
	}
	_, _ = plain.Write(privateKeyBuff.Bytes())
	plain.Close()
	w.Close()
	defer privateKeyFile.Close()
	privateKeyFile.Write(armoredBuff.Bytes())

	// Create public key file
	publicKeyFile, err := os.Create(pubKeyLoc)
	if err != nil {
		return err
	}
	armoredBuff = bytes.NewBuffer(nil)
	publicKeyBuff := bytes.NewBuffer(nil)
	newPair.Serialize(publicKeyBuff)
	w, err = armor.Encode(armoredBuff, openpgp.PublicKeyType, nil)
	if err != nil {
		fmt.Printf("%v", err)
		return nil
	}
	w.Write(publicKeyBuff.Bytes())
	w.Close()
	defer publicKeyFile.Close()
	publicKeyFile.Write(armoredBuff.Bytes())

	return nil
}

func changeDefaultKey() {

}
