package admin

import (
	"encoding/json"
	"fmt"
	"github.com/mongolar/mongolar/models/elements"
	"github.com/mongolar/mongolar/services"
	"github.com/mongolar/mongolar/wrapper"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

func MenuEditor(w *wrapper.Wrapper) {
	if w.Request.Method != "POST" {
		if len(w.APIParams) == 0 {
			http.Error(w.Writer, "Forbidden", 403)
			return
		}
		e := elements.NewElement()
		err := elements.GetById(w.APIParams[0], &e, w)
		if err != nil {
			errmessage := fmt.Sprintf("Element not found to edit for %s by %s.", w.APIParams[0], w.Request.Host)
			w.SiteConfig.Logger.Error(errmessage)
			services.AddMessage("This element was not found", "Error", w)
		} else {
			w.SetPayload("menu", e.ControllerValues)
			w.SetPayload("title", e.Title)
		}
		w.Serve()
	} else {
		post := make(map[string]interface{})
		err := json.NewDecoder(w.Request.Body).Decode(&post)
		if err != nil {
			errmessage := fmt.Sprintf("Unable to update marshall menu elements by %s: %s", w.Request.Host, err.Error())
			w.SiteConfig.Logger.Error(errmessage)
			services.AddMessage("Unable to save menu element.", "Error", w)
			w.Serve()
			return
		}
		p := bson.M{
			"$set": bson.M{
				"controller_values": post["menu"],
			},
		}
		s := bson.M{"_id": bson.ObjectIdHex(w.APIParams[1])}
		c := w.DbSession.DB("").C("elements")
		err = c.Update(s, p)
		if err != nil {
			errmessage := fmt.Sprintf("Unable to update menu element %s by %s: %s", w.APIParams[0], w.Request.Host, err.Error())
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
