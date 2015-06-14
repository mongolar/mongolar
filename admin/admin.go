// Admin is a series of controllers to manage a Mongolar site

package admin

import (
	"github.com/mongolar/mongolar/controller"
	"github.com/mongolar/mongolar/services"
	"github.com/mongolar/mongolar/user"
	"github.com/mongolar/mongolar/wrapper"
	"net/http"
)

type AdminMap controller.ControllerMap

func GetControllerMap(cm controller.ControllerMap) {
	amap, _ := NewAdmin()
	cm["admin"] = amap.Admin
}

// A series of menu items to render on the admin page
type AdminMenu struct {
	MenuItems []map[string]string `json:"menu_items"`
}

// Build a new admin instance and return AdminMap and menu to be altered
func NewAdmin() (*AdminMap, *AdminMenu) {
	amenu := AdminMenu{
		MenuItems: []map[string]string{
			map[string]string{"title": "Home", "template": "admin/main_content_default.html"},
			map[string]string{"title": "Content", "template": "admin/content_editor.html"},
			map[string]string{"title": "Content Types", "template": "admin/content_types_editor.html"},
		},
	}
	amap := &AdminMap{
		"admin_menu":         amenu.AdminMenu,
		"paths":              AdminPaths,
		"path_elements":      PathElements,
		"path_editor":        PathEditor,
		"element":            Element,
		"element_editor":     ElementEditor,
		"add_child":          AddChild,
		"add_existing_child": AddExistingChild,
		"all_content_types":  GetAllContentTypes,
		"edit_content_type":  EditContentType,
		"delete":             Delete,
		"sort_children":      Sort,
		"content":            ContentEditor,
		"menu":               MenuEditor,
		"content_type":       ContentTypeEditor,
		"orphans":            OrphanElements,
	}
	return amap, &amenu
}

//Main controller for all admin functions
func (a AdminMap) Admin(w *wrapper.Wrapper) {
	if c, ok := a[w.APIParams[0]]; ok {
		if validateAdmin(w) {
			w.Shift()
			c(w)
		}
		return
	} else {
		http.Error(w.Writer, "Forbidden", 403)
		return
	}
}

func validateAdmin(w *wrapper.Wrapper) bool {
	user := new(user.User)
	err := user.Get(w)
	loginurls := make(map[string]string)
	w.SiteConfig.RawConfig.MarshalKey("LoginURLs", &loginurls)
	if err != nil {
		services.Redirect(loginurls["login"], w)
		w.Serve()
		return false
	}
	if user.Roles != nil {
		for _, r := range user.Roles {
			if r == "admin" {
				return true
			}
		}
	}
	services.Redirect(loginurls["access_denied"], w)
	w.Serve()
	return false
}

func (a *AdminMenu) AdminMenu(w *wrapper.Wrapper) {
	w.SetContent(a)
	w.Serve()
	return
}
