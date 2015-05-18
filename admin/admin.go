package admin

import (
	//	"github.com/davecgh/go-spew/spew"
	"github.com/mongolar/mongolar/controller"
	"github.com/mongolar/mongolar/form"
	"github.com/mongolar/mongolar/services"
	"github.com/mongolar/mongolar/url"
	"github.com/mongolar/mongolar/wrapper"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"reflect"
)

type AdminMap controller.ControllerMap

type AdminMenu struct {
	MenuItems map[string]map[string]string `json:"menu_items"`
}

func NewAdmin() (*AdminMap, *AdminMenu) {
	mi := make(map[string]map[string]string)
	amenu := AdminMenu{MenuItems: mi}
	amenu.MenuItems["0"] = map[string]string{"title": "Home", "template": "admin/main_content_default.html"}
	amenu.MenuItems["1"] = map[string]string{"title": "Content", "template": "admin/content_editor.html"}
	amenu.MenuItems["2"] = map[string]string{"title": "Content Types", "template": "admin/content_types_editor.html"}
	amap := &AdminMap{
		"menu":              amenu.AdminMenu,
		"paths":             AdminPaths,
		"path_elements":     PathElements,
		"path_editor":       PathEditor,
		"element":           Element,
		"element_editor":    ElementEditor,
		"add_child":         AddChild,
		"all_content_types": GetAllContentTypes,
		"edit_content_type": EditContentType,
	}
	return amap, &amenu
}

func (a AdminMap) Admin(w *wrapper.Wrapper) {
	u := url.UrlToMap(w.Request.URL.Path)
	if c, ok := a[u[2]]; ok {
		if validateAdmin(w.Session) {
			c(w)
		} else {
			http.Error(w.Writer, "Forbidden", 403)
		}
	} else {
		http.Error(w.Writer, "Forbidden", 403)
		return
	}
}

func validateAdmin(s *wrapper.Session) bool {
	return true
}

func (a *AdminMenu) AdminMenu(w *wrapper.Wrapper) {
	w.SetContent(a)
	w.Serve()
	return
}

func AdminPaths(w *wrapper.Wrapper) {
	//TODO: Log Errors here
	pl, err := controller.PathList(w)
	if err != nil {
		w.SiteConfig.Logger.Error("Error getting path list: " + err.Error())
	} else {
		w.SetContent(pl)
	}
	w.Serve()
}

func PathEditor(w *wrapper.Wrapper) {
	if w.Post == nil {
		ops := []string{"published", "unpublished"}
		f := form.NewForm()
		f.AddText("title", "text").AddLabel("Title")
		f.AddText("path", "text").AddLabel("Path")
		f.AddText("template", "text").AddLabel("Template")
		f.AddCheckBox("wildcard").AddLabel("Wildcard")
		o := make([]map[string]string, 0)
		for _, op := range ops {
			r := map[string]string{
				"name":  op,
				"value": op,
			}
			o = append(o, r)
		}
		f.AddRadio("status", o).AddLabel("Status")
		f.AddText("path_id", "text").Hidden()
		u := url.UrlToMap(w.Request.URL.Path)
		if u[3] != "new" {
			p := controller.NewPath()
			err := p.GetById(u[3], w)
			if err != nil {
				w.SiteConfig.Logger.Error("Path not found to edit for " + u[3] + " by " + w.Request.Host)
				services.AddMessage("This path was not found", "Error", w)
				w.Serve()
			} else {
				f.FormData["wildcard"] = p.Wildcard
				f.FormData["template"] = p.Template
				f.FormData["path"] = p.Path
				f.FormData["status"] = p.Status
				f.FormData["title"] = p.Title
			}
		}
		f.Register(w)
		w.SetTemplate("admin/form.html")
		w.SetPayload("form", f)
		w.Serve()
	} else {
		_, err := form.GetValidRegForm(w.Post["form_id"].(string), w)
		if err != nil {
			w.SiteConfig.Logger.Error("Attempt to access invalid form" + w.Post["form_id"].(string) + " by " + w.Request.Host)
			services.AddMessage("Invalid Form", "Error", w)
			w.Serve()
		} else {
			c := w.DbSession.DB("").C("paths")
			if w.Post["mongolarid"].(string) == "new" {
				var wc bool
				if c, ok := w.Post["wildcard"]; ok {
					wc = c.(bool)
				} else {
					wc = false
				}
				p := controller.Path{
					Wildcard: wc,
					Path:     w.Post["path"].(string),
					Template: w.Post["template"].(string),
					Title:    w.Post["title"].(string),
					Status:   w.Post["status"].(string),
				}
				err := c.Insert(p)
				if err != nil {
					w.SiteConfig.Logger.Error("Unable to save new path by " + w.Request.Host + " : " + err.Error())
					services.AddMessage("There was a problem saving your path.", "Error", w)
				} else {
					services.AddMessage("Your element was saved.", "Success", w)
				}
			} else {
				p := bson.M{
					"$set": bson.M{
						"wildcard": w.Post["wildcard"].(bool),
						"path":     w.Post["path"].(string),
						"template": w.Post["template"].(string),
						"title":    w.Post["title"].(string),
						"status":   w.Post["status"].(string),
					},
				}
				s := bson.M{"_id": bson.ObjectIdHex(w.Post["mongolarid"].(string))}
				err := c.Update(s, p)
				if err != nil {
					w.SiteConfig.Logger.Error("Unable to asave path " + w.Post["mongolarid"].(string) +
						" by " + w.Request.Host + " : " + err.Error())
					services.AddMessage("There was a problem saving your path.", "Error", w)
				} else {
					services.AddMessage("Your path was saved.", "Success", w)
				}
				w.Serve()
				return
			}
		}

	}
}

