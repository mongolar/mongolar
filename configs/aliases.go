package configs

import (
	"fmt"
)

// Simple map for the aliases
type Aliases map[string]string

// Builder function takes a SiteMap and returns the Aliases Map
func NewAliases(sm SitesMap) Aliases {
	a := make(Aliases)
	a.setAliases(&sm)
	return a
}

// Itterate over SiteMap and map all domains to their configs
func (a Aliases) setAliases(sm *SitesMap) {
	for k, s := range *sm {
		//TODO check Aliases length, if 0 fatal error
		for _, alias := range s.Aliases {
			fmt.Printf("Mapping domain  %v to sit configuration %v\n", alias, k)
			a[alias] = k
		}
	}
}
