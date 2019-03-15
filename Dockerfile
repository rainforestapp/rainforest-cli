FROM golang:1.9.7
WORKDIR /go/src/github.com/rainforestapp/rainforest-cli
RUN curl -s https://glide.sh/get | sh
ADD glide.yaml glide.lock ./
RUN glide install
COPY . .
RUN go-wrapper install -ldflags "-X main.build=docker"

ENTRYPOINT ["rainforest-cli", "--skip-update"]
