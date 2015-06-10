package form

import (
	"github.com/mongolar/mongolar/wrapper"
	"gopkg.in/mgo.v2/bson"
	"time"
)

// Basic form structure, required by Formly
type Form struct {
	Fields   []*Field      `json:"formFields"`
	FormData interface{}   `json:"formData"`
	FormId   bson.ObjectId `json:"formId"`
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
