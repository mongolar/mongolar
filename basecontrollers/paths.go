package basecontrollers

import (
	"fmt"
	"github.com/mongolar/mongolar/models/elements"
	"github.com/mongolar/mongolar/models/paths"
	"github.com/mongolar/mongolar/services"
	"github.com/mongolar/mongolar/wrapper"
)

// The controller function to retrieve elements ids from the path
func PathValues(w *wrapper.Wrapper) {
	// Request is never url based so we can Never cache this request.
	// TODO: move no caching into a reusable function.
	w.Writer.Header().Add("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Writer.Header().Add("Pragma", "no-cache")
	w.Writer.Header().Add("Expires", "0")
	p := paths.NewPath()
	u := w.Request.Header.Get("CurrentPath")
	qp, err := p.PathMatch(u, "published", w)
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
	var v []elements.Element
	for _, eid := range p.Elements {
		e := elements.NewElement()
		err = elements.GetById(eid, &e, w)
		if err != nil {
			errmessage := fmt.Sprintf("Content not found %s : %s", eid, err.Error())
			w.SiteConfig.Logger.Error(errmessage)
		} else {
			v = append(v, e)
		}
	}
	w.SetPayload("mongolar_slug", qp)
	w.SetContent(v)
	w.SetTemplate(p.Template)
	w.Serve()
	return
}
