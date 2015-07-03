package admin

import (
	"encoding/json"
	"fmt"
	"github.com/mongolar/mongolar/form"
	"github.com/mongolar/mongolar/models/elements"
	"github.com/mongolar/mongolar/models/paths"
	"github.com/mongolar/mongolar/services"
	"github.com/mongolar/mongolar/wrapper"
	"net/http"
)

// Sort controller to sort path and wrapper children
func Sort(w *wrapper.Wrapper) {
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
	case "wrapper":
		if w.Request.Method != "POST" {
			SortWrapperForm(w)
			return
		}
		SortWrapperSubmit(w)
		return
	case "paths":
		if w.Request.Method != "POST" {
			SortPathForm(w)
			return
		}
		SortPathSubmit(w)
		return
	default:
		http.Error(w.Writer, "Forbidden", 403)
		w.Serve()
	}
	return
}

// Controller for wrapper element sort form.
func SortWrapperForm(w *wrapper.Wrapper) {
	var parentid string
	if len(w.APIParams) > 0 {
		parentid = w.APIParams[0]
	} else {
		http.Error(w.Writer, "Forbidden", 403)
		w.Serve()
		return
	}
	e, err := elements.LoadWrapperElement(parentid, w)
	if err != nil {
		errmessage := fmt.Sprintf("Element not found to sort for %s by %s.", w.APIParams[1], w.Request.Host)
		w.SiteConfig.Logger.Error(errmessage)
		services.AddMessage("This element was not found", "Error", w)
		w.Serve()
		return
	}
	if len(e.Elements) > 0 {
		w.SetPayload("elements", e.Elements)
	} else {
		services.AddMessage("This has no elements assigned yet.", "Error", w)
	}
	w.SetTemplate("admin/element_sorter.html")
	w.Serve()
	return

}

// Controller for path sort form.
func SortPathForm(w *wrapper.Wrapper) {
	var parentid string
	if len(w.APIParams) > 0 {
		parentid = w.APIParams[0]
	} else {
		http.Error(w.Writer, "Forbidden", 403)
		w.Serve()
		return
	}
	p, err := paths.LoadPath(parentid, w)
	if err != nil {
		errmessage := fmt.Sprintf("Path not found to sort for %s by %s.", parentid, w.Request.Host)
		w.SiteConfig.Logger.Error(errmessage)
		services.AddMessage("This element was not found", "Error", w)
		w.Serve()
		return
	} else {
		if len(p.Elements) > 0 {
			w.SetPayload("elements", p.Elements)
		} else {
			services.AddMessage("This has no elements assigned yet.", "Error", w)
		}
		w.SetTemplate("admin/element_sorter.html")
		w.Serve()
		return
	}
}

// Controller for wrapper sort submission.
func SortWrapperSubmit(w *wrapper.Wrapper) {
	var parentid string
	if len(w.APIParams) > 0 {
		parentid = w.APIParams[0]
	} else {
		http.Error(w.Writer, "Forbidden", 403)
		w.Serve()
		return
	}
	wes := elements.NewWrapperElements()
	err := json.NewDecoder(w.Request.Body).Decode(&wes)
	if err != nil {
		errmessage := fmt.Sprintf("Unable to marshall elements %s by %s: %s", parentid, w.Request.Host, err.Error())
		w.SiteConfig.Logger.Error(errmessage)
		services.AddMessage("Unable to save elements.", "Error", w)
		w.Serve()
		return
	}
	we, err := elements.LoadWrapperElement(parentid, w)
	if err != nil {
		errmessage := fmt.Sprintf("Element not found to sort for %s by %s.", parentid, w.Request.Host)
		w.SiteConfig.Logger.Error(errmessage)
		services.AddMessage("This element was not found", "Error", w)
		w.Serve()
		return
	}
	we.WrapperElements = wes
	we.Save(w)
	if err != nil {
		errmessage := fmt.Sprintf("Unable to save wrapper element %s by %s : %s", parentid, w.Request.Host, err.Error())
		w.SiteConfig.Logger.Error(errmessage)
		services.AddMessage("Could not save parent element.", "Error", w)
		w.Serve()
		return
	}
	dynamic := services.Dynamic{
		Target:     parentid,
		Controller: "admin/element",
		Template:   "admin/element.html",
		Id:         parentid,
	}
	services.SetDynamic(dynamic, w)
	services.AddMessage("You elements have been updated.", "Success", w)
	w.Serve()
	return
}

