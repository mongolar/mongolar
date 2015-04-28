// Controller Map is a list of API endpoints that allow the admin to
// compile the api calls needed to render a site with Mongolar.

// A controller is any function that will take a wrapper as an argument.

package controller

import (
	"fmt"
	//"github.com/davecgh/go-spew/spew"
	"github.com/jasonrichardsmith/mongolar/service/redirect"
	"github.com/jasonrichardsmith/mongolar/url"
	"github.com/jasonrichardsmith/mongolar/wrapper"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"path"
	"strings"
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
	ControllerValues map[string]interface{} `bson:"controller_values,omitempty"`
	Controller       string                 `bson:"controller"`
	Template         string                 `bson:"template"`
	DynamicId        string                 `bson:"dynamic_id,omitempty"`
}

// Constructor for elements
func NewElement() Element {
	cv := make(map[string]interface{})
	e := Element{ControllerValues: cv}
	return e
}

// Query one element
func (e *Element) getElement(b bson.M, s *mgo.Session) error {
	se := s.Copy()
	defer se.Close()
	c := se.DB("").C("elements")
	err := c.Find(b).One(&e)
	return err
}

// Get one element given an id
func (e *Element) GetById(i string, s *mgo.Session) error {
	b := bson.M{"_id": bson.ObjectIdHex(i)}
	err := e.getElement(b, s)
	return err
}

// Get one element by id and controller path, most common query because you should validate your controller against the id
func (e *Element) GetValidElement(i string, c string, s *mgo.Session) error {
	b := bson.M{"_id": bson.ObjectIdHex(i), "controller": c}
	err := e.getElement(b, s)
	return err
}

//The designated structure for all elements
type Path struct {
	MongoId  bson.ObjectId `bson:"_id,omitempty"`
	Path     string        `bson:"path"`
	Wildcard bool          `bson:"wildcard"`
	Elements []string      `bson:"elements"`
	Template string        `bson:"template"`
	Status   string        `bson:"status"`
}

// Constructor for elements
func NewPath() Path {
	e := make([]string, 1)
	p := Path{Elements: e}
	return p
}

// Get all Paths
func PathList(d *mgo.Session) {
	pl := make([]Path)
	s := w.SiteConfigDbSession.Copy()
	defer s.Close()
	c := s.DB("").C("paths")
	i := c..Find(nil).Limit(50).Iter()
	err := i.All(&pl)
	if err != nil {
	    return nil, err
	}
	return pl, nil
}

// The controller function to retrieve elements ids from the path
func PathValues(w *wrapper.Wrapper) {
	p := NewPath()
	s := w.SiteConfig.DbSession.Copy()
	defer s.Close()
	c := s.DB("").C("paths")
	u := w.Request.Header.Get("CurrentPath")
	u = "test/path"
	qp, err := p.pathMatch(u, c)
	if err != nil {
		if err.Error() == "not found" {
			if w.SiteConfig.FourOFour != u {
				redirect.Set(w.SiteConfig.FourOFour, w)
				w.Serve()
				return
			} else {
				//TODO: Log error for missing 404 path
				return
			}

		}

	}
	var v []map[string]string
	for _, eid := range p.Elements {
		ev := make(map[string]string)
		e := NewElement()
		//TODO handle error here
		err = e.GetById(eid, w.SiteConfig.DbSession)
		ev["template"] = e.Template
		ev["controller"] = e.Controller
		ev["id"] = eid
		v = append(v, ev)
	}
	w.Writer.Header().Add("QueryParameters", qp)
	w.SetContent(v)
	w.Serve()
}

// Path matching query
func (p *Path) pathMatch(u string, c *mgo.Collection) (string, error) {
	var rejects []string
	w := false
	var err error
	for {
		b := bson.M{"path": u, "wildcard": w}
		err = c.Find(b).One(p)
		w = true
		// If query doesnt return anything
		if err != nil {
			rejects = append([]string{path.Base(u)}, rejects...)
			u = path.Dir(u)
			if u == "." {
				break
			}
			continue
		}
		break
	}
	return strings.Join(rejects, "/"), err
}

// The controller function for Values found in the Site Configuration
func DomainPublicValue(w *wrapper.Wrapper) {
	// Get second value in url path
	p := url.UrlToMap(w.Request.URL.Path)
	v := make(map[string]interface{})
	v[p[1]] = w.SiteConfig.PublicValues[p[1]]
	w.SetContent(v)
	w.Serve()
	return
}

// The controller function for Values found directly in the controller values of the element
func ContentValue(w *wrapper.Wrapper) {
	u := url.UrlToMap(w.Request.URL.Path)
	e := NewElement()
	err := e.GetValidElement(u[1], u[0], w.SiteConfig.DbSession)
	fmt.Println(err)
	//TODO: Log Errors here
	w.SetTemplate(e.Template)
	w.SetDynamicId(e.DynamicId)
	w.SetContent(e.ControllerValues['content'])
	w.Serve()
}

// The controller function for Values found directly in the controller values of the element
func WrapperValue(w *wrapper.Wrapper) {
	u := url.UrlToMap(w.Request.URL.Path)
	e := NewElement()
	err := e.GetValidElement(u[1], u[0], w.SiteConfig.DbSession)
	fmt.Println(err)
	//TODO: Log Errors here
	w.SetTemplate(e.Template)
	w.SetDynamicId(e.DynamicId)
	type es []string
	var v []map[string]string
	for _, eid := range e.ControllerValues["elements"].(es) {
		ev := make(map[string]string)
		e := NewElement()
		//TODO handle error here
		err = e.GetById(eid, w.SiteConfig.DbSession)
		fmt.Println(err)
		ev["template"] = e.Template
		ev["controller"] = e.Controller
		ev["id"] = eid
		v = append(v, ev)
	}
	w.SetContent(v)
	w.Serve()
}

// The controller function for elements that are context specific
func SlugValue(w *wrapper.Wrapper) {
	u := url.UrlToMap(w.Request.URL.Path)
	es := NewElement()
	err := es.GetValidElement(u[1], u[0], w.SiteConfig.DbSession)

	fmt.Println(err)
	//TODO: Log Errors here
	i := es.ControllerValues[w.Request.Header.Get("QueryParameter")]
	e := NewElement()
	err = e.GetById(i.(string), w.SiteConfig.DbSession)
	fmt.Println(err)
	//TODO: Log Errors here
	w.SetTemplate(e.Template)
	w.SetDynamicId(e.DynamicId)
	w.SetContent(e.ControllerValues)
	w.Serve()
}

type Form struct {
	MongoId   bson.ObjectId `bson:"_id"`
	SessionId string        `bson:"session_id"`
	Required  []string      `bson:"required"`
}

func NewForm(f map[string]map[string]interface{}, s string) Form {
	i := bson.NewObjectId()
	r := make([]string, 1)
	for k, ff := range f {
		if _, ok := ff["required"]; ok {
			r = append(r, k)
		}
	}
	fo := Form{MongoId: i, SessionId: s, Required: r}
	return fo
}

func (f Form) Register(s *mgo.Session) error {
	se := s.Copy()
	defer se.Close()
	c := s.DB("").C("forms")
	err := c.Insert(s)
	//TODO Log error
	return err
}

func GetRegisteredForm(i string, s *mgo.Session) (*Form, error) {
	f := new(Form)
	se := s.Copy()
	defer se.Close()
	c := s.DB("").C("forms")
	b := bson.M{bson.ObjectIdHex("_id"): i}
	err := c.Find(b).One(f)
	return f, err
}
