package main

import (
	"encoding/json"
	"hash/fnv"
)

type Store struct {
	Buffer *Buffer
	Disk   *Disk
	Mutex  *myRWMutex
}

func NewStore(bufferSize int, filename string, indexFilename string) *Store {
	buffer := NewBuffer(bufferSize)
	disk, err := NewDisk(filename, indexFilename)

	if err != nil {
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

func (s *Store) Set(key string, value json.RawMessage) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	hash := s.hashKey(key)
	s.Buffer.Put(hash, value)
	s.Disk.Put(hash, value)
}
