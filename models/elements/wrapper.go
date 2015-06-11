package elements

type WrapperElement struct {
	Element  `bson:"_,inline"`
	Elements []string `json:"elements"`
}

func NewWrapperElement() WrapperElement {
	e := NewElement()
	els := make([]string, 0)
	cv := WrapperElement{Element: e, Elements: els}
	return cv
}
