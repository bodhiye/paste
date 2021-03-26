FROM golang:alpine AS builder

ENV GOPROXY=https://goproxy.cn,https://goproxy.io,direct \
    GO111MODULE=on

WORKDIR /go/cache
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . /go/src/paste.org.cn/paste
RUN cd /go/src/paste.org.cn/paste && \
    go install /go/src/paste.org.cn/paste

FROM ubuntu:20.04
LABEL MAINTAINER="叶琼州" \
    EMAIL="yeqiongzhou@whu.edu.cn"

RUN sed -i s:archive.ubuntu.com:mirrors.aliyun.com:g /etc/apt/sources.list \
    && sed -i s:security.ubuntu.com:mirrors.aliyun.com:g /etc/apt/sources.list

ENV DEBIAN_FRONTEND=noninteractive
RUN ln -fs /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN apt-get update \
    && apt install -y --no-install-recommends tzdata \
    && dpkg-reconfigure --frontend noninteractive tzdata \
    && apt-get clean

WORKDIR /root
COPY --from=builder /go/bin/paste /root/paste
COPY config.yaml /root/config.yaml
COPY docker-entrypoint.sh /root/docker-entrypoint.sh
RUN chmod +x /root/docker-entrypoint.sh

EXPOSE 80

ENTRYPOINT ["/root/docker-entrypoint.sh"]
