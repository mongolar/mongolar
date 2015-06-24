package admin

import (
	"fmt"
	"github.com/mongolar/mongolar/form"
	"github.com/mongolar/mongolar/services"
	"github.com/mongolar/mongolar/wrapper"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"reflect"
	"strings"
)

type ContentType struct {
	Form    []*form.Field `bson:"form,omitempty" json:"content_form"`
	Type    string        `bson:"type,omitempty" json:"type"`
	MongoId bson.ObjectId `bson:"_id" json:"id"`
}

func GetContentType(w *wrapper.Wrapper) {
	c := w.DbSession.DB("").C("content_types")
	i := bson.M{"_id": bson.ObjectIdHex(w.APIParams[0])}
	var ct ContentType
	err := c.Find(i).One(&ct)
	if err != nil {
		errmessage := fmt.Sprintf("Content Type not found %s : %s", w.APIParams[0], err.Error())
		w.SiteConfig.Logger.Error(errmessage)
		services.AddMessage("Your content types was not found.", "Error", w)
		w.Serve()
		return
	}
	w.SetPayload("content_type", ct)
	w.Serve()
	return
}

func EditContentType(w *wrapper.Wrapper) {
	if len(w.APIParams) < 1 {
		http.Error(w.Writer, "Forbidden", 403)
		w.Serve()
		return
	}
	if w.Request.Method != "POST" {
		EditContentTypeForm(w)
		return
	}
	EditContentTypeSubmit(w)
	return
}

func EditContentTypeForm(w *wrapper.Wrapper) {
	f := form.NewForm()
	ct := new(ContentType)
	if w.APIParams[0] != "new" {
		c := w.DbSession.DB("").C("content_types")
		i := bson.M{"_id": bson.ObjectIdHex(w.APIParams[0])}
		err := c.Find(i).One(ct)
		if err != nil {
			errmessage := fmt.Sprintf("Content Type not found %s : %s", w.APIParams[0], err.Error())
			w.SiteConfig.Logger.Error(errmessage)
			services.AddMessage("Your content types was not found ", "Error", w)
			w.Serve()
			return
		}
		var elements []map[string]interface{}
		for _, field := range ct.Form {
			element := make(map[string]interface{})
			element["type"] = field.Type
			element["key"] = field.Key
			element["label"] = field.TemplateOptions.Label
			element["placeholder"] = field.TemplateOptions.Placeholder
			element["rows"] = field.TemplateOptions.Rows
			element["cols"] = field.TemplateOptions.Cols
			element["options"] = ""
			for _, opt := range field.TemplateOptions.Options {
				element["options"] = fmt.Sprintf("%s%s|%s\n", element["options"], opt["name"], opt["value"])
			}
			elements = append(elements, element)
		}
		data := make(map[string]interface{})
		data["elements"] = elements
		data["content_type"] = ct.Type
		f.FormData = data
	} else {
		data := make(map[string]interface{})
		fd := make([]map[string]string, 0)
		data["elements"] = fd
		data["content_type"] = ""
		f.FormData = data
	}
	f.AddText("content_type", "text").AddLabel("Content Type Name")
	f.AddRepeatSection("elements", "Add another field", FieldFormGroup())
	f.Register(w)
	w.SetPayload("form", f)
	w.SetTemplate("admin/form.html")
	w.Serve()
	return

}

