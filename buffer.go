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
	cacheSize      int
	cache          map[uint32]*Entry // Simple cahe
	cacheQueue     []*Entry          // Most recent at the front
	WriteBatch     []Operation       // Write buffer
	WriteBatchSize int
	BatchTimer     *time.Timer
	Disk           *Disk
}

type Operation struct {
	Key   uint32
	Value json.RawMessage
}

func NewBuffer(cacheSize int, writeBatchSize int, disk *Disk) *Buffer {
	return &Buffer{
		cacheSize:      cacheSize,
		cache:          make(map[uint32]*Entry),
		cacheQueue:     make([]*Entry, 0, cacheSize),
		WriteBatch:     make([]Operation, 0, writeBatchSize),
		WriteBatchSize: writeBatchSize,
		Disk:           disk,
	}
}

// Timer for batch operations
var BatchTimer *time.Timer

// Duration after which the buffer is flushed to disk
const FlushDuration = 1 * time.Minute

func (b *Buffer) UpdateCache(key uint32, value json.RawMessage) {
	if entry, ok := b.cache[key]; ok {
		entry.value = value
		b.moveToFront(entry)
		return
	}

	if len(b.cacheQueue) == b.cacheSize {
		delete(b.cache, b.cacheQueue[b.cacheSize-1].key)
		b.cacheQueue = b.cacheQueue[:b.cacheSize-1]
	}

	entry := &Entry{key, value}
	b.cache[key] = entry
	b.cacheQueue = append([]*Entry{entry}, b.cacheQueue...)
}

func (b *Buffer) BatchPut(ops []Operation) {
	if len(b.WriteBatch) == 0 && len(ops) > 0 {
		b.BatchTimer = time.AfterFunc(FlushDuration, b.flushBuffer)
	}

	for _, op := range ops {
		b.UpdateCache(op.Key, op.Value)
	}

	b.WriteBatch = append(b.WriteBatch, ops...)

	if len(b.WriteBatch) >= b.WriteBatchSize {
		b.flushBuffer()
	}
}

func (b *Buffer) Put(key uint32, value json.RawMessage) {
	b.UpdateCache(key, value)

	// Add operation to batch buffer
	b.WriteBatch = append(b.WriteBatch, Operation{Key: key, Value: value})

	// If this is the first operation in the buffer, start the timer
	if len(b.WriteBatch) == 1 {
		b.BatchTimer = time.AfterFunc(FlushDuration, b.flushBuffer)
	}

	// If buffer size has reached the maximum, flush to disk
	if len(b.WriteBatch) >= b.WriteBatchSize {
		b.flushBuffer()
	}
}

func (b *Buffer) flushBuffer() {
	// Flush the write buffer to disk
	for _, op := range b.WriteBatch {
		err := b.Disk.Put(op.Key, op.Value)
		if err != nil {
			fmt.Println("Error writing to disk:", err)
			return
		}
	}

	// Clear the buffer and stop the timer after flushing
	b.WriteBatch = []Operation{}
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
	for i, e := range b.cacheQueue {
		if e == entry {
			b.cacheQueue = append(b.cacheQueue[:i], b.cacheQueue[i+1:]...)
			b.cacheQueue = append([]*Entry{entry}, b.cacheQueue...)
			break
		}
	}
}
