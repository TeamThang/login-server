#### 添加提醒邮箱

`uri`: `/v1/notify/email`

`method`: `post`

`request`: json

`备注`: 需要鉴权token和UserID是否匹配

`params`:

arg | type | desc
-- | -- | --
UserID | int | 用户ID，必填
Email | str | 提醒邮箱，必填
Subscribed | bool | 是否订阅，选填，默认:false，订阅之后会收到提醒


`request`
```json
{
	"UserID": 48,
	"Email": "zhaoyi.yuan@bitmain.com",
	"Subscribed": true
}
```
`response`:
```json
statusCode: 201  // CREATED
{
    "message": "",
    "status": 201
}
```

#### 删除提醒邮箱

`uri`: `/v1/notify/email/{EmailID}`

`method`: `delete`

`request`: json

`备注`: 需要鉴权token和UserID是否匹配

`params`:

arg | type | desc
-- | -- | --
EmailID | int | 提醒邮箱ID，必填

`request`
```json
url: http://127.0.0.1:8085/v1/notify/email/5
```
`response`:
```json
statusCode: 204 // NO CONTENT - [DELETE]：用户删除数据成功。
```

#### 提醒邮箱订阅

`uri`: `/v1/notify/email/{EmailID}}/user/{UserID}/sub`

`method`: `post`

`request`: json

`备注`: 需要鉴权token和UserID是否匹配

`params`:

arg | type | desc
-- | -- | --
EmailID | int | 提醒邮箱ID，必填

`request`
```json
url: http://127.0.0.1:8085/v1/notify/email/6/user/50/sub
```
`response`:
```json
statusCode: 201  // CREATED
{
    "message": "",  // 消息
    "status": 201   // 请求状态
}
```

#### 提醒邮箱取消订阅

`uri`: `/v1/notify/email/{EmailID}}/user/{UserID}/unsub`

`method`: `post`

`request`: json

`备注`: 需要鉴权token和UserID是否匹配

`params`:

arg | type | desc
-- | -- | --
EmailID | int | 提醒邮箱ID，必填

`request`
```json
url: http://127.0.0.1:8085/v1/notify/email/6/user/50/unsub
```
`response`:
```json
statusCode: 201  // CREATED
{
    "message": "",  // 消息
    "status": 201   // 请求状态
}
```

#### 查询提醒邮件

`uri`: `/v1/notify/email/user/{UserID}`

`method`: `post`

`request`: json

`备注`: 需要鉴权token和UserID是否匹配

`params`:

arg | type | desc
-- | -- | --
EmailID | int | 提醒邮箱ID，必填

`request`
```json
url: http://127.0.0.1:8085/v1/notify/email/user/50
```
`response`:
statusCode: 200 // [GET]
```json
{
    "data": [
        {
            "ID": 4,  // 邮箱ID
            "UserID": 48,  // 用户ID
            "Email": "zhaoyi.yuan@bitmain.com",  // 提醒邮箱
            "Subscribed": false   // 订阅状态
        },
        {
            "ID": 1,
            "UserID": 48,
            "Email": "yuan.zhaoyi@163.com",
            "Subscribed": false
        }
    ],
    "message": "",
    "status": 200
}
```

#### 发送提醒邮件

`uri`: `/v1/notify/email/user/{UserID}`

`method`: `post`

`request`: json

`备注`: 需要鉴权token和UserID是否匹配

`params`:

arg | type | desc
-- | -- | --
UserID | int | 用户ID，必填
Subject | str | 邮件主题，必填
ContentType | str | 邮件内容格式，必填
Content | str | 邮件，必填

`request`
```json
url: http://127.0.0.1:8085/v1/notify/email/user/50
{
	"Subject": "test",
	"ContentType": "text",
	"Content": "abcdefg"
}
```

`response`:
```json
statusCode: 201  // CREATED
{
    "message": "",  // 消息
    "status": 201   // 请求状态
}
```