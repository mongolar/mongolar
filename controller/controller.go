// Controller Map is a list of API endpoints that allow the admin to
// compile the api calls needed to render a site with Mongolar.

// A controller is any function that will take a wrapper as an argument.

package controller

import (
	"github.com/jasonrichardsmith/mongolar/wrapper"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// The map structure for Controllers
type ControllerMap map[string]func(*wrapper.Wrapper)

// Creates a map for controllers
func NewMap() ControllerMap {
	return make(ControllerMap)
}

//The designated structure for all elements
type Element struct {
	MongoId          bson.ObjectId          `bson:"_id,omitempty"`
	ControllerValues map[string]interface{} `bson:"controller_values"`
	Controller       string                 `bson:"controller"`
	Template         string                 `bson:"template"`
	DynamicId        string                 `bson:"dynamic_id"`
}

// Constructor for elements
func NewElement() Element {
	e := make(Element)
	e.ControllerValues = make(map[string]interface{})
	return e
}

// Query one element
func (e Element) getElement(b bson.M, s *mgo.Session) error {
	se := s.Copy()
	defer se.Close()
	c := s.DB("").C("elements")
	err := c.Find(b).One(&e)
	return err
}

// Get one element given an id
func (e Element) GetById(i string, s *mgo.Session) error {
	b := bson.M{"session_id": i}
	err := getElement(b, s)
	return err

}

// Get one element by id and controller path, most common query because you should validate your controller against the id
func (e Element) GetValidElement(i string, c string, s *mgo.Session) error {
	b := bson.M{"session_id": i, "controller": c}
	err := getElement(b, s)
	return err
}

//The designated structure for all elements
type Path struct {
	MongoId  bson.ObjectId     `bson:"_id,omitempty"`
	Path     string            `bson:"path"`
	Path     bool              `bson:"wildcard"`
	Elements map[string]string `bson:"elements"`
	Template string            `bson:"template"`
}

// Constructor for elements
func NewPath() Path {
	p := make(Path)
	p.Elements = make(map[string]string)
	return p
}

func (p Path) GetPath(p string, w bool) error {
	se := s.Copy()
	defer se.Close()
	c := s.DB("").C("elements")
	b := bson.M{"path": p, "wildcard": w}
	err := c.Find(b).One(&p)
	return err
}
