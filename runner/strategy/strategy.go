package strategy

import (
	"time"
)

// RequestEnqueuer is a chan that signals when a request should be sent to the
// worker pool.
//
// Chan is given out by the Strategy constructor and is closed after all requests
// had been processed.
type RequestEnqueuer <-chan struct{}

// RunnerFunc is a function that launches strategy fanout.
type RunnerFunc func()

// RampUp is a strategy that sends more requests with each tick,
// starting from initial r/s to ceil r/s.
func RampUp(count, initial, ceil uint64) (RunnerFunc, RequestEnqueuer) {
	retChan := make(chan struct{}, ceil)
	return func() {
		defer close(retChan)

		var currentSimCount = initial
		for i := uint64(0); i < count; i++ {
			for j := uint64(0); j < currentSimCount; j++ {
				retChan <- struct{}{}
			}
			if currentSimCount < ceil {
				currentSimCount++
			}
		}
	}, retChan
}

// ConstQPS sends requests at a constant rate per second.
func ConstQPS(total uint64, qps uint64) (RunnerFunc, RequestEnqueuer) {
	throttle := time.Tick(time.Duration(1e6/qps) * time.Microsecond)
	retChan := make(chan struct{}, qps*2)
	return func() {
		defer close(retChan)
		for i := uint64(0); i < total; i++ {
			<-throttle
			retChan <- struct{}{}
		}
	}, retChan
}
