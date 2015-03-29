package path

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func Path(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "path, %s!\n", ps.ByName("fullpath"))
}
