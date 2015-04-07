package responders


import(
	"net/http"
	"encode/json"
)

function ServeJson(w http.ResponseWriter, v interface{}){
  js, err := json.Marshal(v)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  w.Write(js)
}
