package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var srv *http.Server

func startServer(kv *Store) {
	r := gin.Default()

	// Load HTML templates
	r.LoadHTMLGlob("templates/*")

	// Root path should re-direct to console
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/console/")
	})

	// Create a route group for the API
	api := r.Group("/api")
	{
		api.GET("/keys/:key", func(c *gin.Context) {
			key := c.Param("key")
			value, ok := kv.Get(key)
			if ok {
				c.JSON(200, gin.H{"value": value})
			} else {
				c.JSON(404, gin.H{"error": "Key not found"})
			}
		})

		api.POST("/keys/:key", func(c *gin.Context) {
			var body json.RawMessage
			err := c.BindJSON(&body)

			if err != nil {
				c.JSON(400, gin.H{"error": "Bad request"})
				return
			}

			err = kv.Set(c.Param("key"), body)

			if err != nil {
				c.JSON(500, gin.H{"error": "Internal server error"})
				return
			} else {
				c.JSON(200, gin.H{"status": "success"})
			}
		})

		// Variant of POST keys where the key is in the body instead of path
		api.POST("/keys", func(c *gin.Context) {
			var body struct {
				Key   string          `json:"key"`
				Value json.RawMessage `json:"value"`
			}
			err := c.BindJSON(&body)

			if err != nil {
				c.JSON(400, gin.H{"error": "Bad request"})
				return
			}

			err = kv.Set(body.Key, body.Value)

			if err != nil {
				c.JSON(500, gin.H{"error": "Internal server error"})
				return
			} else {
				c.JSON(200, gin.H{"status": "success"})
			}
		})
	}

	// Create a route group for the console
	console := r.Group("/console")
	{
		console.GET("/", func(c *gin.Context) {
			c.HTML(http.StatusOK, "index.html", gin.H{})
		})

		console.POST("/keys", func(c *gin.Context) {
			key := c.PostForm("key")
			value := c.PostForm("value")

			if key == "" || value == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Both key and value are required"})
				return
			}

			// Create a JSON manually
			var jsonValue json.RawMessage
			jsonValue = append(jsonValue, []byte(fmt.Sprintf(`{"value": "%s"}`, value))...)

			err := kv.Set(key, jsonValue)

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
				return
			}

			c.JSON(http.StatusOK, gin.H{"status": "success"})
		})

		console.GET("/keys", func(c *gin.Context) {
			key := c.Query("key")
			value, ok := kv.Get(key)
			if ok {
				c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(fmt.Sprintf("<div>%s</div>", value)))
			} else {
				c.Data(http.StatusNotFound, "text/html; charset=utf-8", []byte("<div>Key not found</div>"))
			}
		})
	}

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
