package controller

type ControllerMap map[string]Controller

func NewMap() ControllerMap {
	return make(ControllerMap)
}

type Controller func(http.ResponseWriter, *http.Request, *site.SiteConfig)
