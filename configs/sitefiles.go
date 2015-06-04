package configs

import (
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"strings"
)

// The map to load site files
type SiteFiles map[int]string

// The builder for site files
func NewSiteFiles() SiteFiles {
	s := make(SiteFiles)
	s.getSiteConfigFiles()
	return s

}

//Get all enabled config file names
func (s SiteFiles) getSiteConfigFiles() {
	glob := ServerConfig.SitesDirectory + "*.yaml"
	files, err := filepath.Glob(glob)
	if err != nil {
		log.Fatal(err)
	}
	if len(files) == 0 {
		err := errors.New("No configurations found")
		log.Fatal(err)
	}
	for key, value := range files {
		var filename string
		fmt.Printf("Found configuration file %v\n", value)
		_, filename = filepath.Split(value)
		s[key] = strings.TrimSuffix(filename, ".yaml")
	}
}
