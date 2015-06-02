package oauthlogin

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mongolar/mongolar/controller"
	"github.com/mongolar/mongolar/wrapper"
	"golang.org/x/oauth2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

type LoginMap struct {
	Controllers controller.ControllerMap
	Logins      map[string]Login
	State       string
}

func NewLoginMap() *LoginMap {
	lmap := &LoginMap{
		Controllers: make(controller.ControllerMap),
	}
	lmap.Controllers["loginurls"] = lmap.LoginUrls
	lmap.Controllers["callback"] = lmap.Callback
	lmap.Controllers["logout"] = lmap.Logout
	github := NewGitHub()
	lmap.Logins = map[string]Login{
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
	for k, l := range w.SiteConfig.OAuthLogins {
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
	if _, ok := w.SiteConfig.OAuthLogins[w.APIParams[0]]; ok {
		if _, ok := lo.Logins[w.APIParams[0]]; ok {
			s := w.Request.FormValue("state")
			sc := w.SiteConfig.OAuthLogins[w.APIParams[0]]
			login := lo.Logins[w.APIParams[0]]
			if lo.State != s {
				errmessage := fmt.Sprintf("Invalid oauth state, expected %s, got %s", lo.State, s)
				w.SiteConfig.Logger.Error(errmessage)
				http.Redirect(w.Writer, w.Request, w.SiteConfig.LoginURLs["failure"], 301)
				return
			}
			login.SetConfig(sc, "", "")
			code := w.Request.FormValue("code")
			token, err := login.GetToken(code)
			if err != nil {
				errmessage := fmt.Sprintf("Exchange() failed with %s", err.Error())
				w.SiteConfig.Logger.Error(errmessage)
				http.Redirect(w.Writer, w.Request, w.SiteConfig.LoginURLs["failure"], 301)
				return
			}
			u := login.GetUser()
			err = u.Set(w)
			if err != nil {
				errmessage := fmt.Sprintf("Unable to set user: %s", err.Error())
				w.SiteConfig.Logger.Error(errmessage)
				http.Redirect(w.Writer, w.Request, w.SiteConfig.LoginURLs["failure"], 301)
				return
			}
			err = w.SetSessionValue("user_id", u.MongoId)
			if err != nil {
				errmessage := fmt.Sprintf("Unable to set user id on session: %s", err.Error())
				w.SiteConfig.Logger.Error(errmessage)
				http.Redirect(w.Writer, w.Request, w.SiteConfig.LoginURLs["failure"], 301)
				return
			}
			err = w.SetSessionValue("token", token)
			if err != nil {
				errmessage := fmt.Sprintf("Unable to set token on session: %s", err.Error())
				w.SiteConfig.Logger.Error(errmessage)
				http.Redirect(w.Writer, w.Request, w.SiteConfig.LoginURLs["failure"], 301)
				return
			}
			http.Redirect(w.Writer, w.Request, w.SiteConfig.LoginURLs["success"], 301)
			return
		}
	}
	http.Error(w.Writer, "Forbidden", 403)
	return
}

type Login interface {
	SetConfig(map[string]string, string, string)
	GetUrl() string
	GetUser() *User
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

func (gh *GitHub) GetUser() *User {
	client := gh.Config.Client(oauth2.NoContext, gh.Token)
	test, _ := client.Get("https://api.github.com/user")
	u := new(User)
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

// move below to new package
type User struct {
	MongoId bson.ObjectId `json:"-" bson:"_id"`
	Email   string        `json:"email" bson:"email"`
	Id      int           `json:"id" bson:"id"`
	Name    string        `json:"login" bson:"name"`
	Type    string        `bson:"type"`
	Roles   []string      `bson:"roles,omitempty"`
}

func (u *User) Set(w *wrapper.Wrapper) error {
	c := w.DbSession.DB("").C("users")
	tmpuser := new(User)
	s := bson.M{"id": u.Id, "type": u.Type}
	err := c.Find(s).One(tmpuser)
	if err != nil {
		if err.Error() == "not found" {
			u.MongoId = bson.NewObjectId()
			err := c.Insert(u)
			if err != nil {
				return err
			}
			return nil
		}
		return err
	}
	u.MongoId = tmpuser.MongoId
	err = c.Update(bson.M{"_id": u.MongoId}, bson.M{"$set": u})
	return err
}

func (u *User) Get(w *wrapper.Wrapper) error {
	c := w.DbSession.DB("").C("users")
	var id bson.M
	w.GetSessionValue("user_id", &id)
	if id == nil {
		err := errors.New("User not found")
		return err
	}
	if user_id, ok := id["user_id"]; ok {
		err := c.Find(bson.M{"_id": user_id}).One(u)
		return err
	} else {
		err := errors.New("User not found")
		return err
	}
}
