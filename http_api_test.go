package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAPI(t *testing.T) {
	// Start the server.
	kv := NewStore(100, "test.db", "test.idx")
	go startServer(kv)
	defer stopServer()

	client := &http.Client{}

	// Test POST /keys/:key
	req, _ := http.NewRequest("POST", "http://localhost:8080/api/keys/testKey", bytes.NewBufferString(`{"value":"testValue"}`))
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
	req, _ = http.NewRequest("GET", "http://localhost:8080/api/keys/testKey", nil)
	resp, err = client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	body, _ = ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	assert.Equal(t, 200, resp.StatusCode)
	assert.Contains(t, string(body), "testValue")
	assert.Contains(t, string(body), "value")
}

func TestAPI_NotJSON(t *testing.T) {
	// Start the server.
	kv := NewStore(100, "test.db", "test.idx")
	go startServer(kv)
	defer stopServer()

	client := &http.Client{}

	// Test POST /keys/:key
	req, _ := http.NewRequest("POST", "http://localhost:8080/api/keys/testKey", bytes.NewBufferString(`testValue`))
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	assert.Equal(t, 400, resp.StatusCode)
	assert.Contains(t, string(body), "Bad request")
}
