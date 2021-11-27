# syntax=docker/dockerfile:1

FROM golang:1.17-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY main.go ./
COPY collectors ./collectors

RUN go build -o /hilo-exporter

ENTRYPOINT ["/hilo-exporter"]

EXPOSE 6464/tcp
