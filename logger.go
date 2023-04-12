package graw

import (
	"io"
	"log"
)

func logger(l *log.Logger) *log.Logger {
	if l == nil {
		return log.New(io.Discard, "", 0)
	}

	return l
}
