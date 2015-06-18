package elements

import (
	"github.com/mongolar/mongolar/wrapper"
)

type ContentValues struct {
	Content map[string]interface{} `bson:"content"`
	Type    string                 `bson:"type"`
}

type ContentElement struct {
	ContentValues `bson:"controller_values" json:"content"`
	Element       `bson:",inline"`
}

func (ce *ContentElement) Save(w *wrapper.Wrapper) error {
	return Save(ce.Element.MongoId, ce, w)
}

func NewContentElement() ContentElement {
	e := NewElement()
	c := make(map[string]interface{})
	contentv := ContentValues{Content: c}
	ce := ContentElement{Element: e, ContentValues: contentv}
	return ce
}

func LoadContentElement(i string, w *wrapper.Wrapper) (ContentElement, error) {
	e := NewContentElement()
	err := GetValidElement(i, "content", &e, w)
	return e, err
}
