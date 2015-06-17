package admin

import (
	"fmt"
	"github.com/mongolar/mongolar/services"
	"github.com/mongolar/mongolar/wrapper"
	"gopkg.in/mgo.v2/bson"
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

func DeletePath(w *wrapper.Wrapper){
	id := w.APIParams[1]
	if !bson.IsObjectIdHex(id){
		errmessage := fmt.Sprintf("Attempt to delete invalid hex %s", id)
		w.SiteConfig.Logger.Error(errmessage)
		services.AddMessage("Invalid Path id.", "Error", w)
		w.Serve()
		return
	}
	c := w.DbSession.DB("").C('paths')
	i := bson.M{"_id": bson.ObjectIdHex(id)}
	err := c.Remove(i)
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

func DeleteElement(w *wrapper.Wrapper){
	id := w.APIParams[1]
	if !bson.IsObjectIdHex(id){
		errmessage := fmt.Sprintf("Attempt to delete invalid hex %s", id)
		w.SiteConfig.Logger.Error(errmessage)
		services.AddMessage("Invalid Element id.", "Error", w)
		w.Serve()
		return
	}
	c := w.DbSession.DB("").C('elements')
	i := bson.M{"_id": bson.ObjectIdHex(id)}
	err := c.Remove(i)
	if err != nil {
		errmessage := fmt.Sprintf("Unable to delete %s : %s", id, err.Error())
		w.SiteConfig.Logger.Error(errmessage)
		services.AddMessage("Unable to delete element.", "Error", w)
		w.Serve()
		return
	}
	s := bson.M{"controller_values.elements": id}
	d := bson.M{"$pull": bson.M{"controller_values.elements": id}}
	_, err := c.UpdateAll(s, d)
	if err != nil {
		errmessage := fmt.Sprintf("Unable to delete reference to element %s : %s", id, err.Error())
		w.SiteConfig.Logger.Error(errmessage)
		services.AddMessage("Unable to delete all references to your element.", "Error", w)
		w.Serve()
		return
	}
	s = bson.M{"elements": id}
	d = bson.M{"$pull": bson.M{"elements": id}}
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
	services.AddMessage("Successfully deleted element", "Success", w)
	w.Serve()
	return
}
