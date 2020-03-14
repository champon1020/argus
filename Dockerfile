FROM golang:latest

WORKDIR /go/src/app

RUN mkdir /go/src/app
RUN go get -d -v ./...
RUN go install -v ./...

CMD ["app"]
