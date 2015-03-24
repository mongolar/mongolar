package main

import (
	//"github.com/davecgh/go-spew/spew"
	"github.com/jasonrichardsmith/mongolar/configs"
)

func main() {
	mo := new(*configs.MongolarSites)
	mo.BuildMongolarSiteConfigs()

}
