package admin

import (
	"encoding/json"
	"fmt"
	"github.com/mongolar/mongolar/form"
	"github.com/mongolar/mongolar/models/elements"
	"github.com/mongolar/mongolar/models/paths"
	"github.com/mongolar/mongolar/services"
	"github.com/mongolar/mongolar/wrapper"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

func Sort(w *wrapper.Wrapper) {
	if w.Request.Method != "POST" {
		if w.APIParams[0] == "paths" {
			p := paths.NewPath()
			err := p.GetById(w.APIParams[1], w)
			if err != nil {
				errmessage := fmt.Sprintf("Path not found to sort for %s by %s", w.APIParams[1], w.Request.Host)
				w.SiteConfig.Logger.Error(errmessage)
				services.AddMessage("This path was not found", "Error", w)
				w.Serve()
				return
			}
			if len(p.Elements) == 0 {
				services.AddMessage("This path has no elements.", "Info", w)
			}
			w.SetPayload("elements", p.Elements)
			w.SetTemplate("admin/element_sorter.html")
			w.Serve()
			return
		} else if w.APIParams[0] == "elements" {
			e := elements.NewWrapperElement()
			err := elements.GetById(w.APIParams[1], &e, w)
			if err != nil {
				errmessage := fmt.Sprintf("Element not found to sort for %s by %s.", w.APIParams[1], w.Request.Host)
				w.SiteConfig.Logger.Error(errmessage)
				services.AddMessage("This element was not found", "Error", w)
				w.Serve()
				return
			} else {
				if len(e.Elements) > 0 {
					w.SetPayload("elements", e.Elements)
				} else {
					services.AddMessage("This has no elements assigned yet.", "Error", w)
				}
				w.SetTemplate("admin/element_sorter.html")
				w.Serve()
				return
			}
		}
		http.Error(w.Writer, "Forbidden", 403)
		return
	} else {
		post := make(map[string]interface{})
		err := json.NewDecoder(w.Request.Body).Decode(&post)
		if err != nil {
			errmessage := fmt.Sprintf("Unable to marshall sort elements %s by %s: %s", w.APIParams[0], w.Request.Host, err.Error())
			w.SiteConfig.Logger.Error(errmessage)
			services.AddMessage("Unable to save elements.", "Error", w)
			w.Serve()
			return
		}
		if w.APIParams[0] == "paths" {
			p := bson.M{
				"$set": bson.M{
					"elements": post["elements"],
				},
			}
			s := bson.M{"_id": bson.ObjectIdHex(w.APIParams[1])}
			c := w.DbSession.DB("").C("paths")
			err := c.Update(s, p)
			if err != nil {
				errmessage := fmt.Sprintf("Unable to update path order %s by %s: %s", w.APIParams[0], w.Request.Host, err.Error())
				w.SiteConfig.Logger.Error(errmessage)
				services.AddMessage("Unable to save elements.", "Error", w)
				w.Serve()
				return
			}
			dynamic := services.Dynamic{
				Target:     "centereditor",
				Controller: "admin/path_elements",
				Template:   "admin/path_elements.html",
				Id:         w.APIParams[1],
			}
			services.SetDynamic(dynamic, w)
			services.AddMessage("You elements have been updated.", "Success", w)
			w.Serve()
			return

		} else if w.APIParams[0] == "elements" {
			post := make(map[string]interface{})
			err := json.NewDecoder(w.Request.Body).Decode(&post)
			if err != nil {
				errmessage := fmt.Sprintf("Unable to marshall element %s by %s: %s", w.APIParams[0], w.Request.Host, err.Error())
				w.SiteConfig.Logger.Error(errmessage)
				services.AddMessage("Unable to save elements.", "Error", w)
				w.Serve()
				return
			}
			p := bson.M{
				"$set": bson.M{
					"controller_values.elements": post["elements"],
				},
			}
			s := bson.M{"_id": bson.ObjectIdHex(w.APIParams[1])}
			c := w.DbSession.DB("").C("elements")
			err = c.Update(s, p)
			if err != nil {
				errmessage := fmt.Sprintf("Unable to update element order %s by %s: %s", w.APIParams[0], w.Request.Host, err.Error())
				w.SiteConfig.Logger.Error(errmessage)
				services.AddMessage("Unable to save elements.", "Error", w)
				w.Serve()
				return
			}
			dynamic := services.Dynamic{
				Target:     w.APIParams[1],
				Controller: "admin/element",
				Template:   "admin/element.html",
				Id:         w.APIParams[1],
			}
			services.SetDynamic(dynamic, w)
			services.AddMessage("You elements have been updated.", "Success", w)
			w.Serve()
			return
		}
		http.Error(w.Writer, "Forbidden", 403)
		return
	}
}

func AddChild(w *wrapper.Wrapper) {
	var parenttype string
	if len(w.APIParams) > 1 {
		parenttype = w.APIParams[0]
	} else {
		http.Error(w.Writer, "Forbidden", 403)
		w.Serve()
		return
	}
	w.Shift()
	switch parenttype {
	case "elements":
		AddWrapperChild(w)
	case "paths":
		AddPathChild(w)
	default:
		http.Error(w.Writer, "Forbidden", 403)
		w.Serve()
	}
	return
}

func AddWrapperChild(w *wrapper.Wrapper) {
	var parentid string
	if len(w.APIParams) > 0 {
		parentid = w.APIParams[0]
	} else {
		http.Error(w.Writer, "Forbidden", 403)
		w.Serve()
		return
	}
	e := elements.NewElement()
	e.Title = "New Element"
	err := e.Save(w)
	if err != nil {
		errmessage := fmt.Sprintf("Unable to create new element  by %s : %s", w.Request.Host, err.Error())
		w.SiteConfig.Logger.Error(errmessage)
		services.AddMessage("Could not create a new element.", "Error", w)
		w.Serve()
		return
	}
	var parent elements.WrapperElement
	parent, err = elements.LoadWrapperElement(parentid, w)
	if err != nil {
		errmessage := fmt.Sprintf("Unable to loap parent element  by %s : %s", w.Request.Host, err.Error())
		w.SiteConfig.Logger.Error(errmessage)
		services.AddMessage("Could not load parent element.", "Error", w)
		w.Serve()
		return
	}
	parent.Elements = append(parent.Elements, e.MongoId.Hex())
	err = parent.Save(w)
	if err != nil {
		errmessage := fmt.Sprintf("Unable to loap parent element  by %s : %s", w.Request.Host, err.Error())
		w.SiteConfig.Logger.Error(errmessage)
		services.AddMessage("Could not load parent element.", "Error", w)
		w.Serve()
		return
	}
	dynamic := services.Dynamic{
		Target:     w.APIParams[1],
		Controller: "admin/element",
		Template:   "admin/element.html",
		Id:         w.APIParams[1],
	}
	services.SetDynamic(dynamic, w)
	services.AddMessage("You have added a new element.", "Success", w)
	w.Serve()
	return

}

func AddPathChild(w *wrapper.Wrapper) {
	var parentid string
	if len(w.APIParams) > 0 {
		parentid = w.APIParams[0]
	} else {
		http.Error(w.Writer, "Forbidden", 403)
		w.Serve()
		return
	}
	e := elements.NewElement()
	e.Title = "New Element"
	err := e.Save(w)
	if err != nil {
		errmessage := fmt.Sprintf("Unable to create new element  by %s : %s", w.Request.Host, err.Error())
		w.SiteConfig.Logger.Error(errmessage)
		services.AddMessage("Could not create a new element.", "Error", w)
		w.Serve()
		return
	}
	var parent paths.Path
	parent, err = paths.LoadPath(parentid, w)
	if err != nil {
		errmessage := fmt.Sprintf("Unable to loap path  by %s : %s", w.Request.Host, err.Error())
		w.SiteConfig.Logger.Error(errmessage)
		services.AddMessage("Could not load parent path.", "Error", w)
		w.Serve()
		return
	}
	parent.Elements = append(parent.Elements, e.MongoId.Hex())
	err = parent.Save(w)
	if err != nil {
		errmessage := fmt.Sprintf("Unable to save path by %s : %s", w.Request.Host, err.Error())
		w.SiteConfig.Logger.Error(errmessage)
		services.AddMessage("Could not add child element.", "Error", w)
		w.Serve()
		return
	}
	dynamic := services.Dynamic{
		Target:     "centereditor",
		Controller: "admin/path_elements",
		Template:   "admin/path_elements.html",
		Id:         w.APIParams[1],
	}
	services.SetDynamic(dynamic, w)
	services.AddMessage("You have added a new element.", "Success", w)
	w.Serve()
	return

}

func AddExistingChild(w *wrapper.Wrapper) {
	if w.Request.Method != "POST" {
		c := w.DbSession.DB("").C("elements")
		var elems []elements.Element
		err := c.Find(nil).Limit(50).Iter().All(&elems)
		if err != nil {
			errmessage := fmt.Sprintf("Unable to retrieve a list of all elements: %s", err.Error())
			w.SiteConfig.Logger.Error(errmessage)
			services.AddMessage("There was a problem retrieving the element list.", "Error", w)
			w.Serve()
			return
		}
		options := make([]map[string]string, 0)
		for _, element := range elems {
			option := map[string]string{"name": element.Title, "value": element.MongoId.Hex()}
			options = append(options, option)
		}
		f := form.NewForm()
		f.AddSelect("element", options).AddLabel("Element").Required()
		element := elements.NewElement()
		if w.APIParams[1] == "elements" {
			err := elements.GetById(w.APIParams[1], &element, w)
			if err != nil {
				errmessage := fmt.Sprintf("Unable to retrieve a parent element: %s", err.Error())
				w.SiteConfig.Logger.Error(errmessage)
				services.AddMessage("There was a problem retrieving your form.", "Error", w)
				w.Serve()
				return
			}
			if element.Controller == "slug" {
				f.AddText("slug", "text").Required()
			}
		}
		f.Register(w)
		w.SetTemplate("admin/form.html")
		w.SetPayload("form", f)
		w.Serve()
		return
	} else {
		post := make(map[string]string)
		err := form.GetValidFormData(w, &post)
		if err != nil {
			return
		}
		c := w.DbSession.DB("").C(w.APIParams[0])
		i := bson.M{"_id": bson.ObjectIdHex(w.APIParams[1])}
		if w.APIParams[0] == "elements" {
			if slug, ok := post["slug"]; ok {
				f := "controller_values"
				p := map[string]string{slug: post["element"]}
				err = c.Update(i, bson.M{"$push": bson.M{f: p}})
			} else {
				f := "controller_values.elements"
				err = c.Update(i, bson.M{"$push": bson.M{f: post["element"]}})
			}
		}
		if w.APIParams[0] == "paths" {
			f := "elements"
			err = c.Update(i, bson.M{"$push": bson.M{f: post["element"]}})
		}
		if err != nil {
			errmessage := fmt.Sprintf("Unable to assign child %s to %s : %s", w.APIParams[2], w.APIParams[1], err.Error())
			w.SiteConfig.Logger.Error(errmessage)
			services.AddMessage("Unable to add child element", "Error", w)
			w.Serve()
			return
		}
		services.AddMessage("Child element added", "Success", w)
		w.Serve()
		return
	}
}
