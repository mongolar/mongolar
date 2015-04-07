package content

import (
	"encoding/json"
	"github.com/jasonrichardsmith/mongolar/responders"
	"net/http"
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
	responders.ServeJson(mc)

}
