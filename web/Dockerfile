FROM node:16-alpine AS builder

WORKDIR /web
COPY web/ .

RUN npm config set registry http://registry.npm.taobao.org
RUN npm install --legacy-peer-deps
RUN npm run build

FROM nginx:alpine
LABEL MAINTAINER="叶琼州" \
    EMAIL="yeqiongzhou@whu.edu.cn"

ENV TZ=Asia/Shanghai

COPY --from=builder /web/paste /usr/share/nginx/html
COPY web/public/conf.d/nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80 443