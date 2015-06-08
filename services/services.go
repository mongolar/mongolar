package services

import (
	"github.com/mongolar/mongolar/wrapper"
)

// Add a message to be served
func AddMessage(t string, s string, w *wrapper.Wrapper) {
	message := map[string]string{"text": t, "severity": s}
	messages, err := w.GetAPayload("mongolar_messages")
	if err != nil {
		messages = []map[string]string{message}
	} else {
		messages = append(messages.([]map[string]string), message)
	}
	w.SetPayload("mongolar_messages", messages)
}

func ClearMessage(w *wrapper.Wrapper) {
	w.DeleteAPayload("mongolar_messages")
}

// Set will take a string of the page to redirect to and the wrapper
func Redirect(r string, w *wrapper.Wrapper) {
	w.SetPayload("mongolar_redirect", r)
}

// The DynamicLoad struct will be marshalled and added to you controller return payload.
type Dynamic struct {
	Target     string `json:"target,omitempty"`   // The dynamic id of the dom element being targetted
	Controller string `json:"controller"`         // the controller to invoke
	Id         string `json:"id,omitempty"`       // An id to pass to the element
	Template   string `json:"template,omitempty"` // The template to use
}

// Once the above structure is created it can be passed to the dynamic function
// to be added to the payload.  It takes the DynamicLoad struct and the wrapper from the controller
func SetDynamic(d Dynamic, w *wrapper.Wrapper) {
	w.SetPayload("mongolar_dynamics", []Dynamic{d})
}
