FROM node:14.17.3-alpine AS builder

WORKDIR /web
COPY web/ .

RUN npm install -g cnpm --registry=https://registry.npm.taobao.org
RUN cnpm install || npm install
RUN npm run build

FROM nginx:alpine
LABEL MAINTAINER="叶琼州" \
    EMAIL="yeqiongzhou@whu.edu.cn"

ENV TZ=Asia/Shanghai

COPY --from=builder /web/paste /usr/share/nginx/html
COPY web/public/conf.d/nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80 443
