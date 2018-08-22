# stage 1 (Build)
FROM golang:1.9 as builder

# install go dep
RUN curl -fsSL -o /usr/local/bin/dep https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64 && chmod +x /usr/local/bin/dep

RUN mkdir -p /go/src/github.com/ramsgoli/uberextbackend
WORKDIR /go/src/github.com/ramsgoli/uberextbackend

# install dependencies
COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure --vendor-only

COPY . ./

# build
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix main.go

# stage 2 (Run)
FROM alpine
RUN apk --no-cache add ca-certificates
WORKDIR /root
COPY --from=builder /go/src/github.com/ramsgoli/uberextbackend/uberextbackend .
COPY --from=builder /go/src/github.com/ramsgoli/uberextbackend/.env .


CMD ["PORT=8000", "./uberextbackend"]
