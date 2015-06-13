package configs

// The definision of the sitemap
type SitesMap map[string]*SiteConfig

// Constructor that builds the sitemap.
func NewSitesMap() SitesMap {
	s := make(SitesMap)
	f := NewSiteFiles()
	s.getSiteConfigs(f)
	return s
}

// Builds all sites based off all found files
func (s SitesMap) getSiteConfigs(f []string) {
	for _, value := range f {
		site := NewSiteConfig(value)
		s[value] = site
	}
}
