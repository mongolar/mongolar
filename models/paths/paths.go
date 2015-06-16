package paths

import (
	"errors"
	"github.com/mongolar/mongolar/wrapper"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"path"
	"strings"
)

type PathElements struct {
	Elements []string `bson:"elements,omitempty,inline" json:"elements,omitempty"`
}

func NewPathElements() PathElements {
	e := make([]string, 0)
	p := PathElements{Elements: e}
	return p
}

type Path struct {
	MongoId  bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Path     string        `bson:"path" json:"path"`
	Wildcard bool          `bson:"wildcard" json:"wildcard"`
	Template string        `bson:"template" json:"template"`
	Status   string        `bson:"status" json:"status"`
	Title    string        `bson:"title" json:"title"`
	PathElements
}

// Constructor for paths
func NewPath() Path {
	e := NewPathElements()
	id := bson.NewObjectId()
	p := Path{MongoId: id, PathElements: e}
	return p
}

// Constructor for existing paths
func LoadPath(i string, w *wrapper.Wrapper) (Path, error) {
	e := NewPathElements()
	p := Path{PathElements: e}
	err := p.GetById(i, w)
	return p, err
}

//Save an element in its current state.
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
	_, err := c.Upsert(p.MongoId, p)
	if err != nil {
		return err
	}
	return nil
}

// Get Path by Id
func (p *Path) GetById(i string, w *wrapper.Wrapper) error {
	if !bson.IsObjectIdHex(i) {
		return errors.New("Invalid Id Hex")
	}
	c := w.DbSession.DB("").C("paths")
	err := c.FindId(bson.ObjectIdHex(i)).One(&p)
	return err
}

// Path matching query
func (p *Path) PathMatch(u string, s string, c *mgo.Collection) (string, error) {
	var rejects []string
	w := false
	var err error
	for {
		b := bson.M{"path": u, "wildcard": w, "status": s}
		err = c.Find(b).One(p)
		w = true
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
