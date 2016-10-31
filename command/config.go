package command

import (
	"github.com/robfig/config"
	"os"
)

type Setup struct {
	pubKey  string
	privKey string
}

// Grabs the configuration from the config file
func getConfig() Setup {

	con := config.NewDefault()
	if _, err := os.Stat(HOME_DIR + "/.mwnn/config"); os.IsNotExist(err) { // If there is no config file make one
		*con = generateConfig(*con)
	} else { // Else assert that it has the correct fields
		con, _ = config.ReadDefault(HOME_DIR + "/.mwnn/config")
		if con.HasSection("Keys") {
			if !con.HasOption("Keys", "Public-Key-Location") {
				con.AddOption("Keys", "Public-Key-Location", HOME_DIR+"/.mwnn/pub")
			}
			if !con.HasOption("Keys", "Private-Key-Location") {
				con.AddOption("Keys", "Private-Key-Location", HOME_DIR+"/.mwnn/priv")
			}
		} else {
			*con = generateConfig(*con)
		}
	}
	// No matter the case, we want to write the config file back to itself
	writeConfig(*con)

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
	return set

}

// Generates the default configuration for the keys
func generateConfig(con config.Config) config.Config {

	con.AddSection("Keys")
	con.AddOption("Keys", "Public-Key-Location", HOME_DIR+"/.mwnn/pub")
	con.AddOption("Keys", "Private-Key-Location", HOME_DIR+"/.mwnn/priv")
	return con

}

// Writes the configuration to the default file
func writeConfig(con config.Config) {

	var modePerm os.FileMode
	modePerm = 0777
	con.WriteFile(HOME_DIR+"/.mwnn/config", modePerm, "## MWNN CONFIG ##")

}
