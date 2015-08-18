package bot

import (
	"time"
)

// Controller defines the interface for bots to interact with the engine. These
// methods are requests to the engine to perform actions on behalf of the bot,
// when it decides it is time.
type Controller interface {
	// SetAlarm configures a delayed event. The name will be passed to the
	// bot's Alarm() method when the delay expires.
	SetAlarm(delay time.Duration, name string)
	// Stop stops the engine execution.
	Stop()
}
