FROM golang:1.12.5

COPY . /go/src/github.com/chanwit/meshrun

WORKDIR /go/src/github.com/chanwit/meshrun

RUN GO111MODULE=on go build -o meshrun github.com/chanwit/meshrun/cmd
