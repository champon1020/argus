FROM golang:1.15-alpine as builder

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
ENV GO111MODULE=on
ENV GOPRIVATE=github.com/champon1020

WORKDIR /go/src/github.com/champon1020/argus
COPY . .
RUN go mod download
RUN go build -o argus ./cmd/main.go
RUN chmod +x argus
