// Server config loads all the universal server settings so we can load site configurations,
// write logs, and bind to the correct port.

package configs

import (
	"github.com/spf13/viper"
	"log"
)

var ServerConfigDirectory string

// Server Config, Includes some base configurations for the Server
// 	Port: Which port should be served
// 	SitesDirectory: Where the individual sites folder is located
// 	Log Directory: The directory where logs are stored
type Server struct {
	Port           string
	SitesDirectory string
	LogDirectory   string
}

// Constructor for server config
func NewServer() *Server {
	s := new(Server)
	s.getServerConfig()
	return s
}

// Marshall server config from Yaml file
func (s *Server) getServerConfig() {
	viper.SetConfigName("mongolar")
	viper.AddConfigPath(ServerConfigDirectory)
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	err = viper.Marshal(s)
	if err != nil {
		log.Fatal(err)
	}
	if s.SitesDirectory == "" {
		log.Fatal("No Mongolar sites directory set.")
	}
	if s.LogDirectory == "" {
		log.Fatal("No Mongolar sites directory set.")
	}
}
