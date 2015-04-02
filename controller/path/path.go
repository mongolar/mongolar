package path

import (
	"github.com/jasonrichardsmith/mongolar/configs/site"
	"github.com/jasonrichardsmith/mongolar/controller"
	"net/http"
)

func GetContent(s *http.Request, s *site.SiteConfig) *controller.ControllerResponse {
	test := 'Test'
	return test
}
