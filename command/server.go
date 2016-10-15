package command

import (
	"github.com/spf13/cobra"
	"github.com/t94j0/mwnn/server"
)

var ServerCmd = &cobra.Command{
	Use:   "server",
	Short: "Start MWNN server",
	Long:  "Start the MWNN server",
	RunE: func(cmd *cobra.Command, args []string) error {
		var port string
		cmd.Flags().StringVarP(&port, "port", "p", "8080", "Set the port of the server")
		server.StartServer(port)
		return nil
	},
}
