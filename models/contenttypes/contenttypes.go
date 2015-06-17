package contenttypes

import (
	"errors"
	"github.com/mongolar/mongolar/form"
	"github.com/mongolar/mongolar/wrapper"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"path"
	"strings"
)

type ContentType struct {
	Form    []*form.Field `bson:"form,omitempty" json:"content_form"`
	Type    string        `bson:"type,omitempty" json:"type"`
	MongoId bson.ObjectId `bson:"_id" json:"id"`
}

// Constructor for paths
func NewContenType() Path {
	f := make([]*form.Field, 0)
	id := bson.NewObjectId()
	ct := ContenType{MongoId: id, Form: f}
	return ct
}

// Constructor for existing paths
func LoadContentType(i string, w *wrapper.Wrapper) (ContenType, error) {
	f := make([]*form.Field, 0)
	ct := ContenType{Form: f}
	err := ct.GetById(i, w)
	return ct, err
}

//Save an element in its current state.
func (ct *ContentType) Save(w *wrapper.Wrapper) error {
	if !p.MongoId.Valid() {
		p.MongoId = bson.NewObjectId()
	}
	if p.Type == "" {
		return errors.New("Type required")
	}
	c := w.DbSession.DB("").C("content_types")
	_, err := c.Upsert(p.MongoId, p)
	if err != nil {
		return err
	}
	return nil
}

// Get Path by Id
func (ct *ContentType) GetById(i string, w *wrapper.Wrapper) error {
	if !bson.IsObjectIdHex(i) {
		return errors.New("Invalid Id Hex")
	}
	c := w.DbSession.DB("").C("content_types")
	err := c.FindId(bson.ObjectIdHex(i)).One(&ct)
	return err
}

// Get all Paths
func ContentTypeList(w *wrapper.Wrapper) ([]ContenType, error) {
	cl := make([]ContenType, 0)
	c := w.DbSession.DB("").C("content_types")
	i := c.Find(nil).Limit(50).Iter()
	err := i.All(&cl)
	if err != nil {
		return nil, err
	}
	return cl, nil
}
