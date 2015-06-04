package form

import (
	"fmt"
	"github.com/mongolar/mongolar/services"
	"github.com/mongolar/mongolar/wrapper"
	"gopkg.in/mgo.v2/bson"
	"time"
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
	fi := &Field{
		Type:            "input",
		Key:             k,
		TemplateOptions: new(TemplateOptions),
	}
	f.Fields = append(f.Fields, fi)
	return fi
}

// Add a text are to form
func (f *Form) AddTextArea(k string) *Field {
	fi := &Field{
		Type:            "textarea",
		Key:             k,
		TemplateOptions: new(TemplateOptions),
	}
	f.Fields = append(f.Fields, fi)
	return fi
}

// Add a checkbox to form
func (f *Form) AddCheckBox(k string) *Field {
	fi := &Field{
		Type:            "checkbox",
		Key:             k,
		TemplateOptions: new(TemplateOptions),
	}
	f.Fields = append(f.Fields, fi)
	return fi
}

// Add a radio button to form
func (f *Form) AddRadio(k string, o []map[string]string) *Field {
	fo := &TemplateOptions{
		Options: o,
	}
	fi := &Field{
		Type:            "radio",
		Key:             k,
		TemplateOptions: fo,
	}
	f.Fields = append(f.Fields, fi)
	return fi
}

// Add a radio button to form
func (f *Form) AddSelect(k string, o []map[string]string) *Field {
	fo := &TemplateOptions{
		Options: o,
	}
	fi := &Field{
		Type:            "select",
		Key:             k,
		TemplateOptions: fo,
	}
	f.Fields = append(f.Fields, fi)
	return fi
}

// Add a radio button to form
func (f *Form) AddRepeatSection(k string, b string, fs []*Field) *Field {
	fo := &TemplateOptions{
		Fields:     fs,
		ButtonText: b,
	}
	fi := &Field{
		Type:            "repeatSection",
		Key:             k,
		TemplateOptions: fo,
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
		Created:    time.Now(),
	}
	c := w.DbSession.DB("").C("form_register")
	err := c.Insert(fr)
	return err
}

// Structure for form registration
type FormRegister struct {
	FormFields []*Field      `bson:"fields"`
	FormId     bson.ObjectId `bson:"_id"`
	SessionId  bson.ObjectId `bson:"session_id"`
	Created    time.Time     `bson:"created"`
}

// Retrieve a previously registered form by id
func GetRegisteredForm(i string, w *wrapper.Wrapper) (*FormRegister, error) {
	fr := new(FormRegister)
	c := w.DbSession.DB("").C("form_register")
	err := c.FindId(bson.ObjectIdHex(i)).One(fr)
	return fr, err
}

// Retrieve valid form based on id and session id
func GetValidRegForm(i string, w *wrapper.Wrapper) (*FormRegister, error) {
	fr := new(FormRegister)
	c := w.DbSession.DB("").C("form_register")
	b := bson.M{"session_id": w.Session.Id, "_id": bson.ObjectIdHex(i)}
	err := c.Find(b).One(fr)
	return fr, err
}

func GetValidRegFormM(i string, w *wrapper.Wrapper) (*FormRegister, error) {
	fr, err := GetValidRegForm(i, w)
	if err != nil {
		errmessage := fmt.Sprintf("Attempt to access invalid form %s by %s.", w.Post["form_id"].(string), w.Request.Host)
		w.SiteConfig.Logger.Error(errmessage)
		services.AddMessage("Invalid Form", "Error", w)
		w.Serve()
	}
	return fr, err
}

// Form fields structure
type Field struct {
	Type            string           `json:"type" bson:"type"`
	Hide            bool             `json:"hide,omitempty" bson:"hide,omitempty"`
	Key             string           `json:"key" bson:"key"`
	TemplateOptions *TemplateOptions `json:"templateOptions" bson:"templateOptions"`
	HideExpression  string           `json:"hideExpression,omitempty" bson:"hideExpression"`
}

type TemplateOptions struct {
	Options     []map[string]string `json:"options,omitempty" bson:"options,omitempty"`
	Label       string              `json:"label,omitempty" bson:"label,omitempty"`
	Required    bool                `json:"required,omitempty" bson:"required,omitempty"`
	Placeholder string              `json:"placeholder,omitempty" bson:"placeholder,omitempty"`
	Rows        int                 `json:"rows,omitempty" bson:"rows,omitempty"`
	Cols        int                 `json:"cols,omitempty" bson:"cols,omitempty"`
	Fields      []*Field            `json:"fields,omitempty" bson:"fields,omitempty"`
	ButtonText  string              `json:"btnText,omitempty" bson:"btnText,omitempty"`
}

// Add label to field
func (f *Field) AddLabel(l string) *Field {
	f.TemplateOptions.Label = l
	return f
}

// Add
func (f *Field) AddPlaceHolder(p string) *Field {
	f.TemplateOptions.Placeholder = p
	return f
}

func (f *Field) AddRowsCols(r int, c int) *Field {
	f.TemplateOptions.Rows = r
	f.TemplateOptions.Cols = c
	return f
}

func (f *Field) AddHideExpression(he string) *Field {
	//TODO: Fix this
	f.HideExpression = he
	return f
}

func (f *Field) Required() *Field {
	f.TemplateOptions.Required = true
	return f
}

func (f *Field) Hidden() *Field {
	f.Hide = true
	return f
}
