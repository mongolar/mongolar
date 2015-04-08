package main

import (
	"github.com/jasonrichardsmith/mongolar/configs"
	"github.com/jasonrichardsmith/mongolar/controller"
	"github.com/jasonrichardsmith/mongolar/controller/domain"
	//	"github.com/jasonrichardsmith/mongolar/router"
)

func main() {
	cm := controller.NewMap()
	cm['test'] := domain.Serve()
	Serve(cm)
}

func Serve(cm controller.ControllerMap) {

	c := configs.New()

	spew.Dump(c)
	//HostSwitch := router.New(sites, aliases, HandlersMap)
	//log.Fatal(http.ListenAndServe(":"+Server.Port, HostSwitch))
}
