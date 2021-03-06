# go-pool-of-workers

This is a simplistic Go implementation of the pool of workers [1].

###### When you might need it:
If you want to handle a number/stream of jobs in parallel but with a limited number of goroutines.

All you need is:
 * a job Runner -- a handler for a specific unit of work
 * an optional callback to handle a result of job 

Code example:

    type Job struct{ result SomeType }
    
    func (r *Job) Run() {
        time.Sleep(200 * time.Millisecond)
        r.result = x
    }
    
    func main() {
        p := pool.NewPool(2, 4) // minWorkers, maxWorkers
        p.FnOnJobResult = func(j Job) {}

        tasks := 10
        for tasksCnt > 0; tasksCnt-- {
            p.Submit(&Job{})
        }
        p.Stop()
    }


###### Design
The key difference from canonical pool [2] is an absence of a common job queue.
Instead, there is a common queue of idle workers [1][2]. 


###### Features / properties

* The producer (submitting a job) is unlocked as soon as a idle worker consumes the job, thus potentially reducing the time producer is blocked.
* Number of workers increases on demand from *minWorkers* to *maxWorkers*  
* You may have as many types of jobs as you like


###### References
1. [Go: Worker Pool vs Pool of Workers](https://medium.com/@hau12a1/go-worker-pool-vs-pool-of-workers-b7c0598b4a67)
2. [Go by Example: Worker Pools](https://gobyexample.com/worker-pools)
3. [Handling 1 Million Requests per Minute with Go](http://marcio.io/2015/07/handling-1-million-requests-per-minute-with-golang/)
