package user

import (
	"errors"
	"github.com/mongolar/mongolar/wrapper"
	"gopkg.in/mgo.v2/bson"
)

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
	return nil
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
