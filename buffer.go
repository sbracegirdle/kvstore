package main

type Entry struct {
	key   uint32
	value string
}

type Buffer struct {
	size  int
	cache map[uint32]*Entry
	queue []*Entry
}

func NewBuffer(size int) *Buffer {
	return &Buffer{
		size:  size,
		cache: make(map[uint32]*Entry),
		queue: make([]*Entry, 0, size),
	}
}

func (b *Buffer) Get(key uint32) (string, bool) {
	if entry, ok := b.cache[key]; ok {
		b.moveToFront(entry)
		return entry.value, true
	}
	return "", false
}

func (b *Buffer) Put(key uint32, value string) {
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

func (b *Buffer) moveToFront(entry *Entry) {
	for i, e := range b.queue {
		if e == entry {
			b.queue = append(b.queue[:i], b.queue[i+1:]...)
			b.queue = append([]*Entry{entry}, b.queue...)
			break
		}
	}
}
