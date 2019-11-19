package pool_of_workers

import (
	"sync"
)

type Job interface {
	Run()
}

type Pool struct {
	maxWorkers  int
	workers     []*Worker
	wg          sync.WaitGroup
	quit        chan bool
	results     chan Job
	idleWorkers chan *Worker
	pendingJobs int
}

type FuncProcessJob func(*Worker) bool
type FuncOnJobResult func(Job) bool

func NewPool(maxWorkers int) *Pool {
	p := &Pool{
		maxWorkers:  maxWorkers,
		workers:     make([]*Worker, 0, maxWorkers),
		quit:        make(chan bool),
		results:     make(chan Job),
		idleWorkers: make(chan *Worker, 100),
		pendingJobs: 0,
	}
	return p
}

func (p *Pool) spawnWorker() {
	p.wg.Add(1)
	w := newWorker(p)
	p.workers = append(p.workers, w)

	id := len(p.workers)
	go w.Run(id)
}

func (p *Pool) joinWorkers() {
	close(p.quit)
	p.wg.Wait()
}

func (p *Pool) Start(funcProcessJob FuncProcessJob, funcOnJobResult FuncOnJobResult) {
	p.spawnWorker()

	for {
		select {
		case res := <-p.results:
			p.pendingJobs--

			if funcOnJobResult(res) && p.pendingJobs == 0 {
				p.joinWorkers()
				return
			}

			// if has more tasks and not maxWorkers exceeded, spawn another 1 worker
			if len(p.workers) < p.maxWorkers {
				p.spawnWorker()
			}

		case w := <-p.idleWorkers:
			if funcProcessJob(w) && p.pendingJobs == 0 {
				p.joinWorkers()
				return
			}
		}
	}
}
