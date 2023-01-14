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
	Name     string `json:"name" binding:"required"`
	Received string `json:"received" binding:"required"`
	Actual   string `json:"actual" binding:"required"`
	Passed   bool   `json:"passed" binding:"required"`
	RunError string `json:"runError" binding:"required"`
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
	Stdin       string
	Tests       []TestResult
	BypassCache bool
}

type actionFunc = func(data WorkData, ch chan<- CodeOutput)

type WorkType struct {
	data   WorkData
	action actionFunc
	ch     chan<- CodeOutput
}

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

func (wp *WorkerPool) SubmitJob(data WorkData, action actionFunc, ch chan<- CodeOutput) {
	work := WorkType{data, action, ch}
	wp.queue.enqueue(work)
}

func poolWorker[T any](wg *sync.WaitGroup, queue *ConcurrentQueue[T], idx int) {
	defer wg.Done()

	for {
		val, err := queue.dequeue()
		if err != nil {
			log.Fatal(err)
			os.Exit(0)
		}

		work := val.(WorkType)
		work.action(work.data, work.ch)
	}
}
