# pool-of-workers

This is simplistic pool of workers implementation for Golang.
All you need to do is:
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
