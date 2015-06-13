// Site Files

package configs

import (
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"strings"
)

// The builder for site files
func NewSiteFiles() []string {
	s := make([]string, 0)
	glob := ServerConfig.SitesDirectory + "*.yaml"
	files, err := filepath.Glob(glob)
	if err != nil {
		log.Fatal(err)
	}
	if len(files) == 0 {
		err := errors.New("No configurations found")
		log.Fatal(err)
	}
	for _, value := range files {
		var filename string
		fmt.Printf("Found configuration file %v\n", value)
		_, filename = filepath.Split(value)
		s = append(s, strings.TrimSuffix(filename, ".yaml"))
	}
	return s
}
