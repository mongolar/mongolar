package main

import (
	"fmt"
	//	"github.com/davecgh/go-spew/spew"
	"github.com/jasonrichardsmith/mongolar/configs"
	"github.com/jasonrichardsmith/mongolar/configs/server"
	"github.com/jasonrichardsmith/mongolar/configs/site"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"log"
	"net/http"
)

var Server *server.MongolarServerConfig
var Sites map[string]*site.MongolarSiteConfig
var Aliases map[string]string

//var Routers map[string]*httprouter.Router
var Routers = make(map[string]*httprouter.Router)
var Switch = make(HostSwitch)

//var Switch HostSwitch

func main() {
	Server, Sites, Aliases = configs.GetAll()
	buildRouters()
	buildHostSwitch()
	log.Fatal(http.ListenAndServe(":"+Server.Port, Switch))

}

func buildRouters() {
	for key, site := range Sites {
		buildRouter(site, key)
	}

}

func buildRouter(s *site.MongolarSiteConfig, key string) {
	router := httprouter.New()
	//router.GET("/mongolar/path/*fullpath", MongolarPath)
	//router.GET("/*", IndexHtml)
	Routers[key] = router

}

func buildHostSwitch() {
	for domain, key := range Aliases {
		Switch[domain] = Routers[key]
	}
}

func MongolarPath(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "path, %s!\n", ps.ByName("fullpath"))
}

func IndexHtml(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	file, _ := ioutil.ReadFile(Sites[Aliases[r.Host]].Directory + "/index.html")
	fmt.Fprint(w, string(file))

}

type HostSwitch map[string]http.Handler

func (hs HostSwitch) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if handler := hs[r.Host]; handler != nil {
		handler.ServeHTTP(w, r)
	} else {
		http.Error(w, "Forbidden", 403) // Or Redirect?
	}
}
