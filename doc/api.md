# Paste API

## 创建分享接口

|Method|接口|说明|
| :--- | :--- | :--- |
| `POST` |/v1/paste|创建一个永久分享链接|
| `POST` |/v1/paste/once|创建一个一次性分享链接，阅后即焚|

### `POST /v1/paste | /v1/paste/once`

**`request`**

|字段|类型|是否必选|说明|
| :--- | :--- | :--- | :--- |
|langtype|string|Yes|文本语言类型，支持常见的编程语言类型|
|content|string|Yes|分享的文本内容，最大支持一万个字符|
|password|string|No|文本密码，可选项|
|expireDate|Date|No|过期时间，可选项|

``` http
POST /v1/paste | v1/paste/once HTTP/1.1
Content-Type: application/json

{
    "langtype": "golang",
    "content": "hello, paste.org.cn!",
    "password": "123456"
}
```

**`response`**

|字段|类型|是否必选|说明|
| :--- | :--- | :--- | :--- |
|code|int|Yes|200: 表示正确|
|key|string|No|分享文本的key，可以用来访问文本内容|
|message|string|No|错误描述信息|

``` http
HTTP/1.1 200 OK
Content-Type: applicatoin/json

{
    "code": 200,
    "key": "test"
}
```

## 获取分享内容接口

### `GET /v1/paste/:key?[password=]`

**`request`**

``` http
GET /v1/paste/:test?password=123456 HTTP/1.1
```

**`response`**

|字段|类型|是否必选|说明|
| :--- | :--- | :--- | :--- |
|code|int|Yes|200: 表示正确|
|langtype|string|No|文本语言类型|
|content|string|No|分享的文本内容|
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
