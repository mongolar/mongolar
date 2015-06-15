package basecontrollers

import (
	"fmt"
	"github.com/mongolar/mongolar/models/elements"
	"github.com/mongolar/mongolar/services"
	"github.com/mongolar/mongolar/wrapper"
	"net/http"
)

// The controller function for Values found directly in the controller values of the element
func ContentValues(w *wrapper.Wrapper) {
	var contentid string
	if len(w.APIParams) > 0 {
		contentid = w.APIParams[0]
	} else {
		http.Error(w.Writer, "Forbidden", 403)
		w.Serve()
		return

	}
	e := elements.NewContentElement()
	err := elements.GetValidElement(contentid, "content", &e, w)
	if err != nil {
		errmessage := fmt.Sprintf("Content not found %s : %s", contentid, err.Error())
		w.SiteConfig.Logger.Error(errmessage)
		services.AddMessage("There was a problem loading some content on your page.", "Error", w)
		w.Serve()
		return
	}
	w.SetContent(e.ContentValues.Content)
	w.Serve()
	return
}