func PathElements(w *wrapper.Wrapper) {
	u := url.UrlToMap(w.Request.URL.Path)
	p := controller.NewPath()
	err := p.GetById(u[3], w)
	if err != nil {
		w.SiteConfig.Logger.Error("Path not found to edit for " + u[3] + " by " + w.Request.Host)
		services.AddMessage("This path was not found", "Error", w)
		w.Serve()
	} else {
		w.SetPayload("path", p.Path)
		w.SetPayload("title", p.Title)
		w.SetPayload("elements", p.Elements)
		w.Serve()
	}

}

func Element(w *wrapper.Wrapper) {
	u := url.UrlToMap(w.Request.URL.Path)
	e := controller.NewElement()
	err := e.GetById(u[3], w)
	if err != nil {
		w.SiteConfig.Logger.Error("Element not found to edit for " + u[3] + " by " + w.Request.Host)
		services.AddMessage("This element was not found", "Error", w)
		w.Serve()
	} else {
		w.SetPayload("id", e.MongoId)
		w.SetPayload("title", e.Title)
		w.SetPayload("controller", e.Controller)
		if c, ok := e.ControllerValues["elements"]; ok {
			w.SetPayload("elements", c)
		}
		w.Serve()
	}
}

func ElementSort(w *wrapper.Wrapper) {
	u := url.UrlToMap(w.Request.URL.Path)
	if w.Post == nil {
		if u[3] == "paths" {
			p := controller.NewPath()
			err := p.GetById(u[4], w)
			if err != nil {
				w.SiteConfig.Logger.Error("Path not found to sort for " + u[4] + " by " + w.Request.Host)
				services.AddMessage("This path was not found", "Error", w)
				w.Serve()
			} else {
				w.SetPayload("elements", p.Elements)
				w.Serve()
				return
			}
		} else if u[3] == "elements" {
			e := controller.NewElement()
			err := e.GetById(u[4], w)
			if err != nil {
				w.SiteConfig.Logger.Error("Element not found to sort for " + u[4] + " by " + w.Request.Host)
				services.AddMessage("This element was not found", "Error", w)
				w.Serve()
				return
			} else {
				if es, ok := e.ControllerValues["elements"]; ok {
					els := reflect.ValueOf(es)
					if els.Len() > 0 {
						w.SetPayload("elements", e.ControllerValues["elements"])
					} else {
						services.AddMessage("This has no elements assigned yet.", "Error", w)
					}
				}
				w.Serve()
				return
			}
		} else {
			//TODO Illegal request
		}
	} else {
		//TODO range over Post?
		if u[3] == "paths" {
			p := bson.M{
				"$set": bson.M{
					"elements": w.Post["elements"].(string),
				},
			}
			s := bson.M{"_id": bson.ObjectIdHex(u[4])}
			c := w.DbSession.DB("").C("paths")
			c.Update(s, p)

		} else if u[3] == "elements" {
			p := bson.M{
				"$set": bson.M{
					"controller_values.elements": w.Post["elements"].(string),
				},
			}
			s := bson.M{"_id": bson.ObjectIdHex(u[4])}
			c := w.DbSession.DB("").C("elements")
			c.Update(s, p)
		} else {
			//TODO Illegal request
		}
	}
}
func ElementEditor(w *wrapper.Wrapper) {
	if w.Post == nil {

		f := form.NewForm()
		f.AddText("title", "text").AddLabel("Title")
		f.AddText("controller", "text").AddLabel("Controller")
		f.AddText("template", "text").AddLabel("Template")
		f.AddText("dynamic_id", "text").AddLabel("Dynamic Id")
		f.AddText("element_id", "text").Hidden()
		u := url.UrlToMap(w.Request.URL.Path)
		if u[3] != "new" {
			e := controller.NewElement()
			err := e.GetById(u[3], w)
			if err != nil {
				w.SiteConfig.Logger.Error("Element not found to edit for " + u[3] + " by " + w.Request.Host)
				services.AddMessage("This element was not found", "Error", w)
				w.Serve()
			} else {
				f.FormData["controller"] = e.Controller
				f.FormData["title"] = e.Title
				f.FormData["template"] = e.Template
				f.FormData["dynamic_id"] = e.DynamicId
			}
		}
		f.Register(w)
		w.SetTemplate("admin/form.html")
		w.SetPayload("form", f)
		w.Serve()
	} else {
		_, err := form.GetValidRegForm(w.Post["form_id"].(string), w)
		if err != nil {
			w.SiteConfig.Logger.Error("Attempt to access invalid form" + w.Post["FormId"].(string) + " by " + w.Request.Host)
			services.AddMessage("Invalid Form", "Error", w)
			w.Serve()
		} else {
			c := w.DbSession.DB("").C("elements")
			if w.Post["mongolarid"].(string) == "new" {
				p := controller.Element{
					Controller: w.Post["controller"].(string),
					DynamicId:  w.Post["dynamic_id"].(string),
					Template:   w.Post["template"].(string),
					Title:      w.Post["title"].(string),
				}
				err := c.Insert(p)
				if err != nil {
					w.SiteConfig.Logger.Error("Unable to save new element by " + w.Request.Host + " : " + err.Error())
					services.AddMessage("There was a problem saving your element.", "Error", w)
				} else {
					services.AddMessage("Your element was saved.", "Success", w)
				}
			} else {
				p := bson.M{
					"$set": bson.M{
						"template":   w.Post["template"].(string),
						"title":      w.Post["title"].(string),
						"dynamic_id": w.Post["dynamic_id"].(string),
						"controller": w.Post["controller"].(string),
					},
				}
				s := bson.M{"_id": bson.ObjectIdHex(w.Post["mongolarid"].(string))}
				err := c.Update(s, p)
				if err != nil {
					w.SiteConfig.Logger.Error("Unable to save element " + w.Post["mongolarid"].(string) +
						" by " + w.Request.Host + " : " + err.Error())
					services.AddMessage("There was a problem saving your element.", "Error", w)
				} else {
					services.AddMessage("Your element was saved.", "Success", w)
				}
			}
		}
	}
	w.Serve()

}

