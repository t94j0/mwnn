package command

import (
	"github.com/spf13/cobra"
	"github.com/t94j0/mwnn/keys"
)

var (
	list   bool
	defKey string
)

func init() {
	KeysCmd.Flags().StringVarP(&defKey, "change-key", "c", "", " Change the default key by supplying the email associated with the key")
	KeysCmd.PersistentFlags().BoolVarP(&list, "list", "l", false, "List your available keys")
}

var KeysCmd = &cobra.Command{
	Use:   "keys",
	Short: "Create and list GPG keys",
	Long:  `Manages operations that have to do with key creation and management`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if list == false {
			if defKey == "" {
				if err := keys.GenerateKeyPair(); err != nil {
					return err
				}
				return nil
			} else {
				if err := keys.ChangeDefaultKey(defKey); err != nil {
					return err
				}
				return nil
			}
		} else {
			keys.ListKeys()
			return nil
		}
		return nil
	},
}
