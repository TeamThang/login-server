#### 用户注册

`uri`: `/v1/account`

`method`: `post`

`request`: json

`说明`: LoginName，Mobile，Email不能和现有数据库中重复

arg1 | arg2 | type | desc
-- | -- | -- | --
BasicInfo | LoginName | str | 登录名，必填
BasicInfo | Password | str | 密码(6-20位)，必填
BasicInfo | Mobile | str | 手机，选填
BasicInfo | Email | str | 邮箱，选填
UserInfo | Name | str | 真实姓名,选填
UserInfo | Age | int | 年龄，选填
UserInfo | Birthday | time | 生日(格式:"2018-04-25 18:55:10")，选填
UserInfo | Address | str | 地址, 选填
UserInfo | Source | str | 用户来源, 选填

`request`
```json
{
    "BasicInfo": {  // BasicInfo必须上报
    	    "LoginName": "test1",
            "Password": "123",
            "Mobile": "1",
            "Email": "1@163.com"
            "Source": "quantity"
    },
    "UserInfo": {  // UserInfo可以不上报
    	"Name": "王霸",
    	"Age": 123,
    	"Birthday": "2018-04-25 18:55:10",
    	"Address": "太阳星星月亮"
    }
}
```
`response`:
```json
statusCode: 201  // CREATED
{
    "UserID": 52,  // 用户唯一ID
    "message": "", // 消息
    "status": 201  // 请求状态
}
```
#### 发送邮件验证码

`url`: `/v1/verify/email/user/{UserID}`

`method`: `get`

`request`: json

`备注`	: 需要鉴权token和UserID是否匹配

arg | type | desc
-- | -- | --
UserID | int | 用户id，必填

`request`
```json
http://127.0.0.1:8085//v1/verify/email/user/50
```
`response`:
```json
statusCode: 200 // [GET]
{
    "message": "",  // 消息
    "status": 200  // 请求状态
}
```

#### 验证邮件验证码

`url`: `/v1/verify/email/user/{UserID}/Code/{Code}`

`method`: `put`

`request`: json

`备注`: 需要鉴权token和UserID是否匹配

`params`:

arg | type | desc
-- | -- | --
UserID | int | 用户id，必填
verifyCode | str | 邮箱验证码，必填

`request`
```json
url: http://127.0.0.1:8085/v1/verify/email/user/48/Code/pLtkM
```
`response`:
statusCode: 201  // CREATED
```json
{
    "message": "",  // 消息
    "status": 201,  // 请求状态
    "result": true  // 验证结果
}
```

#### 用户销户

`uri`: `/v1/account/{UserID}`

`method`: `delete`

`request`: json

`备注`: 需要鉴权token和UserID是否匹配

`params`:

arg | type | desc
-- | -- | --
UserID | int | 用户id，必填

`request`
```json
http://127.0.0.1:8085/v1/account/49
```
`response`:
```json
statusCode: 204 // NO CONTENT - [DELETE]：用户删除数据成功。
```

#### 用户信息修改

`uri`: `/v1/account/{UserID}`

`method`: `put`

`request`: json

`备注`: 需要鉴权token和UserID是否匹配

`说明`: 按上报的参数修改当前值。LoginName，Mobile，Email不能和现有数据库中重复

`params`:

arg1 | arg2 | type | desc
-- | -- | -- | --
UserID | | int | 用户id，必填
BasicInfo | LoginName | str | 登录名，选填
BasicInfo | Mobile | str | 手机，选填
BasicInfo | Email | str | 邮箱，选填
UserInfo | Name | str | 真实姓名,选填
UserInfo | Age | int | 年龄，选填
UserInfo | Birthday | time | 生日(格式:"2018-04-25 18:55:10")，选填
UserInfo | Address | str | 地址, 选填

`request`
```json
{
    "BasicInfo": {  // BasicInfo必须上报
    	    "LoginName": "test1",
            "Password": "123",
            "Mobile": "1",
            "Email": "1@163.com"
    },
    "UserInfo": {  // UserInfo可以不上报
    	"Name": "王霸",
    	"Age": 123,
    	"Birthday": "2018-04-25 18:55:10",
    	"Address": "太阳星星月亮"
    }
}
```
`response`:
```json
statusCode: 201  // CREATED
{
    "message": "",  // 消息
    "status": 201  // 请求状态
}
```

#### 用户密码修改

