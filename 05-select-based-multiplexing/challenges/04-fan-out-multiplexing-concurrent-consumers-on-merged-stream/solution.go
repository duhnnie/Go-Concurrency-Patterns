package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const NUM_MESSAGES = 5
const NUM_WORKERS = 5

type TaggedMessage struct {
	Source string
	Value  string
}

type MessageMerger struct {
	inputChannels map[string]<-chan string
	wg            *sync.WaitGroup
	Output        chan TaggedMessage
	mu            *sync.RWMutex
	ctx           context.Context
	cancel        context.CancelFunc
}

func NewMessageMerger(channels map[string]chan string) *MessageMerger {
	ctx, cancel := context.WithCancel(context.Background())

	messageMerger := &MessageMerger{
		inputChannels: make(map[string]<-chan string),
		wg:            &sync.WaitGroup{},
		Output:        make(chan TaggedMessage, 10),
		mu:            &sync.RWMutex{},
		ctx:           ctx,
		cancel:        cancel,
	}

	for tag, channel := range channels {
		messageMerger.AddChannel(tag, channel)
	}

	go func() {
		messageMerger.wg.Wait()
		close(messageMerger.Output)
	}()

	return messageMerger
}

func (mm *MessageMerger) AddChannel(tag string, channel chan string) {
	if _, ok := mm.inputChannels[tag]; ok {
		panic(fmt.Sprintf("tag \"%s\" is already being used", tag))
	}

	mm.mu.Lock()
	mm.inputChannels[tag] = channel
	mm.wg.Add(1)
	mm.mu.Unlock()

	go func() {
		defer mm.wg.Done()

		for {
			select {
			case <-mm.ctx.Done():
				return
			case value, channelOpen := <-channel:
				if !channelOpen {
					return
				}

				mm.mu.RLock()
				_, stillActive := mm.inputChannels[tag]
				mm.mu.RUnlock()

				if !stillActive {
					return
				}

				mm.Output <- TaggedMessage{
					Source: tag,
					Value:  value,
				}
			}
		}
	}()
}

func (mm *MessageMerger) RemoveChannel(tag string) {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	delete(mm.inputChannels, tag)
}

func (mm *MessageMerger) Close() {
	mm.cancel()
}

func worker(number int, input <-chan TaggedMessage, wg *sync.WaitGroup) {
	defer wg.Done()

	for message := range input {
		fmt.Printf("[Worker #%d] processed: [%s] %s\n", number, message.Source, message.Value)
		time.Sleep(time.Duration(rand.Intn(800)) * time.Millisecond)
	}
}

func startWorkers(input <-chan TaggedMessage, size int) *sync.WaitGroup {
	wg := &sync.WaitGroup{}
	wg.Add(size)

	for i := 0; i < size; i++ {
		go worker(i+1, input, wg)
	}

	return wg
}

func produce(tag string, channel chan<- string, numberMessages int) {
	defer close(channel)

	for i := 0; i < numberMessages; i++ {
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		channel <- fmt.Sprintf("%s message #%d", tag, i+1)
	}
}

func main() {
	sportsChan := make(chan string)
	techChan := make(chan string)
	newsChan := make(chan string)
	gamingChan := make(chan string)

	messageMerger := NewMessageMerger(map[string]chan string{
		"sports": sportsChan,
		"tech":   techChan,
		"news":   newsChan,
	})

	go produce("sports", sportsChan, NUM_MESSAGES)
	go produce("tech", techChan, NUM_MESSAGES)
	go produce("news", newsChan, NUM_MESSAGES)
	go produce("gaming", gamingChan, NUM_MESSAGES)

	workersWG := startWorkers(messageMerger.Output, NUM_WORKERS)
	messageMerger.AddChannel("gaming", gamingChan)
	workersWG.Wait()
	fmt.Println("Done")
}
