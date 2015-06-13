package admin

import (
	"fmt"
	"github.com/mongolar/mongolar/services"
	"github.com/mongolar/mongolar/wrapper"
	"gopkg.in/mgo.v2/bson"
)

func Delete(w *wrapper.Wrapper) {
	c := w.DbSession.DB("").C(w.APIParams[0])
	i := bson.M{"_id": bson.ObjectIdHex(w.APIParams[1])}
	err := c.Remove(i)
	if err != nil {
		errmessage := fmt.Sprintf("Unable to delete %s %s : %s", w.APIParams[0], w.APIParams[1], err.Error())
		w.SiteConfig.Logger.Error(errmessage)
		services.AddMessage("Unable to delete.", "Error", w)
		w.Serve()
		return
	}
	if w.APIParams[0] == "elements" {
		s := bson.M{"controller_values.elements": w.APIParams[1]}
		d := bson.M{"$pull": bson.M{"controller_values.elements": w.APIParams[1]}}
		_, err := c.UpdateAll(s, d)
		if err != nil {
			errmessage := fmt.Sprintf("Unable to delete reference to %s %s : %s", w.APIParams[0], w.APIParams[1], err.Error())
			w.SiteConfig.Logger.Error(errmessage)
			services.AddMessage("Unable to delete all references to your element.", "Error", w)
			w.Serve()
			return
		}
		s = bson.M{"elements": w.APIParams[1]}
		d = bson.M{"$pull": bson.M{"elements": w.APIParams[1]}}
		c := w.DbSession.DB("").C("paths")
		_, err = c.UpdateAll(s, d)
		if err != nil {
			errmessage := fmt.Sprintf("Unable to delete reference to %s %s : %s", w.APIParams[0], w.APIParams[1], err.Error())
			w.SiteConfig.Logger.Error(errmessage)
			services.AddMessage("Unable to delete all references to your element.", "Error", w)
			w.Serve()
			return
		}
		dynamic := services.Dynamic{
			Target:   w.APIParams[1],
			Template: "default.html",
		}
		services.SetDynamic(dynamic, w)
	}
	if w.APIParams[0] == "paths" {
		dynamic := services.Dynamic{
			Target:     "pathbar",
			Controller: "admin/paths",
			Template:   "admin/path_list.html",
		}
		services.SetDynamic(dynamic, w)
	}
	services.AddMessage("Successfully deleted "+w.APIParams[0], "Success", w)
	w.Serve()
	return
}
