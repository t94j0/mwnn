package command

import (
	"os/user"

	"github.com/spf13/cobra"
	"../client"
)

var HOME_DIR string

var ClientCmd = &cobra.Command{
	Use:   "client [server]",
	Short: "Start MWNN client",
	Long:  `Starts the MWNN client and connect to the specified server.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var (
			serviceHost        string
			servicePort        string
			publicKeyLocation  string
			privateKeyLocation string
			logLocation        string
		)

		// Create variable to find home folder to find the ~/.mwnn folder
		currUsr, err := user.Current()
		if err != nil {
			return err
		}
		HOME_DIR = currUsr.HomeDir

		set := getConfig()

		serviceHost = args[0]
		cmd.Flags().StringVarP(&servicePort, "port", "p", "8080", "Port to connect to when trying host")
		cmd.Flags().StringVarP(&publicKeyLocation, "publickey", "u", set.pubKey, "Location of public key")
		cmd.Flags().StringVarP(&privateKeyLocation, "privatekey", "r", set.privKey, "Location of private key")
		cmd.Flags().StringVarP(&logLocation, "log", "l", HOME_DIR+"/.mwnn/log.txt", "Location of log file")
		client.StartClient(serviceHost, servicePort, publicKeyLocation, privateKeyLocation, logLocation)
		return nil
	},
}
