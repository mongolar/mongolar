// Controller Map is a list of API endpoints that allow the admin to
// compile the api calls needed to render a site with Mongolar.

// A controller is any function that will take a wrapper as an argument.

package controller

import (
	"fmt"
	"github.com/mongolar/mongolar/models/elements"
	"github.com/mongolar/mongolar/models/paths"
	"github.com/mongolar/mongolar/services"
	"github.com/mongolar/mongolar/wrapper"
)

// The map structure for Controllers
type ControllerMap map[string]func(*wrapper.Wrapper)

// Creates a map for controllers
func NewMap() ControllerMap {
	return make(ControllerMap)
}

// The controller function to retrieve elements ids from the path
func PathValues(w *wrapper.Wrapper) {
	p := paths.NewPath()
	c := w.DbSession.DB("").C("paths")
	u := w.Request.Header.Get("CurrentPath")
	qp, err := p.PathMatch(u, "published", c)
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
		e := elements.NewElement()
		err = elements.GetById(eid, &e, w)
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
	e := elements.NewContentElement()
	err := elements.GetValidElement(w.APIParams[0], "content", &e, w)
	//return
	if err != nil {
		errmessage := fmt.Sprintf("Content not found %s : %s", w.APIParams[1], err.Error())
		w.SiteConfig.Logger.Error(errmessage)
		services.AddMessage("There was a problem loading some content on your page.", "Error", w)
		w.Serve()
		return
	}
	w.SetTemplate(e.Template)
	w.SetDynamicId(e.DynamicId)
	w.SetContent(e.ContentValues.Content)
	w.Serve()
	return
}

// The controller function for Values found directly in the controller values of the element
func WrapperValues(w *wrapper.Wrapper) {
	e := elements.NewWrapperElement()
	err := elements.GetValidElement(w.APIParams[0], "wrapper", &e, w)
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
	for _, id := range e.Elements {
		ev := make(map[string]string)
		e := elements.NewElement()
		err = elements.GetById(id, &e, w)
		if err != nil {
			errmessage := fmt.Sprintf("Content not found %s : %s", id, err.Error())
			w.SiteConfig.Logger.Error(errmessage)
		} else {
			ev["mongolartemplate"] = e.Template
			ev["mongolartype"] = e.Controller
			ev["mongolarid"] = id
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
	es := elements.NewSlugElement()
	err := elements.GetValidElement(w.APIParams[0], "slug", &es, w)
	if err != nil {
		errmessage := fmt.Sprintf("Content not found %s : %s", w.APIParams[0], err.Error())
		w.SiteConfig.Logger.Error(errmessage)
		services.AddMessage("There was a problem loading some content on your page.", "Error", w)
		w.Serve()
		return
	}
	if _, ok := es.Slugs[w.Request.Header.Get("Slug")]; !ok {
		errmessage := fmt.Sprintf("Slug content not found for query %s", w.Request.Header.Get("QueryParameter"))
		w.SiteConfig.Logger.Error(errmessage)
		services.AddMessage("There was a problem loading some content on your page.", "Error", w)
		w.Serve()
		return
	}
	id := es.Slugs[w.Request.Header.Get("Slug")]
	e := elements.NewContentElement()
	err = elements.GetById(id, &e, w)
	if err != nil {
		errmessage := fmt.Sprintf("Content not found %s : %s", w.APIParams[0], err.Error())
		w.SiteConfig.Logger.Error(errmessage)
		services.AddMessage("There was a problem loading some content on your page.", "Error", w)
		w.Serve()
		return
	}
	w.SetTemplate(e.Template)
	w.SetDynamicId(e.DynamicId)
	w.SetContent(e.ContentValues.Content)
	w.Serve()
}
