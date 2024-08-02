FROM golang:1.22.3-alpine3.19 AS builder

COPY . /app
WORKDIR /app

RUN go mod download
RUN go build -o ./bin/auth cmd/auth/main.go

FROM alpine:3.19

WORKDIR /app
COPY --from=builder /app/bin/auth .
COPY config.yaml config.yaml
ENV CONFIG_FILE_PATH=config.yaml

CMD ["./auth"]