func EditContentTypeSubmit(w *wrapper.Wrapper) {
	post := make(map[string]interface{})
	err := form.GetValidFormData(w, &post)
	if err != nil {
		return
	}
	elements := reflect.ValueOf(post["elements"])
	f := form.NewForm()
	for i := 0; i < elements.Len(); i++ {
		var field *form.Field
		element := elements.Index(i).Interface().(map[string]interface{})
		switch element["type"].(string) {
		case "input":
			field = f.AddText(element["key"].(string), "text")
		case "textarea":
			field = f.AddTextArea(element["key"].(string))
		case "radio":
			values := strings.Split(element["options"].(string), "\n")
			opt := make([]map[string]string, 0)
			for _, value := range values {
				namval := strings.Split(value, "|")
				if len(namval) < 2 {
					errmessage := fmt.Sprintf("Attempt to set incorrect form option %s by %s", value, w.Request.Host)
					w.SiteConfig.Logger.Error(errmessage)
					services.AddMessage("Your options must be of the format Name|Value", "Error", w)
					w.Serve()
					return
				}
				newval := map[string]string{
					"name":  namval[0],
					"value": namval[1],
				}
				opt = append(opt, newval)
			}
			field = f.AddRadio(element["key"].(string), opt)
		case "checkbox":
			field = f.AddCheckBox(element["key"].(string))
		default:
			//TODO messaging and logging
			return
		}
		if _, ok := element["label"]; ok {
			if element["label"].(string) != "" {
				field.AddLabel(element["label"].(string))
			}
		}
		if _, ok := element["placeholder"]; ok {
			if element["placeholder"].(string) != "" {
				field.AddPlaceHolder(element["placeholder"].(string))
			}
		}
		if _, ok := element["rows"]; ok {
			if _, ok := element["cols"]; ok {
				if element["rows"].(float64) != 0 && element["cols"] != 0 {
					field.AddRowsCols(int(element["rows"].(float64)), int(element["cols"].(float64)))
				}
			}
		}
	}

	var id bson.ObjectId
	if post["mongolarid"].(string) == "new" {
		id = bson.NewObjectId()
	} else {
		id = bson.ObjectIdHex(post["mongolarid"].(string))
	}
	ct := ContentType{
		Form:    f.Fields,
		Type:    post["content_type"].(string),
		MongoId: id,
	}
	s := bson.M{"_id": id}
	c := w.DbSession.DB("").C("content_types")
	_, err = c.Upsert(s, ct)
	if err != nil {
		errmessage := fmt.Sprintf("Cannnot save content type %s : %s", post["mongolarid"].(string), err.Error())
		w.SiteConfig.Logger.Error(errmessage)
		services.AddMessage("Unable to save content type.", "Error", w)
		w.Serve()
		return
	}
	services.AddMessage("Content type saved.", "Success", w)
	dynamic := services.Dynamic{
		Target:     "contenttypelist",
		Controller: "admin/all_content_types",
		Template:   "admin/content_type_list.html",
	}
	services.SetDynamic(dynamic, w)
	w.Serve()
	return

}

func GetAllContentTypes(w *wrapper.Wrapper) {
	c := w.DbSession.DB("").C("content_types")
	var cts []ContentType
	err := c.Find(nil).Limit(50).Iter().All(&cts)
	if err != nil {
		errmessage := fmt.Sprintf("Unable to retrieve a list of content types.")
		w.SiteConfig.Logger.Error(errmessage)
		services.AddMessage("Unable to retrieve a list of elements.", "Error", w)
		w.Serve()
		return
	}
	w.SetPayload("content_types", cts)
	w.Serve()
	return
}

func FieldFormGroup() []*form.Field {
	ft := []map[string]string{
		map[string]string{"name": "Text Field", "value": "input"},
		map[string]string{"name": "TextArea Field", "value": "textarea"},
		map[string]string{"name": "Radio Buttons", "value": "radio"},
		map[string]string{"name": "Checkbox", "value": "checkbox"},
	}
	f := form.NewForm()
	f.AddRadio("type", ft).AddLabel("Field Type").Required()
	f.AddText("key", "text").AddLabel("Key").Required()
	f.AddText("label", "text").AddLabel("Label")
	f.AddText("placeholder", "text").AddLabel("Placeholder")
	f.AddTextArea("options").AddLabel("Options")
	f.AddText("cols", "text").AddLabel("Columns")
	f.AddText("rows", "text").AddLabel("Rows")
	return f.Fields
}