func ContentEditor(w *wrapper.Wrapper) {
	u := url.UrlToMap(w.Request.URL.Path)
	if w.Post == nil {
		e := controller.NewElement()
		err := e.GetById(u[3], w)
		if err != nil {
			w.SiteConfig.Logger.Error("Element not found to edit for " + u[3] + " by " + w.Request.Host)
			services.AddMessage("This element was not found", "Error", w)
			w.Serve()
		}
		c := w.DbSession.DB("").C("content_types")
		t := bson.M{"type": e.ControllerValues["type"]}
		var ct ContentType
		err = c.Find(t).One(&ct)
		f := form.NewForm()
		f.Fields = ct.Form.Fields
		f.FormData = e.ControllerValues
		w.SetPayload("form", f)
		w.Serve()
		return
	} else {
		e := bson.M{
			"$set": bson.M{
				"controller_values.content": w.Post,
			},
		}
		s := bson.M{"_id": bson.ObjectIdHex(u[3])}
		c := w.DbSession.DB("").C("elements")
		err := c.Update(s, e)
		if err != nil {
			w.SiteConfig.Logger.Error("Element not saved " + u[3] + " by " + w.Request.Host)
			services.AddMessage("Unable to save element.", "Error", w)
			w.Serve()
		}
	}
}

