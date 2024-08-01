FROM golang:1.22.3-alpine3.19 AS builder

ENV config=docker

WORKDIR /app

COPY . /app

ENV CONFIG_FILE_PATH=/app/config.yaml

RUN go mod tidy && \
    go mod download && \
    go get github.com/githubnemo/CompileDaemon && \
    go install github.com/githubnemo/CompileDaemon

ENTRYPOINT CompileDaemon --build="go build -o bin/auth cmd/auth/main.go" --command=./bin/auth
