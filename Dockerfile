FROM golang:1.24.4 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w" \
    -o GO_TODO-list \
    ./main.go

FROM ubuntu:latest

RUN apt-get update && apt-get install -y \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /app/GO_TODO-list .
COPY --from=builder /app/web ./web/

RUN mkdir -p /app/data

ENV TODO_PORT=7540
ENV TODO_DBFILE=/app/data/scheduler.db

EXPOSE 7540

CMD ["./GO_TODO-list"]