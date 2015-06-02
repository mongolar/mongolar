// Admin is a series of controllers to manage a Mongolar site

package admin

import (
	"fmt"
	"github.com/mongolar/mongolar/controller"
	"github.com/mongolar/mongolar/form"
	"github.com/mongolar/mongolar/oauthlogin"
	"github.com/mongolar/mongolar/services"
	"github.com/mongolar/mongolar/url"
	"github.com/mongolar/mongolar/wrapper"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"reflect"
	"strings"
)

type AdminMap controller.ControllerMap

// A series of menu items to render on the admin page
type AdminMenu struct {
	MenuItems []map[string]string `json:"menu_items"`
}

// Build a new admin instance and return AdminMap and menu to be altered
func NewAdmin() (*AdminMap, *AdminMenu) {
	amenu := AdminMenu{
		MenuItems: []map[string]string{
			map[string]string{"title": "Home", "template": "admin/main_content_default.html"},
			map[string]string{"title": "Content", "template": "admin/content_editor.html"},
			map[string]string{"title": "Content Types", "template": "admin/content_types_editor.html"},
		},
	}
	amap := &AdminMap{
		"menu":               amenu.AdminMenu,
		"paths":              AdminPaths,
		"path_elements":      PathElements,
		"path_editor":        PathEditor,
		"element":            Element,
		"element_editor":     ElementEditor,
		"add_child":          AddChild,
		"add_existing_child": AddExistingChild,
		"all_content_types":  GetAllContentTypes,
		"edit_content_type":  EditContentType,
		"delete":             Delete,
		"sort_children":      Sort,
		"content":            ContentEditor,
		"content_type":       ContentTypeEditor,
	}
	return amap, &amenu
}

//Main controller for all admin functions
func (a AdminMap) Admin(w *wrapper.Wrapper) {
	u := url.UrlToMap(w.Request.URL.Path)
	if c, ok := a[u[2]]; ok {
		if validateAdmin(w) {
			c(w)
		}
		return
	} else {
		http.Error(w.Writer, "Forbidden", 403)
		return
	}
}

func validateAdmin(w *wrapper.Wrapper) bool {
	user := new(oauthlogin.User)
	err := user.Get(w)
	if err != nil {
		services.Redirect(w.SiteConfig.LoginURLs["login"], w)
		w.Serve()
		return false
	}
	if user.Roles != nil {
		for _, r := range user.Roles {
			if r == "admin" {
				return true
			}
		}
	}
	services.Redirect(w.SiteConfig.LoginURLs["access_denied"], w)
	w.Serve()
	return false
}

func (a *AdminMenu) AdminMenu(w *wrapper.Wrapper) {
	w.SetContent(a)
	w.Serve()
	return
}

func AdminPaths(w *wrapper.Wrapper) {
	pl, err := controller.PathList(w)
	if err != nil {
		services.AddMessage("There was an error retrieving your site paths", "Error", w)
		errmessage := fmt.Sprintf("Error getting path list: %s", err.Error())
		w.SiteConfig.Logger.Error(errmessage)
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
				errmessage := fmt.Sprintf("Could not retrieve path %s by %s: %s", u[3], w.Request.Host, err.Error())
				w.SiteConfig.Logger.Error(errmessage)
				services.AddMessage("Error retrieving path information.", "Error", w)
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
		_, err := form.GetValidRegFormM(w.Post["form_id"].(string), w)
		if err != nil {
			return
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
					errmessage := fmt.Sprintf("Unable to save new path by %s: %s", w.Request.Host, err.Error())
					w.SiteConfig.Logger.Error(errmessage)
					services.AddMessage("There was a problem saving your path.", "Error", w)
					w.Serve()
					return
				}
				services.AddMessage("Your path was saved.", "Success", w)
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
					errmessage := fmt.Sprintf("Unable to save path %s by %s: %s", w.Post["mongolarid"].(string),
						w.Request.Host, err.Error())
					w.SiteConfig.Logger.Error(errmessage)
					services.AddMessage("There was a problem saving your path.", "Error", w)
					w.Serve()
					return
				} else {
					services.AddMessage("Your path was saved.", "Success", w)
				}
			}
			dynamic := services.Dynamic{
				Target:     "pathbar",
				Controller: "admin/paths",
				Template:   "admin/path_list.html",
			}
			services.SetDynamic(dynamic, w)
			w.Serve()
			return
		}

	}
}

