# 服务端部署

- `docker run -v /mongodb/data:/data/db -p 27017:27017 -d mongo:latest`
- `docker build -f Dockerfile -t paste:latest .`
- `docker run -p 80:80 --network host -d paste:latest`

　　由于 Mac 环境的 Docker Engine 无法访问容器所在的网络，可以把配置文件的 mgo host 改为 `mongodb://host.docker.internal:27017`。
