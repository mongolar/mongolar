package basecontrollers

import (
	"github.com/mongolar/mongolar/controller"
)

// Base controllers includes all the basic controllers for Mongolar.
// DomainPublicValue - Retrieves value from public site config value.  Used for site wide values.
// Path - Returns the top level elements for a page, which will eventually bootstrap the page building process.
// Content - Returns the Values for elements labeled as content.
// Wrapper - Returns child element ids for elemements assigned as wrapper.
// Slug - Returns a content element for a wildcard path based on the slug value set in thes lug element.
// Menu - Returns a Menu for an element tagged as menu.

func GetControllerMap(cm controller.ControllerMap) {
	cm["domian_public_value"] = DomainPublicValue
	cm["path"] = PathValues
	cm["content"] = ContentValues
	cm["wrapper"] = WrapperValues
	cm["slug"] = SlugValues
	cm["menu"] = MenuValues
}
