package elements

import (
	"errors"
	"github.com/mongolar/mongolar/wrapper"
	"gopkg.in/mgo.v2/bson"
)

//The designated structure for all elements
type Element struct {
	MongoId          bson.ObjectId                     `bson:"_id,omitempty" json:"mongolarid"`
	ControllerValues map[string]map[string]interface{} `bson:"controller_values,inline" json:"-"`
	Controller       string                            `bson:"controller" json:"mongolartype"`
	Template         string                            `bson:"template,omitempty" json:"mongolartemplate"`
	DynamicId        string                            `bson:"dynamic_id,omitempty" json:"mongolardyn,omitempty"`
	Title            string                            `bson:"title" json:"title"`
	Classes          string                            `bson:"classes" json:"classes,omitempty"`
}

// Constructor for elements
func NewElement() Element {
	cv := make(map[string]map[string]interface{})
	id := bson.NewObjectId()
	e := Element{MongoId: id, ControllerValues: cv}
	return e
}

// Constructor for existing paths
func LoadElement(i string, w *wrapper.Wrapper) (Element, error) {
	cv := make(map[string]map[string]interface{})
	e := Element{ControllerValues: cv}
	err := GetById(i, e, w)
	return e, err
}

//Save an element in its current state.
func (e *Element) Save(w *wrapper.Wrapper) error {
	if !e.MongoId.Valid() {
		e.MongoId = bson.NewObjectId()
	}
	if e.Controller == "" {
		return errors.New("Controller required")
	}
	if e.Template == "" {
		return errors.New("Template required")
	}
	c := w.DbSession.DB("").C("elements")
	_, err := c.Upsert(e.MongoId, e)
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
