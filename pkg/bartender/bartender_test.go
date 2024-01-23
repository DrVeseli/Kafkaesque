package bartender

import (
	"kafkaesque/pkg/loader"
	"testing"
)

type TestSubscriber struct {
	ch chan loader.Data
}

func (s *TestSubscriber) Receive(data loader.Data) {
	s.ch <- data
}

func TestBartender(t *testing.T) {
	q := loader.NewQueue()
	b := NewBartender(q)

	s := &TestSubscriber{ch: make(chan loader.Data)}
	b.Subscribe(s)

	testData := loader.Data{Key: "value"}
	q.Enqueue(testData)

	go b.Serve()

	data := <-s.ch
	if data != testData {
		t.Errorf("Data does not match, got %+v, want %+v", data, testData)
	}
}
