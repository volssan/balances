FROM golang:1.22-alpine3.19 as builder

COPY go.mod go.sum /go/src/github.com/volssan/balances/
WORKDIR /go/src/github.com/volssan/balances
RUN go mod download
COPY . /go/src/github.com/volssan/balances
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/balances github.com/volssan/balances


FROM alpine

RUN apk add --no-cache ca-certificates && update-ca-certificates
COPY --from=builder /go/src/github.com/volssan/balances/build/balances /usr/bin/balances

EXPOSE 8080 8080

ENTRYPOINT ["/usr/bin/balances"]
