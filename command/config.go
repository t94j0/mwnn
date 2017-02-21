package command

import (
	"os"

	"github.com/robfig/config"
)

// Setup holds the public key and private key in memory
type Setup struct {
	pubKey  string
	privKey string
}

// Grabs the configuration from the config file '~/.mwnn/config`
func getConfig() (Setup, error) {
	con := config.NewDefault()

	// Make sure the folder ~/.mwnn exists
	if _, err := os.Stat(HOME_DIR + "/.mwnn/"); os.IsNotExist(err) {
		if err := createFolder(HOME_DIR + "/.mwnn"); err != nil {
			return Setup{}, err
		}
	}

	// Make sure file ~/.mwnn/config is there
	if _, err := os.Stat(HOME_DIR + "/.mwnn/config"); os.IsNotExist(err) {
		*con = generateConfig(*con)
	} else {
		// Check the config file has the correct fields
		con, err = config.ReadDefault(HOME_DIR + "/.mwnn/config")
		if err != nil {
			return Setup{}, err
		}
		if con.HasSection("Keys") {
			if !con.HasOption("Keys", "Public-Key-Location") {
				con.AddOption("Keys", "Public-Key-Location",
					HOME_DIR+"/.mwnn/pub")
			}
			if !con.HasOption("Keys", "Private-Key-Location") {
				con.AddOption("Keys", "Private-Key-Location",
					HOME_DIR+"/.mwnn/priv")
			}
		} else {
			*con = generateConfig(*con)
		}
	}
	// Write the config file back to itself
	if err := writeConfig(*con); err != nil {
		return Setup{}, err
	}

	// Get public key and private key from config files
	var pub string
	var priv string
	set := Setup{pubKey: "", privKey: ""}
	pub, err := con.String("Keys", "Public-Key-Location")
	if err != nil {
		pub = HOME_DIR + "/.mwnn/pub"
	}
	err = nil
	priv, err = con.String("Keys", "Public-Key-Location")
	if err != nil {
		priv = HOME_DIR + "/.mwnn/priv"
	}
	set.pubKey, set.privKey = pub, priv

	return set, nil

}

// Generates the default configuration for the keys
func generateConfig(con config.Config) config.Config {
	con.AddSection("Keys")
	con.AddOption("Keys", "Public-Key-Location", HOME_DIR+"/.mwnn/pub")
	con.AddOption("Keys", "Private-Key-Location", HOME_DIR+"/.mwnn/priv")
	return con

}

// Writes the configuration to the default file
func writeConfig(con config.Config) error {
	var modePerm os.FileMode
	modePerm = 0600
	if err := con.WriteFile(HOME_DIR+"/.mwnn/config", modePerm, "## MWNN CONFIG ##"); err != nil {
		return err
	}
	return nil
}

func createFolder(path string) error {
	return os.Mkdir(path, 0700)
}
