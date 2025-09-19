package main

import (
	"fmt"
	"sync"
	"time"
)

type Broker struct {
	mu               sync.RWMutex
	topicSubscribers map[string][]chan string
	allSubscribers   []chan string
}

func NewBroker() *Broker {
	return &Broker{
		topicSubscribers: make(map[string][]chan string),
		allSubscribers:   []chan string{},
	}
}

func (b *Broker) Subscribe(topic string, bufferSize int) <-chan string {
	ch := make(chan string, bufferSize)

	b.mu.Lock() // Lock access to subscribers map while we update it
	defer b.mu.Unlock()

	if topic == "all" {
		b.allSubscribers = append(b.allSubscribers, ch)
	} else {
		b.topicSubscribers[topic] = append(b.topicSubscribers[topic], ch)
	}

	return ch
}

func (b *Broker) Publish(topic, message string) {
	b.mu.RLock() // Lock write access to subscribers map while we reading it
	defer b.mu.RUnlock()

	if subs, ok := b.topicSubscribers[topic]; ok {
		for _, sub := range subs {
			// Use select/case to drop messages instead of blocking indefinitely (a common Pub/Sub approach).
			select {
			case sub <- message:
			default:
				// Drop message if subscriber channel is full (avoid blocking)
				fmt.Printf("Warning: dropped message on topic %q (subscriber too slow)\n", topic)
			}
		}
	}

	for _, sub := range b.allSubscribers {
		select {
		case sub <- message:
		default:
			fmt.Println("Warning: dropped message for \"all\"", message)
		}
	}
}

func main() {
	broker := NewBroker()
	sportsSub := broker.Subscribe("sports", 5)
	newsSub := broker.Subscribe("news", 5)
	logger := broker.Subscribe("all", 5)

	go func() {
		for message := range sportsSub {
			fmt.Printf("Sports channel: %s\n", message)
		}
	}()

	go func() {
		for message := range newsSub {
			fmt.Printf("News channel: %s\n", message)
		}
	}()

	go func() {
		for message := range logger {
			fmt.Printf("Log channel: %s\n", message)
		}
	}()

	broker.Publish("sports", "Team A won the match!")
	broker.Publish("news", "Breaking: Go 2.0 announced!")
	broker.Publish("sports", "Bolivia classifies to FIFA World Cup 2026!")
	broker.Publish("news", "Asteroid is coming to Earth!")

	time.Sleep(2 * time.Second)
}
