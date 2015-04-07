package domain

import (
	"github.com/jasonrichardsmith/mongolar/configs/site"
	"github.com/jasonrichardsmith/mongolar/responders/content"
	"github.com/jasonrichardsmith/mongolar/router"
	"net/http"
)

func Serve(w http.ResponseWriter, r *http.Request, s *site.SiteConfig) {
	p := router.URLtoMap(r.URL.Path)
	v := make(map[string]string)
	v[p[1]] = s.PublicValues[p[1]]
	c = content.New(v)
	c.Serve(w)
}
