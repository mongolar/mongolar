package controller

import (
	"github.com/jasonrichardsmith/mongolar/wrapper"
)

type ControllerMap map[string]func(*wrapper.Wrapper)

func NewMap() ControllerMap {
	return make(ControllerMap)
}

type Controller func(*wrapper.Wrapper)
