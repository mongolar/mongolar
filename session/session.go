// Session will establish a user session where developers can store arbitrary user session data
// All session data will be stored in the sessions collection.
// TODO: This package need more development and testing
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
	Id        string       // Session ID
	data      *SessionData // The data to store
	dbSession *mgo.Session // The mg DB session
	query     *mgo.Query   // The query to the session record
}

// The actual data being stored in the db
type SessionData struct {
	MongoId   bson.ObjectId `bson:"_id,omitempty"`
	SessionId string        `bson:"session_id"`
	Data      map[string]interface{}
}

//Session constructor
func New(w http.ResponseWriter, r *http.Request, s *site.SiteConfig) *Session {
	se := new(Session)
	// Duration of expiration, needs to be worked out between cookies and db
	var duration time.Duration = time.Duration(s.SessionExpiration) * time.Hour
	expire := time.Now().Add(duration)
	// Set the cookies
	c, err := r.Cookie("m_session_id")
	if err != nil {
		//need to do something here, not sure what
		fmt.Print(err)
	}
	//  If cookie is not set, set one
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
	//  New or resused cookies will have their expiration refreshed
	c.Expires = expire
	c.RawExpires = expire.Format(time.RFC3339)
	http.SetCookie(w, c)
	// Build the session data
	se.data = &SessionData{
		SessionId: c.Value,
		Data:      make(map[string]interface{}),
	}
	se.Id = c.Value
	// Copy the DB session
	se.dbSession = s.DbSession.Copy()
	se.getQuery(duration)
	se.setDbSession()
	return se
}

// Set the query for this record in the session collection
func (s Session) getQuery(d time.Duration) {
	c := s.dbSession.DB("").C("sessions")
	setCollection(c, d)
	s.query = c.Find(bson.M{"session_id": s.Id})
}

// Insert the db session if it does not exist
func (s Session) setDbSession() {
	c := mgo.Change{
		Update:    s.data,
		Upsert:    true,
		ReturnNew: true,
	}

	s.query.Apply(c, &s.data)
	s.getData()
}

// Get current session data
func (s Session) getData() {
	s.query.One(&s.data)
}

// Close the current session
func (s Session) Close() {
	s.dbSession.Close()
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

// Set value by key
func (s Session) Set(n string, v interface{}) {
	s.getData()
	s.data.Data[n] = v
	s.setDbSession()
}

// Generate a random session key.
func getSessionID() string {
	raw := make([]byte, 30)
	_, err := rand.Read(raw)
	if err != nil {
		fmt.Print(err)
	}
	return hex.EncodeToString(raw)

}

// Set the parameters for the Mongodb collection
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
