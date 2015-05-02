package services

import (
	"github.com/mongolar/mongolar/wrapper"
)

// The basic required Message structure required for messaging.
type Message struct {
	Text     string
	Severity string
}

// Add a message to be served
func Add(m Message, w *wrapper.Wrapper) {
	v, err := w.GetPayload("messages")
	if err != nil {
		v := make([]Message)
	}
	v := append(v, m)
	w.SetPayload(v)
}

// Get all current messages
func Get(w *wrapper.Wrapper) []Message {
	v, err := w.GetPayload("messages")
	if err != nil {
		v := make([]Message)
	}
	return v

}

// Clears all messages
func Clear(w *wrapper.Wrapper) {
	v := make([]Message)
	w.SetPayload("messages", v)
}

// Set will take a string of the page to redirect to and the wrapper
func Set(r string, w *wrapper.Wrapper) {
	w.SetPayload("redirect", r)
}

// The DynamicLoad struct will be marshalled and added to you controller return payload.
type DynamicLoad struct {
	Target     string // The dynamic id of the dom element being targetted
	Controller string // the controller to invoke
	Id         string // An id to pass to the element
	Template   string // The template to use
}

// Once the above structure is created it can be passed to the dynamic function
// to be added to the payload.  It takes the DynamicLoad struct and the wrapper from the controller
func Set(d DynamicLoad, w *wrapper.Wrapper) {
	w.SetPayload("dynamic", d)
}
