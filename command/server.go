package command

import (
	"github.com/spf13/cobra"
	"github.com/t94j0/mwnn/server"
)

var port string

var ServerCmd = &cobra.Command{
	Use:   "server",
	Short: "Start MWNN server",
	Long:  "Start the MWNN server",
	Run: func(cmd *cobra.Command, args []string) {
		server.StartServer(port)
	},
}
