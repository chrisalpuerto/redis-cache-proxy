## Redis Cache Proxy 

Redis cache proxy for hooptuber.com, built in Go
- checks for video metadata inside cache instead of always hitting our API
- saving on API call costs & firestore read costs, especially at large scale