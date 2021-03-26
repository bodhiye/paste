FROM golang:alpine AS builder

ENV GOPROXY=https://goproxy.cn,https://goproxy.io,direct

WORKDIR /go/cache
COPY /server/go.mod .
COPY /server/go.sum .
RUN go mod download

COPY . /go/src/paste.org.cn/paste
RUN cd /go/src/paste.org.cn/paste/server && \
    go install /go/src/paste.org.cn/paste/server

FROM alpine:latest
LABEL MAINTAINER="叶琼州" \
    EMAIL="yeqiongzhou@whu.edu.cn"

RUN apk add --no-cache  gettext tzdata   && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" >  /etc/timezone && \
    date && \
    apk del tzdata

WORKDIR /root
COPY --from=builder /go/bin/server ./
COPY config.yaml ./

EXPOSE 80

ENTRYPOINT ./server
