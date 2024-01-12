package main

// myRWMutex is a struct that represents a read-write mutex.
// It has two channels: mut for write locks and readers for read locks.
type myRWMutex struct {
  mut     chan struct{} // Channel used for write locks.
  readers chan int      // Channel used for read locks.
}

// newMyRWMutex is a constructor function that creates a new myRWMutex.
// It initializes the mut channel with a buffer size of 1 to allow one write lock at a time.
// It also initializes the readers channel with a buffer size of 1 to allow one read lock at a time.
func newMyRWMutex() *myRWMutex {
  return &myRWMutex{
    mut:     make(chan struct{}, 1), // Initialize the mut channel.
    readers: make(chan int, 1),      // Initialize the readers channel.
  }
}

// Lock is a method that acquires a write lock.
// It sends an empty struct to the mut channel.
// If the channel is full, this operation will block until the channel is empty,
// ensuring that only one goroutine can acquire the write lock at a time.
func (m *myRWMutex) Lock() {
  m.mut <- struct{}{} // Acquire the write lock.
}

// Unlock is a method that releases a write lock.
// It receives from the mut channel.
// If the channel is empty, this operation will block until the channel is full,
// ensuring that a goroutine can only release the write lock if it has acquired it.
func (m *myRWMutex) Unlock() {
  <-m.mut // Release the write lock.
}

// RLock is a method that acquires a read lock.
// It sends an integer to the readers channel.
// If the channel is full, this operation will block until the channel is empty,
// ensuring that only one goroutine can acquire the read lock at a time.
func (m *myRWMutex) RLock() {
  m.readers <- 1 // Acquire the read lock.
}

// RUnlock is a method that releases a read lock.
// It receives from the readers channel.
// If the channel is empty, this operation will block until the channel is full,
// ensuring that a goroutine can only release the read lock if it has acquired it.
func (m *myRWMutex) RUnlock() {
  <-m.readers // Release the read lock.
}