func PathElements(w *wrapper.Wrapper) {
	u := url.UrlToMap(w.Request.URL.Path)
	p := controller.NewPath()
	err := p.GetById(u[3], w)
	if err != nil {
		errmessage := fmt.Sprintf("Path not found to edit for %s by %s ", u[3], w.Request.Host)
		w.SiteConfig.Logger.Error(errmessage)
		services.AddMessage("This path was not found", "Error", w)
		w.Serve()
	} else {
		w.SetPayload("path", p.Path)
		w.SetPayload("id", u[3])
		w.SetPayload("title", p.Title)
		w.SetPayload("elements", p.Elements)
		if len(p.Elements) == 0 {
			services.AddMessage("This path has no elements.", "Info", w)
		}
		w.Serve()
	}

}

func Element(w *wrapper.Wrapper) {
	u := url.UrlToMap(w.Request.URL.Path)
	e := controller.NewElement()
	if _, ok := u[3]; !ok {
		http.Error(w.Writer, "Forbidden", 403)
		return
	}
	err := e.GetById(u[3], w)
	if err != nil {
		errmessage := fmt.Sprintf("Element not found to edit for %s by %s.", u[3], w.Request.Host)
		w.SiteConfig.Logger.Error(errmessage)
		services.AddMessage("This element was not found", "Error", w)
	} else {
		w.SetPayload("id", e.MongoId.Hex())
		w.SetPayload("title", e.Title)
		w.SetPayload("controller", e.Controller)
		if e.Controller == "wrapper" {
			w.SetDynamicId(e.MongoId.Hex())
		}
		if c, ok := e.ControllerValues["elements"]; ok {
			w.SetPayload("elements", c)
		}
	}
	w.Serve()
}

