// Controller Map is a list of API endpoints that allow the admin to
// compile the api calls needed to render a site with Mongolar.

// A controller is any function that will take a wrapper as an argument.

package controller

import (
	"fmt"
	"github.com/mongolar/mongolar/services"
	"github.com/mongolar/mongolar/wrapper"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"path"
	"reflect"
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
	MongoId          bson.ObjectId          `bson:"_id,omitempty" json:"id"`
	ControllerValues map[string]interface{} `bson:"controller_values,omitempty" json:"controller_values"`
	Controller       string                 `bson:"controller" json:"controller"`
	Template         string                 `bson:"template,omitempty" json:"template,omitempty"`
	DynamicId        string                 `bson:"dynamic_id,omitempty" json:"dynamic_id"`
	Title            string                 `bson:"title" json:"title"`
	Classes          string                 `bson:"classes"`
}

// Constructor for elements
func NewElement() Element {
	cv := make(map[string]interface{})
	e := Element{ControllerValues: cv}
	return e
}

// Query one element
func (e *Element) GetElement(b bson.M, w *wrapper.Wrapper) error {
	c := w.DbSession.DB("").C("elements")
	err := c.Find(b).One(&e)
	return err
}

// Get one element given an id
func (e *Element) GetById(i string, w *wrapper.Wrapper) error {
	b := bson.M{"_id": bson.ObjectIdHex(i)}
	err := e.GetElement(b, w)
	return err
}

