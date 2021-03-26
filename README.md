# paste——便利贴，在线信息分享平台~

`我们不生产内容，我们只是内容的搬运工`

## 部署

### 本地部署

- `docker run -v /root/data:/data/db -p 27017:27017 -d mongo:latest`
- `docker build -f Dockerfile -t paste:latest .`
- `docker run -p 80:80 --network host -d paste:latest`

　　由于 Mac 环境的 Docker Engine 无法访问容器所在的网络，可以把配置文件的 mgo host 改为 `mongodb://host.docker.internal:27017`。

### Docker-Compose 部署

　　推荐使用这种方式来部署服务，`docker-compose up -d`　执行一键部署服务。
