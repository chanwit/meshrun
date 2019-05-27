FROM golang:1.12.5

COPY . /go/src/github.com/chanwit/meshrun

WORKDIR /go/src/github.com/chanwit/meshrun

RUN GO111MODULE=on CGO_ENABLED=0 \
	go build -a -ldflags '-extldflags "-static"' \
	-o app github.com/chanwit/meshrun/cmd

FROM stratch

COPY --from=0 /go/src/github.com/chanwit/meshrun/app /meshrun

ENTRYPOINT ["/meshrun"]

