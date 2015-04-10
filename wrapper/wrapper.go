package wrapper

import (
	"encoding/json"
	"github.com/jasonrichardsmith/mongolar/configs/site"
	"github.com/jasonrichardsmith/mongolar/session"
	"net/http"
)

type Wrapper struct {
	Writer     http.ResponseWriter
	Request    *http.Request
	SiteConfig *site.SiteConfig
	Session    *session.Session
	Payload    map[string]interface{}
}

func New(w http.ResponseWriter, r *http.Request, s *site.SiteConfig) *Wrapper {
	wr := Wrapper{Writer: w, Request: r, SiteConfig: s}
	wr.Session = session.New(w, r, s)
	wr.Payload = make(map[string]interface{})
	return &wr
}

func (w *Wrapper) SetContent(c map[string]interface{}) {
	w.SetPayload("content", c)
}

func (w *Wrapper) SetPayload(n string, v map[string]interface{}) {
	w.Payload[n] = v
}

func (w *Wrapper) Serve() {
	js, err := json.Marshal(w.Payload)
	if err != nil {
		http.Error(w.Writer, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Writer.Write(js)
}
