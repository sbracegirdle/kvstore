package main

import (
	"encoding/json"
	"fmt"
	"time"
)

type Entry struct {
	key   uint32
	value json.RawMessage
}

type Buffer struct {
	size       int
	cache      map[uint32]*Entry // Simple cahe
	queue      []*Entry          // Most recent at the front
	Batch      []Operation       // Write buffer
	MaxSize    int
	BatchTimer *time.Timer
	Disk       *Disk
}

type Operation struct {
	Key   uint32
	Value json.RawMessage
}

func NewBuffer(size int, maxSize int, disk *Disk) *Buffer {
	return &Buffer{
		size:    size,
		cache:   make(map[uint32]*Entry),
		queue:   make([]*Entry, 0, size),
		Batch:   make([]Operation, 0, maxSize),
		MaxSize: maxSize,
		Disk:    disk,
	}
}

// Timer for batch operations
var BatchTimer *time.Timer

// Duration after which the buffer is flushed to disk
const FlushDuration = 1 * time.Minute

// UpdateCache
func (b *Buffer) UpdateCache(key uint32, value json.RawMessage) {
	if entry, ok := b.cache[key]; ok {
		entry.value = value
		b.moveToFront(entry)
		return
	}

	if len(b.queue) == b.size {
		delete(b.cache, b.queue[b.size-1].key)
		b.queue = b.queue[:b.size-1]
	}

	entry := &Entry{key, value}
	b.cache[key] = entry
	b.queue = append([]*Entry{entry}, b.queue...)
}

func (b *Buffer) Put(key uint32, value json.RawMessage) {
	b.UpdateCache(key, value)

	// Add operation to batch buffer
	b.Batch = append(b.Batch, Operation{Key: key, Value: value})

	// If this is the first operation in the buffer, start the timer
	if len(b.Batch) == 1 {
		b.BatchTimer = time.AfterFunc(FlushDuration, b.flushBuffer)
	}

	// If buffer size has reached the maximum, flush to disk
	if len(b.Batch) >= b.MaxSize {
		b.flushBuffer()
	}
}

func (b *Buffer) flushBuffer() {
	// Flush the write buffer to disk
	for _, op := range b.Batch {
		err := b.Disk.Put(op.Key, op.Value)
		if err != nil {
			fmt.Println("Error writing to disk:", err)
			return
		}
	}

	// Clear the buffer and stop the timer after flushing
	b.Batch = []Operation{}
	if b.BatchTimer != nil {
		b.BatchTimer.Stop()
		b.BatchTimer = nil
	}
}

func (b *Buffer) Get(key uint32) (json.RawMessage, bool) {
	if entry, ok := b.cache[key]; ok {
		b.moveToFront(entry)
		return entry.value, true
	}

	value, ok := b.Disk.Get(key)

	// Update the cache with the value from disk
	if ok {
		b.UpdateCache(key, value)
	}

	return value, ok
}

func (b *Buffer) moveToFront(entry *Entry) {
	for i, e := range b.queue {
		if e == entry {
			b.queue = append(b.queue[:i], b.queue[i+1:]...)
			b.queue = append([]*Entry{entry}, b.queue...)
			break
		}
	}
}
