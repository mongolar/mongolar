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

func SlugUrlEditor(w *wrapper.Wrapper) {
	if len(w.APIParams) < 1 {
		http.Error(w.Writer, "Forbidden", 403)
		w.Serve()
		return
	}
	if w.Request.Method != "POST" {
		SlugUrlEditorForm(w)
		return
	}
	SlugUrlEditorSubmit(w)
	return
}

func SlugUrlEditorForm(w *wrapper.Wrapper) {
	slugid := w.APIParams[0]
	e, err := elements.LoadSlugElement(slugid, w)
	if err != nil {
		errmessage := fmt.Sprintf("Element not found to edit for %s by %s", slugid, w.Request.Host)
		w.SiteConfig.Logger.Error(errmessage)
		services.AddMessage("Unable to load slug parent", "Error", w)
		w.Serve()
		return
	}
	f := form.NewForm()
	data := make(map[string]string)
	for slug, id := range e.Slugs {
		data[id] = slug
		e := elements.NewElement()
		err = elements.GetById(id, &e, w)
		if err != nil {
			errmessage := fmt.Sprintf("Content not found %s : %s", id, err.Error())
			w.SiteConfig.Logger.Error(errmessage)
			services.AddMessage("There was a problem loading some slug elements.", "Error", w)
			w.Serve()
			return
		}
		f.AddText(id, "text").AddLabel(e.Title).Required()

	}
	f.FormData = data
	f.Register(w)
	w.SetTemplate("admin/form.html")
	w.SetPayload("form", f)
	w.Serve()
	return
}

func SlugUrlEditorSubmit(w *wrapper.Wrapper) {
	slugid := w.APIParams[0]
	post := make(map[string]string)
	err := form.GetValidFormData(w, &post)
	if err != nil {
		return
	}
	vals := make(map[string]string)
	for id, slug := range post {
		if bson.IsObjectIdHex(id) {
			vals[slug] = id
		}
	}
	var e elements.SlugElement
	e, err = elements.LoadSlugElement(slugid, w)
	if err != nil {
		errmessage := fmt.Sprintf("Element not found to edit for %s by %s", slugid, w.Request.Host)
		w.SiteConfig.Logger.Error(errmessage)
		services.AddMessage("Unable to load parent slug", "Error", w)
		w.Serve()
		return
	}
	e.Slugs = vals
	err = e.Save(w)
	if err != nil {
		errmessage := fmt.Sprintf("Slugs not saved %s by %s", w.APIParams[0], w.Request.Host)
		w.SiteConfig.Logger.Error(errmessage)
		services.AddMessage("Unable to save slug values.", "Error", w)
		w.Serve()
		return
	}
	services.AddMessage("Slug values updated.", "Success", w)
	w.Serve()
	return
}
