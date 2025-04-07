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

ENV TZ=Asia/Shanghai

WORKDIR /root
COPY --from=builder /go/bin/server ./
COPY /server/config.yaml /root/config.yaml

# 创建日志目录并设置权限
RUN mkdir -p /app/logs && \
    chmod 755 /app/logs

EXPOSE 8000

ENTRYPOINT ./server
