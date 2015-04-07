package router

import (
	//"github.com/davecgh/go-spew/spew"
	"github.com/jasonrichardsmith/mongolar/configs/site"
	"github.com/jasonrichardsmith/mongolar/controller"
	"github.com/jasonrichardsmith/mongolar/router/apiend"
	"github.com/jasonrichardsmith/mongolar/router/jsconfig"
	"net/http"
	"fmt"
)

// The Router should have everything needed to server multiple sites from one go instance
// Aliases will have all domain aliases with the key for the site configuration
// Sites will have all the individual configurations with their key that relates to a Alias
// APIEndPoint is a random string that generates each time a server boots
type Router struct {
	Aliases     map[string]string
	Sites       map[string]*site.SiteConfig
	Controllers map[string]*controller.Controller
	APIEndPoint apiend.APIEndPoint
	Session	    map[string]interface{}
}

// The Constructor for the Router structure
func New(a map[string]string, s map[string]*SiteConfig, c map[string]*controller.Controller, map[string]*) *Router {
	r = new(Router)
	r.Aliases = a
	r.Sites = s
	r.Controllers = c
	r.APIEndPoint = apiend.New()
	r.Session = 
	return r
}

// The Serve HTTP method to qualify as a handler interface.
func (ro Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// Does domain exist
	if val, ok := ro.Alias[r.Host]; ok {

		pathvalues := UrlToMap(r.URL.Path)

		// Set the the site config to an easy to use value.
		s := ro.Sites[r.Aliases[r.Host]]
		switch pathvalues[0] {

		// Mongolar config js is generated dynamically because it gets passed values from site config and endpoint is variable
		// TODO move this to a controller
		case "mongolar_config.js":
			c := jsconfig.JSConfigs{
				APIEndPoint: ro.APIEndPoint,
				TemplateEndpoint: s.TemplateEndpoint,
				ForeignDomains: s.ForeignDomains,
				s.AngularModules
				}
			c.Serve(w)
		// All static assets bypass AngularJS and get served as files.
		// TODO Move this to a controller
		case "assets":
			directory := s.Directory
			http.FileServer(http.Dir(directory + "/assets"))
		// If path is ApiEndPoint this is an API request.
		case r.APIEndPoint:
			w.Header().Set("Content-Type", "application/json")
			fmt.Print('[')
			if val, ok := ro.Controllers[pathvalues[1]]; ok {
				ro.Controllers[pathvalues[1]](r, w, s)
			} else {
				http.Error(w, "Forbidden", 403)
				return
			}

			fmt.Print(']')
		// All other traffic will be handled by the AngularJs router
		default:
			directory := ro.Sites[r.Aliases[ro.Host]].Directory
			http.ServeFile(w, r, directory+"/index.html")
		}

	} else {
		// Domain was not found
		http.Error(w, "Not Found", 404) // Or Redirect?
	}
	return

}

func UrlToMap( u string ) map[int]string {
	// Split the path values
	urlpath = split("\\", u)

	// Map the values as key store values
	pathvalues := make([]string, len(urlpath))
	i := 0
	for k := range mymap {
		pathvalues[i] = k
		i++
	}
	return pathvalues
}
