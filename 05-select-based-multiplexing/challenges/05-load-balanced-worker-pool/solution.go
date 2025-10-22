package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const MIN_JOB_PROCESSING_TIME = 300
const MAX_JOB_PROCESSING_TIME = 1000
const WORKERS_NUM = 3
const JOBS_NUM = 15

type Worker struct {
	id        int
	inputChan chan string
}

func NewWorker(id int, readyChan chan *Worker, wg *sync.WaitGroup) *Worker {
	c := make(chan string)
	w := &Worker{id, c}

	go func() {
		readyChan <- w

		for job := range c {
			fmt.Printf("worker-%d received %s\n", id, job)
			time.Sleep(time.Duration(rand.Int31n(MAX_JOB_PROCESSING_TIME-MIN_JOB_PROCESSING_TIME)+MIN_JOB_PROCESSING_TIME) * time.Millisecond)
			wg.Done()
			readyChan <- w
		}
	}()

	return w
}

type Dispatcher struct {
	jobsChan         <-chan string
	readyWorkersChan <-chan *Worker
	ctx              context.Context
}

func NewDispatcher(jobsChan chan string, readyChan chan *Worker, ctx context.Context) *Dispatcher {
	go func() {
		job := ""
		open := true

		for {
			select {
			case <-ctx.Done():
				return
			case job, open = <-jobsChan:
				if !open {
					return
				}

				select {
				case <-ctx.Done():
					return
				case worker := <-readyChan:
					fmt.Printf("worker-%d ready\n", worker.id)
					worker.inputChan <- job
				}
			}
		}
	}()

	return &Dispatcher{
		jobsChan:         jobsChan,
		readyWorkersChan: readyChan,
		ctx:              ctx,
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	readyChan := make(chan *Worker)
	wg := &sync.WaitGroup{}

	for i := 0; i < WORKERS_NUM; i++ {
		_ = NewWorker(i+1, readyChan, wg)
	}

	jobsChan := make(chan string)
	_ = NewDispatcher(jobsChan, readyChan, ctx)
	wg.Add(JOBS_NUM)

	go func() {
		for i := 0; i < JOBS_NUM; i++ {
			jobsChan <- fmt.Sprintf("job %d", i+1)
		}
		close(jobsChan)
	}()

	wg.Wait()
	cancel() // Not necessary since dispatcher also stops at closing jobsChan channel
}
