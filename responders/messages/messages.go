package messages

import (
	"encoding/json"
	"github.com/jasonrichardsmith/mongolar/responders"
	"net/http"
)

type Messages struct {
	mongular_messages []map[string]string
}

func New(ms []map[string]string) {
	m = new(Messages)
	m.mongular_messages = ms
	return ms
}

func (ms Content) Serve(w http.ResponseWriter) {
	responders.ServeJson(ms)

}
