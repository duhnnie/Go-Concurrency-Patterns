package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const NUM_MESSAGES = 5

type TaggedMessage struct {
	Source string
	Value  string
}

type MessageMerger struct {
	inputChannels map[string]<-chan string
	wg            *sync.WaitGroup
	Output        chan TaggedMessage
	mu            *sync.RWMutex
}

func NewMessageMerger(channels map[string]chan string) *MessageMerger {
	messageMerger := &MessageMerger{
		inputChannels: make(map[string]<-chan string),
		wg:            &sync.WaitGroup{},
		Output:        make(chan TaggedMessage, 10),
		mu:            &sync.RWMutex{},
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

		for value := range channel {
			mm.mu.RLock()
			_, ok := mm.inputChannels[tag]
			mm.mu.RUnlock()

			if !ok {
				break
			}

			mm.Output <- TaggedMessage{
				Source: tag,
				Value:  value,
			}
		}
	}()
}

func (mm *MessageMerger) RemoveChannel(tag string) {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	delete(mm.inputChannels, tag)
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

	go func() {
		defer close(sportsChan)

		for i := 0; i < NUM_MESSAGES; i++ {
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
			sportsChan <- fmt.Sprintf("sports message #%d", i+1)
		}
	}()

	go func() {
		defer close(techChan)

		for i := 0; i < NUM_MESSAGES; i++ {
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
			techChan <- fmt.Sprintf("tech message #%d", i+1)
		}
	}()

	go func() {
		defer close(newsChan)

		for i := 0; i < NUM_MESSAGES; i++ {
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
			newsChan <- fmt.Sprintf("news message #%d", i+1)
		}
	}()

	go func() {
		defer close(gamingChan)

		for i := 0; i < NUM_MESSAGES; i++ {
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
			gamingChan <- fmt.Sprintf("gaming message #%d", i+1)
		}
	}()

	printedMessages := 0

	for taggedMessage := range messageMerger.Output {
		fmt.Printf("[%s] %s\n", taggedMessage.Source, taggedMessage.Value)
		printedMessages++

		if printedMessages == 5 {
			messageMerger.RemoveChannel("news")
			fmt.Println("News channel removed!")

			fmt.Println("gaming channel added")
			messageMerger.AddChannel("gaming", gamingChan)
		}
	}

}
