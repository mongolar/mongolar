package wrapper

import (
	"github.com/jasonrichardsmith/mongolar/configs/site"
	"github.com/jasonrichardsmith/mongolar/session"
	"http/net"
)

type Wrapper struct {
	Writer     http.ResponseWriter
	Request    *http.Request
	SiteConfig *site.SiteConfig
	Session    *session.Session
}

func New(w http.ResponseWriter, r *http.Request, s *site.SiteConfig) {
	wr = Wrapper{Writer: w, Request: r, SiteConfig: s}
	session.New(&wr)

}
