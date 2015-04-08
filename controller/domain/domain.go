package domain

import (
	"github.com/jasonrichardsmith/mongolar/router"
	"github.com/jasonrichardsmith/mongolar/wrapper"
)

func Serve(w *wrapper.Wrapper) {
	p := router.URLtoMap(w.Request.URL.Path)
	v := make(map[string]string)
	v[p[1]] = w.SiteConfig.PublicValues[p[1]]
	c = w.setContent(v)
	w.Serve(w)
}
