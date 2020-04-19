package pool

import (
	"sync"
)

type Job interface {
	Run()
}

type Pool struct {
	minWorkers  int
	maxWorkers  int
	workers     []*Worker
	wg          sync.WaitGroup
	quit        chan bool
	results     chan Job
	idleWorkers chan *Worker
	pendingJobs int

	FnOnJobResult FuncOnJobResult
}

type FuncProcessJob func(*Worker)
type FuncOnJobResult func(Job)

func NewPool(minWorkers, maxWorkers int) *Pool {
	p := &Pool{
		minWorkers:  minWorkers,
		maxWorkers:  maxWorkers,
		workers:     make([]*Worker, 0, maxWorkers),
		quit:        make(chan bool),
		results:     make(chan Job),
		idleWorkers: make(chan *Worker, 100),
		pendingJobs: 0,
	}
	p.minWorkers = max(1, p.minWorkers)
	p.minWorkers = min(p.minWorkers, p.maxWorkers)

	return p
}

func (p *Pool) spawnWorker() {
	p.wg.Add(1)
	w := newWorker(p)
	p.workers = append(p.workers, w)

	id := len(p.workers)
	go w.run(id)
}

func (p *Pool) joinWorkers() {
	close(p.quit)
	p.wg.Wait()
}

func (p *Pool) Submit(j Job) {
	if j != nil {
		for len(p.workers) < p.minWorkers {
			p.spawnWorker()
		}

		// if no spare workers and not maxWorkers exceeded, spawn 1 another worker
		if p.pendingJobs == len(p.workers) && len(p.workers) < p.maxWorkers {
			p.spawnWorker()
		}
	}

	for {
		select {
		case res := <-p.results:
			p.pendingJobs--
			if p.FnOnJobResult != nil {
				p.FnOnJobResult(res)
			}

		case w := <-p.idleWorkers:
			if j != nil {
				p.pendingJobs++
				w.SubmitJob(j)
				return
			}
			if p.pendingJobs == 0 {
				p.joinWorkers()
				return
			}
		}
	}
}

func (p *Pool) Stop() {
	p.Submit(nil)
}
