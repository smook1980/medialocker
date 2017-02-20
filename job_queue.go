package medialocker

// See https://github.com/bamzi/jobrunner
// https://godoc.org/github.com/qor/worker#QorJobInterface
import (
	"sync"
	"sync/atomic"

	"github.com/labstack/gommon/log"
	"context"
)

type Worker interface {
	Work(*JobQueue) error
}

// JobQueue
type JobQueue struct {
	context   context.Context
	wg        sync.WaitGroup
	input     chan Worker
	output    chan Worker
	processed int32
	errors    int32
	Results   <-chan Worker
}

// Queue Worker job
func (jq *JobQueue) Queue(work ...Worker) {
	for _, w := range work {
		jq.input <- w
	}

	return
}

func NewJobQueue(context context.Context, workers int, buffer int) *JobQueue {
	bufferSize := buffer + workers
	input := make(chan Worker, bufferSize)
	output := make(chan Worker, bufferSize)
	jq := &JobQueue{context: context, input: input, output: output}

	for x := 0; x < workers; x++ {
		jq.wg.Add(1)
		go func() {
			defer jq.wg.Done()

			select {
			case w := <-input:
				if err := w.Work(jq); err != nil {
					log.Error(err)
					atomic.AddInt32(&jq.errors, 1)
				}

				atomic.AddInt32(&jq.processed, 1)
			case <-jq.context.Done():
				break
			}
		}()
	}

	go func() {
		jq.wg.Wait()
		close(output)
	}()

	return jq
}
