package server

import (
	"github.com/spf13/viper"
)

const (
	SERVER_CONFIG = "/etc/mongolar/"
)

// Server Config, only port setting right now, but it will probably grow
type Server struct {
	Port string
}

// Constructor
func New() *Server {
	s := new(Server)
	s.getServerConfig()
	return s
}

// Build from config file
func (s *Server) getServerConfig() {
	viper.SetConfigName("mongolar")
	viper.AddConfigPath(SERVER_CONFIG)
	viper.ReadInConfig()
	viper.Marshal(s)
}
