package configs

import (
	//"github.com/davecgh/go-spew/spew"
	"github.com/jasonrichardsmith/mongolar/configs/server"
	"github.com/jasonrichardsmith/mongolar/configs/site"
	"path/filepath"
	"strings"
)

const (
	SITES_DIRECTORY = "/etc/mongolar/enabled/"
	SERVER_CONFIG   = "/etc/mongolar/"
)

func GetAll() (*server.ServerConfig, map[string]*site.SiteConfig, map[string]string) {
	server := GetServerConfig()
	sites := GetSiteConfigs()
	aliases := GetAliasesArray(sites)
	return server, sites, aliases
}

func GetServerConfig() *server.ServerConfig {
	return server.New()
}

func GetSiteConfigs() map[string]*site.SiteConfig {
	var Sites = make(map[string]*site.SiteConfig)
	fs := getSiteConfigFileNames()
	for _, value := range fs {
		s := site.New(value)
		Sites[value] = s
	}
	return Sites
}

// Set alias array
func GetAliasesArray(ms map[string]*site.SiteConfig) map[string]string {
	a := make(map[string]string)
	for k, s := range ms {
		for _, alias := range s.Aliases {
			a[alias] = k
		}
	}
	return a
}

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
