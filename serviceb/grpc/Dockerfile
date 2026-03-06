# syntax=docker/dockerfile:1
FROM golang:1.25 AS builder

WORKDIR /app

# cache modules
COPY go.mod go.sum ./
RUN go mod download

# copy the full module source (protos + serviceb are needed)
COPY protos/ protos/
COPY serviceb/ serviceb/

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o /serviceb ./serviceb/cmd

FROM alpine:3.18
RUN apk add --no-cache ca-certificates

# copy the built binary from the builder stage
COPY --from=builder /serviceb /serviceb

EXPOSE 50051
ENTRYPOINT ["/serviceb"]
