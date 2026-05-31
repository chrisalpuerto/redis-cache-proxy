package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	redis "github.com/redis/go-redis/v9"
)

func (p *Proxy) GetMetadataHandler(w http.ResponseWriter, r *http.Request) {

	// get the videoId from the URL path
	videoID := r.URL.Path[len("/metadata/"):]
	if videoID == "" {
		http.Error(w, `{"ok": false}`, http.StatusBadRequest)
		return
	}

	//define unique cache key for videoId
	cacheKey := fmt.Sprintf("metadata:v1:%s", videoID)
	//attempt to get video metadata from redis cache
	val, err := p.redis.Get(ctx, cacheKey).Result()

	w.Header().Set("Content-Type", "application/json")

	if err == nil {
		// CACHE HIT: videoId found in cache, return cached metadata
		w.Header().Set("X-Cache", "HIT")
		fmt.Fprint(w, val)
		return
	}

	if err != redis.Nil {
		// unexpected error when accessing Redis
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]bool{"ok": false})
		return
	}
	// cache miss
	w.Header().Set("X-Cache", "MISS")
	json.NewEncoder(w).Encode(map[string]bool{"ok": false})
}
