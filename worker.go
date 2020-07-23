package pool

type Worker struct {
	pool    *Pool
	jobChan chan Job
}

func newWorker(p *Pool) *Worker {
	return &Worker{
		pool:    p,
		jobChan: make(chan Job),
	}
}

// Used by user-s code to submit a task to a worker
func (w *Worker) SubmitJob(j Job) {
	w.jobChan <- j
}

// Used by pool to spawn a worker
func (w *Worker) run(id int) {
	for {
		select {
		case w.pool.idleWorkers <- w:
			println("w> idle")

		case <-w.pool.quit:
			println("w> quit")
			w.pool.wg.Done()
			return
		}

		select {
		case job := <-w.jobChan:
			println("w> j", job)
			job.Run()
			w.pool.results <- job

		case <-w.pool.quit:
			println("w> quit 2")

			w.pool.wg.Done()
			return
		}
	}
}
