package basecontrollers

import (
	"github.com/mongolar/mongolar/wrapper"
	"net/http"
)

// The controller function for Values found in the Site Configuration
func DomainPublicValue(w *wrapper.Wrapper) {
	var valuekey string
	if len(w.APIParams) > 0 {
		valuekey = w.APIParams[0]
	} else {
		http.Error(w.Writer, "Forbidden", 403)
		w.Serve()
		return

	}
	if value, ok := w.SiteConfig.PublicValues[valuekey]; ok {
		w.SetPayload("domain_value", value)
	} else {
		http.Error(w.Writer, "Forbidden", 403)
	}
	w.Serve()
	return
}
