package server

import (
	"github.com/spf13/viper"
)

const (
	SERVER_CONFIG = "/etc/mongolar/"
)

// Server Config
type MongolarServerConfig struct {
	Port int
}

// Constructor
func NewMongolarServerConfig() *MongolarServerConfig {
	var s MongolarServerConfig
	s.getServerConfig()
	return &s
}

// Build from config
func (s *MongolarServerConfig) getServerConfig() {
	viper.SetConfigName("mongolar")
	viper.AddConfigPath(SERVER_CONFIG)
	viper.ReadInConfig()
	viper.Marshal(s)
}
