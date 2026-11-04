# Resilient Worker Pool

Implement `RunWorkerPool` to process jobs from a channel using a fixed number of workers. Each job receives a context and returns an error if it fails.

The pool must:

* Run no more than `workers` jobs concurrently and reject a nonpositive worker count.
* Cancel the pool when a job returns an error, stop picking up new jobs, and return the first recorded job error.
* Wait for all currently executing jobs to finish before returning.
* Return the parent context's error if it is canceled and no job error was recorded.
* Avoid leaking its workers, and leave ownership of the jobs channel to the caller.

A worker checks for cancellation after receiving a job, but cancellation can still race with starting that job. Jobs should observe their context so they can finish promptly when canceled.

## Tags
`Concurrency`
