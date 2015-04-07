package controller

type ControllerMap map[string]Controller

func NewMap() ControllerMap {
	return make(ControllerMap)
}

type Controller func(http.ResponseWriter, *http.Request, *site.SiteConfig)

type Content struct {
	mongular_content map[string]interface{}
}

func NewContent(c map[string]interface{}) {
	mc = new(Content)
	mc.mongular_content = c
	return mc
}

func (mc Content) write(w http.ResponseWriter) {
	//ToDo Marshall and serve
}
