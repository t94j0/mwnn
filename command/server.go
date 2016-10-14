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
		cmd.Flags().StringVarP(&port, "port", "p", "6666", "Set the port of the server")
		if err := server.StartServer(port); err != nil {
			return err
		}
		return nil
	},
}
