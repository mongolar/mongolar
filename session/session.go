// Session will establish a user session where developers can store arbitrary user session data
// All session data will be stored in the sessions collection.
// TODO: This package need more development and testing
package session

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	log "github.com/Sirupsen/logrus"
	"github.com/mongolar/mongolar/configs"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"time"
)

// Wrapper structure for Sessions
type Session struct {
	Id         string // Session ID
	Cookie     *http.Cookie
	data       *SessionData // The data to store
	dbSession  *mgo.Session // The mg DB session
	collection *mgo.Collection
}

// The actual data being stored in the db
type SessionData struct {
	SessionId string                 `bson:"session_id"`
	Data      map[string]interface{} `bson:"data"`
}

//Session constructor
func New(w http.ResponseWriter, r *http.Request, s *configs.SiteConfig) *Session {
	se := new(Session)
	// Duration of expiration, needs to be worked out between cookies and db
	var duration time.Duration = time.Duration(s.SessionExpiration * time.Hour)
	expire := time.Now().Add(duration)
	// Set the cookies
	c, err := r.Cookie("m_session_id")
	if err != nil {
		//need to do something here, not sure what
		s.Logger.WithFields(log.Fields{"Remote": r.RemoteAddr}).Warn(err.Error())
	}
	//  If cookie is not set, set one
	if c == nil {
		id, err := getSessionID()
		if err != nil {
			s.Logger.WithFields(log.Fields{"Remote": r.RemoteAddr}).Warn(err.Error())
		}
		c = &http.Cookie{
			Name:   "m_session_id",
			Value:  id,
			Path:   "/",
			Domain: r.Host,
		}
	}
	//  New or reused cookies will have their expiration refreshed
	c.Expires = expire
	c.RawExpires = expire.Format(time.RFC3339)
	se.Cookie = c
	// Build the session data
	se.data = &SessionData{
		SessionId: c.Value,
		Data:      make(map[string]interface{}),
	}
	se.Id = c.Value
	// Copy the DB session
	se.dbSession = s.DbSession.Copy()
	se.collection = se.dbSession.DB("").C("sessions")
	err = setCollection(se.collection, duration)
	if err != nil {
		s.Logger.WithFields(log.Fields{"Remote": r.RemoteAddr}).Error(err.Error())
	}
	http.SetCookie(w, se.Cookie)

	err = se.getData()
	if err != nil {
		se.setSession()
	}
	return se
}

// Get current session data
func (s Session) getData() error {
	err := s.collection.Find(bson.M{"session_id": s.Id}).One(&s.data)
	return err
}

// Get current session data
func (s Session) setSession() {
	s.collection.Insert(bson.M{"session_id": s.Id}, s.data)
}

// Get a session value by key.
func (s Session) Get(n string) (v interface{}, err error) {
	s.getData()
	if v, ok := s.data.Data[n]; ok {
		return v, nil
	} else {
		err := errors.New("No Value")
		return nil, err
	}
}

// Get a session value by key.
func (s Session) Set(n string, v interface{}) {
	k := "data." + n
	d := bson.M{
		"$set": bson.M{
			k: v,
		},
	}
	s.collection.Update(bson.M{"session_id": s.Id}, d)
}

// Close the current session
func (s Session) Close() {
	s.dbSession.Close()
}

// Generate a random session key.
func getSessionID() (string, error) {
	raw := make([]byte, 30)
	_, err := rand.Read(raw)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(raw), nil

}

// Set the parameters for the Mongodb collection
func setCollection(c *mgo.Collection, d time.Duration) error {
	i := mgo.Index{
		Key:         []string{"session_id"},
		Unique:      true,
		DropDups:    true,
		Background:  true,
		Sparse:      false,
		ExpireAfter: d,
	}
	err := c.EnsureIndex(i)
	return err
}
