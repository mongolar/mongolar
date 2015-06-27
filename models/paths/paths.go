package paths

import (
	"errors"
	"github.com/mongolar/mongolar/wrapper"
	"gopkg.in/mgo.v2/bson"
	"path"
	"strings"
)

// Path Structure
type Path struct {
	MongoId      bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Path         string        `bson:"path" json:"path"`
	Wildcard     bool          `bson:"wildcard" json:"wildcard"`
	Template     string        `bson:"template" json:"template"`
	Status       string        `bson:"status" json:"status"`
	Title        string        `bson:"title" json:"title"`
	PathElements `bson:",inline,omitempty"`
}

// PathElements Structure to define the Elements in a path for easy json and bson marshalling.
type PathElements struct {
	Elements []string `bson:"elements,omitempty" json:"elements,omitempty"`
}

// Constructor for paths
func NewPath() Path {
	e := NewPathElements()
	id := bson.NewObjectId()
	p := Path{MongoId: id, PathElements: e}
	return p
}

// Constructor for path elements.
func NewPathElements() PathElements {
	e := make([]string, 0)
	p := PathElements{Elements: e}
	return p
}

// Construct and Get Path by ID
func LoadPath(i string, w *wrapper.Wrapper) (Path, error) {
	e := NewPathElements()
	p := Path{PathElements: e}
	err := p.GetById(i, w)
	return p, err
}

//Save a path in its current state.
func (p *Path) Save(w *wrapper.Wrapper) error {
	if !p.MongoId.Valid() {
		p.MongoId = bson.NewObjectId()
	}
	if p.Path == "" {
		return errors.New("Path required")
	}
	if p.Template == "" {
		return errors.New("Template required")
	}
	if p.Status == "" {
		return errors.New("Status required")
	}
	c := w.DbSession.DB("").C("paths")
	_, err := c.Upsert(bson.M{"_id": p.MongoId}, p)
	if err != nil {
		return err
	}
	return nil
}

// Get Path by Id on constructed Path
func (p *Path) GetById(i string, w *wrapper.Wrapper) error {
	if !bson.IsObjectIdHex(i) {
		return errors.New("Invalid Id Hex")
	}
	c := w.DbSession.DB("").C("paths")
	err := c.FindId(bson.ObjectIdHex(i)).One(&p)
	return err
}

// Given a URL and status this query will attempt to find a matching path.
// First the query qill attempt to implicitly match the url without a wildcard
// Then it will attempt to match based on the same path value having a wildcard
// After that it will remove sections of the url each time looking for a wildcard match
// If it does not find any  matches it retrns the last error.
func (p *Path) PathMatch(u string, s string, w *wrapper.Wrapper) (string, error) {
	c := w.DbSession.DB("").C("paths")
	var rejects []string
	wildcard := false
	var err error
	for {
		b := bson.M{"path": u, "wildcard": wildcard, "status": s}
		err = c.Find(b).One(p)
		wildcard = true
		// If query doesnt return anything
		if err != nil {
			rejects = append([]string{path.Base(u)}, rejects...)
			u = path.Dir(u)
			if u == "/" {
				break
			}
			continue
		}
		break
	}
	return strings.Join(rejects, "/"), err
}

// Add a child element to an existing Path
func AddChild(pathid string, elementid string, w *wrapper.Wrapper) error {
	if !bson.IsObjectIdHex(elementid) {
		return errors.New("Invalid Hex Id")
	}
	p, err := LoadPath(pathid, w)
	if err != nil {
		return err
	}
	p.Elements = append(p.Elements, elementid)
	return p.Save(w)
}

// Delete a path by id.
func Delete(id string, w *wrapper.Wrapper) error {
	if !bson.IsObjectIdHex(id) {
		return errors.New("Invalid Invalid Hex")
	}
	c := w.DbSession.DB("").C("paths")
	i := bson.M{"_id": bson.ObjectIdHex(id)}
	return c.Remove(i)
}

// Delete all references to a child element in all paths by id.
func DeleteAllChild(id string, w *wrapper.Wrapper) error {
	if !bson.IsObjectIdHex(id) {
		return errors.New("Invalid Invalid Hex")
	}
	s := bson.M{"elements": id}
	d := bson.M{"$pull": bson.M{"elements": id}}
	c := w.DbSession.DB("").C("paths")
	_, err := c.UpdateAll(s, d)
	return err

}

// Get all Paths
func PathList(w *wrapper.Wrapper) ([]Path, error) {
	pl := make([]Path, 0)
	c := w.DbSession.DB("").C("paths")
	i := c.Find(nil).Limit(50).Iter()
	err := i.All(&pl)
	if err != nil {
		return nil, err
	}
	return pl, nil
}
