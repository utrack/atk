package runner

import "time"

type Worker struct {
	reqChan <-chan struct{}
	rspChan chan<- Result

	requester func() (int, error)
}

func SpawnWorkerPool(reqChan <-chan struct{}, requester func() (int, error), count int) <-chan Result {
	ret := make(chan Result, count)
	for i := 0; i < count; i++ {
		w := NewWorker(reqChan, ret, requester)
		w.Start()
	}

	return ret
}

func NewWorker(reqChan <-chan struct{}, rspChan chan<- Result, requester func() (int, error)) *Worker {
	ret := &Worker{
		reqChan:   reqChan,
		rspChan:   rspChan,
		requester: requester,
	}
	return ret
}

func (w *Worker) Start() {
	go w.runProc()
}

func (w *Worker) runProc() {
	for range w.reqChan {
		started := time.Now()
		code, err := w.requester()
		w.rspChan <- Result{
			Err:     err,
			Code:    code,
			Elapsed: time.Since(started),
		}
	}
}
