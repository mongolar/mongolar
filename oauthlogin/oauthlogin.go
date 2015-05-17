package oauthlogin
import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"github.com/davecgh/go-spew/spew"
	"github.com/mongolar/mongolar/controller"
	"github.com/mongolar/mongolar/session"
	"github.com/mongolar/mongolar/url"
	"github.com/mongolar/mongolar/wrapper"
	"golang.org/x/oauth2"
	"net/http"
	"strings"
)

		//Build example
		//conf := &oauth2.Config{
		//	ClientID:     l["client_id"],
		//	ClientSecret: l["client_secret"],
		//	Scopes:       s,
		///	Endpoint: oauth2.Endpoint{
		//		AuthURL:  l["auth_url"],
		//		TokenURL: l["token_url"],
		//	},
		//	RedirectURL: "http://" + w.Request.Host + "/" + u[0] + "/" + u[1] + "/callback/" + k,
		//}
	//TODO find common oauth values
	//test, err1 := client.Get("https://api.github.com/user")
	//buf := new(bytes.Buffer)
	//buf.ReadFrom(test.Body)
	//resp := buf.String()
	//spew.Dump(err1)
	//spew.Dump(resp)
	//http.Redirect(w.Writer, w.Request, w.SiteConfig.LoginSuccess, 301)



type LoginMap struct {
	Controllers controller.ControllerMap
	Logins      map[string]Login
}

func NewLoginMap() *LoginMap {
	lmap := &Login{
		Controllers: make(controller.ControllerMap),
	}
	lmap.Controllers["loginurls"] = lmap.LoginUrls
	lmap.Controllers["callback"] = lmap.Callback
	lmap.Controllers["logout"] = lmap.Logout
	return lmap
}

func (l *LoginMap) Login(w *wrapper.Wrapper) {
	u := url.UrlToMap(w.Request.URL.Path)
	if c, ok := l.Controllers[u[2]]; ok {
		if c, ok := w.SiteConfig.Logins[u[3]]; ok {
			if c, ok := l.Logins[u[3]]; ok {
				c(w)
				return
			}
		}
	}
	http.Error(w.Writer, "Forbidden", 403)
	return
}

func (lo *LoginMap) LoginUrls(w *wrapper.Wrapper) {
	u := url.UrlToMap(w.Request.URL.Path)
	us := make(map[string]map[string]string)
	for k, l := range w.SiteConfig.Logins {
		conf := lo.Logins.BuildConfig(l)
		conf.RedirectURL = "http://" + w.Request.Host + "/" + u[0] + "/" + u[1] + "/callback/" + k,
		u := conf.AuthCodeURL(lo.State, oauth2.AccessTypeOffline)
		m := map[string]string{"url": u, "login_text": l["login_text"]}
		us[k] = m
	}
	w.SetPayload("login_links", us)
	w.Serve()
	return
}

func (l *LoginMap) Logout(w *wrapper.Wrapper) {

}

func (l *LoginMap) Callback(w *wrapper.Wrapper) {
	u := url.UrlToMap(w.Request.URL.Path)
	s := w.Request.FormValue("state")
	sc := w.SiteConfig.Login[u[3]]
	if s != sc["state"] {
		w.SiteConfig.Logger.Error("invalid oauth state, expected " + sc["state"] + ", got " + s)
		http.Redirect(w.Writer, w.Request, w.SiteConfig.LoginFailure, 301)
		return
	}
	conf := l.Logins.BuildConfig(sc)
	conf.RedirectURL = "http://" + w.Request.Host + "/" + u[0] + "/" + u[1] + "/callback/" + u[3] ,
	co := w.Request.FormValue("code")
	t, err := conf.Exchange(oauth2.NoContext, co)
	if err != nil {
		w.SiteConfig.Logger.Error("Exchange() failed with " + err.Error())
		http.Redirect(w.Writer, w.Request, w.SiteConfig.LoginFailure, 301)
		return
	}
	client := conf.Client(oauth2.NoContext, t)
	id := l.Logins[u[3]].GetId(client)
	//TODO link session to user
}



type Login interface {
	GetUrl() string
	GetID() string
	GetToken() string
	ValidateLogin() string
}


type LoginStructure struct{
	AuthURL: string
	TokenURL: string
	Scope: string
	Id: string
	Secret:	string
}

type GitHub LoginStructure

func NewGitHub(map[string]string){
	
}
