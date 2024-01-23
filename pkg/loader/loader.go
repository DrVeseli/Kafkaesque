package loader

import (
	"encoding/json"
	"fmt"
)

// Data represents the data structure of your JSON data.
type Data struct {
	Key string `json:"key"`
}

// Queue represents a simple queue using a channel.
type Queue struct {
	ch chan Data
}

// NewQueue creates a new Queue.
func NewQueue() *Queue {
	return &Queue{ch: make(chan Data, 100)} // Buffer size of 100.
}
func (q *Queue) Len() int {
	return len(q.ch)
}

// Enqueue adds an item to the queue.
func (q *Queue) Enqueue(data Data) {
	q.ch <- data
}

// Dequeue removes an item from the queue.
func (q *Queue) Dequeue() (Data, bool) {
	data, ok := <-q.ch
	return data, ok
}

// Loader takes JSON data as a byte slice, unmarshals it into a Data struct,
// and enqueues it into the provided Queue.
func Loader(q *Queue, jsonData []byte) {
	var data Data
	if err := json.Unmarshal(jsonData, &data); err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return
	}

	q.Enqueue(data)
}
