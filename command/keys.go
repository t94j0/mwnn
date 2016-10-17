package command

import (
	"github.com/spf13/cobra"
	"github.com/t94j0/mwnn/keys"
)

var (
	keysPublicKeyLocation  string
	keysPrivateKeyLocation string
	username               string
	email                  string
)

func init() {
	set := getConfig()

	KeysCmd.Flags().StringVarP(&keysPublicKeyLocation, "publickey", "u", set.pubKey, "Location of public key")
	KeysCmd.Flags().StringVarP(&keysPrivateKeyLocation, "privatekey", "r", set.privKey, "Location of private key")
	KeysCmd.Flags().StringVarP(&username, "name", "n", "", "Username of new GPG key")
	KeysCmd.Flags().StringVarP(&email, "email", "e", "", "Email of new GPG key")
}

var KeysCmd = &cobra.Command{
	Use:   "keys",
	Short: "Create and list GPG keys",
	Long:  `Manages opeerations that have to do with key creation and management`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := keys.GenerateKeyPair(keysPublicKeyLocation, keysPrivateKeyLocation, username, email); err != nil {
			return err
		}
		return nil
	},
}
