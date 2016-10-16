package command

import(
	"github.com/spf13/cobra"
	"os"
	"os/user"
)

var HOME_DIR string

func Execute() {
	var rootCmd = &cobra.Command{Use: "mwnn"}

	// Create variable to find home folder to find the ~/.mwnn folder
	currUsr, err := user.Current()
	if err != nil {
		os.Exit(1)
	}
	HOME_DIR = currUsr.HomeDir


	set := getConfig()

	ClientCmd.Flags().StringVarP(&servicePort, "port", "p", "8080", "Port to connect to when trying host")
	ClientCmd.Flags().StringVarP(&publicKeyLocation, "publickey", "u", set.pubKey, "Location of public key")
	ClientCmd.Flags().StringVarP(&privateKeyLocation, "privatekey", "r", set.privKey, "Location of private key")
	ClientCmd.Flags().StringVarP(&logLocation, "log", "l", HOME_DIR+"/.mwnn/log.txt", "Location of log file")

	ServerCmd.Flags().StringVarP(&port, "port", "p", "8080", "Set the port of the server")

	// I could convert this to a function:
	// https://github.com/spf13/hugo/blob/master/commands/hugo.go#L185
	rootCmd.AddCommand(ServerCmd, ClientCmd)
	if err := rootCmd.Execute(); err != nil {
		// TODO: Add logger event for this
		fmt.Println("Error", err)
	}
}
