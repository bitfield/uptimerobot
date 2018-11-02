FROM golang:1.11-alpine AS builder
WORKDIR /src/
COPY . .
ENV CGO_ENABLED=0
RUN apk --no-cache add git ca-certificates
RUN go test ./...
RUN go build -o /uptimerobot

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /uptimerobot /uptimerobot
ENTRYPOINT ["/uptimerobot"]
