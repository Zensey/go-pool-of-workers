package pool_of_workers

import (
	"testing"
	"time"
)

type job struct{}

func (r *job) Run() {
	time.Sleep(200 * time.Millisecond)
}

func Test_Pool(t *testing.T) {
	p := NewPool(10, 20)

	tasks := 10

	tasksCnt := tasks
	resCnt := 0
	submitCnt := 0

	funcPoolOnResult := func(r Job) {
		resCnt++
		println("on result >", resCnt)
	}

	funcPoolProcessJob := func(w *Worker) {
		if tasksCnt > 0 {
			j := job{}
			w.SubmitJob(&j)

			submitCnt++
			println("submit >", submitCnt)
		}
		tasksCnt--
	}
	p.Start(funcPoolProcessJob, funcPoolOnResult)

	if !(resCnt == submitCnt && resCnt == tasks) {
		t.Fail()
	}
}
