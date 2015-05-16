// Session will establish a user session where developers can store arbitrary user session data
// All session data will be stored in the sessions collection.
// TODO: This package need more development and testing
package session

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

// Wrapper structure for Sessions
type Session struct {
	Id    bson.ObjectId `bson:"_id"`
	Token string        `bson:"token"`
}

//Session constructor
func New(w *wrapper.Wrapper) (*Session, error) {
	se := new(Session)
	// Duration of expiration, needs to be worked out between cookies and db
	var duration time.Duration = time.Duration(w.SitConfig.SessionExpiration * time.Hour)
	expire := time.Now().Add(duration)
	// Set the cookies
	c, err := w.Request.Cookie("m_session_id")
	if err != nil {
		return nil, err
	}
	//  If cookie is not set, set one
	if c == nil {
		id := bson.NewObjectId()
		c = &http.Cookie{
			Name:   "m_session_id",
			Value:  id.Hex(),
			Path:   "/",
			Domain: r.Host,
		}
	}
	//  New or reused cookies will have their expiration refreshed
	c.Expires = expire
	c.RawExpires = expire.Format(time.RFC3339)
	se.Id = c.Value
	http.SetCookie(w.Writer, c)
	return se, nil
}

// Get current session data
func (s Session) setSession() {
	s.collection.Upsert(bson.M{"session_id": s.Id}, s.data)
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
