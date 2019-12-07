
all: build

build:
	go get github.com/rakyll/statik
	statik -src=stacktmpl
	go build -o ctrl cmd/stickycli/main.go
