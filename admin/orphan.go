// Orphan elements are a list of elements that have not been asigned to
// a path, wrapper element, slug element.

package admin

import (
	"fmt"
	"github.com/mongolar/mongolar/models/elements"
	"github.com/mongolar/mongolar/models/paths"
	"github.com/mongolar/mongolar/services"
	"github.com/mongolar/mongolar/wrapper"
	"gopkg.in/mgo.v2/bson"
)

func OrphanElements(w *wrapper.Wrapper) {
	assigned := make([]bson.ObjectId, 0)
	paths, err := paths.PathList(w)
	if err != nil {
		errmessage := fmt.Sprintf("Could not retrieve path elements for orphan list: %s", err.Error())
		w.SiteConfig.Logger.Error(errmessage)
		services.AddMessage("Could not retrieve path elements.", "Error", w)
		w.Serve()
	}
	for _, path := range paths {
		for _, element := range path.Elements {
			id := bson.ObjectIdHex(element)
			assigned = append(assigned, id)
		}
	}
	wrappers := make([]elements.WrapperElement, 0)
	c := w.DbSession.DB("").C("elements")
	s := bson.M{"controller": "wrapper"}
	i := c.Find(s).Limit(50).Iter()
	err = i.All(&wrappers)
	if err != nil {
		errmessage := fmt.Sprintf("Could not retrieve wrapper elements for orphan list: %s", err.Error())
		w.SiteConfig.Logger.Error(errmessage)
		services.AddMessage("Could not retrieve wrapper elements.", "Error", w)
		w.Serve()
	}
	for _, wrapper := range wrappers {
		for _, eid := range wrapper.Elements {
			bsonid := bson.ObjectIdHex(eid)
			assigned = append(assigned, bsonid)
		}
	}
	slugs := make([]elements.SlugElement, 0)
	s = bson.M{"controller": "slug"}
	i = c.Find(s).Limit(50).Iter()
	err = i.All(&slugs)
	if err != nil {
		errmessage := fmt.Sprintf("Could not retrieve slug elements for orphan list: %s", err.Error())
		w.SiteConfig.Logger.Error(errmessage)
		services.AddMessage("Could not retrieve slug elements.", "Error", w)
		w.Serve()
	}
	for _, slug := range slugs {
		for _, eid := range slug.Slugs {
			bsonid := bson.ObjectIdHex(eid)
			assigned = append(assigned, bsonid)
		}
	}
	unassigned := new([]elements.Element)
	s = bson.M{"_id": bson.M{"$nin": assigned}}
	i = c.Find(s).Limit(50).Iter()
	err = i.All(unassigned)
	if err != nil {
		errmessage := fmt.Sprintf("Could not retrieve unassigned elements: %s", err.Error())
		w.SiteConfig.Logger.Error(errmessage)
		services.AddMessage("Could not retrieve unassigned elements.", "Error", w)
		w.Serve()
	}
	w.SetTemplate("admin/orphan_path_elements.html")
	w.SetPayload("elements", unassigned)
	w.Serve()
	return
}
