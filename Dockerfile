FROM golang:1.16.6-alpine3.14 as build-env

ENV GOOS linux

RUN apk add --no-cache git build-base && \
  echo "http://dl-cdn.alpinelinux.org/alpine/edge/main" >> /etc/apk/repositories && \
  echo "http://dl-cdn.alpinelinux.org/alpine/edge/community" >> /etc/apk/repositories && \
  echo "http://dl-cdn.alpinelinux.org/alpine/edge/testing" >> /etc/apk/repositories && \
  apk add --no-cache upx

WORKDIR /go/src/github.com/jonascheng/speedtest-go/
COPY . .

RUN make build && upx ./bin/speedtest-go

FROM alpine:3.14

WORKDIR /app

COPY --from=build-env /go/src/github.com/jonascheng/speedtest-go/bin/speedtest-go .

ENTRYPOINT [ "./speedtest-go" ]