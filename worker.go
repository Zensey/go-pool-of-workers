package pool_of_workers

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

func (w *Worker) SubmitJob(j Job) {
	w.pool.incPendingJobs()
	w.jobChan <- j
}

func (w *Worker) Run(id int) {
	for {
		w.pool.idleWorkers <- w

		select {
		case job := <-w.jobChan:
			job.Run()
			w.pool.results <- job // Pool should read all results. For this we count submitted jobs

		case <-w.pool.quit:
			w.pool.wg.Done()

			return
		}
	}
}