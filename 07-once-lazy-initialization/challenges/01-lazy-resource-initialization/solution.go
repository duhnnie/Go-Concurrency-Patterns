package main

import (
	"fmt"
	"math/rand"
	"sync"
)

type Resource struct {
	ID int
}

var resource *Resource
var once sync.Once

func initResource() {
	resource = &Resource{ID: rand.Intn(999)}
	fmt.Println("Resource initialized!")
}

func GetResource() *Resource {
	once.Do(initResource)
	return resource
}

func main() {
	wg := sync.WaitGroup{}
	wg.Add(10)

	for i := range 10 {
		go func() {
			r := GetResource()
			fmt.Printf("goroutine #%d received resource with ID: %d\n", i, r.ID)
			wg.Done()
		}()
	}

	wg.Wait()
	r := GetResource()
	fmt.Printf("Final resource pointer: %p\n", r)
}
