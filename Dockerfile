FROM golang:1.17.1-alpine AS builder
RUN apk add --no-cache git
WORKDIR /build
COPY . .
RUN go build -ldflags "-X main.build=docker" -o build/rainforest

FROM alpine:3.9
RUN apk add --no-cache ca-certificates
COPY --from=builder /build/build/rainforest /usr/local/bin/rainforest
RUN ln -s /usr/local/bin/rainforest /usr/local/bin/rainforest-cli
ENTRYPOINT ["rainforest", "--skip-update"]
