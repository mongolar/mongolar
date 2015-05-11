package login

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/davecgh/go-spew/spew"
	"github.com/mongolar/mongolar/controller"
	"github.com/mongolar/mongolar/url"
	"github.com/mongolar/mongolar/wrapper"
	"golang.org/x/oauth2"
	"net/http"
	"strings"
)

type Login struct {
	Controllers controller.ControllerMap
	State       string
}

func NewLogin() *Login {
	s, _ := getStateString()
	lmap := &Login{
		Controllers: make(controller.ControllerMap),
		State:       s,
	}
	lmap.Controllers["loginurls"] = lmap.LoginUrls
	lmap.Controllers["callback"] = lmap.Callback
	lmap.Controllers["logout"] = lmap.Logout
	return lmap
}

func (l *Login) Login(w *wrapper.Wrapper) {
	u := url.UrlToMap(w.Request.URL.Path)
	if c, ok := l.Controllers[u[2]]; ok {
		c(w)
	} else {
		http.Error(w.Writer, "Forbidden", 403)
		return
	}
}

func (lo *Login) LoginUrls(w *wrapper.Wrapper) {
	u := url.UrlToMap(w.Request.URL.Path)
	us := make(map[string]map[string]string)
	for k, l := range w.SiteConfig.Logins {
		s := strings.Split(l["scopes"], ",")
		conf := &oauth2.Config{
			ClientID:     l["client_id"],
			ClientSecret: l["client_secret"],
			Scopes:       s,
			Endpoint: oauth2.Endpoint{
				AuthURL:  l["auth_url"],
				TokenURL: l["token_url"],
			},
			RedirectURL: "http://" + w.Request.Host + "/" + u[0] + "/" + u[1] + "/callback/" + k,
		}
		u := conf.AuthCodeURL(lo.State, oauth2.AccessTypeOffline)
		m := map[string]string{"url": u, "login_text": l["login_text"]}
		us[k] = m
	}
	w.SetPayload("login_links", us)
	w.Serve()
	return
}

func (l *Login) Logout(w *wrapper.Wrapper) {

}

func (l *Login) Callback(w *wrapper.Wrapper) {
	u := url.UrlToMap(w.Request.URL.Path)
	s := w.Request.FormValue("state")
	if s != l.State {
		w.SiteConfig.Logger.Error("invalid oauth state, expected " + l.State + ", got " + s)
		http.Redirect(w.Writer, w.Request, w.SiteConfig.LoginFailure, 301)
		return
	}

	li := w.SiteConfig.Logins[u[3]]
	sc := strings.Split(li["scopes"], ",")
	conf := &oauth2.Config{
		ClientID:     li["client_id"],
		ClientSecret: li["client_secret"],
		Scopes:       sc,
		Endpoint: oauth2.Endpoint{
			AuthURL:  li["auth_url"],
			TokenURL: li["token_url"],
		},
		RedirectURL: "http://" + w.Request.Host + "/" + u[0] + "/" + u[1] + "/callback/" + u[3],
	}

	co := w.Request.FormValue("code")
	t, err := conf.Exchange(oauth2.NoContext, co)
	if err != nil {
		w.SiteConfig.Logger.Error("Exchange() failed with " + err.Error())
		http.Redirect(w.Writer, w.Request, w.SiteConfig.LoginFailure, 301)
		return
	}
	client := conf.Client(oauth2.NoContext, t)
	//TODO find common oauth values
	test, _ := client.Get("user")
	spew.Dump(test)
	http.Redirect(w.Writer, w.Request, w.SiteConfig.LoginSuccess, 301)
}

// Generate a random session key.
func getStateString() (string, error) {
	raw := make([]byte, 30)
	_, err := rand.Read(raw)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(raw), nil
}
