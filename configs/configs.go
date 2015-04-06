package configs

import (
	//"github.com/davecgh/go-spew/spew"
	"github.com/jasonrichardsmith/mongolar/configs/aliases"
	"github.com/jasonrichardsmith/mongolar/configs/server"
	"github.com/jasonrichardsmith/mongolar/configs/sites"
)

type Configs struct {
	Server   *server.Server
	SitesMap *sites.SitesMap
	Aliases  *aliases.Aliases
}

func New() *Configs {
	c := new(Configs)
	c.Server = server.New()
	c.SitesMap = sites.New()
	c.Aliases = aliases.New(c.SitesMap)
	return c
}
