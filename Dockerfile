FROM  golang:latest

WORKDIR /go/src/app

ENV GO111MODULE=on

# RUN go get bitbucket.org/liamstask/goose/cmd/goose