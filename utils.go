package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	redis "github.com/redis/go-redis/v9"
)

const metadataTTL = 24 * time.Hour

func (p *Proxy) MetadataHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		p.GetMetadataHandler(w, r)
	case http.MethodPut:
		p.PutMetadataHandler(w, r)
	default:
		w.Header().Set("Allow", "GET, PUT")
		http.Error(w, `{"ok": false}`, http.StatusMethodNotAllowed)
	}
}

func getVideoID(r *http.Request) string {
	return strings.TrimPrefix(r.URL.Path, "/metadata/")
}

func metadataCacheKey(videoID string) string {
	return fmt.Sprintf("metadata:v1:%s", videoID)
}

func (p *Proxy) GetMetadataHandler(w http.ResponseWriter, r *http.Request) {

	// get the videoId from the URL path
	videoID := getVideoID(r)
	if videoID == "" {
		http.Error(w, `{"ok": false}`, http.StatusBadRequest)
		return
	}

	//define unique cache key for videoId
	cacheKey := metadataCacheKey(videoID)
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

func (p *Proxy) PutMetadataHandler(w http.ResponseWriter, r *http.Request) {
	videoID := getVideoID(r)
	if videoID == "" {
		http.Error(w, `{"ok": false}`, http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil || len(bytes.TrimSpace(body)) == 0 {
		http.Error(w, `{"ok": false}`, http.StatusBadRequest)
		return
	}

	// Validate that the payload is JSON before storing it as-is in Redis.
	var payload json.RawMessage
	if err := json.Unmarshal(body, &payload); err != nil {
		http.Error(w, `{"ok": false}`, http.StatusBadRequest)
		return
	}

	cacheKey := metadataCacheKey(videoID)
	if err := p.redis.Set(ctx, cacheKey, body, metadataTTL).Err(); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]bool{"ok": false})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"ok":       true,
		"videoId":  videoID,
		"cacheKey": cacheKey,
		"ttl":      int(metadataTTL.Seconds()),
	})
}
