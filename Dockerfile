FROM golang:1.12.4-alpine3.9 AS builder
RUN apk add --no-cache git
WORKDIR /build
COPY . .
RUN go build -ldflags "-X main.build=docker" -o build/rainforest

FROM alpine:3.9
RUN apk add --no-cache ca-certificates
COPY --from=builder /build/build/rainforest /usr/local/bin/rainforest
ENTRYPOINT ["rainforest", "--skip-update"]
