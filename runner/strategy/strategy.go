package strategy

import "time"

// RequestEnqueueFunc is a function that inserts a request to the
// pipeline of requests waiting to be sent.
type RequestEnqueueFunc func()

// RunnerFunc is a function that launches strategy fanout.
type RunnerFunc func()

// RampUp is a strategy that sends more requests with each tick,
// starting from initial r/s to ceil r/s.
func RampUp(count, initial, ceil uint64, eFunc RequestEnqueueFunc) RunnerFunc {
	return func() {

		var currentSimCount = initial
		for i := uint64(0); i < count; i++ {
			for j := uint64(0); j < currentSimCount; j++ {
				eFunc()
			}
			if currentSimCount < ceil {
				currentSimCount++
			}
		}
	}
}

// ConstQPS sends requests at a constant rate per second.
func ConstQPS(total uint64, qps uint64, eFunc RequestEnqueueFunc) RunnerFunc {
	var throttle = time.Tick(time.Duration(1e6/qps) * time.Microsecond)
	return func() {
		for i := uint64(0); i < total; i++ {
			<-throttle
			eFunc()
		}
	}
}
