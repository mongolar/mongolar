package url

import (
	"strings"
)

// Basic function to break a url to a map
func UrlToMap(u string) map[int]string {
	// Split the path values
	urlpath := strings.Split(u, "/")

	// Map the values as key store values
	pathvalues := make(map[int]string)
	i := 0
	for _, k := range urlpath {
		// The first value always evaluates to empty string so we can disregard
		if k != "" {
			pathvalues[i] = k
			i++
		}
	}
	return pathvalues
}
