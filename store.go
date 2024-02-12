package main

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"os"
)

type Store struct {
	Buffer *Buffer
	Mutex  *myRWMutex
	WALog  *os.File // Write-ahead log
}

type StoreEntry struct {
	Key   uint32
	Value json.RawMessage
}

// Maximum size of the buffer before flushing to disk
const MaxBufferSize = 100

func NewStore(bufferSize int, filename string, indexFilename string) *Store {
	disk, err := NewDisk(filename, indexFilename)
	buffer := NewBuffer(bufferSize, MaxBufferSize, disk)
	waLog, err := os.OpenFile("wa.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		fmt.Println("Error creating disk:", err)
		panic(err)
	}

	mutex := newMyRWMutex()
	return &Store{
		Buffer: buffer,
		Mutex:  mutex,
		WALog:  waLog,
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

	return value, ok
}

func (s *Store) Set(key string, value json.RawMessage) error {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	// Write the operation to the log before applying it to the index
	logEntry := fmt.Sprintf("Set %s %s\n", key, string(value))
	_, err := s.WALog.Write([]byte(logEntry))
	if err != nil {
		return err
	}

	// Write the operation to the buffer
	hash := s.hashKey(key)
	s.Buffer.Put(hash, value)

	return nil
}

func (s *Store) BatchSet(entries []StoreEntry) error {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	// Write the operations to the log before applying them to the index
	for _, entry := range entries {
		logEntry := fmt.Sprintf("Set %d %s\n", entry.Key, string(entry.Value))
		_, err := s.WALog.Write([]byte(logEntry))
		if err != nil {
			return err
		}
	}

	// Write the operations to the buffer
	ops := make([]Operation, len(entries))
	for i, entry := range entries {
		ops[i] = Operation{Key: entry.Key, Value: entry.Value}
	}
	s.Buffer.BatchPut(ops)

	return nil
}