// Controller for path sort submission.
func SortPathSubmit(w *wrapper.Wrapper) {
	var parentid string
	if len(w.APIParams) > 0 {
		parentid = w.APIParams[0]
	} else {
		http.Error(w.Writer, "Forbidden", 403)
		w.Serve()
		return
	}
	pes := paths.NewPathElements()
	err := json.NewDecoder(w.Request.Body).Decode(&pes)
	if err != nil {
		errmessage := fmt.Sprintf("Unable to marshall elements %s by %s: %s", parentid, w.Request.Host, err.Error())
		w.SiteConfig.Logger.Error(errmessage)
		services.AddMessage("Unable to save elements.", "Error", w)
		w.Serve()
		return
	}
	pe, err := paths.LoadPath(parentid, w)
	if err != nil {
		errmessage := fmt.Sprintf("Path not found to sort for %s by %s.", parentid, w.Request.Host)
		w.SiteConfig.Logger.Error(errmessage)
		services.AddMessage("This path was not found", "Error", w)
		w.Serve()
		return
	}
	pe.PathElements = pes
	err = pe.Save(w)
	if err != nil {
		errmessage := fmt.Sprintf("Unable to save path %s by %s : %s", parentid, w.Request.Host, err.Error())
		w.SiteConfig.Logger.Error(errmessage)
		services.AddMessage("Could not save path.", "Error", w)
		w.Serve()
		return
	}
	dynamic := services.Dynamic{
		Target:     "centereditor",
		Controller: "admin/path_elements",
		Template:   "admin/path_elements.html",
		Id:         parentid,
	}
	services.SetDynamic(dynamic, w)
	services.AddMessage("You elements have been updated.", "Success", w)
	w.Serve()
	return
}

// Controller for adding children to wrapper element or path.
func AddChild(w *wrapper.Wrapper) {
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
	case "wrapper":
		AddWrapperChild(w)
	case "slug":
		AddSlugChild(w)
	case "paths":
		AddPathChild(w)
	default:
		http.Error(w.Writer, "Forbidden", 403)
		w.Serve()
	}
	return
}

// Controller for adding children to wrapper element.
func AddWrapperChild(w *wrapper.Wrapper) {
	parentid := w.APIParams[0]
	e := elements.NewElement()
	e.Title = "New Element"
	err := e.Save(w)
	if err != nil {
		errmessage := fmt.Sprintf("Unable to create new element  by %s : %s", w.Request.Host, err.Error())
		w.SiteConfig.Logger.Error(errmessage)
		services.AddMessage("Could not create a new element.", "Error", w)
		w.Serve()
		return
	}
	var parent elements.WrapperElement
	parent, err = elements.LoadWrapperElement(parentid, w)
	if err != nil {
		errmessage := fmt.Sprintf("Unable to loap parent element  by %s : %s", w.Request.Host, err.Error())
		w.SiteConfig.Logger.Error(errmessage)
		services.AddMessage("Could not load parent element.", "Error", w)
		w.Serve()
		return
	}
	parent.Elements = append(parent.Elements, e.MongoId.Hex())
	err = parent.Save(w)
	if err != nil {
		errmessage := fmt.Sprintf("Unable to loap parent element  by %s : %s", w.Request.Host, err.Error())
		w.SiteConfig.Logger.Error(errmessage)
		services.AddMessage("Could not load parent element.", "Error", w)
		w.Serve()
		return
	}
	dynamic := services.Dynamic{
		Target:     parentid,
		Controller: "admin/element",
		Template:   "admin/element.html",
		Id:         parentid,
	}
	services.SetDynamic(dynamic, w)
	services.AddMessage("You have added a new element.", "Success", w)
	w.Serve()
	return

}

