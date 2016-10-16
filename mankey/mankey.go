package mankey

import(
	"fmt"
	"os"
	"strings"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
)

var RecievedBadName = errors.New("Names cannot contain \"()<>|\x00\"")
var RecievedBadEmail = errors.New("Email cannot contain \"()<>|\x00\"")
var RecievedBadComment = errors.New("Comment cannot contain \"()<>|\x00\"")

func generateKeyPair(pubKeyLoc, privKeyLoc, name, email, comment string) error {

	if strings.ContainsAny(name, "()<>|\x00") {
		return RecievedBadName
	}
	openpgp.NewEntity(
}

func editKeyName() {

}

func editKeyComment() {

}

func changeDefaultKey(){

}
