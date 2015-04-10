package logger

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type LogChannel chan Log

type Log struct {
	Severity string
	Message  string
	Category string
	Time     time.Time
	MongoId  bson.ObjectId `bson:"_id,omitempty"`
}

func New(s *mgo.Session) LogChannel {
	l := make(LogChannel)
	ns := s.Copy()
	c := ns.DB("").C("logs")
	dbWriter(l, c)
	return l
}

func dbWriter(lc LogChannel, c *mgo.Collection) {
	for {
		l := <-lc
		c.Insert(l)
	}

}
