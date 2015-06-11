package elements

type SlugElement struct {
	Slugs   map[string]string `bson:"controller_values"`
	Element `bson:"_,inline"`
}

func NewSlugElement() SlugElement {
	e := NewElement()
	cv := make(map[string]string)
	se := SlugElement{Element: e, Slugs: cv}
	return se
}
