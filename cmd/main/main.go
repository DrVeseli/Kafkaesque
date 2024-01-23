package main

import (
	"fmt"
	"io"
	"kafkaesque/pkg/bartender"
	"kafkaesque/pkg/loader"
	"net/http"
)

type Subscriber struct{}

func (s *Subscriber) Receive(data loader.Data) {
	fmt.Println("Received data:", data)
}

func main() {
	q := loader.NewQueue() // Initialize the queue
	b := bartender.NewBartender(q)

	http.HandleFunc("/loader", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading body", http.StatusInternalServerError)
			return
		}

		// Pass the data to the loader module.
		loader.Loader(q, body)
	})

	http.HandleFunc("/subscribe", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		// Add a new subscriber.
		b.Subscribe(&Subscriber{})
	})

	go b.Serve()

	fmt.Println("Starting server on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
