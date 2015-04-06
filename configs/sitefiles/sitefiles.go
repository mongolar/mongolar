package sitefiles

import (
	"path/filepath"
	"strings"
)

const (
	SITES_DIRECTORY = "/etc/mongolar/enabled/"
)

type SiteFiles map[int]string

func New() SiteFiles {
	s := make(SiteFiles)
	s.getSiteConfigFiles()
	return s

}

//Get all enabled config file names
func (s SiteFiles) getSiteConfigFiles() {
	glob := SITES_DIRECTORY + "*.yaml"
	files, _ := filepath.Glob(glob)
	for key, value := range files {
		var filename string
		_, filename = filepath.Split(value)
		s[key] = strings.TrimSuffix(filename, ".yaml")
	}
}
