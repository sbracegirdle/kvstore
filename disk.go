package main

import (
	"encoding/gob"
	"fmt"
	"os"
	"encoding/json"
)

type Disk struct {
	PageSize int
}

type Page struct {
	Data map[uint32]json.RawMessage
}

const prefix = "db_"

func NewDisk(pageSize int) *Disk {
	return &Disk{
		PageSize: pageSize,
	}
}

func (d *Disk) Get(key uint32) (json.RawMessage, bool) {
	file, err := os.Open(fmt.Sprintf("%s%d", prefix, key))
	if err != nil {
		return nil, false
	}
	defer file.Close()

	page := &Page{}
	decoder := gob.NewDecoder(file)
	if err := decoder.Decode(page); err != nil {
		return nil, false
	}

	return page.Data[key], true
}

func (d *Disk) Put(key uint32, value json.RawMessage) {
	file, _ := os.Create(fmt.Sprintf("%s%d", prefix, key))
	defer file.Close()

	page := &Page{
		Data: map[uint32]json.RawMessage{
			key: value,
		},
	}
	encoder := gob.NewEncoder(file)
	encoder.Encode(page)
}
