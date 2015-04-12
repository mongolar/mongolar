package message

import (
	"github.com/jasonrichardsmith/mongolar/wrapper"
)

type Message{
	Target string
	Controller string
	Id	string
	Template string
}

func set(m Message, w *wrapper.Wrapper) {
}
