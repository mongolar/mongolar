// Messages allows the user to pass growl style messages back to the browser. 

package message

import (
	"github.com/jasonrichardsmith/mongolar/wrapper"
)

// The basic required Message structure required for messaging.
type Message{
	Text string
	Severity string
}

func Add(m Message, w *wrapper.Wrapper) {
	//TODO grab a slice of values and add to the slices
}
