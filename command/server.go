package command

import (
	"github.com/spf13/cobra"
	"github.com/t94j0/mwnn/server"
)

var port string

func init() {
	ServerCmd.Flags().StringVarP(&port, "port", "p", "8181", "Set the port of the server")
}

var ServerCmd = &cobra.Command{
	Use:   "server",
	Short: "Start MWNN server",
	Long:  "Start the MWNN server",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := server.StartServer(port); err != nil {
			return err
		}
		return nil
	},
}
