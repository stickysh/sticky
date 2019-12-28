
all: release

release:
	@echo ">> building binaries"
	go get github.com/rakyll/statik
	statik -src=stacktmpl
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o ./stickysrv ./cmd/stickysrv

