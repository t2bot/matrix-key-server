# Build first
FROM golang:1.19-alpine AS builder
RUN apk add --no-cache git musl-dev
COPY . /opt
WORKDIR /opt
RUN go build -v -o bin/matrix-key-server

# The actual image (which is lightweight)
FROM alpine
COPY --from=builder /opt/bin/matrix-key-server /usr/local/bin/
RUN apk add --no-cache \
        su-exec \
        ca-certificates
ENTRYPOINT "/usr/local/bin/matrix-key-server"
EXPOSE 8080
