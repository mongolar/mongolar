package basecontrollers

import (
	"fmt"
	"github.com/mongolar/mongolar/models/elements"
	"github.com/mongolar/mongolar/services"
	"github.com/mongolar/mongolar/wrapper"
	"net/http"
)

// The controller function for Values found directly in the controller values of the element
func WrapperValues(w *wrapper.Wrapper) {
	var wrapid string
	if len(w.APIParams) > 0 {
		wrapid = w.APIParams[0]
	} else {
		http.Error(w.Writer, "Forbidden", 403)
		w.Serve()
		return
	}
	e := elements.NewWrapperElement()
	err := elements.GetValidElement(wrapid, "wrapper", &e, w)
	if err != nil {
		errmessage := fmt.Sprintf("Content not found %s : %s", wrapid, err.Error())
		w.SiteConfig.Logger.Error(errmessage)
		services.AddMessage("There was a problem loading some content on your page.", "Error", w)
		w.Serve()
		return
	}
	var v []elements.Element
	for _, id := range e.Elements {
		e := elements.NewElement()
		err = elements.GetById(id, &e, w)
		if err != nil {
			errmessage := fmt.Sprintf("Content not found %s : %s", id, err.Error())
			w.SiteConfig.Logger.Error(errmessage)
		} else {
			v = append(v, e)
		}
	}
	w.SetDynamicId(e.DynamicId)
	w.SetContent(v)
	w.Serve()
}
