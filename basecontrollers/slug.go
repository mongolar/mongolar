package basecontrollers

import (
	"fmt"
	"github.com/mongolar/mongolar/models/elements"
	"github.com/mongolar/mongolar/services"
	"github.com/mongolar/mongolar/wrapper"
	"net/http"
)

// The controller function for elements that are context specific
func SlugValues(w *wrapper.Wrapper) {
	var slugid string
	if len(w.APIParams) > 0 {
		slugid = w.APIParams[0]
	} else {
		http.Error(w.Writer, "Forbidden", 403)
		w.Serve()
		return

	}
	es := elements.NewSlugElement()
	err := elements.GetValidElement(slugid, "slug", &es, w)
	if err != nil {
		errmessage := fmt.Sprintf("Content not found %s : %s", slugid, err.Error())
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
	return
}
