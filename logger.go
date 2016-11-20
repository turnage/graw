package graw

import (
	"io/ioutil"
	"log"
)

func logger(l *log.Logger) *log.Logger {
	if l == nil {
		return log.New(ioutil.Discard, "", 0)
	}

	return l
}
