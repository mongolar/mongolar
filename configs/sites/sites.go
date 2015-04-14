// Sites just defines the site map and the constructor to make one.
package sites

import (
	"github.com/jasonrichardsmith/mongolar/configs/site"
	"github.com/jasonrichardsmith/mongolar/configs/sitefiles"
)

// The definision of the sitemap
type SitesMap map[string]*site.SiteConfig

// Constructor that builds the sitemap.
func New() SitesMap {
	s := make(SitesMap)
	f := sitefiles.New()
	s.getSiteConfigs(f)
	return s
}

// Builds all sites based off all found files
func (s SitesMap) getSiteConfigs(f sitefiles.SiteFiles) {
	for _, value := range f {
		site := site.New(value)
		s[value] = site
	}
}
