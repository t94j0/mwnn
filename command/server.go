package command

import (
	"github.com/spf13/cobra"
	"github.com/t94j0/mwnn/server"
)

var (
	port       string
	debugLevel int
)

func init() {
	ServerCmd.Flags().StringVarP(&port, "port", "p", "8181", "Set the port of the server")
	ServerCmd.Flags().IntVarP(&debugLevel, "debug", "d", 1, "Set debug level. Level 0 - No debugging. Level 1 - Normal debugging. Level 2 - Print GPG messages.")
}

var ServerCmd = &cobra.Command{
	Use:   "server",
	Short: "Start MWNN server",
	Long:  "Start the MWNN server",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := server.StartServer(port, debugLevel); err != nil {
			return err
		}
		return nil
	},
}
