package service

import (
	"encoding/json"
)

type Service interface {
	getService() json.RawMessage
}
