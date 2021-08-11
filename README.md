<div align=center>
<img src="https://file.paste.org.cn/logo.png" width=300" height="300" />
</div>
<div align=center>
	<img src="https://img.shields.io/badge/golang-1.15.6-blue"/>
	<img src="https://img.shields.io/badge/gin-1.6.3-lightBlue"/>
    <img src="https://img.shields.io/badge/node-14.17.3-green"/>
	<img src="https://img.shields.io/badge/vue-2.6.10-brightgreen"/>
</div>

# `paste`——代码便利贴，在线代码分享平台~

`我们不生产代码，我们只是代码的搬运工`

## 部署

先修改 `paste/web/public/config.json` 配置文件，之后使用 docker-compose 方式来部署容器服务，输入 `docker-compose up -d`　执行一键部署服务。如需开启百度统计，取消注释并替换 `paste/web/public/index.html` 中的百度统计脚本。

## TODO

- [ ] 支持多段代码上传，并添加标题描述信息
- [ ] 支持代码截图上传

# API

## 创建分享接口

|Method|接口|说明|
| :--- | :--- | :--- |
| `POST` |/v1/paste|创建一个自定义过期时间的分享链接|
| `POST` |/v1/paste/once|创建一个一次性分享链接，阅后即焚|

### `POST /v1/paste | /v1/paste/once`

**`request`**

|字段|类型|是否必选|说明|
| :--- | :--- | :--- | :--- |
|langtype|string|Yes|代码语言类型，支持常见的编程语言类型|
|content|string|Yes|分享的代码内容，最大支持十万个字符|
|password|string|No|代码文本密码，可选项|
|expireDate|int|No|过期时间，单位秒，可选项|

``` http
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

|字段|类型|是否必选|说明|
| :--- | :--- | :--- | :--- |
|code|int|Yes|201: 表示成功|
|key|string|No|分享代码文本的key，可以用来访问代码内容|
|message|string|No|错误描述信息|

``` http
HTTP/1.1 200 OK
Content-Type: applicatoin/json

{
    "code": 201,
    "key": "abcd123456"
}
```

## 获取分享内容接口

### `GET /v1/paste/:key?[password=]`

**`request`**

``` http
GET /v1/paste/:abcd123456?password=123456 HTTP/1.1
```

**`response`**

|字段|类型|是否必选|说明|
| :--- | :--- | :--- | :--- |
|code|int|Yes|200: 表示成功|
|langtype|string|No|代码语言类型|
|content|string|No|分享的代码内容|
|message|string|No|错误描述信息|

``` http
HTTP/1.1 200 OK
Content-Type: application/json

{
    "code": 200,
    "langtype": "golang",
    "content": "hello, paste.org.cn!"
}
```

# 感谢

web 前端参考了开源项目 [gin-vue-admin](https://github.com/flipped-aurora/gin-vue-admin) & PasteMeFrontend

# 免责声明

- 本平台只提供代码文本分享，与分享内容均没有任何联系。
- 本平台不对代码数据的存储负责。