func ContentTypeEditor(w *wrapper.Wrapper) {
	u := url.UrlToMap(w.Request.URL.Path)
	if w.Post == nil {
		e := controller.NewElement()
		err := e.GetById(u[3], w)
		if err != nil {
			w.SiteConfig.Logger.Error("Element not found to edit for " + u[3] + " by " + w.Request.Host)
			services.AddMessage("This element was not found", "Error", w)
			w.Serve()
		}
		c := w.DbSession.DB("").C("content_types")
		var cts []ContentType
		c.Find(nil).Limit(50).Iter().All(&cts)
		//for _, ct := range cts {
		//	ct.Type
		//}
		f := form.NewForm()
		//f.FormData = e.ControllerValues["type"]
		w.SetPayload("form", f)
		w.Serve()
		return
	} else {
		e := bson.M{
			"$set": bson.M{
				"controller_values.type": w.Post["type"],
			},
		}
		s := bson.M{"_id": bson.ObjectIdHex(u[3])}
		c := w.DbSession.DB("").C("content_types")
		err := c.Update(s, e)
		if err != nil {
			w.SiteConfig.Logger.Error("Element not saved " + u[3] + " by " + w.Request.Host)
			services.AddMessage("Unable to save element.", "Error", w)
			w.Serve()
		}
	}
}

func AddChild(w *wrapper.Wrapper) {
	u := url.UrlToMap(w.Request.URL.Path)
	e := controller.NewElement()
	e.MongoId = bson.NewObjectId()
	e.Title = "New Element"
	c := w.DbSession.DB("").C("elements")
	err := c.Insert(e)
	if err != nil {
		w.SiteConfig.Logger.Error("Unable to create new element  by " + w.Request.Host + " : " + err.Error())
		services.AddMessage("Could not create a new element.", "Error", w)
		w.Serve()
		return
	}
	c = w.DbSession.DB("").C(u[3])
	f := ""
	if u[3] == "elements" {
		f = "controller_values.elements"
	} else if u[3] == "paths" {
		f = "elements"
	} else {
		w.SiteConfig.Logger.Error("Unable to add child element " + e.MongoId.Hex() + " to " + u[4] + " : " + err.Error())
		services.AddMessage("There was a problem, your elemeent was created but was not assigned to your "+u[3]+".", "Error", w)
		w.Serve()
		return

	}
	i := bson.M{"_id": bson.ObjectIdHex(u[4])}
	err = c.Update(i, bson.M{"$push": bson.M{f: e.MongoId.Hex()}})
	if err != nil {
		w.SiteConfig.Logger.Error("Unable to add child element " + e.MongoId.Hex() + " to " + u[4] + " : " + err.Error())
		services.AddMessage("There was a problem, your elemeent was created but was not assigned to your "+u[3]+".", "Error", w)
		w.Serve()
		return
	}
	services.AddMessage("You have added a new element.", "Success", w)
	w.Serve()
	return

}

func AllElements(w *wrapper.Wrapper) {
	c := w.DbSession.DB("").C("elements")
	var es []controller.Element
	c.Find(nil).Limit(50).Iter().All(&es)
	w.SetPayload("elements", es)
	w.Serve()
	return
}

func AddExistingChild(w *wrapper.Wrapper) {
	u := url.UrlToMap(w.Request.URL.Path)
	c := w.DbSession.DB("").C(u[3])
	i := bson.M{"_id": bson.ObjectIdHex(u[4])}
	var f string
	if u[3] == "elements" {
		f = "controller_values.elements"
	}
	if u[3] == "paths" {
		f = "elements"
	}
	err := c.Update(i, bson.M{"$push": bson.M{f: u[5]}})
	if err != nil {
		w.SiteConfig.Logger.Error("Unable to assign child " + u[5] + " to " + u[4] + " : " + err.Error())
		services.AddMessage("Unable to add child element", "Error", w)
		w.Serve()
		return
	}
	services.AddMessage("Child element added", "Error", w)
}

