package admin

import {
	"github.com/jasonrichardsmith/mongolar/controller"
	"github.com/jasonrichardsmith/mongolar/form"
	"github.com/jasonrichardsmith/mongolar/url"
	"github.com/jasonrichardsmith/mongolar/session"
	"github.com/jasonrichardsmith/mongolar/service/messages"
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
	if w.Post == nil {
		ops := {"published", "unpublished"}
		f = form.NewForm()
		f.AddText("path").AddLabel("Path")
		f.AddText("template").AddLabel("Template")
		f.AddCheckBox("wildcard").AddLabel("Wildcard")
		o := make([]map[string]string, 1)
		for _, op := range ops {
			r := map[string]string{
				'name': op
				'value': op
			}
			o = append(o, r)
		}
		f.AddRadio("status", o).AddLabel("Status")
		f.AddText("path_id").Hidden()
		u := url.UrlToMap(w.Request.URL.Path)
		if u[2] != 'new' {
			p := controller.NewPath()
			err := p.GetById(u[2], w.SiteConfig.DbSession)
			if err != nil {
				w.SiteConfig.Logger.Error("Path not found to edit for " +  $u[2] + " by " w.Request.Host)
				m := messages.Message{Text: "This path was not found", Severity: "Error"}
				messages.Set(m, w)
				w.Serve()
			} else {
				f.FormData['wildcard'] = p.Wildcard
				f.FormData['template'] = p.Template
				f.FormData['path'] = p.Path
				f.FormData['status'] = p.Status
			}
		}
		w.SetContent(f)
		w.Serve()
	} else {
		_, err := form.GetValidRegForm(w.Post['FormId'], w.Session, w.SiteConfig.DbSession)
		if  err != nil {
			w.SiteConfig.Logger.Error("Attempt to access invalid form" + w.Post['FormId'] + " by " w.Request.Host)
			m := messages.Message{Text: "Invalid Form"}
			messages.Set(m, w)
			w.Serve()
		} else {
			//update/save path here
		}

	}
}

func ElementEditor(w *wrapper.Wrapper) {
	        //ControllerValues map[string]interface{} `bson:"controller_values,omitempty"`
		//Controller       string                 `bson:"controller"`
		//Template         string                 `bson:"template"`
		//DynamicId        string                 `bson:"dynamic_id,omitempty"`
	if w.Post == nil {

		f = form.NewForm()
		f.AddText("controller").AddLabel("Controller")
		f.AddText("template").AddLabel("Template")
		f.AddCheckBox("dynamic_id").AddLabel("Dynamic Id")
		f.AddText("element_id").Hidden()
		u := url.UrlToMap(w.Request.URL.Path)
		if u[2] != 'new' {
			e := controller.NewElement()
			err := e.GetById(u[2], w.SiteConfig.DbSession)
			if err != nil {
				w.SiteConfig.Logger.Error("Element not found to edit for " +  $u[2] + " by " w.Request.Host)
				m := messages.Message{Text: "This element was not found", Severity: "Error"}
				messages.Set(m, w)
				w.Serve()
			} else {
				f.FormData['controller'] = e.Controller
				f.FormData['template'] = e.Template
				f.FormData['dynamic_id'] = e.DynamicId
			}
		}
		w.SetContent(f)
		w.Serve()
	} else {
		_, err := form.GetValidRegForm(w.Post['FormId'], w.Session, w.SiteConfig.DbSession)
		if  err != nil {
			w.SiteConfig.Logger.Error("Attempt to access invalid form" + w.Post['FormId'] + " by " w.Request.Host)
			m := messages.Message{Text: "Invalid Form"}
			messages.Set(m, w)
			w.Serve()
		} else {
		  //update save element here
		}
	}


}

func WrapperEditor(w *wrapper.Wrapper) {
	
}


