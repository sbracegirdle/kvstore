package main

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
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
		fmt.Println("Error opening file:", err)
		return nil, err
	}

	indexFile, err := os.OpenFile(indexFilename, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("Error opening index file:", err)
		return nil, err
	}

	var index *IndexTree
	stat, err := indexFile.Stat()
	if err != nil {
		fmt.Println("Error getting index file info:", err)
		return nil, err
	}

	if stat.Size() == 0 {
		index = createIndexTree([]IndexValue{}, 3)
	} else {
		indexFile.Seek(0, 0)
		index = new(IndexTree)
		decoder := gob.NewDecoder(indexFile)
		if err := decoder.Decode(index); err != nil {
			fmt.Println("Error decoding index:", err)
			return nil, err
		}
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
		fmt.Println("Error encoding record:", err)
		return err
	}

	d.Index.Insert(&IndexValue{
		Key: key,
		Pos: position,
	})

	d.Index.Print()

	// Seek to the beginning of the file
	_, err := d.IndexFile.Seek(0, 0)
	if err != nil {
		return err
	}

	// Truncate the file to 0 length
	err = d.IndexFile.Truncate(0)
	if err != nil {
		return err
	}

	// Create a new encoder and encode the index
	encoder = gob.NewEncoder(d.IndexFile)
	if err := encoder.Encode(d.Index); err != nil {
		fmt.Println("Error encoding index:", err)
		return err
	}

	return nil
}
