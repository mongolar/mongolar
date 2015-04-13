package aliases

import (
	"github.com/jasonrichardsmith/mongolar/configs/sites"
)

type Aliases map[string]string

func New(sm sites.SitesMap) Aliases {
	a := make(Aliases)
	a.setAliases(&sm)
	return a
}

// Set alias array
func (a Aliases) setAliases(sm *sites.SitesMap) {
	for k, s := range *sm {
		for _, alias := range s.Aliases {
			a[alias] = k
		}
	}
}
