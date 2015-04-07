package content

import (
	"encoding/json"
)

type Content struct {
	mongular_content map[string]interface{}
}

func New(c map[string]interface{}) {
	mc = new(Content)
	mc.mongular_content = c
	return mc
}

func (mc Content) Serve(w http.ResponseWriter) {
	//ToDo Marshall and serve
	j := json.Marshal(mc)

}
