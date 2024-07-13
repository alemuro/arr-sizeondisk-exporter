############################
# Compile Golang service   #
############################

FROM golang:alpine AS builder
RUN \
    apk update && \
    apk add --no-cache git ca-certificates && \
    update-ca-certificates

WORKDIR /go/src/app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -a -installsuffix cgo -o exporter cmd/cli/main.go

############################
# Optimize resulting image #
############################

FROM scratch

LABEL org.opencontainers.image.source https://github.com/alemuro/arr-sizeondisk-exporter

EXPOSE 9101
COPY --from=builder /go/src/app/exporter /go/bin/exporter
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
CMD ["/go/bin/exporter"]