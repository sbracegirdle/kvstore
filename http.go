package main

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
)

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

	r.Run() // listen and serve on 0.0.0.0:8080
}
