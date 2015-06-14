package paths

import (
	"errors"
	"github.com/mongolar/mongolar/wrapper"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"path"
	"strings"
)

type Path struct {
	MongoId  bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Path     string        `bson:"path" json:"path"`
	Wildcard bool          `bson:"wildcard" json:"wildcard"`
	Elements []string      `bson:"elements,omitempty" json:"elements,omitempty"`
	Template string        `bson:"template" json:"template"`
	Status   string        `bson:"status" json:"status"`
	Title    string        `bson:"title" json:"title"`
}

// Constructor for elements
func NewPath() Path {
	e := make([]string, 0)
	p := Path{Elements: e}
	return p
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
