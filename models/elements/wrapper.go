package elements

import (
	"github.com/mongolar/mongolar/wrapper"
)

type WrapperElements struct {
	Elements []string `bson:"elements" json:"elements"`
}

func NewWrapperElements() WrapperElements {
	els := make([]string, 0)
	wels := WrapperElements{els}
	return wels
}

type WrapperElement struct {
	WrapperElements `bson:"controller_values" json:"content"`
	Element         `bson:",inline"`
}

func (we *WrapperElement) Save(w *wrapper.Wrapper) error {
	return Save(we.Element.MongoId, we, w)
}

func NewWrapperElement() WrapperElement {
	e := NewElement()
	els := NewWrapperElements()
	cv := WrapperElement{els, e}
	return cv
}

func LoadWrapperElement(i string, w *wrapper.Wrapper) (WrapperElement, error) {
	e := NewWrapperElement()
	err := GetById(i, &e, w)
	return e, err
}
