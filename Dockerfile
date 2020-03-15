FROM golang:latest

WORKDIR /go/src/github.com/champon1020/argus

COPY . /go/src/github.com/champon1020/argus

RUN go get -d -v ./...
RUN go build -o cmd/argus cmd/main.go

CMD ./cmd/argus stg
