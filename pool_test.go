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
	results := 0

	p.FnOnJobResult = func(r Job) {
		results++
		println("on result >", results)
	}

	for i := 0; i < tasks; i++ {
		p.Submit(&job{})
	}
	p.Stop()

	if results != tasks {
		t.Fail()
	}
}

func Test_PoolEmpty(t *testing.T) {
	p := NewPool(2, 8)

	p.FnOnJobResult = func(r Job) {
		println("on result >")
	}
	p.Stop()
}
