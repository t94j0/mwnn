package keys

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/howeyc/gopass"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
)

var RecievedBadName = errors.New("Names cannot contain \"()<>|\x00\"")
var RecievedBadEmail = errors.New("Email cannot contain \"()<>|\x00\"")
var RecievedBadComment = errors.New("Comment cannot contain \"()<>|\x00\"")
var PasswordMismatch = errors.New("Input passwords do not match")

func generateKeyPair(pubKeyLoc, privKeyLoc, name, comment, email string) error {

	if strings.ContainsAny(name, "()<>|\x00") {
		return RecievedBadName
	}
	if strings.ContainsAny(email, "()<>|\x00") {
		return RecievedBadEmail
	}
	if strings.ContainsAny(comment, "()<>|\x00") {
		return RecievedBadComment
	}
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
	newPair, err := openpgp.NewEntity(name, comment, email, nil)
	if err != nil {
		return err
	}
	// TODO Encode keys with armor.Encode() and write to files
	for _, identity := range newPair.Identities {
		if err := identity.SelfSignature.SignUserId(identity.UserId.Id, newPair.PrimaryKey, newPair.PrivateKey, nil); err != nil {
			return err
		}
	}

	// Create public key file
	publicKeyFile, err := os.Create(pubKeyLoc)
	if err != nil {
		return err
	}
	armor.Encode(publicKeyFile, openpgp.PublicKeyType, nil)
	defer publicKeyFile.Close()
	newPair.Serialize(publicKeyFile)

	// Create private key file
	privateKeyFile, err := os.Create(privKeyLoc)
	if err != nil {
		return err
	}
	armor.Encode(privateKeyFile, openpgp.PrivateKeyType, nil)
	defer privateKeyFile.Close()
	newPair.Serialize(privateKeyFile)

	return nil
}

func editKeyName() {

}

func editKeyComment() {

}

func changeDefaultKey() {

}
