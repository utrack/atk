package runner

import "time"

type Result struct {
	Err     error
	Code    int
	At      time.Time
	Elapsed time.Duration
}
