package elements

import (
	"github.com/mongolar/mongolar/wrapper"
)

type MenuItem struct {
	Title    string      `bson:"title"`
	Url      string      `bson:"url"`
	Children interface{} `bson:"menu_items"`
}

type MenuElement struct {
	MenuItems []MenuItem `bson:"controller_values" json:"menu_items"`
	Element   `bson:"_,inline"`
}

func (me *MenuElement) Save(w *wrapper.Wrapper) error {
	return Save(me.Element.MongoId, me, w)
}

func NewMenuElement() MenuElement {
	e := NewElement()
	menuitems := make([]MenuItem, 0)
	me := MenuElement{Element: e, MenuItems: menuitems}
	return me
}

func LoadMenuElement(i string, w *wrapper.Wrapper) (MenuElement, error) {
	e := NewMenuElement()
	err := GetValidElement(i, "menu", &e, w)
	return e, err
}
