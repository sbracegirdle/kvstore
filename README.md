# KVStore

This is a Key-Value Store developed in Go. This project is for learning purposes and is not suitable for production use.

Features:

- [x] Hash map data structure
- [x] LRU cache buffering
- [x] HTTP API
- [x] Paged I/O
- [x] Document storage
- [x] Safe concurrency (mutex etc)
- [x] Single file storage
- [x] Web Console: Basic Web UI for querying and retrieving data.
- [x] Write Ahead Log: Add support for writing all data operations to a log
- [x] Indexes: Implement a way to retrieve records performantly from a non-primary key
- [x] Batch Operations: Add support for batch get/set operations.

Roadmap:

- [ ] Transactions: Implement transactions to allow multiple operations to be executed atomically.
- [ ] Compression: Add data compression to save storage space.
- [ ] Encryption: Implement data encryption for security.
- [ ] Replication: Add support for data replication across multiple nodes.
- [ ] Sharding: Implement sharding to distribute data across multiple nodes.
- [ ] Query Language: Implement a simple query language for complex retrievals.
- [ ] Access Control: Add user authentication and access control.
- [ ] Telemetry: Emit OpenTelemetry metrics and traces
- [ ] Query Language: Implement a simple query language for complex retrievals.


## Files

- `main.go`: This is the main entry point of the application, which initialises a store and starts the HTTP api.
- `disk.go`: This file contains the `Disk` struct and its methods. The `Disk` struct represents a disk where the key-value pairs are stored. It has methods for getting and putting data on the disk.
- `index.go`: B-tree based index for finding the position of a record from its key.
- `store.go`: This file contains the Store struct and its methods. The Store struct represents a key-value store that uses a buffer and a disk for storage. It has methods for setting and getting key-value pairs. The Set method stores the key-value pair in both the buffer and the disk. The Get method first tries to get the value from the buffer. If it's not in the buffer, it tries to get it from the disk and if successful, puts it in the buffer for future access.
- `buffer.go`: This file contains the Buffer struct and its methods. The Buffer struct represents a buffer that stores a certain number of key-value pairs in memory for quick access. It has methods for getting and putting data in the buffer. If the buffer is full and a new key-value pair needs to be put in the buffer, it removes the least recently used (LRU cache) key-value pair before putting the new one.
- `http.go`: This file contains the startServer function which starts an HTTP server. The server has two routes: a GET route for getting the value of a key and a POST route for setting the value of a key. The server uses the Store to get and set the key-value pairs.

## Usage (as a library)

To use the key-value store, you can create a new disk with a specified page size:

```go
kv := NewStore(100, 1000)
```

You can then put a key-value pair on the disk:

```go
kv.Set("myKey", "myValue")
```

And get a value from the disk by its key:

```go
value, ok := kv.Get(tt.key)
```

## Usage (http API)

The HTTP API provides two endpoints: a GET endpoint for retrieving the value of a key and a POST endpoint for setting the value of a key.

To retrieve the value of a key, you can use the following curl command:

```sh
curl http://localhost:8080/keys/your_key
```

Replace `your_key`` with the key you want to retrieve. The server will return the value of the key.

To set the value of a key, you can use the following curl command:

```sh
curl -X POST -H "Content-Type: application/json" -d '{"Value":"your_value"}' http://localhost:8080/keys/your_key
```

Replace your_key with the key you want to set and your_value with the value you want to set. The server will store the key-value pair and return a confirmation message.

Please note that the server must be running for these commands to work.

## Running the server

To run the K/V store, including http server, navigate to the kvstore directory and run the main.go file:

```sh
go run .
```

## Running the tests


```sh
go test
```