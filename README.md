# paste——便利贴，在线信息分享平台~

`我们不生产内容，我们只是内容的搬运工`

## 本地部署

- `docker run -v /root/data:/data/db -p 27017:27017 -d mongo:latest`
- `docker build -f Dockerfile -t paste:latest .`
- `docker run -p 80:80 --network host -d paste:latest`
