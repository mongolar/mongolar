package logger

import (
	"github.com/Sirupsen/logrus"
	"log"
	"os"
	"path/filepath"
)

// Wrapper for builidng logrus logger
func New(f string) *logrus.Logger {
	d := filepath.Dir(f)
	_, err := os.Stat(d)
	if err != nil {
		os.MkdirAll(d, 770)
	}
	_, err0 := os.Stat(f)
	var fi *os.File
	var err1 error
	if err0 == nil {
		fi, err1 = os.OpenFile(f, os.O_RDWR|os.O_APPEND, 770)
	} else {
		fi, err1 = os.Create(f)
	}
	if err1 != nil {
		log.Fatal(err1.Error())
	}
	l := logrus.New()
	l.Out = fi
	return l
}
