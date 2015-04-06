package mongolarserver

import (
	"github.com/jasonrichardsmith/mongolar/configs"
	"github.com/jasonrichardsmith/mongolar/router"
)

func Serve(cm ControllerMap) {

	c := configs.New()

	spew.Dump(c)
	HostSwitch := router.New(sites, aliases, HandlersMap)
	//log.Fatal(http.ListenAndServe(":"+Server.Port, HostSwitch))
}
