FROM nginx:alpine
LABEL MAINTAINER="叶琼州" \
    EMAIL="yeqiongzhou@whu.edu.cn"

ENV TZ=Asia/Shanghai

COPY web/public/conf.d/nginx.conf /etc/nginx/
COPY web/paste/ /usr/share/nginx/html/

EXPOSE 80
