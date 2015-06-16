package elements

type ContentValues struct {
	Content map[string]interface{} `bson:"content"`
	Type    string                 `bson:"type"`
}

type ContentElement struct {
	ContentValues `bson:"controller_values" json:"content"`
	Element       `bson:"_,inline"`
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
	err := GetById(i, e, w)
	return e, err
}
