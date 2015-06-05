// Wrapper defines the struct to be passed back to the controller.
// It contains the entirety of the response and performs all marshalling and write operations
package wrapper

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mongolar/mongolar/configs"
	"gopkg.in/mgo.v2"
	"net/http"
	"strings"
)

// Wrapper structure required to be passed back to the Controller
type Wrapper struct {
	Writer     http.ResponseWriter    // The response writer
	Request    *http.Request          // The request
	Post       map[string]interface{} // Post data from AngularJS
	SiteConfig *configs.SiteConfig    // The configuration for the site being accessed
	Session    *Session               // Session for user
	Payload    map[string]interface{} // This is the sum of the payload that will be returned to the user
	DbSession  *mgo.Session           // The master MongoDb session that gets copied
	APIParams  []string
}

//Constructor for the Wrapper
func New(w http.ResponseWriter, r *http.Request, s *configs.SiteConfig) *Wrapper {
	wr := Wrapper{Writer: w, Request: r, SiteConfig: s}
	var err error
	if r.Method == "POST" {
		wr.Post, err = formPostData(r)
		if err != nil {
			errmessage := fmt.Sprintf("Could not load Post Data: %s", err.Error())
			wr.SiteConfig.Logger.Error(errmessage)
		}
	}
	wr.DbSession = s.DbSession.Copy()
	//Get session
	err = wr.NewSession()
	if err != nil {
		errmessage := fmt.Sprintf("Unable to create new session: %s", err.Error())
		wr.SiteConfig.Logger.Error(errmessage)
	}
	err = wr.SetSession()
	if err != nil {
		errmessage := fmt.Sprintf("Unable to save session to db sessions: %s", err.Error())
		wr.SiteConfig.Logger.Error(errmessage)
	}
	// Define payload
	wr.Payload = make(map[string]interface{})
	wr.APIParams = strings.Split(r.URL.Path, "/")
	wr.Shift()
	return &wr
}

// Shift API Params over by one
func (w *Wrapper) Shift() {
	w.APIParams = w.APIParams[1:]
}

// Load post data from AngulaJS
func formPostData(r *http.Request) (map[string]interface{}, error) {
	p := make(map[string]interface{})
	err := json.NewDecoder(r.Body).Decode(&p)
	return p, err
}

// Helper function for the controller to easily add its final content to the Payload
func (w *Wrapper) SetContent(c interface{}) {
	w.SetPayload("content", c)
}

// Helper function for the controller to easily add its final content to the Payload
func (w *Wrapper) SetTemplate(t string) {
	w.SetPayload("mongolartemplate", t)
}

// Helper function for the controller to easily add its final content to the Payload
func (w *Wrapper) SetClasses(c string) {
	w.SetPayload("mongolarclasses", c)
}

// Helper function for the controller to easily add its final content to the Payload
func (w *Wrapper) SetDynamicId(i string) {
	w.SetPayload("mongolardyn", i)
}

// Sets payload based on a keyvalue
func (w *Wrapper) SetPayload(n string, v interface{}) {
	w.Payload[n] = v
}

// Gets payload based on a keyvalue
func (w *Wrapper) GetAPayload(n string) (interface{}, error) {
	if v, ok := w.Payload[n]; ok {
		return v, nil
	}
	err := errors.New("Payload value not set")
	return nil, err
}

// Gets payload based on a keyvalue
func (w *Wrapper) DeleteAPayload(n string) {
	delete(w.Payload, n)
}

// Gets entire payload
func (w *Wrapper) GetPayload() map[string]interface{} {
	return w.Payload
}

// The final serve function.  This will marshall the payload and serve it to the user.
func (w *Wrapper) Serve() {
	js, err := json.Marshal(w.Payload)
	if err != nil {
		http.Error(w.Writer, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Writer.Write(js)
	w.DbSession.Close()
	return
}
