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
	var slugid string
	if len(w.APIParams) > 0 {
		slugid = w.APIParams[0]
	} else {
		http.Error(w.Writer, "Forbidden", 403)
		w.Serve()
		return
	}

	if w.Request.Method != "POST" {
		e := elements.NewSlugElement()
		err := elements.GetById(slugid, &e, w)
		if err != nil {
			errmessage := fmt.Sprintf("Element not found to edit for %s by %s", slugid, w.Request.Host)
			w.SiteConfig.Logger.Error(errmessage)
			services.AddMessage("This element was not found", "Error", w)
			w.Serve()
			return
		}
		f := form.NewForm()
		data := make(map[string]string)
		for slug, id := range e.Slugs {
			f.AddText(id, "text")
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
	} else {
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
		cv := bson.M{"controller_values": vals}
		s := bson.M{"_id": bson.ObjectIdHex(post["mongolarid"])}
		c := w.DbSession.DB("").C("elements")
		err = c.Update(s, cv)
		if err != nil {
			errmessage := fmt.Sprintf("Element not saved %s by %s", w.APIParams[0], w.Request.Host)
			w.SiteConfig.Logger.Error(errmessage)
			services.AddMessage("Unable to save element.", "Error", w)
			w.Serve()
			return
		}
		services.AddMessage("Element content type saved.", "Success", w)
		w.Serve()
		return
	}
}
