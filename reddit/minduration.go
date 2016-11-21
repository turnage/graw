package reddit

import (
	"time"
)

// lol no generics
func maxOf(a, b time.Duration) time.Duration {
	// lol no ternary statements
	if a > b {
		return a
	}

	return b
}
