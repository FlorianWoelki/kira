package pool

import (
	"log"
	"os"
	"sync"
	"time"
)

type CodeOutput struct {
	User          User
	TempDirName   string
	CompileResult string
	CompileError  string
	CompileTime   time.Duration
	RunResult     string
	RunError      string
	RunTime       time.Duration
}

type WorkerPool struct {
	queue        *ConcurrentQueue[WorkType]
	nWorkers     int
	workingGroup *sync.WaitGroup
}

type WorkType struct {
	lang        string
	code        string
	test        string
	bypassCache bool
	action      func(lang, code string, tests string, bypassCache bool, ch chan<- CodeOutput)
	ch          chan<- CodeOutput
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

func (wp *WorkerPool) SubmitJob(lang, code string, test string, bypassCache bool, action func(lang, code string, tests string, bypassCache bool, ch chan<- CodeOutput), ch chan<- CodeOutput) {
	work := WorkType{
		lang:        lang,
		code:        code,
		test:        test,
		ch:          ch,
		bypassCache: bypassCache,
		action:      action,
	}
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
		work.action(work.lang, work.code, work.test, work.bypassCache, work.ch)
	}
}