func Delete(w *wrapper.Wrapper) {
	u := url.UrlToMap(w.Request.URL.Path)
	c := w.DbSession.DB("").C(u[3])
	i := bson.M{"_id": bson.ObjectIdHex(u[4])}
	err := c.Remove(i)
	if err != nil {
		w.SiteConfig.Logger.Error("Unable to delete " + u[3] + " " + u[4] + " : " + err.Error())
		services.AddMessage("Unable to delete "+u[3], "Error", w)
		w.Serve()
		return
	}
	if u[3] == "elements" {
		s := bson.M{"controller_values.elements": u[4]}
		d := bson.M{"$pull": bson.M{"controller_values.elements": u[4]}}
		_, err := c.UpdateAll(s, d)
		if err != nil {
			w.SiteConfig.Logger.Error("Unable to delete reference to " + u[3] + " " + u[4] + " : " + err.Error())
			services.AddMessage("Unable to delete reference to "+u[3], "Error", w)
			w.Serve()
			return
		}
		s = bson.M{"elements": u[4]}
		d = bson.M{"$pull": bson.M{"elements": u[4]}}
		c := w.DbSession.DB("").C("paths")
		_, err = c.UpdateAll(s, d)
		if err != nil {
			w.SiteConfig.Logger.Error("Unable to delete reference to " + u[3] + " " + u[4] + " : " + err.Error())
			services.AddMessage("Unable to delete reference to "+u[3], "Error", w)
			w.Serve()
			return
		}
	}
	services.AddMessage("Successfully deleted "+u[3], "Success", w)
	w.Serve()
	return
}

type ContentType struct {
	Form    form.Form     `bson:"form" json:"content_form"`
	Type    string        `bson:"type" json:"type"`
	MongoId bson.ObjectId `bson:"_id" json:"id"`
}

func GetContentType(w *wrapper.Wrapper) {
	u := url.UrlToMap(w.Request.URL.Path)
	c := w.DbSession.DB("").C("content_types")
	i := bson.M{"_id": bson.ObjectIdHex(u[3])}
	var ct ContentType
	err := c.Find(i).One(&ct)
	if err != nil {
		w.SiteConfig.Logger.Error("Content Type not found " + u[3] + " : " + err.Error())
		services.AddMessage("Your content types was not found "+u[3], "Error", w)
		w.Serve()
		return
	}
	w.SetPayload("content_type", ct)
	w.Serve()
	return
}

func EditContentType(w *wrapper.Wrapper) {
	u := url.UrlToMap(w.Request.URL.Path)
	c := w.DbSession.DB("").C("content_types")
	i := bson.M{"_id": bson.ObjectIdHex(u[3])}
	var ct ContentType
	err := c.Find(i).One(&ct)
	if err != nil {
		w.SiteConfig.Logger.Error("Content Type not found " + u[3] + " : " + err.Error())
		services.AddMessage("Your content types was not found "+u[3], "Error", w)
		w.Serve()
		return
	}
	f := form.NewForm()
	f.AddRepeatSection("elements", "Add another field", FieldFormGroup())
	test := []map[string]string{map[string]string{"key": "test"}}
	f.FormData["elements"] = test
	w.SetPayload("form", f)
	w.SetTemplate("admin/form.html")
	w.Serve()
	return
}

func GetAllContentTypes(w *wrapper.Wrapper) {
	c := w.DbSession.DB("").C("content_types")
	var cts []ContentType
	c.Find(nil).Limit(50).Iter().All(&cts)
	w.SetPayload("content_types", cts)
	w.Serve()
	return
}

func FieldFormGroup() []*form.Field {
	ft := []map[string]string{
		map[string]string{"name": "Text Field", "value": "text"},
		map[string]string{"name": "TextArea Field", "value": "textarea"},
		map[string]string{"name": "Radio Buttons", "value": "radio"},
		map[string]string{"name": "Checkbox", "value": "checkbox"},
	}
	f := form.NewForm()
	f.AddRadio("field_type", ft).AddLabel("Field Type").Required()
	f.AddText("key", "text").AddLabel("Key").Required()
	f.AddText("label", "text").AddLabel("Label")
	f.AddText("placeholder", "text").AddLabel("Placeholder")
	f.AddTextArea("options").AddLabel("Options").AddHideExpression("form.FormData.field_type != 'radio'")
	f.AddText("cols", "text").AddLabel("Columns").AddHideExpression("form.FormData.field_type != 'textarea'")
	f.AddText("rows", "text").AddLabel("Rows").AddHideExpression("form.FormData.field_type != 'textarea'")
	return f.Fields
}
