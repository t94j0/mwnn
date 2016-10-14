package command

import "github.com/spf13/cobra"

func Execute() {
	var rootCmd = &cobra.Command{Use: "mwnn"}

	// I could convert this to a function:
	// https://github.com/spf13/hugo/blob/master/commands/hugo.go#L185
	rootCmd.AddCommand(ServerCmd, ClientCmd)
	rootCmd.Execute()
}
