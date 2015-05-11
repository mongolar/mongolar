package router

import (
	"github.com/mongolar/mongolar/configs"
	"github.com/mongolar/mongolar/controller"
	"github.com/mongolar/mongolar/router/jsconfig"
	"github.com/mongolar/mongolar/url"
	"github.com/mongolar/mongolar/wrapper"
	"net/http"
	"sort"
)

// The Router should have everything needed to server multiple sites from one go instance
// Aliases will have all domain aliases with the key for the site configuration
// Sites will have all the individual configurations with their key that relates to a Alias
// APIEndPoint is a random string that generates each time a server boots and defines
// where all API calls will take place.
type Router struct {
	Aliases     configs.Aliases
	Sites       configs.SitesMap
	Controllers controller.ControllerMap
}

// The Constructor for the Router structure
func New(a configs.Aliases, s configs.SitesMap, c controller.ControllerMap) *Router {
	r := new(Router)
	r.Aliases = a
	r.Sites = s
	r.Controllers = c
	return r
}

// The Serve HTTP method to qualify as a handler interface.
func (ro Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// Does domain exist
	if d, ok := ro.Aliases[r.Host]; ok {

		pathvalues := url.UrlToMap(r.URL.Path)

		// Set the the site config to an easy to use value.
		s := ro.Sites[d]
		switch pathvalues[0] {
		// Mongolar config js is generated dynamically because it gets passed values from site config and endpoint is variable
		// TODO move this to a controller
		case "mongolar_config.js":
			c := jsconfig.JsConfigs{
				APIEndPoint:      s.APIEndPoint,
				TemplateEndpoint: s.TemplateEndpoint,
				ForeignDomains:   s.ForeignDomains,
				AngularModules:   s.AngularModules,
			}
			c.Serve(w)

		// All static assets bypass AngularJS and get served as files.
		// TODO Move this to a controller
		case "assets":
			d := s.Directory
			http.ServeFile(w, r, d+"/"+r.URL.Path[1:])

		// If path is ApiEndPoint this is an API request.
		case s.APIEndPoint:
			i := sort.SearchStrings(s.Controllers, pathvalues[1])
			if s.Controllers[i] == pathvalues[1] {
				w.Header().Set("Content-Type", "application/json")
				// Build a wrapper for the controller
				wr := wrapper.New(w, r, s)
				//If the controller exists call it
				if c, ok := ro.Controllers[pathvalues[1]]; ok {
					c(wr)
					return
				} else {
					http.Error(w, "Forbidden", 403)
					return
				}
			} else {
				http.Error(w, "Forbidden", 403)
				return
			}

		// All other traffic will be handled by the AngularJs router
		default:
			d := s.Directory
			http.ServeFile(w, r, d)
			return
		}

	} else {
		// Domain was not found
		http.Error(w, "Not Found", 404) // Or Redirect?
	}
	return

}
