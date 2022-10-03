package pool

import (
	"log"
	"os"
	"sync"
)

type CodeOutput struct {
	User        User
	TempDirName string
	Result      string
	Error       string
}

type WorkerPool struct {
	queue        *ConcurrentQueue[WorkType]
	nWorkers     int
	queueSize    int
	workingGroup *sync.WaitGroup
}

type WorkType struct {
	lang   string
	code   string
	action func(lang, code string, ch chan<- CodeOutput)
	ch     chan<- CodeOutput
}

func NewWorkerPool(nWorkers, queueSize int) *WorkerPool {
	var wg sync.WaitGroup

	queue := NewConcurrentQueue[WorkType](uint32(queueSize))

	for idx := 0; idx < nWorkers; idx++ {
		wg.Add(1)
		go poolWorker(&wg, queue, idx)
	}

	return &WorkerPool{
		queue:        queue,
		workingGroup: &wg,
		nWorkers:     nWorkers,
		queueSize:    queueSize,
	}
}

func (wp *WorkerPool) SubmitJob(lang, code string, action func(lang, code string, ch chan<- CodeOutput), ch chan<- CodeOutput) {
	work := WorkType{
		lang:   lang,
		code:   code,
		ch:     ch,
		action: action,
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
		work.action(work.lang, work.code, work.ch)
	}
}
