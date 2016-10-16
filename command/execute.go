package command

import(
	"github.com/spf13/cobra"
	"os"
	"os/user"
)

var HOME_DIR string

func Execute() {
	var rootCmd = &cobra.Command{Use: "mwnn"}

	// I could convert this to a function:
	// https://github.com/spf13/hugo/blob/master/commands/hugo.go#L185
	rootCmd.AddCommand(ServerCmd, ClientCmd)
	if err := rootCmd.Execute(); err != nil {
		// TODO: Add logger event for this
		fmt.Println("Error", err)
	}
}
