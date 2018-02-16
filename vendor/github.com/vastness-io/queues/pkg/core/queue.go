package core

// Queue is a abstract data type.
type Queue interface {
	Size() int64 // Size of the current queue.

	Enqueue(interface{}) // Add node to Tail of the queue.

	Dequeue() interface{} // Removes node from the Head position.
}

// Queue is a abstract data type which blocks on dequeue if the queue is empty.
type BlockingQueue interface {
	Size() int64 // Size of the current queue.

	Enqueue(interface{}) // Add node to Tail of the queue.

	Dequeue() (interface{}, bool) // Removes node from the Head position.

	ShutDown() // Signals the queue to shutdown. This is necessary to notify go routines which are currently blocked on trying to dequeue.
}