`uri`: `/v1/pwd/user/{UserID}`

`method`: `put`

`request`: json

`备注`: 需要鉴权token和UserID是否匹配

`params`

arg | type | desc
-- | -- | -- | --
UserID | int | 用户id，必填
OldPW | str | 当前密码，必填
NewPW |  str | 新密码，必填


`request`
```json
url: http://127.0.0.1:8085/v1/pwd/user/50
{
	"OldPW": "123",
	"NewPW": "456"
}
```
`response`:
```json
statusCode: 201  // CREATED
{
    "message": "",  // 消息
    "status": 201  // 请求状态
}
```

#### 用户信息查询

`uri`: `/v1/account/{UserID}`

`method`: `get`

`request`: uri

`备注`: 需要鉴权token和UserID是否匹配

`params`:

arg | type | desc
-- | -- | -- | --
UserID | int | 用户id，必填

`request`
```
url: http://127.0.0.1:8085/v1/account/50
```
`response`:
```json
statusCode: 200 // [GET]
{
    "data": {
        "BasicInfo": {
            "ID": 47,
            "LoginName": "test1",
            "Mobile": "1",
            "Email": "1@163.com",
            "CreatedAt": "2018-04-26 15:14:14",  // 创建时间
            "UpdatedAt": "2018-04-26 15:14:40",  // 修改时间
            "DeletedAt": "",                     // 软删除时间
            "LoginTime": "0001-01-01 08:05:43",  // 本次登录时间
            "LastLoginTime": "0001-01-01 08:05:43",  // 上次登录时间
            "LoginCount": 0  // 登录次数
        },
        "UserInfo": {
            "Name": "王霸",
            "Age": 123,
            "Birthday": "2018-04-26 02:55:10",
            "Address": "太阳星星月亮",
            "Other": null
        }
    },
    "message": "",
    "status": 200
}
```

#### 账户登录

`url`: `/v1/login`

`method`: `post`

`request`: json

`说明1`: 登录顺序: LoginName->Mobile->Email;这三个至少上报一个

`说明2`: 鉴权顺序: Server->Name，"all"为所有权限，必须配置并赋权给用户

arg | type | desc
-- | -- | --
LoginName | int | 登录名，选填
Mobile | str | 登录手机，选填
Email | str | 登录邮箱，选填
Password | str | 登录密码，必填
Right | dict | 权限，Server:登录服务;Name:配置的权限名


`request`
```json
{
	"Password": "456",
	"LoginName": "test3",
	"Right":{"Server":"all","Name": "all"}
}
```
`response`:
```json
statusCode: 201  // CREATED
{
    "message": "login success",
    "status": 201
}
Cookie:
session_id: ae979a3b-3773-488a-b3c6-b4fb5e31e5c0  // 登录服务的token(uuid)
login_server: all  // 成功登录的服务
```

#### 账户注销

`url`: `/v1/logout`

`method`: `get`

`request`: json


arg | type | desc
-- | -- | --

`request`
```json
Cookie:
session_id=a1132ce8-757f-4c18-876b-423641c299b0
```
`response`:
```json
statusCode: 200 // [GET]
{
    "message": "注销成功",
    "status": 200
}
```

#### 获当前登录的用户信息

`url`: `/v1/get_user_info`

`method`: `get`

`request`: json


arg | type | desc
-- | -- | --

`request`
```json
Cookie:
session_id=ae979a3b-3773-488a-b3c6-b4fb5e31e5c0
```
`response`:
```json
statusCode: 200 // [GET]
{
    "data": {
        "BasicInfo": {
            "ID": 47,
            "LoginName": "test1",
            "Mobile": "1",
            "Email": "1@163.com",
            "Source": "",
            "CreatedAt": "2018-04-26 15:14:14",
            "UpdatedAt": "2018-05-01 11:03:16",
            "DeletedAt": "",
            "LoginTime": "2018-05-01 11:03:16",
            "LoginCount": 14
        },
        "UserInfo": {
            "Name": "王霸",
            "Age": 123,
            "Birthday": "2018-04-26 02:55:10",
            "Address": "太阳星星月亮",
            "Other": null
        }
    },
    "message": "",
    "status": 200
}
Cookie:
session_id: ae979a3b-3773-488a-b3c6-b4fb5e31e5c0  // 登录服务的token(uuid)
server: login  // 成功登录的服务
```