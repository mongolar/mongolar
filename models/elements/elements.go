package elements

import (
	"errors"
	"github.com/mongolar/mongolar/wrapper"
	"gopkg.in/mgo.v2/bson"
)

//The designated structure for all elements
type Element struct {
	MongoId    bson.ObjectId `bson:"_id,omitempty" json:"mongolarid"`
	Controller string        `bson:"controller" json:"mongolartype"`
	Template   string        `bson:"template,omitempty" json:"mongolartemplate"`
	DynamicId  string        `bson:"dynamic_id,omitempty" json:"mongolardyn,omitempty"`
	Title      string        `bson:"title" json:"title"`
	Classes    string        `bson:"classes" json:"mongolarclasses,omitempty"`
}

func (e *Element) Save(w *wrapper.Wrapper) error {
	return Save(e.MongoId, e, w)
}

// Constructor for elements
func NewElement() Element {
	id := bson.NewObjectId()
	e := Element{MongoId: id}
	return e
}

// Constructor for existing paths
func LoadElement(i string, w *wrapper.Wrapper) (Element, error) {
	var e Element
	err := GetById(i, &e, w)
	return e, err
}

//Save an element in its current state.
func Save(id bson.ObjectId, v interface{}, w *wrapper.Wrapper) error {
	c := w.DbSession.DB("").C("elements")
	_, err := c.Upsert(bson.M{"_id": id}, v)
	if err != nil {
		return err
	}
	return nil
}

// Query one element
func GetElement(b bson.M, v interface{}, w *wrapper.Wrapper) error {
	c := w.DbSession.DB("").C("elements")
	err := c.Find(b).One(v)
	return err
}

// Get one element given an id
func GetById(i string, v interface{}, w *wrapper.Wrapper) error {
	if !bson.IsObjectIdHex(i) {
		return errors.New("Invalid Id Hex")
	}
	b := bson.M{"_id": bson.ObjectIdHex(i)}
	err := GetElement(b, v, w)
	return err
}

// Get one element by id and controller path, most common query because you should validate your controller against the id
func GetValidElement(i string, c string, v interface{}, w *wrapper.Wrapper) error {
	if !bson.IsObjectIdHex(i) {
		return errors.New("Invalid Id Hex")
	}
	b := bson.M{"_id": bson.ObjectIdHex(i), "controller": c}
	err := GetElement(b, v, w)
	return err
}

func Delete(id string, w *wrapper.Wrapper) error {
	if !bson.IsObjectIdHex(id) {
		return errors.New("Invalid Invalid Hex")
	}
	c := w.DbSession.DB("").C("elements")
	i := bson.M{"_id": bson.ObjectIdHex(id)}
	return c.Remove(i)
}

// Get all Elements
func ElementList(w *wrapper.Wrapper) ([]Element, error) {
	el := make([]Element, 0)
	c := w.DbSession.DB("").C("elements")
	i := c.Find(nil).Limit(50).Iter()
	err := i.All(&el)
	if err != nil {
		return nil, err
	}
	return el, nil
}
