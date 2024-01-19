package main

import (
	"encoding/gob"
	"encoding/json"
	"os"
)

type Disk struct {
	Index     *IndexTree
	IndexFile *os.File
	File      *os.File
}

type Record struct {
	Key  uint32
	Data json.RawMessage
}

func NewDisk(filename string, indexFilename string) (*Disk, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	indexFile, err := os.Open(indexFilename)
	if err != nil {
		return nil, err
	}

	index := createIndexTree([]IndexValue{}, 3)
	decoder := gob.NewDecoder(indexFile)
	if err := decoder.Decode(index); err != nil {
		return nil, err
	}

	return &Disk{
		Index:     index,
		IndexFile: indexFile,
		File:      file,
	}, nil
}

func (d *Disk) Get(key uint32) (json.RawMessage, bool) {
	pos, success := d.Index.Get(key)
	if !success {
		return nil, false
	}

	d.File.Seek(pos, 0)

	record := &Record{}
	decoder := gob.NewDecoder(d.File)
	if err := decoder.Decode(record); err != nil {
		return nil, false
	}

	return record.Data, true
}

func (d *Disk) Put(key uint32, data json.RawMessage) error {
	record := &Record{
		Key:  key,
		Data: data,
	}

	position, _ := d.File.Seek(0, 2)
	encoder := gob.NewEncoder(d.File)
	if err := encoder.Encode(record); err != nil {
		return err
	}

	d.Index.Insert(&IndexValue{
		key: key,
		pos: position,
	})

	encoder = gob.NewEncoder(d.IndexFile)
	if err := encoder.Encode(d.Index); err != nil {
		return err
	}

	return nil
}
