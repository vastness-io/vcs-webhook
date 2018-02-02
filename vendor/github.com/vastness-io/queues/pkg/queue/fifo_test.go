package queue

import (
	"fmt"
	"testing"
	"time"
)

func TestFIFOQueue(t *testing.T) {
	q := NewFIFOQueue()
	count := 10000
	var first, tail int
	for i := 0; i < count; i++ {
		q.Enqueue(i)
		tail++
		if q.Size() <= 0 {
			t.Error("Shouldn't be empty")
		}
		if q.Size() != tail {
			t.Error("Should be equal")
		}

		v, shutdown := q.Dequeue()

		if shutdown {
			t.Error("Should be false")
		}

		if v.(int) != first {
			t.Error("Should be the same element")
		}
		tail--
		first++
	}

	if q.Size() != 0 {
		t.Error("Meant to be an empty queue")
	}

	signalCh := make(chan struct{})

	go func() {
		_, shutdown := q.Dequeue()

		if shutdown {
			signalCh <- struct{}{}
		}
	}()

	q.ShutDown()

	select {

	case <-time.After(30 * time.Second):
		t.Fatal("should have shutdown")

	case <-signalCh:
		fmt.Println("signal received")
	}

}
