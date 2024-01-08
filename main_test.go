package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSet(t *testing.T) {
	kv := NewStore(100, 1000)

	tests := []struct {
		key   string
		value string
	}{
		{"key1", "value1"},
		{"key2", "value2"},
		{"key3", "value3"},
		{"key4", "value4"},
		{"key5", "value5"},
		{"key6", "value6"},
		{"key7", "value7"},
		{"key8", "value8"},
		{"key9", "value9"},
		{"key10", "value10"},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			kv.Set(tt.key, tt.value)

			got, _ := kv.Get(tt.key)
			if got != tt.value {
				t.Errorf("Set(%q) = %v, want %v", tt.key, got, tt.value)
			}
		})
	}
}

func TestSetOverwrite(t *testing.T) {
	kv := NewStore(100, 1000)
	key := "key1"
	value1 := "value1"
	value2 := "value2"

	kv.Set(key, value1)
	kv.Set(key, value2)

	got, _ := kv.Get(key)
	if got != value2 {
		t.Errorf("Set(%q) = %v, want %v", key, got, value2)
	}
}

func TestGetNonExistentKey(t *testing.T) {
	kv := NewStore(100, 1000)
	key := "nonexistent"

	_, ok := kv.Get(key)
	if ok {
		t.Errorf("Get(%q) = %v, want %v", key, ok, false)
	}
}

func TestSetGetLargeData(t *testing.T) {
	kv := NewStore(100, 1000)
	key := "large"
	value := strings.Repeat("a", 1<<20) // 1 MiB

	kv.Set(key, value)
	got, _ := kv.Get(key)

	if got != value {
		t.Errorf("Set(%q) = %v, want %v", key, got, value)
	}
}

func TestAPI(t *testing.T) {
	// Start the server.
	kv := NewStore(100, 1000)
	go startServer(kv)

	client := &http.Client{}

	// Test POST /keys/:key
	req, _ := http.NewRequest("POST", "http://localhost:8080/keys/testKey", bytes.NewBufferString(`{"value":"testValue"}`))
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	assert.Equal(t, 200, resp.StatusCode)
	assert.Contains(t, string(body), "success")

	// Test GET /keys/:key
	req, _ = http.NewRequest("GET", "http://localhost:8080/keys/testKey", nil)
	resp, err = client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	body, _ = ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	assert.Equal(t, 200, resp.StatusCode)
	assert.Contains(t, string(body), "testValue")
}
