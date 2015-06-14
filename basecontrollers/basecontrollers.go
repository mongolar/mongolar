package basecontrollers

import (
	"github.com/mongolar/mongolar/controller"
)

func GetControllerMap(cm controller.ControllerMap) {
	cm["domian_public_value"] = DomainPublicValue
	cm["path"] = PathValues
	cm["content"] = ContentValues
	cm["wrapper"] = WrapperValues
	cm["slug"] = SlugValues
}
