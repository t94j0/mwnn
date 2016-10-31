package keys

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/howeyc/gopass"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	//	"golang.org/x/crypto/openpgp/packet"
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
	/*
		keyBuff := bytes.NewBuffer(nil)
		k, _ := packet.SerializeSymmetricKeyEncrypted(keyBuff, privateKeyPass, nil)
		symKeyBuff := bytes.NewBuffer(nil)
		packWrite, _ := packet.SerializeSymmetricallyEncrypted(symKeyBuff, packet.CipherAES128, k, nil)
	*/
	// Create private key file
	privateKeyFile, err := os.Create(privKeyLoc)
	if err != nil {
		return err
	}
	armoredBuff := bytes.NewBuffer(nil)
	finalPrivKey := bytes.NewBuffer(nil)
	newPair.SerializePrivate(finalPrivKey, nil)

	w, err := armor.Encode(armoredBuff, openpgp.PrivateKeyType, nil)
	if err != nil {
		fmt.Printf("%v", err)
		return nil
	}
	packWrite, err := openpgp.SymmetricallyEncrypt(w, privateKeyPass, nil, nil)

	//	privateKeyBuff := make([]byte, 10000)
	//	finalPrivKey := bytes.NewBuffer(nil)
	packWrite.Write(finalPrivKey.Bytes())
	packWrite.Close()
	//	finalPrivKey := bytes.NewBuffer(privateKeyBuff)
	//	_, _ = w.Write(symKeyBuff.Bytes())
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
