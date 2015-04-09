package domain

import (
	"github.com/jasonrichardsmith/mongolar/router"
	"github.com/jasonrichardsmith/mongolar/wrapper"
)

func Serve(w *wrapper.Wrapper) {
	p := router.UrlToMap(w.Request.URL.Path)
	v := make(map[string]interface{})
	v[p[1]] = w.SiteConfig.PublicValues[p[1]]
	w.SetContent(v)
	w.Serve()
}
