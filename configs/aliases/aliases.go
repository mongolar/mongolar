// This builds a map of aliases that have a key that is the domain
// that references the site config that applies to it

package aliases

import (
	"fmt"
	"github.com/jasonrichardsmith/mongolar/configs/sites"
)

// Simple map for the aliases
type Aliases map[string]string

// Builder function takes a SiteMap and returns the Aliases Map
func New(sm sites.SitesMap) Aliases {
	a := make(Aliases)
	a.setAliases(&sm)
	return a
}

// Itterate over SiteMap and map all domains to their configs
func (a Aliases) setAliases(sm *sites.SitesMap) {
	for k, s := range *sm {
		//TODO check Aliases length, if 0 fatal error
		for _, alias := range s.Aliases {
			fmt.Printf("Mapping domain  %v to sit configuration %v\n", alias, k)
			a[alias] = k
		}
	}
}
