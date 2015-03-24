package configs

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/viper"
	"path/filepath"
	"strings"
)

const (
	// Server Paths
	SERVER_CONFIG   = "/etc/mongular/"
	SITES_DIRECTORY = "/etc/mongular/enabled/"
)

// Individual Site Configuration Type
type MongolarSiteConfig struct {
	MongoDb   map[string]string
	Domain    string
	Directory string
	Aliases   map[string]string
}

// Wrapper type for entire server config structure
type MongolarSites struct {
	SiteConfigs map[string]MongolarSiteConfig
	Aliases     map[string]string
}

/*
Ignore this stuff for now
// Server Config
type MongolarServerConfig struct {
	Port int
}
*/
// Wrapper to load all site configurations
func (ms *MongolarSites) BuildMongolarSiteConfigs() {
	files := getSiteConfigFileNames()
	for _, value := range files {
		var site MongolarSiteConfig
		site.getSiteConfig(value)
		ms.SiteConfigs[value] = site
		spew.Dump(site)
		spew.Dump(value)
	}
}

/*
Ignore this stuff for now
// Set alias array
func (s *MongolarSites) getAliasesArray(file string) {
	for _, value := range s.SiteConfigs {
		s.Aliases[file] = value.Domain
		for _, alias := range value.Aliases {
			s.Aliases[file] = alias
		}
	}
}
*/
// Get server config

// Get one site configuration and assign it to the structure
func (s MongolarSiteConfig) getSiteConfig(filename string) {
	viper.SetConfigName(filename)
	viper.AddConfigPath(SITES_DIRECTORY)
	viper.ReadInConfig()
	viper.Marshal(&s)
}

/*
Ignore this stuff for now
func (s *MongolarServerConfig) getServerConfig() {
	viper.SetConfigName("mongolar")
	viper.AddConfigPath(SERVER_CONFIG)
	viper.ReadInConfig()
	viper.Marshal(&s)
}
*/

//Get all enabled config file names
func getSiteConfigFileNames() (files []string) {
	glob := SITES_DIRECTORY + "*.yaml"
	files, _ = filepath.Glob(glob)
	for key, value := range files {
		var filename string
		_, filename = filepath.Split(value)
		files[key] = strings.TrimSuffix(filename, ".yaml")
	}
	return files
}
