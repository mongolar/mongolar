package admin

import (
	"fmt"
	"github.com/mongolar/mongolar/form"
	"github.com/mongolar/mongolar/models/elements"
	"github.com/mongolar/mongolar/services"
	"github.com/mongolar/mongolar/wrapper"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

func ContentTypeEditor(w *wrapper.Wrapper) {
	var elementid string
	if len(w.APIParams) == 0 {
		http.Error(w.Writer, "Forbidden", 403)
		w.Serve()
		return
	}
	if w.Request.Method != "POST" {
		ContentTypeEditorForm(w)
		return
	}
	ContentTypeEditorSubmit(w)
	return
}

func ContentTypeEditorForm(w *wrapper.Wrapper) {
	elementid = w.APIParams[0]
	e, err := elements.LoadContentElement(elementid, w)
	if err != nil {
		errmessage := fmt.Sprintf("Element not found to edit for %s by %s", elementid, w.Request.Host)
		w.SiteConfig.Logger.Error(errmessage)
		services.AddMessage("This element was not found", "Error", w)
		w.Serve()
		return
	}

}
func ContentTypeEditorSubmit(w *wrapper.Wrapper) {
	elementid = w.APIParams[0]
	e, err := elements.LoadContentElement(elementid, w)
	if err != nil {
		errmessage := fmt.Sprintf("Element not found to edit for %s by %s", elementid, w.Request.Host)
		w.SiteConfig.Logger.Error(errmessage)
		services.AddMessage("This element was not found", "Error", w)
		w.Serve()
		return
	}

}
func ContentTypeEditor(w *wrapper.Wrapper) {
	var elementid string
	if len(w.APIParams) > 0 {
		elementid = w.APIParams[0]
	} else {
		http.Error(w.Writer, "Forbidden", 403)
		w.Serve()
		return
	}
	if w.Request.Method != "POST" {
		e := elements.NewContentElement()
		err := elements.GetById(elementid, &e, w)
		if err != nil {
			errmessage := fmt.Sprintf("Element not found to edit for %s by %s", elementid, w.Request.Host)
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
		data := map[string]string{"type": e.ContentValues.Type}
		f.FormData = data
		f.Register(w)
		w.SetTemplate("admin/form.html")
		w.SetPayload("form", f)
	} else {
		post := make(map[string]string)
		err := form.GetValidFormData(w, &post)
		if err != nil {
			return
		}
		e := bson.M{
			"$set": bson.M{
				"controller_values.type": post["type"],
			},
		}
		s := bson.M{"_id": bson.ObjectIdHex(post["mongolarid"])}
		c := w.DbSession.DB("").C("elements")
		err = c.Update(s, e)
		if err != nil {
			errmessage := fmt.Sprintf("Element not saved %s by %s", w.APIParams[0], w.Request.Host)
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
	var elementid string
	if len(w.APIParams) > 0 {
		elementid = w.APIParams[0]
	} else {
		http.Error(w.Writer, "Forbidden", 403)
		w.Serve()
		return
	}
	e := elements.NewContentElement()
	err := elements.GetById(elementid, &e, w)
	if err != nil {
		errmessage := fmt.Sprintf("Element not found to edit for %s by %s", elementid, w.Request.Host)
		w.SiteConfig.Logger.Error(errmessage)
		services.AddMessage("This element was not found", "Error", w)
		w.Serve()
		return
	}
	if e.ContentValues.Type == "" {
		errmessage := fmt.Sprintf("No content type set for %s", elementid, w.Request.Host)
		w.SiteConfig.Logger.Error(errmessage)
		services.AddMessage("This element doesn't have a content type set.  Set a content type to edit values.", "Error", w)
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
	if w.Request.Method != "POST" {
		if e.Controller != "content" {
			http.Error(w.Writer, "Forbidden", 403)
			return
		}
		f := form.NewForm()
		f.Fields = ct.Form
		f.FormData = e.ContentValues.Content
		f.Register(w)
		w.SetTemplate("admin/form.html")
		w.SetPayload("form", f)
	} else {
		post := make(map[string]interface{})
		err := form.GetValidFormData(w, &post)
		if err != nil {
			return
		}

		content_values := make(map[string]string)
		for _, field := range ct.Form {
			content_values[field.Key] = post[field.Key].(string)
		}
		e := bson.M{
			"$set": bson.M{
				"controller_values.content": content_values,
			},
		}
		s := bson.M{"_id": bson.ObjectIdHex(post["mongolarid"].(string))}
		c := w.DbSession.DB("").C("elements")
		err = c.Update(s, e)
		if err != nil {
			errmessage := fmt.Sprintf("Element not saved %s by %s", w.APIParams[0], w.Request.Host)
			w.SiteConfig.Logger.Error(errmessage)
			services.AddMessage("Unable to save element.", "Error", w)
		} else {
			services.AddMessage("Element content saved.", "Success", w)
			dynamic := services.Dynamic{
				Target:     post["mongolarid"].(string),
				Id:         post["mongolarid"].(string),
				Controller: "admin/element",
				Template:   "admin/element.html",
			}
			services.SetDynamic(dynamic, w)
		}
	}
	w.Serve()
	return
}
