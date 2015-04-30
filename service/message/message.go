// Messages allows the user to pass growl style messages back to the browser.

package message

import (
	"github.com/mongolar/mongolar/wrapper"
)

// The basic required Message structure required for messaging.
type Message struct {
	Text     string
	Severity string
}

// Add a message to be served
func Add(m Message, w *wrapper.Wrapper) {
	v, err := w.GetPayload("messages")
	if err != nil {
		v := make([]Message)
	}
	v := append(v, m)
	w.SetPayload(v)
}

// Get all current messages
func Get(w *wrapper.Wrapper) []Message {
	v, err := w.GetPayload("messages")
	if err != nil {
		v := make([]Message)
	}
	return v

}

// Clears all messages
func Clear(w *wrapper.Wrapper) {
	v := make([]Message)
	w.SetPayload("messages", v)
}
