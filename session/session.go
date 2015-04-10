package session

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/jasonrichardsmith/mongolar/configs/site"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"time"
)

// Wrapper structure for Sessions
type Session struct {
	Id        string
	data      *SessionData
	dbSession *mgo.Session
	query     *mgo.Query
}

// The actual data being stored
type SessionData struct {
	MongoId   bson.ObjectId `bson:"_id,omitempty"`
	SessionId string        `bson:"session_id"`
	Data      map[string]interface{}
}

//Builds the session
func New(w http.ResponseWriter, r *http.Request, s *site.SiteConfig) *Session {
	se := new(Session)
	var duration time.Duration = time.Duration(s.SessionExpiration) * time.Hour
	expire := time.Now().Add(duration)
	c, err := r.Cookie("m_session_id")
	fmt.Print(err)
	if c == nil {
		id := getSessionID()
		c = &http.Cookie{
			Name:     "m_session_id",
			Value:    id,
			Path:     "/",
			Domain:   r.Host,
			MaxAge:   0,
			Secure:   true,
			HttpOnly: true,
			Raw:      "m_session_id=" + id,
			Unparsed: []string{"m_session_id=" + id},
		}
	}
	c.Expires = expire
	c.RawExpires = expire.Format(time.RFC3339)
	http.SetCookie(w, c)

	se.data = &SessionData{
		SessionId: c.Value,
		Data:      make(map[string]interface{}),
	}
	se.Id = c.Value
	se.dbSession = s.DbSession.Copy()
	se.getQuery(duration)
	se.setDbSession()
	return se
}

func (s Session) getQuery(d time.Duration) {
	c := s.dbSession.DB("").C("sessions")
	setCollection(c, d)
	s.query = c.Find(bson.M{"session_id": s.Id})
}

func (s Session) setDbSession() {
	c := mgo.Change{
		Update:    s.data,
		Upsert:    true,
		ReturnNew: true,
	}

	s.query.Apply(c, &s.data)
	s.getData()
}

func (s Session) getData() {
	s.query.One(&s.data)
}

func (s Session) Close() {
	s.dbSession.Close()
}

func (s Session) Get(n string) (v interface{}, err error) {
	s.getData()
	if v, ok := s.data.Data[n]; ok {
		return v, nil
	} else {
		err := errors.New("No Value")
		return nil, err
	}
}

func (s Session) Set(n string, v interface{}) {
	s.getData()
	s.data.Data[n] = v
	s.setDbSession()
}

func getSessionID() string {
	raw := make([]byte, 30)
	_, err := rand.Read(raw)
	if err != nil {
		fmt.Print(err)
	}
	return hex.EncodeToString(raw)

}

func setCollection(c *mgo.Collection, d time.Duration) {
	i := mgo.Index{
		Key:         []string{"SessionId"},
		Unique:      true,
		DropDups:    true,
		Background:  true,
		Sparse:      false,
		ExpireAfter: d,
	}
	err := c.EnsureIndex(i)
	fmt.Print(err)
}
