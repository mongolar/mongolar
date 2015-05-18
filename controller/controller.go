// Controller Map is a list of API endpoints that allow the admin to
// compile the api calls needed to render a site with Mongolar.

// A controller is any function that will take a wrapper as an argument.

package controller

import (
	//	"github.com/davecgh/go-spew/spew"
	"github.com/mongolar/mongolar/services"
	"github.com/mongolar/mongolar/url"
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
	Template         string                 `bson:"template" json:"template"`
	DynamicId        string                 `bson:"dynamic_id,omitempty" json:"dynamic_id"`
	Title            string                 `bson:"title"`
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
	e := make([]string, 1)
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
	pl := make([]Path, 1)
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
	qp, err := p.pathMatch(u, c)
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
			w.SiteConfig.Logger.Error("Content not found " + eid + " by " + w.Request.Host)
		}
		ev["mongolartemplate"] = e.Template
		ev["mongolartype"] = e.Controller
		ev["mongolarid"] = eid
		v = append(v, ev)
	}
	w.Writer.Header().Add("QueryParameters", qp)
	w.SetContent(v)
	w.SetTemplate(p.Template)
	w.Serve()
	return
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
	// Get second value in url path
	p := url.UrlToMap(w.Request.URL.Path)
	v := make(map[string]interface{})
	v[p[2]] = w.SiteConfig.PublicValues[p[2]]
	w.SetContent(v)
	w.Serve()
	return
}

// The controller function for Values found directly in the controller values of the element
func ContentValues(w *wrapper.Wrapper) {
	u := url.UrlToMap(w.Request.URL.Path)
	e := NewElement()
	err := e.GetValidElement(u[2], u[1], w)
	if err != nil {
		w.SiteConfig.Logger.Error("Content not found " + u[2] + " by " + w.Request.Host)
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
	w.SiteConfig.Logger.Error("Content not found " + u[2] + " by " + w.Request.Host)
	services.AddMessage("There was a problem loading some content on your page.", "Error", w)
	w.Serve()
	return
}

// The controller function for Values found directly in the controller values of the element
func WrapperValues(w *wrapper.Wrapper) {
	u := url.UrlToMap(w.Request.URL.Path)
	e := NewElement()
	err := e.GetValidElement(u[2], u[1], w)
	if err != nil {
		w.SiteConfig.Logger.Error("Content not found " + u[2] + " by " + w.Request.Host)
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
			w.SiteConfig.Logger.Error("Content not found " + eid + " by " + w.Request.Host)
		} else {
			ev["mongolartemplate"] = e.Template
			ev["mongolartype"] = e.Controller
			ev["mongolarid"] = eid
			v = append(v, ev)
		}
	}
	w.SetContent(v)
	w.Serve()
}

// The controller function for elements that are context specific
func SlugValues(w *wrapper.Wrapper) {
	u := url.UrlToMap(w.Request.URL.Path)
	es := NewElement()
	err := es.GetValidElement(u[2], u[1], w)
	if err != nil {
		w.SiteConfig.Logger.Error("Content not found " + u[2] + " by " + w.Request.Host)
		services.AddMessage("There was a problem loading some content on your page.", "Error", w)
		w.Serve()
		return
	}
	i := es.ControllerValues[w.Request.Header.Get("QueryParameter")]
	e := NewElement()
	err = e.GetById(i.(string), w)
	if err != nil {
		w.SiteConfig.Logger.Error("Content not found " + u[2] + " by " + w.Request.Host)
		services.AddMessage("There was a problem loading some content on your page.", "Error", w)
		w.Serve()
		return
	}
	w.SetTemplate(e.Template)
	w.SetDynamicId(e.DynamicId)
	w.SetContent(e.ControllerValues)
	w.Serve()
}
