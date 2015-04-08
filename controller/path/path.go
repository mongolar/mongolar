package path

import (
	"github.com/jasonrichardsmith/mongolar/wrapper"
)
func Serve(w *wrapper.Wrapper) {
	v := make(map[string]string)
	v['test'] = "Test"
	w.SetContent(v)
	w.Serve()
	return
}
