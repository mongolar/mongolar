// Controller Map is a list of API endpoints that allow the admin to
// compile the api calls needed to render a site with Mongolar.

// A controller is any function that will take a wrapper as an argument.

package controller

import (
	"github.com/mongolar/mongolar/wrapper"
)

// The map structure for Controllers
type ControllerMap map[string]func(*wrapper.Wrapper)

// Creates a map for controllers
func NewMap() ControllerMap {
	return make(ControllerMap)
}
