
### 依赖包
```
go get github.com/gin-gonic/gin
go get gopkg.in/mgo.v2
go get -u github.com/jinzhu/gorm
go get github.com/lib/pq
go get github.com/apsdehal/go-logger
go get -u github.com/satori/go.uuid
```

### 目录结构

```shell
├── bin           # go编译目录
├── src  # 代码目录
    ├── github.com        # 第三方包
    ├── server            # 服务代码，handler和web服务启动
    │   ├── account       # 用户账户管理
    │   ├── authen        # 权限管理
    │   ├── login         # 登录接口
    │   ├── notify        # 提醒模块，当前只支持邮件提醒
    │   ├── verify        # 账户验证模块
    │   ├── util          # 公共模块
    │   ├── conf          # 配置文件，启动时根据配置文件初始化并赋值
    │   ├── log           # 日志模块
    └── ├── db            # 数据库配置
        ├── base
        │   ├── middleware  # 接口鉴权中间件
        │   ├── router      # 接口路由，Restful风格
        │   └── base        # handler公用函数
        └── main            # 服务启动主目录，需要定义配置文件路径
```

### 配置文件结构
配置文件目录下需包含:
- 基本文件: config.json

```json
{
  "Redis": {
    "Host": "127.0.0.1",
    "Port": 6379,
    "PassWord": "",
    "DB": 2
  },
  "Postgre": {
    "Host": "",
    "Port": 9969,
    "User": "",
    "DbName": "",
    "PassWord": ""
  },
  "Email": {
    "Host": "smtp.qq.com",
    "Port": 465,
    "User": "",
    "PassWord": "",
    "Sender": ""
  }
}
```

- 启动配置： server.json

```json
{
  "LogLevel": "debug",
  "LogPath": "",
  "TCPAddr": "127.0.0.1:8085",
  "MaxConnNum": 20000
}
```

### 接口文档:
[http://47.98.55.223:8095/](http://47.98.55.223:8095/)