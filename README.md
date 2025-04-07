<div align=center>
<img src="web/public/logo.png" width=300" height="300" />
</div>
<div align=center>
	<img src="https://img.shields.io/badge/golang-1.15.6-blue"/>
	<img src="https://img.shields.io/badge/gin-1.6.3-lightBlue"/>
    <img src="https://img.shields.io/badge/node-14.17.3-green"/>
	<img src="https://img.shields.io/badge/vue-2.6.10-brightgreen"/>
</div>

# `paste`——代码便利贴，在线代码分享平台~

`我们不生产代码，我们只是代码的搬运工`

## 项目概述

`Paste` 是一个在线代码分享平台，允许用户上传代码片段并生成可分享的链接。平台支持多种编程语言，并提供如下特性：
- 自定义代码过期时间
- 密码保护
- 一次性阅读（阅后即焚）
该项目使用前后端分离架构，后端采用 `Go` 语言开发，前端基于 `Vue.js` 框架构建。数据存储使用 `MongoDB`。整个应用通过 `Docker` 容器化部署，使用 `docker-compose` 简化部署流程。

## 技术栈

**后端**  
- `Go 1.15.6` 及以上版本
- `Gin` (Web 框架)
- `MongoDB` (数据库)
- `logrus` (日志库)
- `viper` (配置管理)

**前端**
`Vue.js 2.6.10`
`Node.js 14.17.3`
`axios` (HTTP 客户端)
`vue-router` (前端路由)

**腾讯云 COS 生命周期管理**

系统利用腾讯云 COS 提供的生命周期管理功能，对上传的图片对象自动进行过期处理：

- **自动初始化**：应用启动时，自动根据配置设置 COS 存储桶的生命周期规则
- **目录分类管理**：根据不同的过期时间将文件存储在不同的目录，对每个目录应用独立的生命周期规则
- **精确过期控制**：根据用户选择的过期时间，自动匹配最合适的存储目录和生命周期规则
- **优势**：不消耗应用服务器资源，可靠性高，由云服务提供商保障

目录结构和生命周期规则示例：

```
/expires/1h/          - 1小时内过期的内容（保留1天）
/expires/1d/          - 1天内过期的内容（保留2天）
/expires/1w/          - 1周内过期的内容（保留9天）
/expires/1m/          - 1个月内过期的内容（保留37天）
/expires/1y/          - 1年内过期的内容（保留372天）
```

## 项目部署

**部署架构**
项目使用 `Docker Compose` 部署三个主要服务：
1. **mongo**: `MongoDB` 数据库服务，用于存储代码片段
2. **server**: `Go` 后端 `API` 服务，提供 `API` 接口
3. **web**: 前端 `Vue.js` 应用，通过 `Nginx` 提供静态资源服务

各服务通过内部网络相互通信，对外暴露的端口包括：
- `80/443` 端口：`Web` 前端服务 `（HTTP/HTTPS）`
- `8000` 端口：后端 `API` 服务
- `27017` 端口：`MongoDB` 数据库服务

**环境要求**：
- `Docker`
- `Docker Compose`

**部署步骤**
- 如果需要，可以修改 `paste/web/public/config.json` 配置文件
- 使用 `docker-compose` 方式来部署容器服务，执行 `docker-compose up -d`　执行一键部署服务。
- 如需开启百度统计，取消注释并替换 `paste/web/public/index.html` 中的百度统计脚本
- 如需开启 `https` 访问，需要先上传 `Nginx` 服务器类型 `SSL` 证书到 `paste/web/public/conf.d` 目录下，并修改 `nginx.conf` 配置文件。

