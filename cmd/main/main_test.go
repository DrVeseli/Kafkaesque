package main

import (
	"bytes"
	"io"
	"kafkaesque/pkg/bartender"
	"kafkaesque/pkg/loader"
	"net/http"
	"net/http/httptest"
	"testing"
)

type TestSubscriber struct{}

func (s *TestSubscriber) Receive(data loader.Data) {
	// You can add code here to check the received data.
}

func TestEndpoints(t *testing.T) {
	q := loader.NewQueue() // Initialize the queue
	b := bartender.NewBartender(q)

	loaderHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

	subscribeHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		// Add a new subscriber.
		b.Subscribe(&TestSubscriber{})
	})

	// Test the /loader endpoint.
	{
		server := httptest.NewServer(loaderHandler)
		defer server.Close()

		resp, err := http.Post(server.URL+"/loader", "application/json", bytes.NewBuffer([]byte(`{"key": "value"}`)))
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status OK; got %v", resp.Status)
		}
	}

	// Test the /subscribe endpoint.
	{
		server := httptest.NewServer(subscribeHandler)
		defer server.Close()

		resp, err := http.Post(server.URL+"/subscribe", "application/json", nil)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status OK; got %v", resp.Status)
		}
	}
}
