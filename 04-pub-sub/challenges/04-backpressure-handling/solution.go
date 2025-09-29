package main

// Implement the graceful cancellation with context.Context

import (
	"fmt"
	"sync"
	"time"
)

type Strategy string

const (
	BLOCK = Strategy("block")
	DROP  = Strategy("drop")
	EVICT = Strategy("evict")
)

type Subscriber struct {
	strategy Strategy
	c        chan string
}

func NewSubscriber(strategy Strategy, bufferSize int) *Subscriber {
	return &Subscriber{
		strategy: strategy,
		c:        make(chan string, bufferSize),
	}
}

type Broker struct {
	subscribers map[string][]*Subscriber
	mu          *sync.RWMutex
}

func NewBroker() *Broker {
	return &Broker{
		make(map[string][]*Subscriber),
		&sync.RWMutex{},
	}
}

func (b *Broker) Subscribe(topic string, bufferSize int, strategy Strategy) <-chan string {
	b.mu.Lock()
	defer b.mu.Unlock()

	sub := NewSubscriber(strategy, bufferSize)

	if _, ok := b.subscribers[topic]; !ok {
		b.subscribers[topic] = []*Subscriber{}
	}

	b.subscribers[topic] = append(b.subscribers[topic], sub)
	return sub.c
}

func (b *Broker) deliver(sub *Subscriber, i int, message string) {
	fmt.Printf("Publishing message to subscriber %d: %s\n", i, message)

	select {
	case sub.c <- message:
	default:
		switch sub.strategy {
		case BLOCK:
			fmt.Printf("Blocking until subscriber %d has space: %s\n", i, message)
			sub.c <- message
			fmt.Printf("Message delivered to subscriber %d after blocking: %s\n", i, message)
		case DROP:
			fmt.Printf("Dropping newest message for subscriber %d: %s\n", i, message)
		case EVICT:
			select {
			case <-sub.c:
				fmt.Printf("Evicting oldest message for subscriber %d\n", i)
			default:
				fmt.Printf("Buffer empty at eviction for subscriber %d (unexpected)\n", i)
			}
			sub.c <- message
		}
	}
}

func (b *Broker) Publish(topic, message string) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if subs, ok := b.subscribers[topic]; ok {
		for i, sub := range subs {
			b.deliver(sub, i, message)
		}
	}
}

func (b *Broker) Unsubscribe(subChan chan string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	for topic, subscribers := range b.subscribers {
		for i, sub := range subscribers {
			if sub.c == subChan {
				b.subscribers[topic] = append(b.subscribers[topic][0:i], b.subscribers[topic][i+1:]...)
				close(sub.c)
				return
			}
		}
	}
}

func main() {
	messages := map[string][]string{
		"sports": {
			"Team A won the match!",
			"Team B lost the match!",
			"Team C won the match!",
			"Team D lost the match!",
			"Team E won the match!",
			"Team F lost the match!",
			"Team G won the match!",
			"Team H lost the match!",
			"Team I won the match!",
			"Team J lost the match!",
			"Team K won the match!",
			"Team L lost the match!",
			"Team M won the match!",
			"Team N lost the match!",
			"Team O won the match!",
			"Team P lost the match!",
			"Team Q won the match!",
			"Team R lost the match!",
		},
	}

	broker := NewBroker()

	blockingSub := broker.Subscribe("sports", 3, BLOCK)
	droppingSub := broker.Subscribe("sports", 3, DROP)
	evictSub := broker.Subscribe("sports", 3, EVICT)

	go func() {
		for message := range blockingSub {
			fmt.Printf("[BlockingSub] received: %s\n", message)
			time.Sleep(1000 * time.Millisecond) // Reads messages every second
		}
	}()

	go func() {
		for message := range droppingSub {
			fmt.Printf("[DroppingSub] received: %s\n", message)
			time.Sleep(1000 * time.Millisecond) // Reads messages every second
		}
	}()

	go func() {
		for message := range evictSub {
			fmt.Printf("[EvictSub] received: %s\n", message)
			time.Sleep(1000 * time.Millisecond) // Reads messages every second
		}
	}()

	for _, message := range messages["sports"] {
		broker.Publish("sports", message)
		time.Sleep(100 * time.Millisecond) // Publish messages every 300ms
	}

	time.Sleep(20 * time.Second)
}
