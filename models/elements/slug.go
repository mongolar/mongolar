package elements

import (
	"errors"
	"github.com/mongolar/mongolar/wrapper"
	"gopkg.in/mgo.v2/bson"
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

func SlugDeleteAllChild(id string, w *wrapper.Wrapper) error {
	if !bson.IsObjectIdHex(id) {
		return errors.New("Invalid Invalid Hex")
	}
	slugelements := make([]SlugElement, 0)
	c := w.DbSession.DB("").C("elements")
	i := c.Find(bson.M{"controller": "slug"}).Limit(50).Iter()
	err := i.All(&slugelements)
	if err != nil {
		if err.Error() == "not found" {
			return nil
		}
		return err
	}
	for _, element := range slugelements {
		change := false
		for slug, eid := range element.Slugs {
			if id == eid {
				delete(element.Slugs, slug)
				change = true
			}
		}
		if change {
			err := element.Save(w)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
