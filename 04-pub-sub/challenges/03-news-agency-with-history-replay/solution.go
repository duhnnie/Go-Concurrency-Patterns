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

// Ring Buffer
type RingBuffer[T any] struct {
	items      []T
	current    int
	currentLen int
	mu         sync.RWMutex
}

func NewRingBuffer[T any](size int) *RingBuffer[T] {
	return &RingBuffer[T]{
		items:   make([]T, size),
		current: -1,
	}
}

func (r *RingBuffer[T]) Insert(item T) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.current = (r.current + 1) % len(r.items)
	r.items[r.current] = item

	if r.currentLen < len(r.items) {
		r.currentLen++
	}
}

func (r *RingBuffer[T]) GetList() []T {
	r.mu.RLock()
	defer r.mu.RUnlock()

	sliceLen := len(r.items)
	slice := make([]T, sliceLen)
	slice = append(r.items[(r.current+1)%sliceLen:r.currentLen], r.items[0:r.current+1]...)
	return slice
}

// News Agency
type NewsAgency struct {
	topicSubscribers map[Topic][]chan string
	allSubscribers   []chan string
	topicMessages    map[Topic]*RingBuffer[string]
	mu               sync.RWMutex
}

func NewNewsAgency() *NewsAgency {
	return &NewsAgency{
		topicSubscribers: make(map[Topic][]chan string),
		allSubscribers:   make([]chan string, 0),
		topicMessages:    make(map[Topic]*RingBuffer[string]),
	}
}

func (na *NewsAgency) Subscribe(topic Topic) chan string {
	subscriberChan := make(chan string, 10)

	na.mu.Lock()
	na.topicSubscribers[topic] = append(na.topicSubscribers[topic], subscriberChan)
	ring := na.topicMessages[topic]
	na.mu.Unlock()

	// replay last messages if available
	if ring != nil {
		for _, msg := range ring.GetList() {
			subscriberChan <- msg
		}
	}

	return subscriberChan
}

func (na *NewsAgency) SubscribeAll() chan string {
	subscriberChan := make(chan string, 10)
	na.mu.Lock()
	na.allSubscribers = append(na.allSubscribers, subscriberChan)
	na.mu.Unlock()
	return subscriberChan
}

func (na *NewsAgency) Publish(topic Topic, message string) {
	// store in history
	na.mu.Lock()
	if _, ok := na.topicMessages[topic]; !ok {
		na.topicMessages[topic] = NewRingBuffer[string](5)
	}
	na.topicMessages[topic].Insert(message)

	// copy subscribers
	topicSubs := append([]chan string(nil), na.topicSubscribers[topic]...)
	allSubs := append([]chan string(nil), na.allSubscribers...)
	na.mu.Unlock()

	// fan-out to topic subscribers
	for _, sub := range topicSubs {
		go func(ch chan string) {
			select {
			case ch <- message:
			default:
				fmt.Printf("⚠️ Dropping %s message %q (slow consumer)\n", topic, message)
			}
		}(sub)
	}

	// fan-out to "all" subscribers
	for _, sub := range allSubs {
		go func(ch chan string) {
			select {
			case ch <- message:
			default:
				fmt.Printf("⚠️ Dropping broadcast message %q (slow consumer)\n", message)
			}
		}(sub)
	}
}

func reporter(topic Topic, news []string, agency *NewsAgency, wg *sync.WaitGroup) {
	defer wg.Done()
	for _, item := range news {
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(500)+300))
		agency.Publish(topic, item)
	}
}

func main() {
	allNews := map[Topic][]string{
		SPORTS: {
			"Team A won the match!",
			"Team B lost the match!",
			"Team C won the match!",
			"Team D lost the match!",
			"Team E won the match!",
			"Team F lost the match!",
		},
	}

	agency := NewNewsAgency()
	var wg sync.WaitGroup
	var wg2 sync.WaitGroup

	// Subscriber A joins immediately
	sportsSubA := agency.Subscribe(SPORTS)
	wg2.Add(1)

	go func() {
		count := 0
		for msg := range sportsSubA {
			fmt.Println("📺 Sports A:", msg)
			count++
			if count == 3 {
				wg2.Done() // after 3 messages, we start B
			}
		}
	}()

	// Subscriber B joins late, should replay history
	go func() {
		wg2.Wait()
		sportsSubB := agency.Subscribe(SPORTS)
		for msg := range sportsSubB {
			fmt.Println("📺 Sports B:", msg)
		}
	}()

	// Reporters publish news
	for topic, items := range allNews {
		wg.Add(1)
		go reporter(topic, items, agency, &wg)
	}

	wg.Wait()
	time.Sleep(1 * time.Second)
}
