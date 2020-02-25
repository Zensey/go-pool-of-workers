package pool_of_workers

import (
	"github.com/cznic/mathutil"
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
	return p
}

func (p *Pool) incPendingJobs() {
	p.pendingJobs++
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

func (p *Pool) Start(funcProcessJob FuncProcessJob, funcOnJobResult FuncOnJobResult) {
	p.minWorkers = mathutil.Max(1, p.minWorkers)
	p.minWorkers = mathutil.Min(p.minWorkers, p.maxWorkers)

	for i := 0; i < p.minWorkers; i++ {
		p.spawnWorker()
	}

	for {
		select {
		case res := <-p.results:
			p.pendingJobs--

			funcOnJobResult(res)

			// if has more tasks and not maxWorkers exceeded, spawn 1 another worker
			if len(p.workers) < p.maxWorkers {
				p.spawnWorker()
			}

		case w := <-p.idleWorkers:
			funcProcessJob(w)
			if p.pendingJobs == 0 {
				p.joinWorkers()
				return
			}
		}
	}
}
