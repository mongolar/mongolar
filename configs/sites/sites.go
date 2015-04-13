package sites

import (
	"github.com/jasonrichardsmith/mongolar/configs/site"
	"github.com/jasonrichardsmith/mongolar/configs/sitefiles"
)

type SitesMap map[string]*site.SiteConfig

func New() SitesMap {
	s := make(SitesMap)
	f := sitefiles.New()
	s.getSiteConfigs(f)
	return s
}

func (s SitesMap) getSiteConfigs(f sitefiles.SiteFiles) {
	for _, value := range f {
		site := site.New(value)
		s[value] = site
	}
}
