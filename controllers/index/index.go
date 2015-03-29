package index

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	//"io/ioutil"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//file, _ := ioutil.ReadFile(site.Directory + "/index.html")
	file := "test"
	fmt.Fprint(w, string(file))

}
