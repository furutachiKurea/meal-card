# 三层架构模块职责划分

## 作用
- 将 HTTP 处理、业务逻辑、数据访问三类关注点分层隔离
- 每层只与相邻层通信，禁止跨层直接调用

## 职责边界

### Handler 层（handler/）
- 负责：HTTP 请求接收、参数绑定与格式校验、调用 Service、将结果序列化为 JSON、将 BizError 映射为对应 HTTP 状态码
- 不负责：任何业务规则判断、数据库操作

文件：
- `card_handler.go`：发卡、查卡、存款、消费、挂失、取消挂失、注销
- `statistics_handler.go`：7 项统计接口
- `window_handler.go`：窗口列表、创建窗口

统一错误处理：`handleError` 通过 `errors.As` 识别 `BizError`，调用 `bizErrStatus` 映射到 404/409/400，非 BizError 统一返回 500。

---

### Service 层（service/）
- 负责：业务规则校验、状态流转、跨实体协调、定义业务错误码
- 不负责：HTTP 协议细节、直接执行 SQL

文件：
- `card_service.go`：所有卡片生命周期业务
- `statistics_service.go`：7 项统计逻辑
- `window_service.go`：窗口管理

`BizError` 结构体定义在此层：
```go
type BizError struct {
    Code    string
    Message string
}
```

---

### Repository 层（repository/）
- 负责：GORM CRUD 封装、统计聚合 SQL
- 不负责：业务规则、状态流转

文件：
- `card_repository.go`：卡片、持卡人、存款记录、消费记录的增删改查及聚合统计
- `window_repository.go`：窗口增删改查

---

## 错误码与 HTTP 状态码映射

| 错误码 | HTTP 状态码 | 说明 |
|--------|------------|------|
| CARD_NOT_FOUND | 404 | 卡不存在 |
| WINDOW_NOT_FOUND | 404 | 窗口不存在 |
| CARD_ALREADY_ACTIVE | 409 | 该证件号已有 active 卡 |
| CARD_NOT_ACTIVE | 409 | 卡不是 active 状态 |
| CARD_NOT_LOST | 409 | 卡不是 lost 状态 |
| CARD_CANCELLED | 409 | 卡已注销 |
| CARD_LOST | 409 | 卡已挂失 |
| CARD_ALREADY_CANCELLED | 409 | 卡已经是 cancelled 状态 |
| INSUFFICIENT_BALANCE | 409 | 余额不足 |
| INVALID_AMOUNT | 400 | 金额非法 |
| VALIDATION_ERROR | 400 | 参数校验失败 |
| （其他） | 500 | 服务端内部错误 |

## 路由结构（16 个接口，前缀 /api）

| 方法 | 路径 | Handler |
|------|------|---------|
| POST | /cards | CardHandler.IssueCard |
| GET | /cards/:id | CardHandler.GetCard |
| POST | /cards/:id/deposits | CardHandler.Deposit |
| POST | /cards/:id/transactions | CardHandler.CreateTransaction |
| PUT | /cards/:id/loss-report | CardHandler.ReportLoss |
| DELETE | /cards/:id/loss-report | CardHandler.CancelLossReport |
| POST | /cards/:id/cancellation | CardHandler.CancelCard |
| GET | /statistics/meal-revenue | StatisticsHandler.GetMealRevenue |
| GET | /statistics/window-revenue | StatisticsHandler.GetWindowRevenue |
| GET | /statistics/deposit-details | StatisticsHandler.GetDepositDetails |
| GET | /statistics/deposit-summary | StatisticsHandler.GetDepositSummary |
| GET | /statistics/active-balance | StatisticsHandler.GetActiveBalance |
| GET | /statistics/daily-report | StatisticsHandler.GetDailyReport |
| GET | /statistics/yearly-report | StatisticsHandler.GetYearlyReport |
| GET | /windows | WindowHandler.ListWindows |
| POST | /windows | WindowHandler.CreateWindow |

中间件：zerolog 请求日志、CORS（允许所有来源，允许 GET/POST/PUT/DELETE/OPTIONS）
