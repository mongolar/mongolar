package router

import (
	"github.com/davecgh/go-spew/spew"
	//"fmt"
	"github.com/jasonrichardsmith/mongolar/configs/aliases"
	"github.com/jasonrichardsmith/mongolar/configs/sites"
	"github.com/jasonrichardsmith/mongolar/controller"
	"github.com/jasonrichardsmith/mongolar/router/apiend"
	"github.com/jasonrichardsmith/mongolar/router/jsconfig"
	"github.com/jasonrichardsmith/mongolar/wrapper"
	"net/http"
	"strings"
)

// The Router should have everything needed to server multiple sites from one go instance
// Aliases will have all domain aliases with the key for the site configuration
// Sites will have all the individual configurations with their key that relates to a Alias
// APIEndPoint is a random string that generates each time a server boots
type Router struct {
	Aliases     aliases.Aliases
	Sites       sites.SitesMap
	Controllers controller.ControllerMap
	APIEndPoint string
}

// The Constructor for the Router structure
func New(a aliases.Aliases, s sites.SitesMap, c controller.ControllerMap) *Router {
	r := new(Router)
	r.Aliases = a
	r.Sites = s
	r.Controllers = c
	r.APIEndPoint = apiend.New()
	return r
}

// The Serve HTTP method to qualify as a handler interface.
func (ro Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// Does domain exist
	if d, ok := ro.Aliases[r.Host]; ok {
		spew.Dump(d)

		pathvalues := UrlToMap(r.URL.Path)

		// Set the the site config to an easy to use value.
		s := ro.Sites[d]
		switch pathvalues[1] {

		// Mongolar config js is generated dynamically because it gets passed values from site config and endpoint is variable
		// TODO move this to a controller
		case "mongolar_config.js":
			c := jsconfig.JsConfigs{
				APIEndPoint:      ro.APIEndPoint,
				TemplateEndpoint: s.TemplateEndpoint,
				ForeignDomains:   s.ForeignDomains,
				AngularModules:   s.AngularModules,
			}
			c.Serve(w)

		// All static assets bypass AngularJS and get served as files.
		// TODO Move this to a controller
		case "assets":
			directory := s.Directory
			http.FileServer(http.Dir(directory + "/assets"))

		// If path is ApiEndPoint this is an API request.
		case ro.APIEndPoint:
			w.Header().Set("Content-Type", "application/json")
			wr := wrapper.New(w, r, s)
			if c, ok := ro.Controllers[pathvalues[1]]; ok {
				c(wr)
			} else {
				http.Error(w, "Forbidden", 403)
				return
			}

		// All other traffic will be handled by the AngularJs router
		default:
			directory := s.Directory
			http.ServeFile(w, r, directory+"/index.html")
		}

	} else {
		// Domain was not found
		http.Error(w, "Not Found", 404) // Or Redirect?
	}
	return

}

func UrlToMap(u string) map[int]string {
	// Split the path values
	urlpath := strings.Split(u, "/")

	// Map the values as key store values
	pathvalues := make(map[int]string, len(urlpath))
	i := 0
	for _, k := range urlpath {
		pathvalues[i] = k
		i++
	}
	return pathvalues
}