func AddSlugChild(w *wrapper.Wrapper) {
	parentid := w.APIParams[0]
	e := elements.NewElement()
	e.Title = "New Element"
	e.Controller = "content"
	err := e.Save(w)
	if err != nil {
		errmessage := fmt.Sprintf("Unable to create new element  by %s : %s", w.Request.Host, err.Error())
		w.SiteConfig.Logger.Error(errmessage)
		services.AddMessage("Could not create a new element.", "Error", w)
		w.Serve()
		return
	}
	var parent elements.SlugElement
	parent, err = elements.LoadSlugElement(parentid, w)
	if err != nil {
		errmessage := fmt.Sprintf("Unable to loap parent element  by %s : %s", w.Request.Host, err.Error())
		w.SiteConfig.Logger.Error(errmessage)
		services.AddMessage("Could not load parent element.", "Error", w)
		w.Serve()
		return
	}
	if parent.Slugs == nil {
		slug := map[string]string{e.MongoId.Hex(): e.MongoId.Hex()}
		parent.Slugs = slug
	} else {
		parent.Slugs[e.MongoId.Hex()] = e.MongoId.Hex()
	}
	err = parent.Save(w)
	if err != nil {
		errmessage := fmt.Sprintf("Unable to save parent element  by %s : %s", w.Request.Host, err.Error())
		w.SiteConfig.Logger.Error(errmessage)
		services.AddMessage("Could not save parent element.", "Error", w)
		w.Serve()
		return
	}
	dynamic := services.Dynamic{
		Target:     parentid,
		Controller: "admin/element",
		Template:   "admin/element.html",
		Id:         parentid,
	}
	services.SetDynamic(dynamic, w)
	services.AddMessage("You have added a new element.", "Success", w)
	w.Serve()

}

// Controller for adding children to path.
func AddPathChild(w *wrapper.Wrapper) {
	parentid := w.APIParams[0]
	e := elements.NewElement()
	e.Title = "New Element"
	err := e.Save(w)
	if err != nil {
		errmessage := fmt.Sprintf("Unable to create new element  by %s : %s", w.Request.Host, err.Error())
		w.SiteConfig.Logger.Error(errmessage)
		services.AddMessage("Could not create a new element.", "Error", w)
		w.Serve()
		return
	}
	var parent paths.Path
	parent, err = paths.LoadPath(parentid, w)
	if err != nil {
		errmessage := fmt.Sprintf("Unable to loap path  by %s : %s", w.Request.Host, err.Error())
		w.SiteConfig.Logger.Error(errmessage)
		services.AddMessage("Could not load parent path.", "Error", w)
		w.Serve()
		return
	}
	parent.Elements = append(parent.Elements, e.MongoId.Hex())
	err = parent.Save(w)
	if err != nil {
		errmessage := fmt.Sprintf("Unable to save path by %s : %s", w.Request.Host, err.Error())
		w.SiteConfig.Logger.Error(errmessage)
		services.AddMessage("Could not add child element.", "Error", w)
		w.Serve()
		return
	}
	dynamic := services.Dynamic{
		Target:     "centereditor",
		Controller: "admin/path_elements",
		Template:   "admin/path_elements.html",
		Id:         parentid,
	}
	services.SetDynamic(dynamic, w)
	services.AddMessage("You have added a new element.", "Success", w)
	w.Serve()
	return
}

