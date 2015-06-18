package elements

import (
	"github.com/mongolar/mongolar/wrapper"
)

type SlugElement struct {
	Slugs   map[string]string `bson:"controller_values"`
	Element `bson:"_,inline"`
}

func (se *SlugElement) Save(w *wrapper.Wrapper) error {
	return Save(se.Element.MongoId, se, w)
}

func NewSlugElement() SlugElement {
	e := NewElement()
	cv := make(map[string]string)
	se := SlugElement{Element: e, Slugs: cv}
	return se
}

func LoadSlugElement(i string, w *wrapper.Wrapper) (SlugElement, error) {
	e := NewSlugElement()
	err := GetValidElement(i, "slug", &e, w)
	return e, err
}
