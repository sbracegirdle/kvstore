package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var srv *http.Server

func startServer(kv *Store) {
	r := gin.Default()

	r.GET("/keys/:key", func(c *gin.Context) {
		key := c.Param("key")
		value, ok := kv.Get(key)
		if ok {
			c.JSON(200, gin.H{"value": value})
		} else {
			c.JSON(404, gin.H{"error": "Key not found"})
		}
	})

	r.POST("/keys/:key", func(c *gin.Context) {
		var body json.RawMessage
		err := c.BindJSON(&body)

		if err != nil {
			c.JSON(400, gin.H{"error": "Bad request"})
			return
		}

		kv.Set(c.Param("key"), body)
		c.JSON(200, gin.H{"status": "success"})
	})

	srv = &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()
}

func stopServer() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		panic(err) // failure/timeout shutting down the server gracefully
	}
}
