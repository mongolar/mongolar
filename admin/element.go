package admin

import (
	"fmt"
	"github.com/mongolar/mongolar/form"
	"github.com/mongolar/mongolar/models/elements"
	"github.com/mongolar/mongolar/services"
	"github.com/mongolar/mongolar/wrapper"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"unicode"
)

func Element(w *wrapper.Wrapper) {
	e := elements.NewElement()
	if len(w.APIParams) == 0 {
		http.Error(w.Writer, "Forbidden", 403)
		return
	}
	err := elements.GetById(w.APIParams[0], &e, w)
	if err != nil {
		errmessage := fmt.Sprintf("Element not found to edit for %s by %s.", w.APIParams[0], w.Request.Host)
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

func ElementEditor(w *wrapper.Wrapper) {
	if w.Request.Method != "POST" {
		f := form.NewForm()
		f.AddText("title", "text").AddLabel("Title")
		op := make([]map[string]string, 0)
		for _, ec := range w.SiteConfig.ElementControllers {
			uc := []rune(ec)
			uc[0] = unicode.ToUpper(uc[0])
			name := string(uc)
			op = append(op, map[string]string{"name": name, "value": ec})
		}
		f.AddSelect("controller", op).AddLabel("Controller")
		f.AddText("template", "text").AddLabel("Template")
		f.AddText("dynamic_id", "text").AddLabel("Dynamic Id")
		f.AddText("classes", "text").AddLabel("Classes")
		f.AddText("element_id", "text").Hidden()
		if w.APIParams[0] != "new" {
			e := elements.NewElement()
			err := elements.GetById(w.APIParams[0], &e, w)
			if err != nil {
				errmessage := fmt.Sprintf("Element not found to edit for %s by %s", w.APIParams[0], w.Request.Host)
				w.SiteConfig.Logger.Error(errmessage)
				services.AddMessage("This element was not found", "Error", w)
				w.Serve()
				return
			}
			f.FormData = e
		}
		f.Register(w)
		w.SetTemplate("admin/form.html")
		w.SetPayload("form", f)
	} else {
		post := make(map[string]string)
		err := form.GetValidFormData(w, &post)
		if err != nil {
			return
		} else {
			c := w.DbSession.DB("").C("elements")
			if post["mongolarid"] == "new" {
				p := elements.Element{
					Controller: post["controller"],
					DynamicId:  post["dynamic_id"],
					Template:   post["template"],
					Title:      post["title"],
					Classes:    post["classes"],
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
						"template":   post["template"],
						"title":      post["title"],
						"dynamic_id": post["dynamic_id"],
						"controller": post["controller"],
						"classes":    post["classes"],
					},
				}
				s := bson.M{"_id": bson.ObjectIdHex(post["mongolarid"])}
				err := c.Update(s, p)
				if err != nil {
					errmessage := fmt.Sprintf("Unable to save element %s by %s : %s",
						post["mongolarid"], w.Request.Host, err.Error())
					w.SiteConfig.Logger.Error(errmessage)
					services.AddMessage("There was a problem saving your element.", "Error", w)
				} else {
					services.AddMessage("Your element was saved.", "Success", w)
					dynamic := services.Dynamic{
						Target:     post["mongolarid"],
						Id:         post["mongolarid"],
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

func AllElements(w *wrapper.Wrapper) {
	c := w.DbSession.DB("").C("elements")
	var es []elements.Element
	err := c.Find(nil).Limit(50).Iter().All(&es)
	if err != nil {
		errmessage := fmt.Sprintf("Unable to retrieve a list of all elements: %s", err.Error())
		w.SiteConfig.Logger.Error(errmessage)
		services.AddMessage("There was a problem retrieving the element list.", "Error", w)
		w.Serve()
		return
	}
	w.SetPayload("elements", es)
	w.Serve()
	return
}
