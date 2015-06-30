package admin

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/mongolar/mongolar/models/elements"
	"github.com/mongolar/mongolar/models/paths"
	"github.com/mongolar/mongolar/services"
	"github.com/mongolar/mongolar/wrapper"
	"net/http"
)

func Delete(w *wrapper.Wrapper) {
	var parenttype string
	if len(w.APIParams) > 1 {
		parenttype = w.APIParams[0]
	} else {
		http.Error(w.Writer, "Forbidden", 403)
		w.Serve()
		return
	}
	w.Shift()
	switch parenttype {
	case "elements":
		DeleteElement(w)
		return
	case "paths":
		DeletePath(w)
		return
	default:
		http.Error(w.Writer, "Forbidden", 403)
		w.Serve()
	}
	return
}

func DeletePath(w *wrapper.Wrapper) {
	id := w.APIParams[0]
	err := paths.Delete(id, w)
	if err != nil {
		errmessage := fmt.Sprintf("Unable to delete path %s : %s", id, err.Error())
		w.SiteConfig.Logger.Error(errmessage)
		services.AddMessage("Unable to delete.", "Error", w)
		w.Serve()
		return
	}
	dynamic := services.Dynamic{
		Target:     "pathbar",
		Controller: "admin/paths",
		Template:   "admin/path_list.html",
	}
	services.SetDynamic(dynamic, w)
	services.AddMessage("Successfully deleted path", "Success", w)
	w.Serve()
	return
}

func DeleteElement(w *wrapper.Wrapper) {
	id := w.APIParams[0]
	err := elements.Delete(id, w)
	if err != nil {
		errmessage := fmt.Sprintf("Unable to delete %s : %s", id, err.Error())
		w.SiteConfig.Logger.Error(errmessage)
		services.AddMessage("Unable to delete element.", "Error", w)
	}
	err = elements.WrapperDeleteAllChild(id, w)
	if err != nil {
		errmessage := fmt.Sprintf("Unable to delete reference to element %s : %s", id, err.Error())
		w.SiteConfig.Logger.Error(errmessage)
		services.AddMessage("Unable to delete all references to your element.", "Error", w)
	}
	err = elements.SlugDeleteAllChild(id, w)
	spew.Dump(id)
	if err != nil {
		errmessage := fmt.Sprintf("Unable to delete reference to element %s : %s", id, err.Error())
		w.SiteConfig.Logger.Error(errmessage)
		services.AddMessage("Unable to delete all references to your element.", "Error", w)
	}
	err = paths.DeleteAllChild(id, w)
	if err != nil {
		errmessage := fmt.Sprintf("Unable to delete reference to %s : %s", id, err.Error())
		w.SiteConfig.Logger.Error(errmessage)
		services.AddMessage("Unable to delete all references to your element.", "Error", w)
	}
	dynamic := services.Dynamic{
		Target:   id,
		Template: "default.html",
	}
	services.SetDynamic(dynamic, w)
	services.AddMessage("Successfully deleted element", "Success", w)
	w.Serve()
	return
}
