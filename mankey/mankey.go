package mankey

import(
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
	privateKeyPass, err = gopass.GetPasswd()
	if err != nil {
		return err
	}
	fmt.Printf("Re-enter Private Key Password: ")
	privateKeyPass1, err = gopass.GetPasswd()
	if err != nil {
		return err
	}
	if privateKeyPass != privateKeyPass1 {
		return PasswordMismatch
	}
	if newPair, err := openpgp.NewEntity(name, comment, email, nil); err != nil {
		return err
	}
	// TODO Encode keys with armor.Encode() and write to files
}

func editKeyName() {

}

func editKeyComment() {

}

func changeDefaultKey(){

}
