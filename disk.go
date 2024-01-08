package main

import (
	"encoding/gob"
	"fmt"
	"os"
)

type Disk struct {
	PageSize int
}

type Page struct {
	Data map[uint32]string
}

const prefix = "db_"

func NewDisk(pageSize int) *Disk {
	return &Disk{
		PageSize: pageSize,
	}
}

func (d *Disk) Get(key uint32) (string, bool) {
	file, err := os.Open(fmt.Sprintf("%s%d", prefix, key))
	if err != nil {
		return "", false
	}
	defer file.Close()

	page := &Page{}
	decoder := gob.NewDecoder(file)
	if err := decoder.Decode(page); err != nil {
		return "", false
	}

	return page.Data[key], true
}

func (d *Disk) Put(key uint32, value string) {
	file, _ := os.Create(fmt.Sprintf("%s%d", prefix, key))
	defer file.Close()

	page := &Page{
		Data: map[uint32]string{
			key: value,
		},
	}
	encoder := gob.NewEncoder(file)
	encoder.Encode(page)
}
