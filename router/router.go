package router

import (
	//"github.com/davecgh/go-spew/spew"
	"github.com/jasonrichardsmith/mongolar/configs/site"
	"net/http"
	"text/template"
)

const (
	MongolarScript = `mongular.constant('mongularConfig', {
    			mongular_url: '/{{ .APIEndPoint }}/',
    			templates_url: '/{{ .TemplateEndpoint }}'}
		);
		mongular.config(function($sceDelegateProvider) {
   			$sceDelegateProvider.resourceUrlWhitelist([
     				// Allow same origin resource loads.
     				'self',
				{{ range .ForeignDomains }}
     				'http://{{.}}/**',
				{{ end }}
   			]);
 		});`
)

type Router struct {
	Aliases     map[string]string
	Sites       map[string]*site.SiteConfig
	Controllers map[string]*controller.Controller
	APIEndPoint string
}

func NewRouter() {
	//TODO: build the constructor
}

type MongolarJSConfig struct {
	APIEndPoint      string
	TemplateEndpoint string
	ForeignDomains   []string
}

func (ro Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if val, ok := ro.Alias[r.Host]; ok {
		urlpath = split("\\", r.URL.Path)
		pathvalues := make([]string, len(urlpath))
		i := 0
		for k := range mymap {
			pathvalues[i] = k
			i++
		}
		s := ro.Sites[r.Aliases[r.Host]]
		switch pathvalues[0] {
		case "mongolar_config.js":
			c := MongolarJSConfig{APIEndPoint: ro.APIEndPoint, TemplateEndpoint: s.TemplateEndpoint, ForeignDomains: s.ForeignDomains}
			c.serveMongolarConfig(w)
		case "assets":
			directory := ro.Sites[r.Aliases[r.Host]].Directory
			http.FileServer(http.Dir(directory + "/assets"))

		case r.APIEndPoint:
			if val, ok := ro.Controllers[pathvalues[1]]; ok {
				content := ro.Controllers[pathvalues[1]].getContent(r, ro.Sites[r.Aliases[ro.Host]])
			} else {
				http.Error(w, "Forbidden", 403)
			}
		default:
			directory := ro.Sites[r.Aliases[ro.Host]].Directory
			http.ServeFile(w, r, directory+"/index.html")
		}

	} else {
		http.Error(w, "Not Found", 404) // Or Redirect?
	}

}

func (c *MongolarJSConfig) serveMongolarConfig(c MongolarJSConfig, w http.ResponseWriter) {
	t := template.New("Mongolar Config JS")
	t, err := t.Parse(MongolarScript)
	err = t.Execute(w, c)
}

func rand_str(size int) string {
	alphanum := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, size)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(bytes)
}
