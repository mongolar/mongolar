package elements

import (
	"errors"
	"github.com/mongolar/mongolar/wrapper"
	"gopkg.in/mgo.v2/bson"
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
	err := GetValidElement(i, "wrapper", &e, w)
	return e, err
}

func WrapperDeleteAllChild(id string, w *wrapper.Wrapper) error {
	if !bson.IsObjectIdHex(id) {
		return errors.New("Invalid Invalid Hex")
	}
	c := w.DbSession.DB("").C("elements")
	s := bson.M{"controller_values.elements": id}
	d := bson.M{"$pull": bson.M{"controller_values.elements": id}}
	_, err := c.UpdateAll(s, d)
	return err
}
