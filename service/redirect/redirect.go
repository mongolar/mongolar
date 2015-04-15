// The redirect service allows you to redirect the page the user is on
package redirect

import (
	"github.com/jasonrichardsmith/mongolar/wrapper"
)

// Set will take a string of the page to redirect to and the wrapper
func Set(r string, w *wrapper.Wrapper) {
	w.SetPayload("redirect", string)
}
