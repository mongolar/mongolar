// The path controller bootstraps a page and loads the direct children elements in the of the page itself,
// and sets the page template.
package path

import (
	"github.com/jasonrichardsmith/mongolar/wrapper"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)



func PathServe(w *wrapper.Wrapper) {
	c := getCollection(w.SiteConfig.DbSession, "path")
	p := w.Request.Header.Get('CurrentPath')
	ps := p.Split("/")
	ifgetPath(c, p, false)
	for i := 0; i <10; i++{
	}
	w.SetContent(v)
	w.Serve()
	return
}

func getCollection(s *mgo.Session, c string) *mgo.Collection {
	se = s.Copy()
	c := se.DB("").C(c)
	return c
}

func getPath(c *mgo.Collection, p string, w bool) {
	return c.Find(bson.M{"path": p, "wildcard": w}).One()
}

//
func pathMatch(dir string) (string, string) {
	var rejects []string
	for {
		//TODO perform mongodb query here
			
		// If query doesnt return anything
		if ?{
			rejects = append([]string{path.Base(dir)}, rejects...)
			dir = path.Dir(dir)
			continue
		}
		break
	}
	return dir, strings.Join(rejects, "/")
}
