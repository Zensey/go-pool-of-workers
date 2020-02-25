# pool-of-workers

This is a simplistic pool of workers [1] implementation for Golang.
The key difference is an absence of a common job queue as there is only a queue of idle workers.

All you need to use it is:
 * a job executor (Runner), a specific code which handles a unit of work
 * a callback, providing workers with a job-units (fnJobProvider)

Code example:

    type Runner struct{}
    
    func (r *Runner) Run() {
        time.Sleep(200 * time.Millisecond)
    }
    
    func main() {
        p := NewPool(10, 20)
        tasks := 10
    
        fnOnResult   := func(r Job) {}
        fnJobProvider := func(w *Worker) {
            if tasksCnt > 0 {
                j := Runner{}
                w.SubmitJob(&j)
            }
            tasksCnt--
        }
        p.Start(fnJobProvider, fnOnResult)
    }

References
1. [Go: Worker Pool vs Pool of Workers](https://medium.com/@hau12a1/go-worker-pool-vs-pool-of-workers-b7c0598b4a67)
2. [Handling 1 Million Requests per Minute with Go](http://marcio.io/2015/07/handling-1-million-requests-per-minute-with-golang/)