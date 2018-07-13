FROM golang:1.10-alpine as builder
WORKDIR /go/src/github.com/bitfield/uptimerobot/
ADD . .
ENV CGO_ENABLED=0
RUN apk --no-cache add git ca-certificates && \
    go get -t . && \
    go test && \
    go build -o /uptimerobot

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /uptimerobot /uptimerobot
ENTRYPOINT ["/uptimerobot"]
