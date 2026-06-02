package main

import (
	"context"
	"log"
	"net/http"
	"os"

	redis "github.com/redis/go-redis/v9"
)

var ctx = context.Background()

// Proxy holds shared dependencies for the cache proxy server
type Proxy struct {
	redis *redis.Client
}

func main() {
	// Connect to Redis using REDIS_URL env var (differs in prod vs local)
	opt, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		log.Fatalf("failed to parse REDIS_URL: %v", err)
	}
	client := redis.NewClient(opt)

	proxy := &Proxy{redis: client}

	// /metadata/{videoId}
	http.HandleFunc("/metadata/", proxy.MetadataHandler)
	// Cloud Run injects PORT; default to 8080 for local dev
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Starting server on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
