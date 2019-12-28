FROM golang:1.12.11-alpine3.9 AS builder

RUN apk add --update --no-cache ca-certificates git make=4.2.1-r2

WORKDIR /build

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN make release

FROM golang:1.12.11-alpine3.9

RUN apk --no-cache add ca-certificates git

WORKDIR /sticky

RUN mkdir git
RUN mkdir blueprint
RUN mkdir actions
RUN mkdir bin

COPY --from=builder /build/stickysrv bin/

RUN chmod +x bin/stickysrv

EXPOSE 6060

ENTRYPOINT ["./bin/stickysrv", "-dir=/sticky"]

