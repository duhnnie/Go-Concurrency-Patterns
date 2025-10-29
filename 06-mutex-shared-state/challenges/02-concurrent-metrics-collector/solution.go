package main

import (
	"fmt"
	"math/rand/v2"
	"sync"
	"time"
)

type Metrics struct {
	Requests     int
	Errors       int
	TotalLatency time.Duration
	sync.RWMutex
}

func (m *Metrics) RecordRequest(latency time.Duration) {
	m.Lock()
	defer m.Unlock()
	m.Requests++
	m.TotalLatency += latency
}

func (m *Metrics) RecordError() {
	m.Lock()
	defer m.Unlock()
	m.Errors++
}

func (m *Metrics) AverageLatency() time.Duration {
	m.RLock()
	defer m.RUnlock()
	if m.Requests == 0 {
		return 0
	}
	return m.TotalLatency / time.Duration(m.Requests)
}

const (
	WORKERS_NUM        = 10
	WORKER_REQUEST_MIN = 10
	WORKER_REQUEST_MAX = 30
	WORKER_SLEEP_MIN   = 100 * time.Millisecond
	WORKER_SLEEP_MAX   = 500 * time.Millisecond
	ERROR_RATE         = 0.15
)

func worker(m *Metrics, wg *sync.WaitGroup) {
	defer wg.Done()
	reqCount := rand.IntN(WORKER_REQUEST_MAX-WORKER_REQUEST_MIN) + WORKER_REQUEST_MIN

	for i := 0; i < reqCount; i++ {
		start := time.Now()
		time.Sleep(WORKER_SLEEP_MIN + time.Duration(rand.IntN(int(WORKER_SLEEP_MAX-WORKER_SLEEP_MIN))))
		if rand.Float64() < ERROR_RATE {
			m.RecordError()
		} else {
			m.RecordRequest(time.Since(start))
		}
	}
}

func main() {
	metrics := &Metrics{}
	wg := &sync.WaitGroup{}
	wg.Add(WORKERS_NUM)

	for i := 0; i < WORKERS_NUM; i++ {
		go worker(metrics, wg)
	}

	wg.Wait()
	fmt.Printf("Total requests: %d\n", metrics.Requests)
	fmt.Printf("Total errors: %d\n", metrics.Errors)
	fmt.Printf("Average latency: %dms\n", metrics.AverageLatency()/time.Millisecond)
}
