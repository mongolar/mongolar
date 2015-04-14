// The dynamic load will add a dynamic loader to your payload.
// This allows you to target any dom element to reload with new controllers and or ids
package dynamicload

import (
	"github.com/jasonrichardsmith/mongolar/wrapper"
)

// The DynamicLoad struct will be marshalled and added to you controller return payload.
type DynamicLoad struct {
	Target     string // The dynamic id of the dom element being targetted
	Controller string // the controller to invoke
	Id         string // An id to pass to the element
	Template   string // The template to use
}

// Once the above structure is created it can be passed to the dynamic function
// to be added to the payload.  It takes the DynamicLoad struct and the wrapper from the controller
func dynamic(d DynamicLoad, w *wrapper.Wrapper) {
	w.SetPayload("dynamic", d)
}
