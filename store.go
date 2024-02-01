package main

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"time"
)

type Store struct {
	Buffer *Buffer
	Disk   *Disk
	Mutex  *myRWMutex
}

type Operation struct {
	Key   uint32
	Value json.RawMessage
}

// Buffer for batch operations
var BatchBuffer []Operation

// Maximum size of the buffer before flushing to disk
const MaxBufferSize = 100

// Timer for batch operations
var BatchTimer *time.Timer

// Duration after which the buffer is flushed to disk
const FlushDuration = 1 * time.Minute

func NewStore(bufferSize int, filename string, indexFilename string) *Store {
	buffer := NewBuffer(bufferSize)
	disk, err := NewDisk(filename, indexFilename)

	if err != nil {
		fmt.Println("Error creating disk:", err)
		panic(err)
	}

	mutex := newMyRWMutex()
	return &Store{
		Buffer: buffer,
		Disk:   disk,
		Mutex:  mutex,
	}
}

func (s *Store) hashKey(key string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(key))
	return h.Sum32()
}

func (s *Store) Get(key string) (json.RawMessage, bool) {
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()

	hash := s.hashKey(key)
	value, ok := s.Buffer.Get(hash)
	if ok {
		return value, ok
	}
	value, ok = s.Disk.Get(hash)
	if ok {
		s.Buffer.Put(hash, value)
	}
	return value, ok
}

func (s *Store) Set(key string, value json.RawMessage) error {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	hash := s.hashKey(key)
	s.Buffer.Put(hash, value)

	// Add operation to batch buffer
	BatchBuffer = append(BatchBuffer, Operation{Key: hash, Value: value})

	// If this is the first operation in the buffer, start the timer
	if len(BatchBuffer) == 1 {
		BatchTimer = time.AfterFunc(FlushDuration, s.flushBuffer)
	}

	// If buffer size has reached the maximum, flush to disk
	if len(BatchBuffer) >= MaxBufferSize {
		s.flushBuffer()
	}

	return nil
}

func (s *Store) flushBuffer() {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	for _, op := range BatchBuffer {
		err := s.Disk.Put(op.Key, op.Value)
		if err != nil {
			fmt.Println("Error writing to disk:", err)
			return
		}
	}
	// Clear the buffer and stop the timer after flushing
	BatchBuffer = []Operation{}
	if BatchTimer != nil {
		BatchTimer.Stop()
		BatchTimer = nil
	}
}
