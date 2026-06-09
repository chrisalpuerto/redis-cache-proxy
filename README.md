## Redis Cache Proxy 

Redis cache proxy for hooptuber.com, built in Go
- checks for video metadata inside cache instead of always hitting our API
- saving on API call costs & firestore read costs, especially at large scale

### Endpoints

`/metadata/{videoId}` uses Redis keys in the format `metadata:v1:{videoId}` with a TTL of 24 hours.

#### `GET /metadata/{videoId}`

Returns cached metadata for the given `videoId`.

- Request body: none
- `200 OK` on cache hit with `Content-Type: application/json` and `X-Cache: HIT`
- `200 OK` on cache miss with `Content-Type: application/json`, `X-Cache: MISS`, and `{"ok": false}`
- `400 Bad Request` if `videoId` is missing
- `405 Method Not Allowed` for unsupported methods with `Allow: GET, PUT`
- `500 Internal Server Error` if Redis returns an unexpected error


#### `PUT /metadata/{videoId}`

Stores the request body as cached metadata for the given `videoId`.

- Request body: non-empty valid JSON
- `200 OK` with `{"ok": true, "videoId": "...", "cacheKey": "...", "ttl": 86400}` on success
- `400 Bad Request` if `videoId` is missing, the body is empty, or the body is not valid JSON
- `405 Method Not Allowed` for unsupported methods with `Allow: GET, PUT`
- `500 Internal Server Error` if Redis returns an unexpected error
