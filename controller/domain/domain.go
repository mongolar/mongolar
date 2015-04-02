package domain

import (
	"github.com/jasonrichardsmith/mongolar/configs/site"
	"github.com/jasonrichardsmith/mongolar/controller"
	"net/http"
)

func GetContent(s *http.Request, s *site.SiteConfig) *controller.ControllerResponse {
	//file, _ := ioutil.ReadFile(site.Directory + "/index.html")
	file := "test"
	return file
}
