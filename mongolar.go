package main

import (
	"github.com/mongolar/mongolar/admin"
	"github.com/mongolar/mongolar/configs"
	"github.com/mongolar/mongolar/controller"
	"github.com/mongolar/mongolar/oauthlogin"
	//"github.com/davecgh/go-spew/spew"
	"github.com/mongolar/mongolar/router"
	"gopkg.in/mgo.v2"
	"net/http"
	"time"
)

func main() {
	amap, _ := admin.NewAdmin()
	lmap := oauthlogin.NewLoginMap()
	cm := controller.NewMap()
	cm["domian_public_value"] = controller.DomainPublicValue
	cm["path"] = controller.PathValues
	cm["content"] = controller.ContentValues
	cm["wrapper"] = controller.WrapperValues
	cm["slug"] = controller.SlugValues
	cm["admin"] = amap.Admin
	cm["login"] = lmap.Login
	Serve(cm)
}

func Serve(cm controller.ControllerMap) {
	c, port := configs.New()
	EnsureIndexes(c)
	HostSwitch := router.New(c.Aliases, c.SitesMap, cm)
	http.ListenAndServe(":"+port, HostSwitch)
}

func EnsureIndexes(configs *configs.Configs) {
	for _, site_config := range configs.SitesMap {
		db_session := site_config.DbSession.Copy()
		defer db_session.Close()
		var duration time.Duration = time.Duration(site_config.SessionExpiration * time.Hour)
		i := mgo.Index{
			Key:         []string{"updated"},
			Unique:      false,
			DropDups:    false,
			Background:  true,
			Sparse:      false,
			ExpireAfter: duration,
		}
		c := db_session.DB("").C("sessions")
		c.EnsureIndex(i)
		i = mgo.Index{
			Key:        []string{"path", "wildcard"},
			Unique:     true,
			DropDups:   true,
			Background: true,
			Sparse:     false,
		}
		c = db_session.DB("").C("paths")
		c.EnsureIndex(i)
		i = mgo.Index{
			Key:        []string{"id", "type"},
			Unique:     true,
			DropDups:   true,
			Background: true,
			Sparse:     false,
		}
		c = db_session.DB("").C("users")
		c.EnsureIndex(i)
	}
}
