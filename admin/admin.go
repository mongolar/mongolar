package admin

import {
	"github.com/jasonrichardsmith/mongolar/controller"
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
		w.SiteConfig.Logger.Error("Error getting path list")
	}
	else {
		w.SetContent(e.ControllerValues)
	}
	w.Serve()
}


func EditPath(w *wrapper.Wrapper) {

}

func ElementEditor(w *wrapper.Wrapper) {
	
}

func WrapperEditor(w *wrapper.Wrapper) {

}


