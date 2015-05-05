package main

import (
	"github.com/mongolar/mongolar/admin"
	"github.com/mongolar/mongolar/configs"
	"github.com/mongolar/mongolar/controller"
	"github.com/mongolar/mongolar/router"
	"net/http"
)

func main() {
	am := admin.NewAdmin()
	cm := controller.NewMap()
	cm["domian_public_value"] = controller.DomainPublicValue
	cm["path"] = controller.PathValues
	cm["content"] = controller.ContentValues
	cm["wrapper"] = controller.WrapperValues
	cm["slug"] = controller.SlugValues
	cm["admin"] = am.Admin
	Serve(cm)
}

func Serve(cm controller.ControllerMap) {

	c := configs.New()
	//	spew.Dump(c)
	HostSwitch := router.New(c.Aliases, c.SitesMap, cm)
	http.ListenAndServe(":"+c.Server.Port, HostSwitch)
}
