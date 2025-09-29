package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Topic string

const (
	SPORTS   = Topic("Sports")
	TECH     = Topic("Tech")
	POLITICS = Topic("Politics")
	WORLD    = Topic("World")
)

type NewsAgency struct {
	topicSubscribers map[Topic][]chan<- string
	allSubscribers   []chan<- string
	mu               *sync.RWMutex
}

func NewNewsAgency() *NewsAgency {
	return &NewsAgency{
		topicSubscribers: make(map[Topic][]chan<- string),
		allSubscribers:   make([]chan<- string, 0),
		mu:               &sync.RWMutex{},
	}
}

func (na *NewsAgency) Subscribe(topic Topic) chan string {
	subcriberChan := make(chan string, 10)
	na.mu.Lock()
	defer na.mu.Unlock()
	na.topicSubscribers[topic] = append(na.topicSubscribers[topic], subcriberChan)
	return subcriberChan
}

func (na *NewsAgency) SubscribeAll() chan string {
	subcriberChan := make(chan string, 10)
	na.mu.Lock()
	defer na.mu.Unlock()
	na.allSubscribers = append(na.allSubscribers, subcriberChan)
	return subcriberChan
}

func (na *NewsAgency) Publish(topic Topic, message string) {
	na.mu.RLock()
	defer na.mu.RUnlock()

	for _, subscriber := range na.topicSubscribers[topic] {
		select {
		case subscriber <- message:
		default:
			fmt.Printf("Warning - dropping %s message: %q\n", topic, message)
		}
	}

	for _, allSubscriber := range na.allSubscribers {
		select {
		case allSubscriber <- message:
		default:
			fmt.Printf("Warning - dropping message for ALL consumer: %q\n", message)
		}
	}
}

func (na *NewsAgency) Unsubscribe(subscriber chan string) {
	na.mu.Lock()
	defer na.mu.Unlock()

	for i, subs := range na.allSubscribers {
		if subs == subscriber {
			fmt.Println("All subscriber is unscribing...")
			na.allSubscribers = append(na.allSubscribers[:i], na.allSubscribers[i+1:]...)
			close(subscriber)
			return
		}
	}

	for topic, subscribers := range na.topicSubscribers {
		for i, subs := range subscribers {
			if subs == subscriber {
				fmt.Printf("%s subscriber is unscribing...\n", topic)
				na.topicSubscribers[topic] = append(na.topicSubscribers[topic][:i], na.topicSubscribers[topic][i+1:]...)
				close(subscriber)
				return
			}
		}
	}
}

func reporter(topic Topic, news []string, agency *NewsAgency, wg *sync.WaitGroup) {
	for _, newItem := range news {
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000)+500))
		agency.Publish(topic, newItem)
	}

	wg.Done()
}

func main() {
	allNews := map[Topic][]string{
		SPORTS: {
			"Team A won the match!",
			"Team B lost the match!",
		},
		TECH: {
			"New AI chip released by major company",
			"Quantum breakthrough in labs",
		},
		POLITICS: {
			"Election results announced",
		},
		WORLD: {
			"Giant dinosaur destroys Tokyo",
		},
	}

	agency := NewNewsAgency()
	wg := &sync.WaitGroup{}

	subscribers := map[Topic]chan string{
		SPORTS:   agency.Subscribe(SPORTS),
		TECH:     agency.Subscribe(TECH),
		POLITICS: agency.Subscribe(POLITICS),
		WORLD:    agency.Subscribe(WORLD),
	}

	allSubcriber := agency.SubscribeAll()

	for topic, subscriber := range subscribers {
		go func(topic Topic, subscriber chan string) {
			for message := range subscriber {
				fmt.Printf("%s: %s\n", string(topic), message)

				if topic == SPORTS {
					agency.Unsubscribe(subscriber) // Unsubscribe sport subscriber
				}
			}
		}(topic, subscriber)
	}

	go func(subscriber chan string) {
		for message := range subscriber {
			fmt.Printf("ALL: %s\n", message)
		}
	}(allSubcriber)

	// Start reporters
	for topic, newsTopic := range allNews {
		wg.Add(1)
		go reporter(topic, newsTopic, agency, wg)
	}

	wg.Wait()
	time.Sleep(500 * time.Millisecond)
}
