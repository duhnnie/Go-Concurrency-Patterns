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

const MIN_JOB_PRODUCING_TIME = 100
const MAX_JOB_PRODUCING_TIME = 500

const WORKERS_NUM = 3
const JOBS_NUM = 5
const JOB_CHANNEL_BUFFER_SIZE = 3

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
			fmt.Printf("[worker] worker-%d received %s\n", id, job)
			time.Sleep(time.Duration(rand.Int31n(MAX_JOB_PROCESSING_TIME-MIN_JOB_PROCESSING_TIME)+MIN_JOB_PROCESSING_TIME) * time.Millisecond)
			fmt.Printf("[worker] worker-%d completed! %s\n", id, job)
			wg.Done()
			fmt.Printf("[worker] worker-%d ready\n", id)
			readyChan <- w
		}
	}()

	return w
}

type Dispatcher struct {
	highPriorityJobs   <-chan string
	mediumPriorityJobs <-chan string
	lowPriorityJobs    <-chan string
	readyWorkersChan   <-chan *Worker
	ctx                context.Context // user-supplied context, so user can cancel dispatcher execution
	wg                 *sync.WaitGroup
	doneChannel        chan struct{} // Notifies when dispatcher is done
}

func (d *Dispatcher) Run() {
	highOpen := true
	mediumOpen := true
	lowOpen := true

	go func() {
		defer func() {
			d.wg.Wait()
			d.doneChannel <- struct{}{}
		}()

		for {
			job := ""

			if highOpen {
				select {
				case <-d.ctx.Done():
					return
				case job, highOpen = <-d.highPriorityJobs:
				default:
				}
			}

			if job == "" && mediumOpen {
				select {
				case <-d.ctx.Done():
					return
				case job, mediumOpen = <-d.mediumPriorityJobs:
				default:
				}
			}

			if job == "" && lowOpen {
				select {
				case <-d.ctx.Done():
					return
				case job, lowOpen = <-d.lowPriorityJobs:
				default:
				}
			}

			if job != "" {
				select {
				case <-d.ctx.Done():
					return
				case worker := <-d.readyWorkersChan:
					d.wg.Add(1)
					fmt.Printf("[dispatcher] assigning \"%s\" to \"worker-%d\"\n", job, worker.id)
					worker.inputChan <- job
				}
			} else {
				time.Sleep(10 * time.Millisecond) // Prevent Spinning When Idle, To reduce CPU usage
			}

			if !(highOpen || mediumOpen || lowOpen) {
				return
			}
		}
	}()
}

func (d *Dispatcher) Done() <-chan struct{} {
	return d.doneChannel
}

func NewDispatcher(highPriorityJobs, mediumPriorityJobs, lowPriorityJobs chan string, readyChan chan *Worker, ctx context.Context, wg *sync.WaitGroup) *Dispatcher {
	return &Dispatcher{
		highPriorityJobs:   highPriorityJobs,
		mediumPriorityJobs: mediumPriorityJobs,
		lowPriorityJobs:    lowPriorityJobs,
		readyWorkersChan:   readyChan,
		ctx:                ctx,
		wg:                 wg,
		doneChannel:        make(chan struct{}),
	}
}

func produce(id string, jobChannel chan<- string, numJobs int) {
	for i := 0; i < numJobs; i++ {
		job := fmt.Sprintf("%s #%d", id, i+1)
		fmt.Printf("[producer] produced: %s\n", job)
		jobChannel <- job
		time.Sleep(time.Duration(rand.Intn(MAX_JOB_PRODUCING_TIME-MIN_JOB_PRODUCING_TIME)+MIN_JOB_PRODUCING_TIME) * time.Millisecond)
	}

	close(jobChannel)
}

func main() {
	ctx, _ := context.WithCancel(context.Background())
	// ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	readyChan := make(chan *Worker)
	wg := &sync.WaitGroup{}

	for i := 0; i < WORKERS_NUM; i++ {
		_ = NewWorker(i+1, readyChan, wg)
	}

	highPriorityJobs := make(chan string, JOB_CHANNEL_BUFFER_SIZE)
	mediumPriorityJobs := make(chan string, JOB_CHANNEL_BUFFER_SIZE)
	lowPriorityJobs := make(chan string, JOB_CHANNEL_BUFFER_SIZE)

	dispatcher := NewDispatcher(highPriorityJobs, mediumPriorityJobs, lowPriorityJobs, readyChan, ctx, wg)
	dispatcher.Run()

	go produce("🔴", highPriorityJobs, JOBS_NUM)
	go produce("🟡", mediumPriorityJobs, JOBS_NUM)
	go produce("🟢", lowPriorityJobs, JOBS_NUM)

	<-dispatcher.Done()
}
