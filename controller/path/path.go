package path

import (
	"github.com/jasonrichardsmith/mongolar/configs/site"
	"github.com/jasonrichardsmith/mongolar/controller"
	"net/http"
)
func Serve(w http.ResponseWriter, r *http.Request, s *site.SiteConfig) {
	v := make(map[string]string)
	v['test'] = "Test"
	c = content.New(v)
	c.Serve()
	return
}
