package command

import "fmt"

var HOME_DIR string

func Execute() {
	//var rootCmd = &cobra.Command{Use: "mwnn"}
	var rootCmd = ClientCmd
	// I could convert this to a function:
	// https://github.com/spf13/hugo/blob/master/commands/hugo.go#L185
	rootCmd.AddCommand(ServerCmd)
	if err := rootCmd.Execute(); err != nil {
		// TODO: Add logger event for this
		fmt.Println("Error", err)
	}
}
