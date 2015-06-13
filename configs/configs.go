// Configs is the main point of entry for configuring server config, site configs and aliases

package configs

import (
	"log"
	"os"
)

var ServerConfig *Server

func init() {
	ServerConfigDirectory = os.Getenv("MONGOLAR_SERVER_CONFIG")
	if ServerConfigDirectory == "" {
		ServerConfigDirectory = "/etc/mongolar/"
	}
	info, err := os.Stat(ServerConfigDirectory)
	if err != nil {
		log.Fatal(err)
	}
	if info.IsDir() {
		return
	}
	log.Fatal(err)
}

// Wrapper structure for SitesMapAliases
type Configs struct {
	SitesMap SitesMap
	Aliases  Aliases
}

// Constructor for Configs structure
func New() (*Configs, string) {
	ServerConfig = NewServer()
	c := new(Configs)
	c.SitesMap = NewSitesMap()
	c.Aliases = NewAliases(c.SitesMap)
	return c, ServerConfig.Port
}
