package contenttypes

import (
	"errors"
	"github.com/mongolar/mongolar/form"
	"github.com/mongolar/mongolar/wrapper"
	"gopkg.in/mgo.v2/bson"
)

type ContentType struct {
	Form    []*form.Field `bson:"form,omitempty" json:"content_form,omitempty"`
	Type    string        `bson:"type,omitempty" json:"type,omitempty"`
	MongoId bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
}

func NewContentType() ContentType {
	f := make([]*form.Field, 0)
	id := bson.NewObjectId()
	ct := ContentType{MongoId: id, Form: f}
	return ct
}

func LoadContentType(i string, w *wrapper.Wrapper) (ContentType, error) {
	f := make([]*form.Field, 0)
	ct := ContentType{Form: f}
	err := ct.GetById(i, w)
	return ct, err
}
func LoadContentTypeT(t string, w *wrapper.Wrapper) (ContentType, error) {
	f := make([]*form.Field, 0)
	ct := ContentType{Form: f}
	err := ct.GetByType(t, w)
	return ct, err
}

func (ct *ContentType) Save(w *wrapper.Wrapper) error {
	if !ct.MongoId.Valid() {
		ct.MongoId = bson.NewObjectId()
	}
	if ct.Type == "" {
		return errors.New("Type required")
	}
	c := w.DbSession.DB("").C("content_types")
	_, err := c.Upsert(ct.MongoId, ct)
	if err != nil {
		return err
	}
	return nil
}

func (ct *ContentType) GetById(i string, w *wrapper.Wrapper) error {
	if !bson.IsObjectIdHex(i) {
		return errors.New("Invalid Id Hex")
	}
	c := w.DbSession.DB("").C("content_types")
	err := c.FindId(bson.ObjectIdHex(i)).One(&ct)
	return err
}

func (ct *ContentType) GetByType(t string, w *wrapper.Wrapper) error {
	c := w.DbSession.DB("").C("content_types")
	err := c.Find(bson.M{"type": t}).One(&ct)
	return err
}

func AllContentTypes(w *wrapper.Wrapper) ([]ContentType, error) {
	cl := make([]ContentType, 0)
	c := w.DbSession.DB("").C("content_types")
	i := c.Find(nil).Limit(50).Iter()
	err := i.All(&cl)
	if err != nil {
		return nil, err
	}
	return cl, nil
}
