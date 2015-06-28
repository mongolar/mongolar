package admin

import (
	"encoding/json"
	"fmt"
	"github.com/mongolar/mongolar/models/elements"
	"github.com/mongolar/mongolar/services"
	"github.com/mongolar/mongolar/wrapper"
	"net/http"
)

func MenuEditor(w *wrapper.Wrapper) {
	if len(w.APIParams) == 0 {
		http.Error(w.Writer, "Forbidden", 403)
		return
	}
	menuid := w.APIParams[0]
	if w.Request.Method != "POST" {
		e, err := elements.LoadMenuElement(menuid, w)
		if err != nil {
			errmessage := fmt.Sprintf("Element not found to edit for %s by %s.", menuid, w.Request.Host)
			w.SiteConfig.Logger.Error(errmessage)
			services.AddMessage("This element was not found", "Error", w)
			w.Serve()
			return
		}
		if e.MenuItems == nil {
			items := make(map[string][]map[string]string)
			items["menu_items"] = make([]map[string]string, 0)
			w.SetPayload("menu", items)
		} else {
			w.SetPayload("menu", e)
		}
		w.SetPayload("title", e.Title)
		w.SetTemplate("admin/menu_editor.html")
		w.Serve()
		return
	} else {
		e, err := elements.LoadMenuElement(menuid, w)
		if err != nil {
			errmessage := fmt.Sprintf("Element not found to edit for %s by %s.", menuid, w.Request.Host)
			w.SiteConfig.Logger.Error(errmessage)
			services.AddMessage("This element was not found", "Error", w)
			w.Serve()
			return
		}
		err = json.NewDecoder(w.Request.Body).Decode(&e)
		if err != nil {
			errmessage := fmt.Sprintf("Unable to update marshall menu elements by %s: %s", w.Request.Host, err.Error())
			w.SiteConfig.Logger.Error(errmessage)
			services.AddMessage("Unable to save menu element.", "Error", w)
			w.Serve()
			return
		}
		err = e.Save(w)
		if err != nil {
			errmessage := fmt.Sprintf("Unable to update menu element %s by %s: %s", menuid, w.Request.Host, err.Error())
			w.SiteConfig.Logger.Error(errmessage)
			services.AddMessage("Unable to save menu element.", "Error", w)
			w.Serve()
			return
		}
		dynamic := services.Dynamic{
			Target:     "righteditor",
			Controller: "",
			Template:   "",
			Id:         "",
		}
		services.SetDynamic(dynamic, w)
		services.AddMessage("You menu element have been updated.", "Success", w)
		w.Serve()
	}

}
