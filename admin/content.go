package admin

import (
	"fmt"
	"github.com/mongolar/mongolar/form"
	"github.com/mongolar/mongolar/models/contenttypes"
	"github.com/mongolar/mongolar/models/elements"
	"github.com/mongolar/mongolar/services"
	"github.com/mongolar/mongolar/wrapper"
	"net/http"
)

// Controller to change content type for content type element.
func ContentTypeEditor(w *wrapper.Wrapper) {
	if len(w.APIParams) == 0 {
		http.Error(w.Writer, "Forbidden", 403)
		w.Serve()
		return
	}
	if w.Request.Method != "POST" {
		ContentTypeEditorForm(w)
		return
	}
	ContentTypeEditorSubmit(w)
	return
}

// Controller to ddisplay form to change content type for conteent element
func ContentTypeEditorForm(w *wrapper.Wrapper) {
	elementid := w.APIParams[0]
	e, err := elements.LoadContentElement(elementid, w)
	if err != nil {
		errmessage := fmt.Sprintf("Element not found to edit for %s by %s", elementid, w.Request.Host)
		w.SiteConfig.Logger.Error(errmessage)
		services.AddMessage("This element was not found", "Error", w)
		w.Serve()
		return
	}
	cts := make([]contenttypes.ContentType, 0)
	cts, err = contenttypes.AllContentTypes(w)
	if err != nil {
		errmessage := fmt.Sprintf("Unable to query all Content Types: %s", err.Error())
		w.SiteConfig.Logger.Error(errmessage)
		services.AddMessage("Unable to retrieve content types.", "Error", w)
		w.Serve()
		return
	}
	opts := make([]map[string]string, 0)
	for _, ct := range cts {
		opt := map[string]string{
			"name":  ct.Type,
			"value": ct.Type,
		}
		opts = append(opts, opt)
	}
	f := form.NewForm()
	f.AddSelect("type", opts)
	data := map[string]string{"type": e.ContentValues.Type}
	f.FormData = data
	f.Register(w)
	w.SetTemplate("admin/form.html")
	w.SetPayload("form", f)
	w.Serve()
	return

}

// Controller to handle submission for content type change form.
func ContentTypeEditorSubmit(w *wrapper.Wrapper) {
	elementid := w.APIParams[0]
	e, err := elements.LoadContentElement(elementid, w)
	if err != nil {
		errmessage := fmt.Sprintf("Element not found to edit for %s by %s", elementid, w.Request.Host)
		w.SiteConfig.Logger.Error(errmessage)
		services.AddMessage("This element was not found", "Error", w)
		w.Serve()
		return
	}
	type Post struct {
		Type string `json:"type"`
	}
	var post Post
	err = form.GetValidFormData(w, &post)
	if err != nil {
		return
	}
	e.ContentValues.Type = post.Type
	err = e.Save(w)
	if err != nil {
		errmessage := fmt.Sprintf("Element not saved %s by %s", w.APIParams[0], w.Request.Host)
		w.SiteConfig.Logger.Error(errmessage)
		services.AddMessage("Unable to save element.", "Error", w)
	} else {
		services.AddMessage("Element content type saved.", "Success", w)
	}
	w.Serve()
	return
}

// Controller to edit content in element.
func ContentEditor(w *wrapper.Wrapper) {
	if len(w.APIParams) == 0 {
		http.Error(w.Writer, "Forbidden", 403)
		w.Serve()
		return
	}
	if w.Request.Method != "POST" {
		ContentEditorForm(w)
		return
	}
	ContentEditorSubmit(w)
	return
}

// Controller for content editing form.
func ContentEditorForm(w *wrapper.Wrapper) {
	elementid := w.APIParams[0]
	e, err := elements.LoadContentElement(elementid, w)
	if err != nil {
		errmessage := fmt.Sprintf("Element not found to edit for %s by %s", elementid, w.Request.Host)
		w.SiteConfig.Logger.Error(errmessage)
		services.AddMessage("This element was not found", "Error", w)
		w.Serve()
		return
	}
	if e.ContentValues.Type == "" {
		errmessage := fmt.Sprintf("No content type set for %s", elementid, w.Request.Host)
		w.SiteConfig.Logger.Error(errmessage)
		services.AddMessage("This element doesn't have a content type set.  Set a content type to edit values.", "Error", w)
		w.Serve()
		return
	}
	var ct contenttypes.ContentType
	ct, err = contenttypes.LoadContentTypeT(e.ContentValues.Type, w)
	if err != nil {
		errmessage := fmt.Sprintf("Unable to find content type %s : %s", e.ContentValues.Type, err.Error())
		w.SiteConfig.Logger.Error(errmessage)
		services.AddMessage("Unable to find content type.", "Error", w)
		w.Serve()
		return
	}
	f := form.NewForm()
	f.Fields = ct.Form
	f.FormData = e.ContentValues.Content
	f.Register(w)
	w.SetTemplate("admin/form.html")
	w.SetPayload("form", f)
	w.Serve()
	return
}

// Controller to handle content editor submission.
func ContentEditorSubmit(w *wrapper.Wrapper) {
	elementid := w.APIParams[0]
	e, err := elements.LoadContentElement(elementid, w)
	post := make(map[string]interface{})
	err = form.GetValidFormData(w, &post)
	if err != nil {
		return
	}
	e.ContentValues.Content = post
	delete(e.ContentValues.Content, "mongolartype")
	delete(e.ContentValues.Content, "mongolarid")
	err = e.Save(w)
	if err != nil {
		errmessage := fmt.Sprintf("Element not saved %s by %s", w.APIParams[0], w.Request.Host)
		w.SiteConfig.Logger.Error(errmessage)
		services.AddMessage("Unable to save element.", "Error", w)
		w.Serve()
		return
	}
	services.AddMessage("Element content saved.", "Success", w)
	dynamic := services.Dynamic{
		Target:     elementid,
		Id:         elementid,
		Controller: "admin/element",
		Template:   "admin/element.html",
	}
	services.SetDynamic(dynamic, w)
	w.Serve()
	return

}
