package wrapper

import (
	"gopkg.in/mgo.v2/bson"
	"net"
	"net/http"
	"time"
)

// Wrapper structure for Sessions
type Session struct {
	Id      bson.ObjectId `bson:"_id"`
	Updated time.Time     `bson:"updated"`
}

//Session constructor
func (w *Wrapper) NewSession() error {
	se := new(Session)
	// Duration of expiration, needs to be worked out between cookies and db
	var duration time.Duration = time.Duration(w.SiteConfig.SessionExpiration * time.Hour)
	expire := time.Now().Add(duration)
	// Set the cookies
	c, err := w.Request.Cookie("m_session_id")
	if err != nil && err.Error() != "http: named cookie not present" {
		return err
	}
	//  If cookie is not set, set one
	if c == nil {
		host, _, _ := net.SplitHostPort(w.Request.Host)
		se.Id = bson.NewObjectId()
		c = &http.Cookie{
			Name:   "m_session_id",
			Value:  se.Id.Hex(),
			Path:   "/",
			Domain: host,
		}
	} else {
		se.Id = bson.ObjectIdHex(c.Value)
	}
	//  New or reused cookies will have their expiration refreshed
	c.Expires = expire
	c.RawExpires = expire.Format(time.RFC3339)
	http.SetCookie(w.Writer, c)
	w.Session = se
	return nil
}

func (w *Wrapper) SetSession() error {
	w.Session.Updated = time.Now()
	c := w.DbSession.DB("").C("sessions")
	_, err := c.Upsert(bson.M{"_id": w.Session.Id}, bson.M{"$set": bson.M{"_id": w.Session.Id, "updated": w.Session.Updated}})
	if err != nil {
		return err
	}
	return nil
}

// Get current session data
func (w *Wrapper) SetSessionValue(k string, v interface{}) error {
	c := w.DbSession.DB("").C("sessions")
	err := c.Update(bson.M{"_id": w.Session.Id}, bson.M{"$set": bson.M{k: v}})
	if err != nil {
		return err
	}
	return nil
}

// Get a session value by key.
func (w *Wrapper) GetSessionValue(n string, i interface{}) error {
	c := w.DbSession.DB("").C("sessions")
	err := c.Find(bson.M{"_id": w.Session.Id}).Select(bson.M{n: 1}).One(i)
	if err != nil {
		return err
	}
	return nil
}
