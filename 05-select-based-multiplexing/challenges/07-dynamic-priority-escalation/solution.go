package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

const MIN_JOB_PROCESSING_TIME = 300
const MAX_JOB_PROCESSING_TIME = 1000

const MIN_JOB_PRODUCING_TIME = 100
const MAX_JOB_PRODUCING_TIME = 500

const WORKERS_NUM = 3
const JOBS_NUM = 5
const JOB_CHANNEL_BUFFER_SIZE = 1

const LOW_TO_MEDIUM_ESCALATION_MS = 800
const MEDIUM_TO_HIGH_ESCALATION_MS = 300

type Job struct {
	ID         string
	Payload    string
	EnqueuedAt time.Time
}

type MessageQueue struct {
	mu    *sync.RWMutex
	queue []*Job
}

func NewMessageQueue() *MessageQueue {
	return &MessageQueue{
		mu:    &sync.RWMutex{},
		queue: make([]*Job, 0),
	}
}

func (mq *MessageQueue) Enqueue(jobs ...*Job) {
	mq.mu.Lock()
	defer mq.mu.Unlock()

	for _, job := range jobs {
		job.EnqueuedAt = time.Now()
		mq.queue = append(mq.queue, job)
	}
}

func (mq *MessageQueue) Dequeue() (*Job, bool) {
	mq.mu.Lock()
	defer mq.mu.Unlock()

	if len(mq.queue) == 0 {
		return nil, false
	}

	job := mq.queue[0]
	mq.queue = mq.queue[1:]
	return job, true
}

func (mq *MessageQueue) Len() int {
	mq.mu.RLock()
	defer mq.mu.RUnlock()
	return len(mq.queue)
}

func (mq *MessageQueue) RemoveOlderThan(duration time.Duration) []*Job {
	toRemove := []*Job{}
	indexesToRemove := []int{}
	mq.mu.Lock()
	defer mq.mu.Unlock()

	for i, job := range mq.queue {
		// NOTE: assume that jobs are ordered by creation time
		if time.Since(job.EnqueuedAt) < duration {
			break
		}

		indexesToRemove = append(indexesToRemove, i)
		toRemove = append(toRemove, mq.queue[i])
	}

	removed := 0

	for _, index := range indexesToRemove {
		mq.queue = append(mq.queue[:index-removed], mq.queue[index-removed+1:]...)
		removed += 1
	}

	return toRemove
}

type Worker struct {
	id        int
	inputChan chan *Job
}

func NewWorker(id int, readyChan chan *Worker, wg *sync.WaitGroup) *Worker {
	c := make(chan *Job)
	w := &Worker{id, c}

	go func() {
		readyChan <- w

		for job := range c {
			fmt.Printf("[worker] worker-%d received %s\n", id, job.ID)
			time.Sleep(time.Duration(rand.Int31n(MAX_JOB_PROCESSING_TIME-MIN_JOB_PROCESSING_TIME)+MIN_JOB_PROCESSING_TIME) * time.Millisecond)
			fmt.Printf("[worker] worker-%d completed! %s\n", id, job.ID)
			wg.Done()
			fmt.Printf("[worker] worker-%d ready\n", id)
			readyChan <- w
		}
	}()

	return w
}

type Dispatcher struct {
	inputJobs        [3]<-chan *Job
	queues           [3]*MessageQueue
	readyWorkersChan <-chan *Worker
	ctx              context.Context // user-supplied context, so user can cancel dispatcher execution
	wg               *sync.WaitGroup
	doneChannel      chan struct{} // Notifies when dispatcher is done
	completedCount   int32
}

