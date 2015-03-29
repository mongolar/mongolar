package server

import (
	"github.com/spf13/viper"
)

const (
	SERVER_CONFIG = "/etc/mongolar/"
)

// Server Config
type ServerConfig struct {
	Port string
}

// Constructor
func New() *ServerConfig {
	var s ServerConfig
	s.getServerConfig()
	return &s
}

// Build from config
func (s *ServerConfig) getServerConfig() {
	viper.SetConfigName("mongolar")
	viper.AddConfigPath(SERVER_CONFIG)
	viper.ReadInConfig()
	viper.Marshal(s)
}
