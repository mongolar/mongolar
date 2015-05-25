package configs

import (
	"github.com/spf13/viper"
	"log"
)

// Server Config, only port setting right now, but it will probably grow
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

// Build from config file
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
