FROM golang:alpine AS builder

ENV GOPROXY=https://goproxy.cn,https://goproxy.io,direct \
    GO111MODULE=on

WORKDIR /go/cache
COPY go.mod .
COPY go.sum .
RUN go mod download

WORKDIR /go/src/paste.org.cn/paste
COPY . .
RUN go build -o server

FROM alpine:latest
LABEL MAINTAINER="叶琼州" \
    EMAIL="yeqiongzhou@whu.edu.cn"

RUN apk add --no-cache  gettext tzdata   && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" >  /etc/timezone && \
    date && \
    apk del tzdata

WORKDIR /root
COPY --from=builder /go/src/paste.org.cn/paste/server ./
COPY config.yaml ./

EXPOSE 80

ENTRYPOINT ./server
