package main

import (
	"fmt"
	"sync"
	"time"
)

type Broker struct {
	mu               sync.RWMutex
	topicSubscribers map[string][]chan<- string
	allSubscribers   []chan<- string
}

func NewBroker() *Broker {
	return &Broker{
		topicSubscribers: make(map[string][]chan<- string),
		allSubscribers:   []chan<- string{},
	}
}

func (b *Broker) Subscribe(topic string) <-chan string {
	subscriberChan := make(chan string)

	b.mu.Lock() // Lock access to subscribers maps for cases of concurrent access
	defer b.mu.Unlock()

	if topic == "all" {
		b.allSubscribers = append(b.allSubscribers, subscriberChan)
	} else {
		b.topicSubscribers[topic] = append(b.topicSubscribers[topic], subscriberChan)
	}

	return subscriberChan
}

func (b *Broker) Publish(topic, message string) {
	b.mu.RLock() // Lock write access while reading all subscribers for deliver messages
	defer b.mu.RUnlock()

	if subscribers, ok := b.topicSubscribers[topic]; ok {
		for _, subscriber := range subscribers {
			select {
			case subscriber <- message:
			default:
				// Drops message if subscriber channel is full (avoid blocking)
				fmt.Println("Warning: dropped message for topic", topic)
			}

		}
	}

	for _, allSubscriber := range b.allSubscribers {
		select {
		case allSubscriber <- message:
		default:
			fmt.Println("Warning: dropped message for subscribers to \"all\"")
		}

	}
}

func main() {
	broker := NewBroker()
	sportsSub := broker.Subscribe("sports")
	_ = broker.Subscribe("news")
	logger := broker.Subscribe("all")

	go func() {
		for message := range sportsSub {
			fmt.Printf("Sports channel: %s\n", message)
		}
	}()

	// go func() {
	// 	for message := range newsSub {
	// 		fmt.Printf("News channel: %s\n", message)
	// 	}
	// }()

	go func() {
		for message := range logger {
			fmt.Printf("Log channel: %s\n", message)
		}
	}()

	broker.Publish("sports", "Team A won the match!")
	broker.Publish("news", "Breaking: Go 2.0 announced!")
	// broker.Publish("news", "Evo Morales goes to jail!")

	time.Sleep(10 * time.Second)
}
