package controller

import (
	"encoding/json"
)

type Controller interface {
	GetContent(*http.Request, *site.SiteConfig) *ControllerResponse
}

type ControllerResponse struct {
	Content   json.RawMessage
	Template  string
	DynamicId string
}

func New() *ControllerResponse {
	return new(ControllerResponse)
}
