package pool

import (
	"testing"
	"time"
)

type job struct{}

func (r *job) Run() {
	time.Sleep(200 * time.Millisecond)
}

func Test_Pool(t *testing.T) {
	p := NewPool(2, 8)

	tasks := 10
	tasksCnt := tasks
	resCnt := 0
	submitCnt := 0

	p.FnOnJobResult = func(r Job) {
		resCnt++
		println("on result >", resCnt)
	}

	for ; tasksCnt > 0; tasksCnt-- {
		p.Submit(&job{})

		submitCnt++
		println("submit >", submitCnt)
	}
	p.Stop()

	if !(resCnt == submitCnt && resCnt == tasks) {
		t.Fail()
	}
}