func Sort(w *wrapper.Wrapper) {
	u := url.UrlToMap(w.Request.URL.Path)
	if w.Post == nil {
		if u[3] == "paths" {
			p := controller.NewPath()
			err := p.GetById(u[4], w)
			if err != nil {
				errmessage := fmt.Sprintf("Path not found to sort for %s by %s", u[4], w.Request.Host)
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
		} else if u[3] == "elements" {
			e := controller.NewElement()
			err := e.GetById(u[4], w)
			if err != nil {
				errmessage := fmt.Sprintf("Element not found to sort for %s by %s.", u[4], w.Request.Host)
				w.SiteConfig.Logger.Error(errmessage)
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
		if u[3] == "paths" {
			p := bson.M{
				"$set": bson.M{
					"elements": w.Post["elements"],
				},
			}
			s := bson.M{"_id": bson.ObjectIdHex(u[4])}
			c := w.DbSession.DB("").C("paths")
			err := c.Update(s, p)
			if err != nil {
				errmessage := fmt.Sprintf("Unable to update path order %s by %s: %s", u[3], w.Request.Host, err.Error())
				w.SiteConfig.Logger.Error(errmessage)
				services.AddMessage("Unable to save elements.", "Error", w)
				w.Serve()
				return
			}
			dynamic := services.Dynamic{
				Target:     "centereditor",
				Controller: "admin/path_elements",
				Template:   "admin/path_elements.html",
				Id:         u[4],
			}
			services.SetDynamic(dynamic, w)
			services.AddMessage("You elements have been updated.", "Success", w)
			w.Serve()
			return

		} else if u[3] == "elements" {
			p := bson.M{
				"$set": bson.M{
					"controller_values.elements": w.Post["elements"],
				},
			}
			s := bson.M{"_id": bson.ObjectIdHex(u[4])}
			c := w.DbSession.DB("").C("elements")
			err := c.Update(s, p)
			if err != nil {
				errmessage := fmt.Sprintf("Unable to update element order %s by %s: %s", u[3], w.Request.Host, err.Error())
				w.SiteConfig.Logger.Error(errmessage)
				services.AddMessage("Unable to save elements.", "Error", w)
				w.Serve()
				return
			}
			dynamic := services.Dynamic{
				Target:     u[4],
				Controller: "admin/element",
				Template:   "admin/element.html",
				Id:         u[4],
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
func ElementEditor(w *wrapper.Wrapper) {
	if w.Post == nil {

		f := form.NewForm()
		f.AddText("title", "text").AddLabel("Title")
		f.AddText("controller", "text").AddLabel("Controller")
		f.AddText("template", "text").AddLabel("Template")
		f.AddText("dynamic_id", "text").AddLabel("Dynamic Id")
		f.AddText("classes", "text").AddLabel("Classes")
		f.AddText("element_id", "text").Hidden()
		u := url.UrlToMap(w.Request.URL.Path)
		if u[3] != "new" {
			e := controller.NewElement()
			err := e.GetById(u[3], w)
			if err != nil {
				errmessage := fmt.Sprintf("Element not found to edit for %s by %s", u[3], w.Request.Host)
				w.SiteConfig.Logger.Error(errmessage)
				services.AddMessage("This element was not found", "Error", w)
				w.Serve()
				return
			}
			f.FormData["controller"] = e.Controller
			f.FormData["title"] = e.Title
			f.FormData["template"] = e.Template
			f.FormData["dynamic_id"] = e.DynamicId
			f.FormData["classes"] = e.Classes
		}
		f.Register(w)
		w.SetTemplate("admin/form.html")
		w.SetPayload("form", f)
	} else {
		_, err := form.GetValidRegFormM(w.Post["form_id"].(string), w)
		if err != nil {
			return
		} else {
			c := w.DbSession.DB("").C("elements")
			if w.Post["mongolarid"].(string) == "new" {
				p := controller.Element{
					Controller: w.Post["controller"].(string),
					DynamicId:  w.Post["dynamic_id"].(string),
					Template:   w.Post["template"].(string),
					Title:      w.Post["title"].(string),
					Classes:    w.Post["classes"].(string),
				}
				err := c.Insert(p)
				if err != nil {
					errmessage := fmt.Sprintf("Unable to save new element by %s : %s", w.Request.Host, err.Error())
					w.SiteConfig.Logger.Error(errmessage)
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
						"classes":    w.Post["classes"].(string),
					},
				}
				s := bson.M{"_id": bson.ObjectIdHex(w.Post["mongolarid"].(string))}
				err := c.Update(s, p)
				if err != nil {
					errmessage := fmt.Sprintf("Unable to save element %s by %s : %s",
						w.Post["mongolarid"].(string), w.Request.Host, err.Error())
					w.SiteConfig.Logger.Error(errmessage)
					services.AddMessage("There was a problem saving your element.", "Error", w)
				} else {
					services.AddMessage("Your element was saved.", "Success", w)
					dynamic := services.Dynamic{
						Target:     w.Post["mongolarid"].(string),
						Id:         w.Post["mongolarid"].(string),
						Controller: "admin/element",
						Template:   "admin/element.html",
					}
					services.SetDynamic(dynamic, w)
				}
			}
		}
	}
	w.Serve()
	return
}

func ContentTypeEditor(w *wrapper.Wrapper) {
	u := url.UrlToMap(w.Request.URL.Path)
	if w.Post == nil {
		e := controller.NewElement()
		err := e.GetById(u[3], w)
		if err != nil {
			errmessage := fmt.Sprintf("Element not found to edit for %s by %s", u[3], w.Request.Host)
			w.SiteConfig.Logger.Error(errmessage)
			services.AddMessage("This element was not found", "Error", w)
			w.Serve()
			return
		}
		c := w.DbSession.DB("").C("content_types")
		var cts []ContentType
		err = c.Find(nil).Limit(50).Iter().All(&cts)
		if err != nil {
			errmessage := fmt.Sprintf("Unable to query all Content Types: %s", err.Error())
			w.SiteConfig.Logger.Error(errmessage)
			services.AddMessage("Unable to retrieve content types.", "Error", w)
			w.Serve()
			return
		}
		f := form.NewForm()
		opts := make([]map[string]string, 0)
		for _, ct := range cts {
			opt := map[string]string{
				"name":  ct.Type,
				"value": ct.Type,
			}
			opts = append(opts, opt)
		}
		f.AddSelect("type", opts)
		if t, ok := e.ControllerValues["type"]; ok {
			f.FormData["type"] = t.(string)
		} else {
			f.FormData["type"] = ""
		}
		f.Register(w)
		w.SetTemplate("admin/form.html")
		w.SetPayload("form", f)
	} else {
		_, err := form.GetValidRegFormM(w.Post["form_id"].(string), w)
		if err != nil {
			return
		}
		e := bson.M{
			"$set": bson.M{
				"controller_values.type": w.Post["type"],
			},
		}
		s := bson.M{"_id": bson.ObjectIdHex(w.Post["mongolarid"].(string))}
		c := w.DbSession.DB("").C("elements")
		err = c.Update(s, e)
		if err != nil {
			errmessage := fmt.Sprintf("Element not saved %s by %s", u[3], w.Request.Host)
			w.SiteConfig.Logger.Error(errmessage)
			services.AddMessage("Unable to save element.", "Error", w)
		} else {
			services.AddMessage("Element content type saved.", "Success", w)
		}
	}
	w.Serve()
	return
}

func ContentEditor(w *wrapper.Wrapper) {
	u := url.UrlToMap(w.Request.URL.Path)
	e := controller.NewElement()
	err := e.GetById(u[3], w)
	if err != nil {
		errmessage := fmt.Sprintf("Element not found to edit for %s by %s", u[3], w.Request.Host)
		w.SiteConfig.Logger.Error(errmessage)
		services.AddMessage("This element was not found", "Error", w)
		w.Serve()
		return
	}
	c := w.DbSession.DB("").C("content_types")
	var ct ContentType
	s := bson.M{"type": e.ControllerValues["type"]}
	err = c.Find(s).One(&ct)
	if err != nil {
		errmessage := fmt.Sprintf("Unable to find content type %s : %s", e.ControllerValues["type"], err.Error())
		w.SiteConfig.Logger.Error(errmessage)
		services.AddMessage("Unable to find content type.", "Error", w)
		w.Serve()
		return
	}
	if w.Post == nil {
		if e.Controller != "content" {
			http.Error(w.Writer, "Forbidden", 403)
			return
		}
		f := form.NewForm()
		f.Fields = ct.Form
		if content, ok := e.ControllerValues["content"]; ok {
			f.FormData = content.(map[string]interface{})
		} else {
			f.FormData = make(map[string]interface{})
		}
		f.Register(w)
		w.SetPayload("form", f)
	} else {
		_, err := form.GetValidRegFormM(w.Post["form_id"].(string), w)
		if err != nil {
			return
		}

		content_values := make(map[string]string)
		for _, field := range ct.Form {
			content_values[field.Key] = w.Post[field.Key].(string)
		}
		e := bson.M{
			"$set": bson.M{
				"controller_values.content": content_values,
			},
		}
		s := bson.M{"_id": bson.ObjectIdHex(w.Post["mongolarid"].(string))}
		c := w.DbSession.DB("").C("elements")
		err = c.Update(s, e)
		if err != nil {
			errmessage := fmt.Sprintf("Element not saved %s by %s", u[3], w.Request.Host)
			w.SiteConfig.Logger.Error(errmessage)
			services.AddMessage("Unable to save element.", "Error", w)
		} else {
			services.AddMessage("Element content saved.", "Success", w)
			dynamic := services.Dynamic{
				Target:     w.Post["mongolarid"].(string),
				Id:         w.Post["mongolarid"].(string),
				Controller: "admin/element",
				Template:   "admin/element.html",
			}
			services.SetDynamic(dynamic, w)
		}
	}
	w.Serve()
	return
}

func AddChild(w *wrapper.Wrapper) {
	u := url.UrlToMap(w.Request.URL.Path)
	e := controller.NewElement()
	e.MongoId = bson.NewObjectId()
	e.Title = "New Element"
	c := w.DbSession.DB("").C("elements")
	err := c.Insert(e)
	if err != nil {
		errmessage := fmt.Sprintf("Unable to create new element  by %s : %s", w.Request.Host, err.Error())
		w.SiteConfig.Logger.Error(errmessage)
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
		errmessage := fmt.Sprintf("Invalid parent item type %s by %s", u[3], w.Request.Host)
		w.SiteConfig.Logger.Error(errmessage)
		message := fmt.Sprintf("Attempt to assign child to illegal parent %s.", u[3])
		services.AddMessage(message, "Error", w)
		w.Serve()
		return

	}
	i := bson.M{"_id": bson.ObjectIdHex(u[4])}
	err = c.Update(i, bson.M{"$push": bson.M{f: e.MongoId.Hex()}})
	if err != nil {
		errmessage := fmt.Sprintf("Unable to add child element %s to %s : %s", e.MongoId.Hex(), u[4], err.Error())
		w.SiteConfig.Logger.Error(errmessage)
		message := fmt.Sprintf("There was a problem, your elemeent was created but was not assigned to your %s.", u[3])
		services.AddMessage(message, "Error", w)
		w.Serve()
		return
	}
	var dynamic services.Dynamic
	if u[3] == "elements" {
		dynamic = services.Dynamic{
			Target:     u[4],
			Controller: "admin/element",
			Template:   "admin/element.html",
			Id:         u[4],
		}
	} else if u[3] == "paths" {
		dynamic = services.Dynamic{
			Target:     "centereditor",
			Controller: "admin/path_elements",
			Template:   "admin/path_elements.html",
			Id:         u[4],
		}
	}
	services.SetDynamic(dynamic, w)
	services.AddMessage("You have added a new element.", "Success", w)
	w.Serve()
	return

}

func AllElements(w *wrapper.Wrapper) {
	c := w.DbSession.DB("").C("elements")
	var es []controller.Element
	err := c.Find(nil).Limit(50).Iter().All(&es)
	if err != nil {
		errmessage := fmt.Sprintf("Unable to retrieve a list of all elements: %s", err.Error())
		w.SiteConfig.Logger.Error(errmessage)
		message := fmt.Sprintf("There was a problem retrieving the element list.")
		services.AddMessage(message, "Error", w)
	}
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
		errmessage := fmt.Sprintf("Unable to assign child %s to %s : %s", u[5], u[4], err.Error())
		w.SiteConfig.Logger.Error(errmessage)
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
		errmessage := fmt.Sprintf("Unable to delete %s %s : %s", u[3], u[4], err.Error())
		w.SiteConfig.Logger.Error(errmessage)
		services.AddMessage("Unable to delete.", "Error", w)
		w.Serve()
		return
	}
	if u[3] == "elements" {
		s := bson.M{"controller_values.elements": u[4]}
		d := bson.M{"$pull": bson.M{"controller_values.elements": u[4]}}
		_, err := c.UpdateAll(s, d)
		if err != nil {
			errmessage := fmt.Sprintf("Unable to delete reference to %s %s : %s", u[3], u[4], err.Error())
			w.SiteConfig.Logger.Error(errmessage)
			services.AddMessage("Unable to delete all references to your element.", "Error", w)
			w.Serve()
			return
		}
		s = bson.M{"elements": u[4]}
		d = bson.M{"$pull": bson.M{"elements": u[4]}}
		c := w.DbSession.DB("").C("paths")
		_, err = c.UpdateAll(s, d)
		if err != nil {
			errmessage := fmt.Sprintf("Unable to delete reference to %s %s : %s", u[3], u[4], err.Error())
			w.SiteConfig.Logger.Error(errmessage)
			services.AddMessage("Unable to delete all references to your element.", "Error", w)
			w.Serve()
			return
		}
		dynamic := services.Dynamic{
			Target:   u[4],
			Template: "default.html",
		}
		services.SetDynamic(dynamic, w)
	}
	if u[3] == "paths" {
		dynamic := services.Dynamic{
			Target:     "pathbar",
			Controller: "admin/paths",
			Template:   "admin/path_list.html",
		}
		services.SetDynamic(dynamic, w)
	}
	services.AddMessage("Successfully deleted "+u[3], "Success", w)
	w.Serve()
	return
}

type ContentType struct {
	Form    []*form.Field `bson:"form,omitempty" json:"content_form"`
	Type    string        `bson:"type,omitempty" json:"type"`
	MongoId bson.ObjectId `bson:"_id" json:"id"`
}

func GetContentType(w *wrapper.Wrapper) {
	u := url.UrlToMap(w.Request.URL.Path)
	c := w.DbSession.DB("").C("content_types")
	i := bson.M{"_id": bson.ObjectIdHex(u[3])}
	var ct ContentType
	err := c.Find(i).One(&ct)
	if err != nil {
		errmessage := fmt.Sprintf("Content Type not found %s : %s", u[3], err.Error())
		w.SiteConfig.Logger.Error(errmessage)
		services.AddMessage("Your content types was not found.", "Error", w)
		w.Serve()
		return
	}
	w.SetPayload("content_type", ct)
	w.Serve()
	return
}

func EditContentType(w *wrapper.Wrapper) {
	if w.Post == nil {
		u := url.UrlToMap(w.Request.URL.Path)
		f := form.NewForm()
		ct := new(ContentType)
		if u[3] != "new" {
			c := w.DbSession.DB("").C("content_types")
			i := bson.M{"_id": bson.ObjectIdHex(u[3])}
			err := c.Find(i).One(ct)
			if err != nil {
				errmessage := fmt.Sprintf("Content Type not found %s : %s", u[3], err.Error())
				w.SiteConfig.Logger.Error(errmessage)
				services.AddMessage("Your content types was not found ", "Error", w)
				w.Serve()
				return
			}
			var elements []map[string]interface{}
			for _, field := range ct.Form {
				element := make(map[string]interface{})
				element["type"] = field.Type
				element["key"] = field.Key
				element["label"] = field.TemplateOptions.Label
				element["placeholder"] = field.TemplateOptions.Placeholder
				element["rows"] = field.TemplateOptions.Rows
				element["cols"] = field.TemplateOptions.Cols
				element["options"] = ""
				for _, opt := range field.TemplateOptions.Options {
					element["options"] = fmt.Sprintf("%s%s|%s\n", element["options"], opt["name"], opt["value"])
				}
				elements = append(elements, element)
			}
			f.FormData["elements"] = elements
			f.FormData["content_type"] = ct.Type
		} else {
			fd := make([]map[string]string, 0)
			f.FormData["elements"] = fd
			f.FormData["content_type"] = ""
		}
		f.AddText("content_type", "text").AddLabel("Content Type Name")
		f.AddRepeatSection("elements", "Add another field", FieldFormGroup())
		f.Register(w)
		w.SetPayload("form", f)
		w.SetTemplate("admin/form.html")
		w.Serve()
		return
	} else {
		_, err := form.GetValidRegFormM(w.Post["form_id"].(string), w)
		if err != nil {
			return
		}
		elements := reflect.ValueOf(w.Post["elements"])
		f := form.NewForm()
		for i := 0; i < elements.Len(); i++ {
			var field *form.Field
			element := elements.Index(i).Interface().(map[string]interface{})
			switch element["type"].(string) {
			case "input":
				field = f.AddText(element["key"].(string), "text")
			case "textarea":
				field = f.AddTextArea(element["key"].(string))
			case "radio":
				values := strings.Split(element["options"].(string), "\n")
				opt := make([]map[string]string, 0)
				for _, value := range values {
					namval := strings.Split(value, "|")
					if len(namval) < 2 {
						errmessage := fmt.Sprintf("Attempt to set incorrect form option %s by %s", value, w.Request.Host)
						w.SiteConfig.Logger.Error(errmessage)
						services.AddMessage("Your options must be of the format Name|Value", "Error", w)
						w.Serve()
						return
					}
					newval := map[string]string{
						"name":  namval[0],
						"value": namval[1],
					}
					opt = append(opt, newval)
				}
				field = f.AddRadio(element["key"].(string), opt)
			case "checkbox":
				field = f.AddCheckBox(element["key"].(string))
			default:
				//TODO messaging and logging
				return
			}
			if _, ok := element["label"]; ok {
				if element["label"].(string) != "" {
					field.AddLabel(element["label"].(string))
				}
			}
			if _, ok := element["placeholder"]; ok {
				if element["placeholder"].(string) != "" {
					field.AddPlaceHolder(element["placeholder"].(string))
				}
			}
			if _, ok := element["rows"]; ok {
				if _, ok := element["cols"]; ok {
					if element["rows"].(float64) != 0 && element["cols"] != 0 {
						field.AddRowsCols(int(element["rows"].(float64)), int(element["cols"].(float64)))
					}
				}
			}
		}

		var id bson.ObjectId
		if w.Post["mongolarid"].(string) == "new" {
			id = bson.NewObjectId()
		} else {
			id = bson.ObjectIdHex(w.Post["mongolarid"].(string))
		}
		ct := ContentType{
			Form:    f.Fields,
			Type:    w.Post["content_type"].(string),
			MongoId: id,
		}
		s := bson.M{"_id": id}
		c := w.DbSession.DB("").C("content_types")
		_, err = c.Upsert(s, ct)
		if err != nil {
			errmessage := fmt.Sprintf("Cannnot save content type %s : %s", w.Post["mongolarid"].(string), err.Error())
			w.SiteConfig.Logger.Error(errmessage)
			services.AddMessage("Unable to save content type.", "Error", w)
			w.Serve()
			return
		}
		services.AddMessage("Content type saved.", "Success", w)
		dynamic := services.Dynamic{
			Target:     "contenttypelist",
			Controller: "admin/all_content_types",
			Template:   "admin/content_type_list.html",
		}
		services.SetDynamic(dynamic, w)
		w.Serve()
		return
	}
}

func GetAllContentTypes(w *wrapper.Wrapper) {
	c := w.DbSession.DB("").C("content_types")
	var cts []ContentType
	err := c.Find(nil).Limit(50).Iter().All(&cts)
	if err != nil {
		errmessage := fmt.Sprintf("Unable to retrieve a list of content types.")
		w.SiteConfig.Logger.Error(errmessage)
		services.AddMessage("Unable to retrieve a list of elements.", "Error", w)
		w.Serve()
		return
	}
	w.SetPayload("content_types", cts)
	w.Serve()
	return
}

func FieldFormGroup() []*form.Field {
	ft := []map[string]string{
		map[string]string{"name": "Text Field", "value": "input"},
		map[string]string{"name": "TextArea Field", "value": "textarea"},
		map[string]string{"name": "Radio Buttons", "value": "radio"},
		map[string]string{"name": "Checkbox", "value": "checkbox"},
	}
	f := form.NewForm()
	f.AddRadio("type", ft).AddLabel("Field Type").Required()
	f.AddText("key", "text").AddLabel("Key").Required()
	f.AddText("label", "text").AddLabel("Label")
	f.AddText("placeholder", "text").AddLabel("Placeholder")
	f.AddTextArea("options").AddLabel("Options")
	f.AddText("cols", "text").AddLabel("Columns")
	f.AddText("rows", "text").AddLabel("Rows")
	return f.Fields
}
