# 后端

Go + Echo + GORM + SQLite

## 环境要求

- Go 1.26+

## 启动

```bash
# 安装依赖
go mod download

# 启动服务（监听 :8080，首次运行自动创建 meal_card.db）
go run main.go
```

## 测试

```bash
go test ./...
```

## 项目结构

```
backend/
├── main.go              # 入口
├── db/init.go           # 数据库初始化
├── router/router.go     # 路由注册
├── handler/             # HTTP 层
├── service/             # 业务逻辑层（含单元测试）
├── repository/          # 数据访问层
└── model/               # GORM 数据模型
```

数据库文件 `meal_card.db` 在启动时自动创建于当前目录，删除后重启即可重置数据。