// Controller for adding existing element to wrapper element or path.
func AddExistingChild(w *wrapper.Wrapper) {
	var parenttype string
	if len(w.APIParams) > 1 {
		parenttype = w.APIParams[0]
	} else {
		http.Error(w.Writer, "Forbidden", 403)
		w.Serve()
		return
	}
	if w.Request.Method != "POST" {
		AddExistingChildForm(w)
		return
	}
	w.Shift()
	switch parenttype {
	case "wrapper":
		AddExistingWrapperSubmit(w)
		return
	case "paths":
		AddExistingPathSubmit(w)
		return
	default:
		http.Error(w.Writer, "Forbidden", 403)
		w.Serve()
	}
	return
}

// Controller for adding existing element form.
func AddExistingChildForm(w *wrapper.Wrapper) {
	elems, err := elements.ElementList(w)
	if err != nil {
		errmessage := fmt.Sprintf("Unable to retrieve a list of all elements: %s", err.Error())
		w.SiteConfig.Logger.Error(errmessage)
		services.AddMessage("There was a problem retrieving the element list.", "Error", w)
		w.Serve()
		return
	}
	options := make([]map[string]string, 0)
	for _, element := range elems {
		option := map[string]string{"name": element.Title, "value": element.MongoId.Hex(), "group": element.Controller}
		options = append(options, option)
	}
	f := form.NewForm()
	f.AddSelect("element", options).AddLabel("Element").Required()
	f.Register(w)
	w.SetTemplate("admin/form.html")
	w.SetPayload("form", f)
	w.Serve()
	return

}

// Controller for adding existing element to wrapper element form submission.
func AddExistingWrapperSubmit(w *wrapper.Wrapper) {
	parentid := w.APIParams[0]
	var post map[string]string
	err := form.GetValidFormData(w, &post)
	if err != nil {
		return
	}
	var parent elements.WrapperElement
	parent, err = elements.LoadWrapperElement(parentid, w)
	if err != nil {
		errmessage := fmt.Sprintf("Unable to loap parent element  by %s : %s", w.Request.Host, err.Error())
		w.SiteConfig.Logger.Error(errmessage)
		services.AddMessage("Could not load parent element.", "Error", w)
		w.Serve()
		return
	}
	parent.Elements = append(parent.Elements, post["element"])
	err = parent.Save(w)
	if err != nil {
		errmessage := fmt.Sprintf("Unable to save parent element  by %s : %s", w.Request.Host, err.Error())
		w.SiteConfig.Logger.Error(errmessage)
		services.AddMessage("Could not save parent element.", "Error", w)
		w.Serve()
		return
	}
	dynamic := services.Dynamic{
		Target:     parentid,
		Controller: "admin/element",
		Template:   "admin/element.html",
		Id:         parentid,
	}
	services.SetDynamic(dynamic, w)
	services.AddMessage("You have added an existing element.", "Success", w)
	w.Serve()
}

// Controller for adding existing element to path form submission.
func AddExistingPathSubmit(w *wrapper.Wrapper) {
	parentid := w.APIParams[0]
	var post map[string]string
	err := form.GetValidFormData(w, &post)
	if err != nil {
		return
	}
	var parent paths.Path
	parent, err = paths.LoadPath(parentid, w)
	if err != nil {
		errmessage := fmt.Sprintf("Unable to loap parent element  by %s : %s", w.Request.Host, err.Error())
		w.SiteConfig.Logger.Error(errmessage)
		services.AddMessage("Could not load parent element.", "Error", w)
		w.Serve()
		return
	}
	parent.Elements = append(parent.Elements, post["element"])
	err = parent.Save(w)
	if err != nil {
		errmessage := fmt.Sprintf("Unable to save parent element  by %s : %s", w.Request.Host, err.Error())
		w.SiteConfig.Logger.Error(errmessage)
		services.AddMessage("Could not save parent element.", "Error", w)
		w.Serve()
		return
	}
	dynamic := services.Dynamic{
		Target:     "centereditor",
		Controller: "admin/path_elements",
		Template:   "admin/path_elements.html",
		Id:         parentid,
	}
	services.SetDynamic(dynamic, w)
	services.AddMessage("You have added an existing element.", "Success", w)
	w.Serve()
}
