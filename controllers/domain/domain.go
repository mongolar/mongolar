package domain

import (
	"github.com/jasonrichardsmith/mongolar/configs/site"
	"github.com/jasonrichardsmith/mongolar/router"
	"net/http"
)

func Serve(w http.ResponseWriter, r *http.Request, s *site.SiteConfig) {
	p := router.URLtoMap(r.URL.Path)
	v = s.DomainPublic[p[1]]

	fmt.Fprintf(w, "Hello, %q")
}