`docker-compose up -d`命令执行过后，直接访问：[http://127.0.0.1:80](http://127.0.0.1:80) 。

## 后端设计

**目录结构**
```text
server/
├── db/               # 数据库操作相关代码
├── middleware/       # 中间件代码
├── proto/            # 协议定义
├── router/           # 路由配置
├── service/          # 业务逻辑实现
├── util/             # 工具函数
├── config.yaml       # 配置文件
├── Dockerfile        # 服务容器化配置
├── go.mod            # Go 模块依赖
├── go.sum            # Go 模块依赖校验
└── main.go           # 程序入口
```

### 核心模块

#### 数据模型（db）
数据库模块主要负责与 `MongoDB` 的交互，核心数据结构为 `PasteEntry`：
```go
type PasteEntry struct {
    Key       string    `json:"key" bson:"key"`             // 唯一标识
    Langtype  string    `json:"langtype" bson:"langtype"`   // 代码语言类型
    Content   string    `json:"content" bson:"content"`     // 代码内容
    Password  string    `json:"-" bson:"password,omitempty"` // 密码保护
    ClientIP  string    `json:"-" bson:"client_ip"`         // 客户端 IP
    Once      bool      `json:"-" bson:"once,omitempty"`    // 是否一次性阅读
    CreatedAt time.Time `json:"-" bson:"created_at"`        // 创建时间
    ExpireAt  time.Time `json:"-" bson:"expire_at,omitempty"` // 过期时间
}
```

#### 服务层（service）
服务层实现了具体的业务逻辑，包括以下主要功能：
- `PostPaste`: 创建普通代码分享
- `PostPasteOnce`: 创建一次性代码分享
- `GetPaste`: 获取代码内容

#### 路由层（router）
路由层定义了 `API` 接口路径映射：
- `POST /v1/paste`: 创建常规代码分享
- `POST /v1/paste/once`: 创建一次性代码分享
- `GET /v1/paste/:key`: 获取代码内容
- `ANY /health`: 健康检查接口

#### 中间件 (middleware)
- 中间件提供了如下功能：
- 日志记录
- 请求 ID 生成
- 异常恢复

#### 配置管理
系统通过 `config.yaml` 配置文件管理各种设置：
- 日志级别
- 服务器地址和端口
- `MongoDB` 连接信息  

## API接口

### 创建代码分享接口

| Method   | 接口           | 说明                             |
| :------- | :------------- | :------------------------------- |
| `POST` | /v1/paste      | 创建一个自定义过期时间的分享链接 |
| `POST` | /v1/paste/once | 创建一个一次性分享链接，阅后即焚 |

#### `POST /v1/paste | /v1/paste/once`

**`request`**

| 字段       | 类型   | 是否必选 | 说明                                 |
| :--------- | :----- | :------- | :----------------------------------- |
| langtype   | string | Yes      | 代码语言类型，支持常见的编程语言类型 |
| content    | string | Yes      | 分享的代码内容，最大支持十万个字符   |
| password   | string | No       | 代码文本密码，可选项                 |
| expireDate | int    | No       | 过期时间，单位秒，可选项             |

```http
POST /v1/paste | v1/paste/once HTTP/1.1
Content-Type: application/json

{
    "langtype": "golang",
    "content": "hello, paste.org.cn!",
    "password": "123456",
    "expireDate": 3600 // 一小时后过期
}
```

**`response`**

| 字段    | 类型   | 是否必选 | 说明                                    |
| :------ | :----- | :------- | :-------------------------------------- |
| code    | int    | Yes      | 201: 表示成功                           |
| key     | string | No       | 分享代码文本的key，可以用来访问代码内容 |
| message | string | No       | 错误描述信息                            |

```http
HTTP/1.1 200 OK
Content-Type: applicatoin/json

{
    "code": 201,
    "key": "abcd123456"
}
```

### 获取代码分享内容接口

#### `GET /v1/paste/:key?[password=]`

**`request`**

```http
GET /v1/paste/:abcd123456?password=123456 HTTP/1.1
```

**`response`**

| 字段     | 类型   | 是否必选 | 说明           |
| :------- | :----- | :------- | :------------- |
| code     | int    | Yes      | 200: 表示成功  |
| langtype | string | No       | 代码语言类型   |
| content  | string | No       | 分享的代码内容 |
| message  | string | No       | 错误描述信息   |

```http
HTTP/1.1 200 OK
Content-Type: application/json

{
    "code": 200,
    "langtype": "golang",
    "content": "hello, paste.org.cn!"
}
```

## TODO

- [ ] 支持多段代码上传，并添加标题描述信息
- [ ] 支持代码截图上传 

## 感谢

web 前端参考了开源项目 [gin-vue-admin](https://github.com/flipped-aurora/gin-vue-admin)

## 免责声明

- 本平台只提供代码文本分享，与分享内容均没有任何联系。
- 本平台不对代码数据的存储负责。
