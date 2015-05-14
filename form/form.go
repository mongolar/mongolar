package form

import (
	"github.com/mongolar/mongolar/wrapper"
	"gopkg.in/mgo.v2/bson"
)

// Basic form structure, required by Formly
type Form struct {
	Fields   []*Field               `json:"formFields"`
	FormData map[string]interface{} `json:"formData"`
	FormId   bson.ObjectId          `json:"formId"`
}

// Constructor for form
func NewForm() *Form {
	fd := make(map[string]interface{})
	fi := make([]*Field, 0)
	f := Form{
		Fields:   fi,
		FormData: fd,
		FormId:   bson.NewObjectId(),
	}
	return &f
}

// Add a text field to form
func (f *Form) AddText(k string, t string) *Field {
	to := map[string]interface{}{"type": t, "label": ""}
	fi := &Field{
		Type:            "input",
		Key:             k,
		TemplateOptions: to,
	}
	f.Fields = append(f.Fields, fi)
	return fi
}

// Add a text are to form
func (f *Form) AddTextArea(k string) *Field {
	to := map[string]interface{}{"label": ""}
	fi := &Field{
		Type:            "textarea",
		Key:             k,
		TemplateOptions: to,
	}
	f.Fields = append(f.Fields, fi)
	return fi
}

// Add a checkbox to form
func (f *Form) AddCheckBox(k string) *Field {
	to := map[string]interface{}{"label": ""}
	fi := &Field{
		Type:            "checkbox",
		Key:             k,
		TemplateOptions: to,
	}
	f.Fields = append(f.Fields, fi)
	return fi
}

// Add a radio button to form
func (f *Form) AddRadio(k string, o []map[string]string) *Field {
	to := map[string]interface{}{"options": o, "label": ""}
	fi := &Field{
		Type:            "radio",
		Key:             k,
		TemplateOptions: to,
	}
	f.Fields = append(f.Fields, fi)
	return fi
}

// Add a radio button to form
func (f *Form) AddSelect(k string, o []map[string]string) *Field {
	to := map[string]interface{}{"options": o, "label": ""}
	fi := &Field{
		Type:            "radio",
		Key:             k,
		TemplateOptions: to,
	}
	f.Fields = append(f.Fields, fi)
	return fi
}

// Add a radio button to form
func (f *Form) AddRepeatSection(k string, b string, fs []*Field) *Field {
	to := map[string]interface{}{"fields": fs, "label": "", "btnText": b}
	fi := &Field{
		Type:            "repeatSection",
		Key:             k,
		TemplateOptions: to,
	}
	f.Fields = append(f.Fields, fi)
	return fi
}

// Register the form in the database
func (f *Form) Register(w *wrapper.Wrapper) error {
	fr := FormRegister{
		FormFields: f.Fields,
		FormId:     f.FormId,
		SessionId:  w.Session.Id,
	}
	se := w.SiteConfig.DbSession.Copy()
	defer se.Close()
	c := se.DB("").C("form_register")
	err := c.Insert(fr)
	return err
}

// Structure for form registration
type FormRegister struct {
	FormFields []*Field      `bson:"fields"`
	FormId     bson.ObjectId `bson:"_id"`
	SessionId  string        `bson:"session_id"`
}

// Retrieve a previously registered form by id
func GetRegisteredForm(i string, w *wrapper.Wrapper) (*FormRegister, error) {
	fr := new(FormRegister)
	se := w.SiteConfig.DbSession.Copy()
	defer se.Close()
	c := se.DB("").C("form_register")
	err := c.FindId(bson.ObjectIdHex(i)).One(fr)
	return fr, err
}

// Retrieve valid form based on id and session id
func GetValidRegForm(i string, w *wrapper.Wrapper) (*FormRegister, error) {
	fr := new(FormRegister)
	se := w.SiteConfig.DbSession.Copy()
	defer se.Close()
	c := se.DB("").C("form_register")
	b := bson.M{"session_id": w.Session.Id, "_id": bson.ObjectIdHex(i)}
	err := c.Find(b).One(fr)
	return fr, err
}

// Form fields structure
type Field struct {
	Type            string                 `json:"type" bson:"type"`
	Hide            bool                   `json:"hide,omitempty" bson:"hide,omitempty"`
	Key             string                 `json:"key" bson:"key"`
	TemplateOptions map[string]interface{} `json:"templateOptions" bson:"template_options"`
	HideExpression  string                 `json:"hideExpression,omitempty" bson:"hide_expression"`
}

// Add label to field
func (f *Field) AddLabel(l string) *Field {
	f.TemplateOptions["label"] = l
	return f
}

// Add
func (f *Field) AddPlaceHolder(p string) *Field {
	f.TemplateOptions["placeholder"] = p
	return f
}

func (f *Field) AddRowsCols(r int, c int) *Field {
	f.TemplateOptions["rows"] = r
	f.TemplateOptions["cols"] = c
	return f
}

func (f *Field) AddHideExpression(he string) *Field {
	//f.HideExpression = he
	return f
}

func (f *Field) Required() *Field {
	f.TemplateOptions["required"] = true
	return f
}

func (f *Field) Hidden() *Field {
	f.Hide = true
	return f
}
