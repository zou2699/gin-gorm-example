FROM golang:1.11

WORKDIR /go/src/app

COPY . .
COPY ./templates /go/bin/

RUN go get -d -v ./...
RUN go install -v ./...
RUN rm -rf /go/src/*

WORKDIR /go/bin/
CMD ["./app"]
