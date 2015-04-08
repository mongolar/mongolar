package controller

import (
	"github.com/jasonrichardsmith/mongolar/wrapper"
)

type ControllerMap map[string]Controller

func NewMap() ControllerMap {
	return make(ControllerMap)
}

type Controller func(*wrapper.Wrapper)
