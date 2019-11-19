package pool_of_workers

import (
	"testing"
)

type job struct{}

func (r *job) Run() {}

func Test_Pool(t *testing.T) {
	p := NewPool(2)

	tasks := 20
	resultsCnt := tasks

	funcPoolOnResult := func(r Job) bool {
		resultsCnt--
		return resultsCnt == 0
	}

	funcPoolProcessJob := func(w *Worker) bool {
		if tasks > 0 {
			j := job{}

			w.SubmitJob(&j)
		}
		joinAll := tasks <= 0
		tasks--
		return joinAll
	}

	p.Start(funcPoolProcessJob, funcPoolOnResult)
}