// Get one element by id and controller path, most common query because you should validate your controller against the id
func (e *Element) GetValidElement(i string, c string, w *wrapper.Wrapper) error {
	b := bson.M{"_id": bson.ObjectIdHex(i), "controller": c}
	err := e.GetElement(b, w)
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

//The designated structure for all elements
type Path struct {
	MongoId  bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Path     string        `bson:"path" json:"path"`
	Wildcard bool          `bson:"wildcard" json:"wildcard"`
	Elements []string      `bson:"elements,omitempty" json:"elements"`
	Template string        `bson:"template" json:"template"`
	Status   string        `bson:"status" json:"status"`
	Title    string        `bson:"title" json:"title"`
}

// Constructor for elements
func NewPath() Path {
	e := make([]string, 0)
	p := Path{Elements: e}
	return p
}

// Get Path by Id
func (p *Path) GetById(i string, w *wrapper.Wrapper) error {
	c := w.DbSession.DB("").C("paths")
	err := c.FindId(bson.ObjectIdHex(i)).One(&p)
	return err
}

// Get all Paths
func PathList(w *wrapper.Wrapper) ([]Path, error) {
	pl := make([]Path, 0)
	c := w.DbSession.DB("").C("paths")
	i := c.Find(nil).Limit(50).Iter()
	err := i.All(&pl)
	if err != nil {
		return nil, err
	}
	return pl, nil
}

// The controller function to retrieve elements ids from the path
func PathValues(w *wrapper.Wrapper) {
	p := NewPath()
	c := w.DbSession.DB("").C("paths")
	u := w.Request.Header.Get("CurrentPath")
	qp, err := p.pathMatch(u, "published", c)
	if err != nil {
		if err.Error() == "not found" {
			if "/"+w.SiteConfig.FourOFour != u {
				services.Redirect("/"+w.SiteConfig.FourOFour, w)
				w.Serve()
				return
			} else {
				services.AddMessage("There was a problem with the system.", "Error", w)
				w.Serve()
				return
			}

		}

	}
	var v []map[string]string
	for _, eid := range p.Elements {
		ev := make(map[string]string)
		e := NewElement()
		err = e.GetById(eid, w)
		if err != nil {
			errmessage := fmt.Sprintf("Content not found %s : %s", eid, err.Error())
			w.SiteConfig.Logger.Error(errmessage)
		}
		ev["mongolartemplate"] = e.Template
		ev["mongolartype"] = e.Controller
		ev["mongolarclasses"] = e.Classes
		ev["mongolarid"] = eid
		v = append(v, ev)
	}
	w.SetPayload("mongolar_slug", qp)
	w.SetContent(v)
	w.SetTemplate(p.Template)
	w.Serve()
	return
}

// Path matching query
func (p *Path) pathMatch(u string, s string, c *mgo.Collection) (string, error) {
	var rejects []string
	w := false
	var err error
	for {
		b := bson.M{"path": u, "wildcard": w, "status": s}
		err = c.Find(b).One(p)
		w = true
		// If query doesnt return anything
		if err != nil {
			rejects = append([]string{path.Base(u)}, rejects...)
			u = path.Dir(u)
			if u == "/" {
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
	v := make(map[string]interface{})
	v[w.APIParams[0]] = w.SiteConfig.PublicValues[w.APIParams[0]]
	w.SetContent(v)
	w.Serve()
	return
}

// The controller function for Values found directly in the controller values of the element
func ContentValues(w *wrapper.Wrapper) {
	e := NewElement()
	err := e.GetValidElement(w.APIParams[0], "content", w)
	if err != nil {
		errmessage := fmt.Sprintf("Content not found %s : %s", w.APIParams[1], err.Error())
		w.SiteConfig.Logger.Error(errmessage)
		services.AddMessage("There was a problem loading some content on your page.", "Error", w)
		w.Serve()
		return
	}
	if val, ok := e.ControllerValues["content"]; ok {
		w.SetTemplate(e.Template)
		w.SetDynamicId(e.DynamicId)
		w.SetContent(val)
		w.Serve()
		return
	}
	errmessage := fmt.Sprintf("Content not found %s", w.APIParams[0])
	w.SiteConfig.Logger.Error(errmessage)
	services.AddMessage("There was a problem loading some content on your page.", "Error", w)
	w.Serve()
	return
}

// The controller function for Values found directly in the controller values of the element
func WrapperValues(w *wrapper.Wrapper) {
	e := NewElement()
	err := e.GetValidElement(w.APIParams[0], "wrapper", w)
	if err != nil {
		errmessage := fmt.Sprintf("Content not found %s : %s", w.APIParams[0], err.Error())
		w.SiteConfig.Logger.Error(errmessage)
		services.AddMessage("There was a problem loading some content on your page.", "Error", w)
		w.Serve()
		return
	}
	w.SetTemplate(e.Template)
	w.SetDynamicId(e.DynamicId)
	var v []map[string]string
	es := reflect.ValueOf(e.ControllerValues["elements"])
	for i := 0; i < es.Len(); i++ {
		el := es.Index(i)
		eid := el.Interface().(string)
		ev := make(map[string]string)
		e := NewElement()
		err = e.GetById(eid, w)
		if err != nil {
			errmessage := fmt.Sprintf("Content not found %s : %s", eid, err.Error())
			w.SiteConfig.Logger.Error(errmessage)
		} else {
			ev["mongolartemplate"] = e.Template
			ev["mongolartype"] = e.Controller
			ev["mongolarid"] = eid
			ev["mongolarclasses"] = e.Classes
			v = append(v, ev)
		}
	}
	w.SetClasses(e.Classes)
	w.SetContent(v)
	w.Serve()
}

// The controller function for elements that are context specific
func SlugValues(w *wrapper.Wrapper) {
	es := NewElement()
	err := es.GetValidElement(w.APIParams[0], "slug", w)
	if err != nil {
		errmessage := fmt.Sprintf("Content not found %s : %s", w.APIParams[0], err.Error())
		w.SiteConfig.Logger.Error(errmessage)
		services.AddMessage("There was a problem loading some content on your page.", "Error", w)
		w.Serve()
		return
	}
	if _, ok := es.ControllerValues[w.Request.Header.Get("Slug")]; !ok {
		errmessage := fmt.Sprintf("Slug content not found for query %s", w.Request.Header.Get("QueryParameter"))
		w.SiteConfig.Logger.Error(errmessage)
		services.AddMessage("There was a problem loading some content on your page.", "Error", w)
		w.Serve()
		return
	}
	i := es.ControllerValues[w.Request.Header.Get("Slug")]
	e := NewElement()
	err = e.GetById(i.(string), w)
	if err != nil {
		errmessage := fmt.Sprintf("Content not found %s : %s", w.APIParams[0], err.Error())
		w.SiteConfig.Logger.Error(errmessage)
		services.AddMessage("There was a problem loading some content on your page.", "Error", w)
		w.Serve()
		return
	}
	w.SetTemplate(e.Template)
	w.SetDynamicId(e.DynamicId)
	w.SetContent(e.ControllerValues["content"])
	w.Serve()
}
