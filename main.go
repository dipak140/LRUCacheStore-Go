package main

import (
	"log"
	"time"

	"github.com/dipak140/LRUCacheStore-Go/internal"
	"github.com/gin-gonic/gin"
)

func main() {
	cache := internal.NewLRUCache(1024, 5*time.Second)

	log.Print(cache)

	r := gin.Default()

	r.GET("/cache/:key", func(ctx *gin.Context) {
		internal.GetFromCache(ctx, cache)
	})

	r.POST("/cache", func(ctx *gin.Context) {
		internal.PutToCache(ctx, cache)
	})

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}
