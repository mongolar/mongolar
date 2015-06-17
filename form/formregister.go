package form

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mongolar/mongolar/services"
	"github.com/mongolar/mongolar/wrapper"
	"gopkg.in/mgo.v2/bson"
	"time"
)

// Structure for form registration
type FormRegister struct {
	FormFields []*Field      `bson:"fields"`
	FormId     bson.ObjectId `bson:"_id"`
	SessionId  bson.ObjectId `bson:"session_id"`
	Created    time.Time     `bson:"created"`
}

func GetValidFormData(w *wrapper.Wrapper, post interface{}) error {
	p := make([]byte, w.Request.ContentLength)
	_, err := w.Request.Body.Read(p)
	if err != nil {
		errmessage := fmt.Sprintf("Error processing post values %s: %s", w.Request.Host, err.Error())
		w.SiteConfig.Logger.Error(errmessage)
		services.AddMessage("There was an issue processing your form.", "Error", w)
		w.Serve()
		return errors.New("Could not marshall Post values")
	}
	data := make(map[string]interface{})
	err = json.Unmarshal(p, &data)
	if err != nil {
		errmessage := fmt.Sprintf("Error processing post values %s: %s", w.Request.Host, err.Error())
		w.SiteConfig.Logger.Error(errmessage)
		services.AddMessage("There was an issue processing your form.", "Error", w)
		w.Serve()
		return errors.New("Could not marshall Post values")
	}
	register, reg_err := GetFormRegister(data["form_id"].(string), w)
	if reg_err != nil {
		errmessage := fmt.Sprintf("Invalid or expired form %s: %s", w.Request.Host, err.Error())
		w.SiteConfig.Logger.Error(errmessage)
		services.AddMessage("Your form was expired, please try again.", "Error", w)
		w.Serve()
		return errors.New("Inalid or expired form")
	}
	err = json.Unmarshal(p, post)
	if err != nil {
		errmessage := fmt.Sprintf("Error processing post values %s: %s", w.Request.Host, err.Error())
		w.SiteConfig.Logger.Error(errmessage)
		services.AddMessage("There was an issue processing your form.", "Error", w)
		w.Serve()
		return errors.New("Could not marshall Post values")
	}
	missing := register.ValidateRequired(data)
	if len(missing) > 0 {
		for _, label := range missing {

			message := fmt.Sprintf("%s is required.", label)
			services.AddMessage(message, "Error", w)
		}
		w.Serve()
		return errors.New("missing fields")
	}
	return nil

}

// Check for missing required fields.
func (fr *FormRegister) ValidateRequired(data map[string]interface{}) map[string]string {
	missing := make(map[string]string)
	for _, f := range fr.FormFields {
		if f.TemplateOptions.Required == true {
			if _, ok := data[f.Key]; !ok {
				missing[f.Key] = f.TemplateOptions.Label
			}
		}
	}
	return missing
}

// Retrieve a previously registered form by id
func GetFormRegister(i string, w *wrapper.Wrapper) (*FormRegister, error) {
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
