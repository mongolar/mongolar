package admin

import (
	"fmt"
	"github.com/mongolar/mongolar/form"
	"github.com/mongolar/mongolar/models/paths"
	"github.com/mongolar/mongolar/services"
	"github.com/mongolar/mongolar/wrapper"
	"gopkg.in/mgo.v2/bson"
)

func AdminPaths(w *wrapper.Wrapper) {
	pl, err := paths.PathList(w)
	if err != nil {
		services.AddMessage("There was an error retrieving your site paths", "Error", w)
		errmessage := fmt.Sprintf("Error getting path list: %s", err.Error())
		w.SiteConfig.Logger.Error(errmessage)
	} else {
		w.SetContent(pl)
	}
	w.Serve()
}

func PathEditor(w *wrapper.Wrapper) {
	if w.Request.Method != "POST" {
		ops := []string{"published", "unpublished"}
		f := form.NewForm()
		f.AddText("title", "text").AddLabel("Title").Required()
		f.AddText("path", "text").AddLabel("Path").Required()
		f.AddText("template", "text").AddLabel("Template").Required()
		f.AddCheckBox("wildcard").AddLabel("Wildcard")
		o := make([]map[string]string, 0)
		for _, op := range ops {
			r := map[string]string{
				"name":  op,
				"value": op,
			}
			o = append(o, r)
		}
		f.AddRadio("status", o).AddLabel("Status").Required()
		f.AddText("path_id", "text").Hidden()
		if w.APIParams[0] != "new" {
			p := paths.NewPath()
			err := p.GetById(w.APIParams[0], w)
			if err != nil {
				errmessage := fmt.Sprintf("Could not retrieve path %s by %s: %s", w.APIParams[0], w.Request.Host, err.Error())
				w.SiteConfig.Logger.Error(errmessage)
				services.AddMessage("Error retrieving path information.", "Error", w)
				w.Serve()
			} else {
				f.FormData = p
			}
		}
		f.Register(w)
		w.SetTemplate("admin/form.html")
		w.SetPayload("form", f)
		w.Serve()
	} else {
		type PathPost struct {
			*paths.Path
			Id string `json:"mongolarid"`
		}
		realpath := paths.NewPath()
		path := PathPost{&realpath, ""}
		err := form.GetValidFormData(w, &path)
		if err != nil {
			return
		} else {
			c := w.DbSession.DB("").C("paths")
			if path.Id == "new" {
				err := c.Insert(realpath)
				if err != nil {
					errmessage := fmt.Sprintf("Unable to save new path by %s: %s", w.Request.Host, err.Error())
					w.SiteConfig.Logger.Error(errmessage)
					services.AddMessage("There was a problem saving your path.", "Error", w)
					w.Serve()
					return
				}
				services.AddMessage("Your path was saved.", "Success", w)
			} else {
				p := bson.M{
					"$set": realpath,
				}
				s := bson.M{"_id": bson.ObjectIdHex(path.Id)}
				err := c.Update(s, p)
				if err != nil {
					errmessage := fmt.Sprintf("Unable to save path %s by %s: %s", path.Id,
						w.Request.Host, err.Error())
					w.SiteConfig.Logger.Error(errmessage)
					services.AddMessage("There was a problem saving your path.", "Error", w)
					w.Serve()
					return
				} else {
					services.AddMessage("Your path was saved.", "Success", w)
					dynamic := services.Dynamic{
						Target:     "pathbar",
						Controller: "admin/paths",
						Template:   "admin/path_list.html",
					}
					services.SetDynamic(dynamic, w)
					w.Serve()
					return
				}
			}
		}

	}
}

func PathElements(w *wrapper.Wrapper) {
	p := paths.NewPath()
	err := p.GetById(w.APIParams[0], w)
	if err != nil {
		errmessage := fmt.Sprintf("Path not found to edit for %s by %s ", w.APIParams[0], w.Request.Host)
		w.SiteConfig.Logger.Error(errmessage)
		services.AddMessage("This path was not found", "Error", w)
		w.Serve()
	} else {
		w.SetPayload("id", w.APIParams[0])
		w.SetPayload("path", p.Path)
		w.SetPayload("title", p.Title)
		w.SetPayload("elements", p.Elements)
		if len(p.Elements) == 0 {
			services.AddMessage("This path has no elements.", "Info", w)
		}
		w.Serve()
	}

}
