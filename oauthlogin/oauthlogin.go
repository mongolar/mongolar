package oauthlogin

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/mongolar/mongolar/controller"
	"github.com/mongolar/mongolar/user"
	"github.com/mongolar/mongolar/wrapper"
	"golang.org/x/oauth2"
	"net/http"
)

func GetControllerMap(cm controller.ControllerMap) {
	lmap := NewLoginMap()
	cm["login"] = lmap.Login
}

type LoginMap struct {
	Controllers controller.ControllerMap
	Logins      map[string]OALogin
	State       string
}

func NewLoginMap() LoginMap {
	lmap := LoginMap{
		Controllers: make(controller.ControllerMap),
	}
	lmap.Controllers["loginurls"] = lmap.LoginUrls
	lmap.Controllers["callback"] = lmap.Callback
	lmap.Controllers["logout"] = lmap.Logout
	github := NewGitHub()
	lmap.Logins = map[string]OALogin{
		"github": github,
	}
	lmap.State = StateString()
	return lmap
}

func StateString() string {
	raw := make([]byte, 30)
	rand.Read(raw)
	return hex.EncodeToString(raw)
}

func (l *LoginMap) Login(w *wrapper.Wrapper) {
	if controller, ok := l.Controllers[w.APIParams[0]]; ok {
		w.Shift()
		controller(w)
		return
	}
	http.Error(w.Writer, "Forbidden", 403)
	return
}

func (lo *LoginMap) LoginUrls(w *wrapper.Wrapper) {
	us := make(map[string]map[string]string)
	oauthlogins := make(map[string]map[string]string)
	w.SiteConfig.RawConfig.MarshalKey("OAuthLogins", &oauthlogins)
	for k, l := range oauthlogins {
		c := "http://" + w.Request.Host + "/" + w.SiteConfig.APIEndPoint + "/login/callback/" + k
		login := lo.Logins[k]
		login.SetConfig(l, c, lo.State)
		u := login.GetUrl()
		m := map[string]string{"url": u, "login_text": l["login_text"]}
		us[k] = m
	}
	w.SetPayload("login_links", us)
	w.Serve()
	return
}

func (lo *LoginMap) Logout(w *wrapper.Wrapper) {
	w.SetPayload("logout", "test logout")
	w.Serve()
	return
}

func (lo *LoginMap) Callback(w *wrapper.Wrapper) {
	oauthlogins := make(map[string]map[string]string)
	w.SiteConfig.RawConfig.MarshalKey("OAuthLogins", &oauthlogins)
	loginurls := make(map[string]string)
	w.SiteConfig.RawConfig.MarshalKey("LoginURLs", &loginurls)
	if _, ok := oauthlogins[w.APIParams[0]]; ok {
		if _, ok := lo.Logins[w.APIParams[0]]; ok {
			s := w.Request.FormValue("state")
			sc := oauthlogins[w.APIParams[0]]
			login := lo.Logins[w.APIParams[0]]
			if lo.State != s {
				errmessage := fmt.Sprintf("Invalid oauth state, expected %s, got %s", lo.State, s)
				w.SiteConfig.Logger.Error(errmessage)
				http.Redirect(w.Writer, w.Request, loginurls["failure"], 301)
				return
			}
			login.SetConfig(sc, "", "")
			code := w.Request.FormValue("code")
			token, err := login.GetToken(code)
			if err != nil {
				errmessage := fmt.Sprintf("Exchange() failed with %s", err.Error())
				w.SiteConfig.Logger.Error(errmessage)
				http.Redirect(w.Writer, w.Request, loginurls["failure"], 301)
				return
			}
			u := login.GetUser()
			err = u.Set(w)
			if err != nil {
				errmessage := fmt.Sprintf("Unable to set user: %s", err.Error())
				w.SiteConfig.Logger.Error(errmessage)
				http.Redirect(w.Writer, w.Request, loginurls["failure"], 301)
				return
			}
			err = w.SetSessionValue("user_id", u.MongoId)
			if err != nil {
				errmessage := fmt.Sprintf("Unable to set user id on session: %s", err.Error())
				w.SiteConfig.Logger.Error(errmessage)
				http.Redirect(w.Writer, w.Request, loginurls["failure"], 301)
				return
			}
			err = w.SetSessionValue("token", token)
			if err != nil {
				errmessage := fmt.Sprintf("Unable to set token on session: %s", err.Error())
				w.SiteConfig.Logger.Error(errmessage)
				http.Redirect(w.Writer, w.Request, loginurls["failure"], 301)
				return
			}
			http.Redirect(w.Writer, w.Request, loginurls["success"], 301)
			return
		}
	}
	http.Error(w.Writer, "Forbidden", 403)
	return
}

type OALogin interface {
	SetConfig(map[string]string, string, string)
	GetUrl() string
	GetUser() *user.User
	GetToken(string) (*oauth2.Token, error)
}

type LoginStructure struct {
	Token    *oauth2.Token
	AuthURL  string
	TokenURL string
	Scope    string
	ClientId string
	Secret   string
	Callback string
	State    string
	Config   *oauth2.Config
}

func (ls *LoginStructure) BuildConfig() {
	ls.Config = &oauth2.Config{
		ClientID:     ls.ClientId,
		ClientSecret: ls.Secret,
		Scopes:       []string{ls.Scope},
		Endpoint: oauth2.Endpoint{
			AuthURL:  ls.AuthURL,
			TokenURL: ls.TokenURL,
		},
		RedirectURL: ls.Callback,
	}
}

type GitHub struct {
	LoginStructure
}

func NewGitHub() *GitHub {
	gh := &GitHub{
		LoginStructure{
			AuthURL:  "https://github.com/login/oauth/authorize",
			TokenURL: "https://github.com/login/oauth/access_token",
			Scope:    "user:email",
		},
	}
	return gh
}

func (gh *GitHub) SetConfig(v map[string]string, c string, s string) {
	gh.ClientId = v["client_id"]
	gh.Secret = v["client_secret"]
	gh.Callback = c
	gh.State = s
	gh.BuildConfig()
}

func (gh *GitHub) GetUrl() string {
	return gh.Config.AuthCodeURL(gh.State, oauth2.AccessTypeOffline)
}

func (gh *GitHub) GetUser() *user.User {
	client := gh.Config.Client(oauth2.NoContext, gh.Token)
	test, _ := client.Get("https://api.github.com/user")
	u := new(user.User)
	json.NewDecoder(test.Body).Decode(u)
	u.Type = "github"
	return u

}

func (gh *GitHub) GetToken(code string) (*oauth2.Token, error) {
	token, err := gh.Config.Exchange(oauth2.NoContext, code)
	if err == nil {
		gh.Token = token
	}
	return token, err
}
