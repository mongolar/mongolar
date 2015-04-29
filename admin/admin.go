package admin

import {
	"github.com/jasonrichardsmith/mongolar/controller"
	"github.com/jasonrichardsmith/mongolar/form"
	"github.com/jasonrichardsmith/mongolar/url"
	"github.com/jasonrichardsmith/mongolar/session"
}

type AdminMap struct controller.ControllerMap


func NewAdmin() *AdminMap {
	am := AdminMap{
		"menu": AdminMenu
		"paths" : AdminPaths
	}
	return am
}

func (a *AdminMap) Admin(w *wrapper.Wrapper) {
        u := url.UrlToMap(w.Request.URL.Path)
	if c, ok := a[u[1]]; ok {
		if validateAdmin() {
			c(w)
		} else {
			http.Error(w.Writer, "Forbidden", 403)
		}
	} else {
		http.Error(w.Writer, "Forbidden", 403)
		return
	}
}

func validateAdmin(s session.Session) bool {
	return true
}

func AdminMenu(w *wrapper.Wrapper) {
	w.SetContent(w.SiteConfig.Misc['AdminMenu'])
	w.Serve()
	return
}


func PathList(w *wrapper.Wrapper) {
	//TODO: Log Errors here
	w.SetTemplate("admin/pathlist")
	pl, err := controller.PathList(w.SiteConfig.DbSession)
	if err != nil {
		w.SiteConfig.Logger.Error("Error getting path list: " + err.Error())
	}
	else {
		w.SetContent(e.ControllerValues)
	}
	w.Serve()
}


func EditPath(w *wrapper.Wrapper) {
	op := {"published", "unpublished"}
	f = form.NewForm()
	f.AddText("path").AddLabel("Path")
	f.AddCheckBox("wildcard").AddLabel("Wildcard")
	f.AddText("template").AddLabel("Template")
	o := make([]map[string]string)

	f.AddRadio("wildcard").AddLabel("Status")

	Path     string        `bson:"path"`
	Wildcard bool          `bson:"wildcard"`
	Elements []string      `bson:"elements"`
	Template string        `bson:"template"`
	Status   string        `bson:"status"`
}

func ElementEditor(w *wrapper.Wrapper) {
	
}

func WrapperEditor(w *wrapper.Wrapper) {

}


