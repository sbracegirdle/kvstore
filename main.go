package main

func main() {
	kv := NewStore(100, 1000)
	startServer(kv)
}
