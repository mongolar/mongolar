package elements

type WrapperElements struct {
	Elements []string `bson:"elements" json:"elements"`
}
type WrapperElement struct {
	WrapperElements `bson:"controller_values" json:"content"`
	Element         `bson:"_,inline"`
}

func NewWrapperElement() WrapperElement {
	e := NewElement()
	els := make([]string, 0)
	cv := WrapperElement{WrapperElements{Elements: els}, e}
	return cv
}
