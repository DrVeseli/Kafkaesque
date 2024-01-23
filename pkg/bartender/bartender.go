package bartender

import (
	"fmt"
	"kafkaesque/pkg/loader"
	"sync"
	"time"
)

type Subscriber interface {
	Receive(data loader.Data)
}

type Bartender struct {
	q           *loader.Queue
	subscribers []Subscriber
	mu          sync.Mutex
}

func NewBartender(q *loader.Queue) *Bartender {
	return &Bartender{
		q: q,
	}
}

func (b *Bartender) Subscribe(s Subscriber) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.subscribers = append(b.subscribers, s)
}

func (b *Bartender) Serve() {
	for {
		data, ok := b.q.Dequeue()
		if !ok {
			time.Sleep(1 * time.Second) // Wait for a while if the queue is empty
			continue
		}

		b.mu.Lock()
		for _, s := range b.subscribers {
			s.Receive(data)
		}
		b.mu.Unlock()
	}
}

func (b *Bartender) Measure() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		b.mu.Lock()
		queueSize := b.q.Len()
		b.mu.Unlock()

		if queueSize >= 50 {
			start := time.Now()
			for i := 0; i < 50; i++ {
				data, ok := b.q.Dequeue()
				if !ok {
					break
				}

				b.mu.Lock()
				for _, s := range b.subscribers {
					s.Receive(data)
				}
				b.mu.Unlock()
			}
			elapsed := time.Since(start)
			fmt.Printf("Time taken to process 50 items: %s\n", elapsed)
		}
	}
}
