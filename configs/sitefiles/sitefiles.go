// Sites files loads all the configs from the directory /etc/mongolar/enabled
// and loads them in a map.

package sitefiles

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"
)

//  Where configs are stored
const (
	SITES_DIRECTORY = "/etc/mongolar/enabled/"
)

// The map to load site files
type SiteFiles map[int]string

// The builder for site files
func New() SiteFiles {
	s := make(SiteFiles)
	s.getSiteConfigFiles()
	return s

}

//Get all enabled config file names
func (s SiteFiles) getSiteConfigFiles() {
	glob := SITES_DIRECTORY + "*.yaml"
	files, err := filepath.Glob(glob)
	if err != nil {
		log.Fatal(err)
	}
	for key, value := range files {
		var filename string
		fmt.Printf("Found configuration file %v\n", value)
		_, filename = filepath.Split(value)
		s[key] = strings.TrimSuffix(filename, ".yaml")
	}
}
