package graw

import (
	"log"

	"github.com/turnage/graw/reddit"
)

func foreman(kill chan<- bool, errs <-chan error, logger *log.Logger) error {
	for err := range errs {
		switch err {
		case nil:
		case reddit.BusyErr:
			logger.Printf("Reddit was busy; staying up.")
		case reddit.GatewayErr:
			logger.Printf("Reddit connection faulted; staying up.")
		default:
			close(kill)
			return err
		}
	}

	return nil
}
