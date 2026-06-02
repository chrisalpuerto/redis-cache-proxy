## Redis Cache Proxy 

Redis cache proxy for hooptuber.com, built in Go
- checks for video metadata inside cache instead of always hitting our API
- saving on API call costs & firestore read costs, especially at large scale

### Endpoints

- `GET /metadata/{videoId}` returns cached metadata if present
- `PUT /metadata/{videoId}` stores JSON metadata in Redis for 24 hours

Example:

```bash
curl -X PUT http://localhost:8080/metadata/abc123 \
  -H "Content-Type: application/json" \
  -d '{"title":"Test video","duration":120}'
```
