// This is the top level configuration loading for Server Config,
// Site Aliases, and Site Configurations.
package configs

import (
	"log"
	"os"
)

// Where the server configuration directory is set
var ServerConfig *Server

// Initialize some basic settings for configuration
func init() {
	// Check to see if environment variable is set
	ServerConfigDirectory = os.Getenv("MONGOLAR_SERVER_CONFIG")
	if ServerConfigDirectory == "" {
		ServerConfigDirectory = "/etc/mongolar/"
	}
	// Does directory exist
	info, err := os.Stat(ServerConfigDirectory)
	if err != nil {
		log.Fatal(err)
	}
	// Is it  a directory
	if info.IsDir() {
		return
	}
	log.Fatal(err)
}

// Wrapper structure for SitesMap and Aliases
type Configs struct {
	SitesMap SitesMap
	Aliases  Aliases
}

// Bootstrap all configurations
func New() (*Configs, string) {
	ServerConfig = NewServer()
	c := new(Configs)
	c.SitesMap = NewSitesMap()
	c.Aliases = NewAliases(c.SitesMap)
	return c, ServerConfig.Port
}
