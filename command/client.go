package command

import (
	"fmt"
	"os/user"

	"github.com/spf13/cobra"
	"github.com/t94j0/mwnn/client"
)

var HOME_DIR string

var (
	ServiceHost        string
	ServicePort        string
	PublicKeyLocation  string
	PrivateKeyLocation string
	LogLocation        string
)

func init() {
	// Create variable to find home folder to find the ~/.mwnn folder
	currUsr, err := user.Current()
	if err != nil {
		// TODO: Create logger event for this
		fmt.Println(err)
	}
	HOME_DIR = currUsr.HomeDir

	set := getConfig()

	ClientCmd.Flags().StringVarP(&ServicePort, "port", "p", "8080", "Port to connect to when trying host")
	ClientCmd.Flags().StringVarP(&PublicKeyLocation, "publickey", "u", set.pubKey, "Location of public key")
	ClientCmd.Flags().StringVarP(&PrivateKeyLocation, "privatekey", "r", set.privKey, "Location of private key")
	ClientCmd.Flags().StringVarP(&LogLocation, "log", "l", HOME_DIR+"/.mwnn/log.txt", "Location of log file")
}

var ClientCmd = &cobra.Command{
	Use:   "client [server]",
	Short: "Start MWNN client",
	Long:  `Starts the MWNN client and connect to the specified server.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		ServiceHost = args[0]

		client.StartClient(ServiceHost, ServicePort, PublicKeyLocation, PrivateKeyLocation, LogLocation)
		return nil
	},
}
