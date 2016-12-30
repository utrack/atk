package reports

import (
	"io"

	"github.com/tsenart/vegeta/lib"
	"github.com/utrack/atk/runner"
)

type Reporter interface {
	Output(runner.Result)
}

type VegetaReporter struct {
	enc vegeta.Encoder
}

// NewVegetaOutput returns encoder that encodes result data
// to tsenart/vegeta-compatible format.
func NewVegetaOutput(out io.Writer) *VegetaReporter {
	enc := vegeta.NewEncoder(out)
	return &VegetaReporter{
		enc: enc,
	}
}

func ReportFromChan(rep Reporter, resultCount uint64, ch <-chan runner.Result) {
	results := uint64(0)
	for res := range ch {
		rep.Output(res)
		results++
		if results == resultCount {
			return
		}
	}
}

func (o *VegetaReporter) Output(r runner.Result) {
	res := &vegeta.Result{
		Code:      uint16(r.Code),
		Timestamp: r.At,
		Latency:   r.Elapsed,
	}
	if r.Err != nil {
		res.Error = r.Err.Error()
	}
	o.enc.Encode(res)
}
