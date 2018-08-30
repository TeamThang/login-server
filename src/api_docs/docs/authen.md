#### 添加自定义权限

`uri`: `/v1/right`

`method`: `post`

`request`: json

`备注`: 需要鉴权token权限

`params`:

arg | type | desc
-- | -- | -- | --
Server | str | 支持的服务，必填
Name | str | 服务对应权限，必填
Desc| str | 权限描述, 选填

```json
{
	"Server":"all",
	"Name":"all",
	"Desc":"所有权限"
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

#### 删除自定义权限

`uri`: `/v1/right/{RightID}`

`method`: `post`

`request`: json

`备注`: 需要鉴权token权限

`params`:

arg | type | desc
-- | -- | -- | --
RightID | str | 权限表ID，必填

`request`
```json

```
`response`:
```json
statusCode: 204 // NO CONTENT - [DELETE]：用户删除数据成功。
```

#### 修改自定义权限

`uri`: `/v1/right/{RightID}`

`method`: `put`

`request`: json

`备注`: 需要鉴权token权限

`params`:

arg | type | desc
-- | -- | -- | --
RightID | str | 权限表ID，必填
Server | str | 支持的服务，选填
Name | str | 服务对应权限，选填
Desc| str | 权限描述, 选填

`request`
```json
url: http://127.0.0.1:8085/v1/right/6
{
	"Server":"all",
	"Name":"all",
	"Desc":"所有的权限"
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

#### 查询所有自定义权限

`uri`: `/v1/right`

`method`: `get`

`request`: url

`params`: null

`备注`: 需要鉴权token权限

`request`
```
url: http://127.0.0.1:8085/v1/right
```
`response`:
```json
statusCode: 200 // [GET]
{
    "data": [
        {
            "ID": 2,
            "Server": "all",
            "Name": "all",
            "Desc": "所有的权限",
            "UserRightRelations": null
        },
        {
            "ID": 3,
            "Server": "login",
            "Name": "all",
            "Desc": "登录服务权限",
            "UserRightRelations": null
        }
    ],
    "message": "",  // 消息
    "status": 200  // 请求状态
}
```

#### 用户绑定权限

`uri`: `/v1/bind_auth`

`method`: `post`

`request`: json

`备注`: 需要鉴权token和UserID是否匹配

`params`:

arg | type | desc
-- | -- | -- | --
UserID  | str | 用户id，必填
RightID | str | 权限ID，必填


`request`
```
{
	"UserID": 50,
	"RightID": 2
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

#### 用户解除绑定权限

`uri`: `/v1/bind_auth/user/{UserID}/right/{RightID}`

`method`: `delete`

`request`: json

`备注`: 需要鉴权token和UserID是否匹配

`params`:

arg | type | desc
-- | -- | -- | --
UserID  | str | 用户id，必填
RightID | str | 权限ID，必填


`request`
```
http://127.0.0.1:8085/v1/bind_auth/user/50/right/3
```
`response`:
```json
statusCode: 204 // NO CONTENT - [DELETE]：用户删除数据成功。
```

### 用户查询绑定权限

`uri`: `/v1/bind_auth/user/{UserID}`

`method`: `post`

`request`: json

`备注`: 需要鉴权token和UserID是否匹配

`params`:

arg | type | desc
-- | -- | -- | --
UserID  | str | 用户id，必填


`request`
```
url: http://127.0.0.1:8085/v1/bind_auth/user/50
```
`response`:
```json
statusCode: 200 // [GET]
{
    "data": {
        "UserID": 50,
        "Rights": [
            {
                "ID": 2,
                "Server": "all",
                "Name": "all",
                "Desc": "所有的权限",
                "UserRightRelations": null
            },
            {
                "ID": 3,
                "Server": "login",
                "Name": "all",
                "Desc": "登录服务权限",
                "UserRightRelations": null
            }
        ]
    },
    "message": "",
    "status": 200
}
```