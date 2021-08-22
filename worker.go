package telegrambot

import "flag"

// read http://marcio.io/2015/07/handling-1-million-requests-per-minute-with-golang/

const defaultJobQueueSize = 100
const defaultWorkersCount = 10

// Each job is an update that should be handled by bot.handleUpdate()
type job struct {
	Update Update
	Bot    *Bot
}

var jobQueue chan job

type worker struct {
	WorkerPool chan chan job
	JobChan    chan job
	stop       chan bool
}

type dispatcher struct {
	WorkerPool chan chan job
	maxWorkers int
}

// initDispatcher Initializes the dispatcher and job queue. Number of workers and size of job queues are parsed from command-line flags.
func initDispatcher() {
	maxJobs := flag.Int("max-jobs", defaultJobQueueSize, "Max capacity for jobs (updates) queue.")
	maxWorkers := flag.Int("max-workers", defaultWorkersCount, "Max number of workers to handle updates concurrently.")
	flag.Parse()

	jobQueue = make(chan job, *maxJobs)
	dispatcher := newDispatcher(*maxWorkers)
	dispatcher.run()
}

func newDispatcher(maxWorkers int) *dispatcher {
	workerPool := make(chan chan job, maxWorkers)
	return &dispatcher{workerPool, maxWorkers}
}

func (d *dispatcher) run() {
	for i := 0; i < d.maxWorkers; i++ {
		worker := newWorker(d.WorkerPool)
		worker.Start()
	}

	go d.dispatch()
}

// dipatch waits for new jobs and assign each job to a new worker. blocks if no new job or free workers are available
func (d *dispatcher) dispatch() {
	for {
		j := <-jobQueue
		go func(j job) {
			jobChan := <-d.WorkerPool
			jobChan <- j
		}(j)
	}
}

func newWorker(workerPool chan chan job) worker {
	return worker{
		WorkerPool: workerPool,
		JobChan:    make(chan job),
		stop:       make(chan bool),
	}
}

// in Start(), Worker waits for new jobs in a new goroutine and calls bot.handleUpdate for each job.
// The goroutine stops (calls return) when worker.stop channel receives "true"
func (w worker) Start() {
	go func() {
		for {
			w.WorkerPool <- w.JobChan
			select {
			case j := <-w.JobChan:
				j.Bot.handleUpdate(j.Update)
			case <-w.stop:
				return
			}
		}
	}()
}

func (w worker) Stop() {
	go func() {
		w.stop <- true
	}()
}
