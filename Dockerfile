FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o cache-proxy .

FROM gcr.io/distroless/static-debian12

WORKDIR /app
COPY --from=builder /app/cache-proxy /app/cache-proxy

ENV PORT=8080
EXPOSE 8080

CMD ["/app/cache-proxy"]