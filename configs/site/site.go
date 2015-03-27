package site

import (
	"github.com/spf13/viper"
)

const (
	SITES_DIRECTORY = "/etc/mongolar/enabled/"
)

// Individual Site Configuration Type
type MongolarSiteConfig struct {
	MongoDb   map[string]string
	Domain    string
	Directory string
	Aliases   []string
}

func NewMongolarSiteConfig(file string) *MongolarSiteConfig {
	s := MongolarSiteConfig{
		MongoDb: make(map[string]string),
	}
	s.getSiteConfig(file)
	return &s
}

// Get one site configuration and marshall it
func (s *MongolarSiteConfig) getSiteConfig(file string) {
	v := viper.New()
	v.SetConfigName(file)
	v.AddConfigPath(SITES_DIRECTORY)
	v.ReadInConfig()
	v.Marshal(s)
}
