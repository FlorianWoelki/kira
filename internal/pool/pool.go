package pool

import (
	"log"
	"os"
	"sync"
)

type Output struct {
	Result string `json:"result"`
	Error  string `json:"error"`
	Time   int64  `json:"time"`
}

type TestResult struct {
	Name     string   `json:"name"`
	Received string   `json:"received"`
	Actual   string   `json:"actual"`
	Stdin    []string `json:"stdin"`
	Passed   bool     `json:"passed"`
	RunError string   `json:"runError"`
}

type TestOutput struct {
	Results []TestResult `json:"results" binding:"required"`
	Time    int64        `json:"time" binding:"required"`
}

type CodeOutput struct {
	User          User
	TempDirName   string
	CompileOutput Output
	RunOutput     Output
	TestOutput    TestOutput
}

type WorkerPool struct {
	queue        *ConcurrentQueue[WorkType]
	nWorkers     int
	workingGroup *sync.WaitGroup
}

type WorkData struct {
	Lang        string
	Code        string
	Stdin       []string
	Tests       []TestResult
	BypassCache bool
}

type ActionOutput struct {
	Once   chan CodeOutput
	Stream chan string
}

type actionFunc = func(data WorkData, output ActionOutput, terminate chan<- bool)

// WorkType represents a unit of work to be executed by the worker pool.
type WorkType struct {
	// data represents the input data for the work unit.
	data WorkData
	// action is the function that performs the actual work.
	action actionFunc
	// actionOutput represents the output of the work unit.
	actionOutput ActionOutput
	// terminate is a channel used to signal the worker to terminate.
	terminate chan<- bool
}

// NewWorkerPool creates a new worker pool instance with the specific number of worker
// goroutines.
func NewWorkerPool(nWorkers int) *WorkerPool {
	var wg sync.WaitGroup

	queue := NewConcurrentQueue[WorkType]()

	for idx := 0; idx < nWorkers; idx++ {
		wg.Add(1)
		go poolWorker(&wg, queue, idx)
	}

	return &WorkerPool{
		queue:        queue,
		workingGroup: &wg,
		nWorkers:     nWorkers,
	}
}

// SubmitJob adds a new work unit to the worker pool.
func (wp *WorkerPool) SubmitJob(data WorkData, action actionFunc, actionOutput ActionOutput, terminate chan<- bool) {
	work := WorkType{data, action, actionOutput, terminate}
	wp.queue.enqueue(work)
}

// poolWorker is a worker goroutine that continually dequeues work units from the
// provided queue and executes them.
func poolWorker[T any](wg *sync.WaitGroup, queue *ConcurrentQueue[T], idx int) {
	defer wg.Done()

	for {
		val, err := queue.dequeue()
		if err != nil {
			log.Fatal(err)
			os.Exit(0)
		}

		// Execute the actual work.
		work := val.(WorkType)
		work.action(work.data, work.actionOutput, work.terminate)
	}
}
