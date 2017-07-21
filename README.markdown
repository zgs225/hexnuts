HEXNUTS
===

配置服务 By Pure Golang. 根据配置生成配置文件。

### Usage

```
Usage:
        hexnuts command [arguments]

Commands:

        server  启动配置服务
        sync    同步配置文件
        monitor 启动监听程序
```

**注意**：监听程序只监听后缀为`.hexnuts`的文件

#### 启动服务

```
hexnuts server -tls -cert ./certs/cert.pem -key ./certs/key.pem -addr :5678 -monitor :5679 -dumps ./hexdumps.db
```

#### 启动同步监听

```
hexnuts monitor -server http://localhost:5678 -monitor.server localhost:5679 -in ./config -out ./config -tls
```

#### 手动同步

```
hexnuts sync -tls -server http://localhost:5678 -in ./config/config.js.hexnuts -out ./config/config.js
```

### API

#### 保存配置

```
POST /set
key=xxxxx
value=xxxxx
```

#### 更新配置

```
POST /update
key=xxxxx
value=xxxxx
```

#### 获取配置

```
GET /get?key=xxxxx
```

#### 删除配置

```
POST /del
key=xxxxx
```