func (d *Dispatcher) Run() {
	for i, queue := range d.queues {
		go func(q *MessageQueue, index int) {
			for job := range d.inputJobs[i] {
				q.Enqueue(job)
			}
			atomic.AddInt32(&d.completedCount, 1)
		}(queue, i)
	}

	// medium to high
	go func() {
		for {
			select {
			case <-d.ctx.Done():
				return
			default:
				toMove := d.queues[1].RemoveOlderThan(MEDIUM_TO_HIGH_ESCALATION_MS * time.Millisecond)
				for _, job := range toMove {
					fmt.Printf("[escalation] Job %s promoted to HIGH\n", job.ID)
					d.queues[0].Enqueue(job)
				}
				time.Sleep(MEDIUM_TO_HIGH_ESCALATION_MS * time.Millisecond)
			}
		}
	}()

	// low to medium
	go func() {
		for {
			select {
			case <-d.ctx.Done():
				return
			default:
				toMove := d.queues[2].RemoveOlderThan(LOW_TO_MEDIUM_ESCALATION_MS * time.Millisecond)
				for _, job := range toMove {
					fmt.Printf("[escalation] Job %s promoted to MEDIUM\n", job.ID)
					d.queues[1].Enqueue(job)
				}
				time.Sleep(LOW_TO_MEDIUM_ESCALATION_MS * time.Millisecond)
			}
		}
	}()

	go func() {
		defer func() {
			d.wg.Wait()
			d.doneChannel <- struct{}{}
		}()

		for {
			var job *Job
			ok := false

			if d.queues[0].Len() > 0 {
				job, ok = d.queues[0].Dequeue()
			}

			if !ok && d.queues[1].Len() > 0 {
				job, ok = d.queues[1].Dequeue()
			}

			if !ok && d.queues[2].Len() > 0 {
				job, ok = d.queues[2].Dequeue()
			}

			if !ok && atomic.LoadInt32(&d.completedCount) == 3 {
				return
			}

			if job == nil {
				time.Sleep(10 * time.Millisecond)
				continue
			}

			select {
			case <-d.ctx.Done():
				return
			case worker := <-d.readyWorkersChan:
				d.wg.Add(1)
				worker.inputChan <- job
			}
		}
	}()
}

func (d *Dispatcher) Done() <-chan struct{} {
	return d.doneChannel
}

func NewDispatcher(highPriorityJobs, mediumPriorityJobs, lowPriorityJobs chan *Job, readyChan chan *Worker, ctx context.Context, wg *sync.WaitGroup) *Dispatcher {
	queues := [3]*MessageQueue{
		NewMessageQueue(),
		NewMessageQueue(),
		NewMessageQueue(),
	}

	inputJobs := [3]<-chan *Job{
		highPriorityJobs,
		mediumPriorityJobs,
		lowPriorityJobs,
	}

	return &Dispatcher{
		inputJobs:        inputJobs,
		queues:           queues,
		readyWorkersChan: readyChan,
		ctx:              ctx,
		wg:               wg,
		doneChannel:      make(chan struct{}),
	}
}

func produce(id string, jobChannel chan<- *Job, numJobs int) {
	for i := 0; i < numJobs; i++ {
		job := Job{
			ID:         fmt.Sprintf("%s-%d", id, i+1),
			Payload:    fmt.Sprintf("Task %s-%d", id, i+1),
			EnqueuedAt: time.Now(),
		}

		fmt.Printf("[producer] produced: %s\n", job.ID)
		jobChannel <- &job
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

	highPriorityJobs := make(chan *Job, JOB_CHANNEL_BUFFER_SIZE)
	mediumPriorityJobs := make(chan *Job, JOB_CHANNEL_BUFFER_SIZE)
	lowPriorityJobs := make(chan *Job, JOB_CHANNEL_BUFFER_SIZE)

	dispatcher := NewDispatcher(highPriorityJobs, mediumPriorityJobs, lowPriorityJobs, readyChan, ctx, wg)
	dispatcher.Run()

	go produce("🔴", highPriorityJobs, JOBS_NUM)
	go produce("🟡", mediumPriorityJobs, JOBS_NUM)
	go produce("🟢", lowPriorityJobs, JOBS_NUM)

	<-dispatcher.Done()
}
