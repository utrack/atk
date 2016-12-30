package runner

import "time"

type Result struct {
	Err     error
	Code    int
	Elapsed time.Duration
}
