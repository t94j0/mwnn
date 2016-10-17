package command

import (
	"fmt"
	"os/user"

	"github.com/spf13/cobra"
	"github.com/t94j0/mwnn/client"
)

var (
	serviceHost        string
	servicePort        string
	publicKeyLocation  string
	privateKeyLocation string
	logLocation        string
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

	ClientCmd.Flags().StringVarP(&serviceHost, "host", "c", "localhost", "Server to connect to")
	ClientCmd.Flags().StringVarP(&servicePort, "port", "p", "8181", "Port to connect to when trying host")
	ClientCmd.Flags().StringVarP(&publicKeyLocation, "publickey", "u", set.pubKey, "Location of public key")
	ClientCmd.Flags().StringVarP(&privateKeyLocation, "privatekey", "r", set.privKey, "Location of private key")
	ClientCmd.Flags().StringVarP(&logLocation, "log", "l", HOME_DIR+"/.mwnn/log.txt", "Location of log file")
}

var ClientCmd = &cobra.Command{
	Use:   "mwnn",
	Short: "Start MWNN client",
	Long:  `Starts the MWNN client and connect to the specified server.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := client.StartClient(serviceHost, servicePort, publicKeyLocation, privateKeyLocation, logLocation); err != nil {
			return err
		}
		return nil
	},
}
