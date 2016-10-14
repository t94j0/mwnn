package command

import (
	"os/user"

	"github.com/spf13/cobra"
	"github.com/t94j0/mwnn/client"
)

var ClientCmd = &cobra.Command{
	Use:   "client",
	Short: "Start MWNN client",
	Long:  "Start the MWNN client",
	RunE: func(cmd *cobra.Command, args []string) error {
		var (
			serviceHost        string
			servicePort        string
			publicKeyLocation  string
			privateKeyLocation string
			logLocation        string
		)

		currUsr, err := user.Current()
		if err != nil {
			return err
		}
		HOME_DIR := currUsr.HomeDir

		cmd.Flags().StringVarP(&serviceHost, "host", "c", "localhost", "Host to connect to")
		cmd.Flags().StringVarP(&servicePort, "port", "p", "6666", "Port to connect to when trying host")
		cmd.Flags().StringVarP(&publicKeyLocation, "publickey", "u", HOME_DIR+"/.mwnn/key.pub", "Location of public key")
		cmd.Flags().StringVarP(&privateKeyLocation, "privatekey", "r", HOME_DIR+"/.mwnn/key.prv", "Location of private key")
		cmd.Flags().StringVarP(&logLocation, "log", "l", HOME_DIR+"/.mwnn/log.txt", "Location of log file")
		client.StartClient(serviceHost, servicePort, publicKeyLocation, privateKeyLocation, logLocation)
		return nil
	},
}
