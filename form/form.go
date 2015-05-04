package form

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Basic form structure, required by Formly
type Form struct {
	Fields   []*Field          `json: "formFields"`
	FormData map[string]string `json: "formData"`
	FormId   string            `json: "formId"`
}

// Constructor for form
func NewForm() *Form {
	fd := make(map[string]string)
	fi := make([]Field, 1)
	f := Form{
		Fields:   fi,
		FormData: fd,
		FormId:   bson.NewObjectId().String(),
	}
	return f
}

// Add a text field to form
func (f *Form) AddText(k string, t string) *Field {
	to := map[string]interface{}{"type": t}
	fi := Field{
		Type:            "text",
		Key:             k,
		TemplateOptions: to,
	}
	f.Fields = merge(f.Fields, fi)
	return fi
}

// Add a text are to form
func (f *Form) AddTextArea(k string) *Field {
	to := make(map[string]interface{})
	fi := Field{
		Type:            "textarea",
		Key:             k,
		TemplateOptions: to,
	}
	f.Fields = merge(f.Fields, fi)
	return fi
}

// Add a checkbox to form
func (f *Form) AddCheckBox(k string) *Field {
	to := make(map[string]interface{})
	fi := Field{
		Type:            "checkbox",
		Key:             k,
		TemplateOptions: to,
	}
	f.Fields = merge(f.Fields, fi)
	return fi
}

// Add a radio button to form
func (f *Form) AddRadio(k string, o []map[string]string) *Field {
	to := map[string]interface{}{"options": o}
	fi := Field{
		Type:            "radio",
		Key:             k,
		TemplateOptions: to,
	}
	f.Fields = merge(f.Fields, fi)
	return fi
}

// Register the form in the database
func (f *Form) Register(s session.Session, ds *mgo.Session) error {
	fr := FormRegister{
		FormFields: f.Fields,
		FormId:     bson.ObjectIdHex(f.FormId),
		SessionId:  s.Id,
	}
	se := ds.Copy()
	defer se.Close()
	c := se.DB("").C("form_register")
	err := c.Insert(fr)
	return err
}

// Structure for form registration
type FormRegister struct {
	FormFields Form          `bson: "fields"`
	FormId     bson.ObjectId `bson: "_id"`
	SessionId  string        `bson: "session_id"`
}

// Retrieve a previously registered form by id
func GetRegisteredForm(i string, s *mgo.Session) (*FormRegister, error) {
	fr := new(FormRegister)
	se := s.Copy()
	defer se.Close()
	c := se.DB("").C("form_register")
	err := c.FindId(bson.ObjectIdHex(i)).One(fr)
	return fr, err
}

// Retrieve valid form based on id and session id
func GetValidRegForm(i string, se session.Session, s *mgo.Session) (*FormRegister, error) {
	fr := new(FormRegister)
	se := s.Copy()
	defer se.Close()
	c := se.DB("").C("form_register")
	b := bson.M{"session_id": se, "_id": bson.ObjectIdHex(i)}
	err = c.Find(b).One(fr)
	return fr, err
}

// Form fields structure
type Field struct {
	Type            string                 `json: "type"`
	Key             string                 `json: "key"`
	TemplateOptions map[string]interface{} `json: "templateOptions"`
}

// Add label to field
func (f *Field) AddLabel(l string) {
	f.TemplateOptions["label"] = l
	return f
}

// Add
func (f *Field) AddPlaceHolder(p string) {
	f.TemplateOptions["placeholder"] = p
	return f
}

func (f *Field) AddRowsCols(r int, c int) {
	f.TemplateOptions["rows"] = r
	f.TemplateOptions["cols"] = c
	return f
}

func (f *Field) Required() {
	f.TemplateOptions["required"] = true
	return f
}

func (f *Field) Hidden() {
	f.TemplateOptions["hidden"] = true
	return f
}
