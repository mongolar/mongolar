// The domain controller will return values in the PublicValues map of the site configuration.
package domain

import (
	"github.com/jasonrichardsmith/mongolar/router"
	"github.com/jasonrichardsmith/mongolar/wrapper"
)

// The controller func that will be invoked by the path being hit
func Serve(w *wrapper.Wrapper) {
	// Get second value in url path
	p := router.UrlToMap(w.Request.URL.Path)
	v := make(map[string]interface{})
	v[p[1]] = w.SiteConfig.PublicValues[p[1]]
	w.SetContent(v)
	w.Serve()
}
