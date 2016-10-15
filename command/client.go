package command

import (
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
var ClientCmd = &cobra.Command{
	Use:   "client [server] [# port] [# publickey] [# privatekey] [# log]",
	Short: "Start MWNN client",
	Long:  `Starts the MWNN client and connect to the specified server.`,
	Run: func(cmd *cobra.Command, args []string) {
		serviceHost = args[0]
		client.StartClient(serviceHost, servicePort, publicKeyLocation, privateKeyLocation, logLocation)
	},
}
