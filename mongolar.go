package main

import (
	"github.com/jasonrichardsmith/mongolar/configs"
	//	"github.com/jasonrichardsmith/mongolar/configs/server"
	//	"github.com/jasonrichardsmith/mongolar/configs/site"
	"github.com/jasonrichardsmith/mongolar/controllers/index"
	"github.com/jasonrichardsmith/mongolar/controllers/path"
	"github.com/jasonrichardsmith/mongolar/router"
	"github.com/julienschmidt/httprouter"
	//	"io/ioutil"
	"log"
	"net/http"
)

var HandlersMap = map[string]httprouter.Handle{
	"/mongolar/path/*fullpath": path.Path,
	"/": index.Index,
}

func main() {
	Server, sites, aliases := configs.GetAll()
	HostSwitch := router.Build(sites, aliases, HandlersMap)
	log.Fatal(http.ListenAndServe(":"+Server.Port, HostSwitch))
}
