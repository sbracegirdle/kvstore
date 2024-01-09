package main

import (
	"hash/fnv"
	"encoding/json"
)

type Store struct {
	Buffer *Buffer
	Disk   *Disk
}

func NewStore(bufferSize, diskSize int) *Store {
	buffer := NewBuffer(bufferSize)
	disk := NewDisk(diskSize)
	return &Store{
		Buffer: buffer,
		Disk:   disk,
	}
}

func (s *Store) hashKey(key string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(key))
	return h.Sum32()
}

func (s *Store) Get(key string) (json.RawMessage, bool) {
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
	hash := s.hashKey(key)
	s.Buffer.Put(hash, value)
	s.Disk.Put(hash, value)
}
