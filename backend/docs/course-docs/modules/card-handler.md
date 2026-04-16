# card_handler HTTP 接口模块

## 作用
- 将 HTTP 请求转换为 CardService 调用，将结果序列化为 JSON 响应
- 提供学籍验证直通接口（ValidateStudent）
- 提供按证件号查当前有效卡接口（GetCardByIDNumber，直接调 CardRepository）

## 职责边界
- 负责：路径/查询参数解析、请求体绑定、错误码→HTTP 状态码映射、响应序列化
- 不负责：业务规则校验（交给 CardService）、SQL 查询（交给 Repository）

## 注入依赖
- CardService：饭卡核心业务
- StudentValidator：学籍验证接口（用于 ValidateStudent handler）
- CardRepository：直接查按证件号查当前有效卡（跳过 service 层，因为没有附加业务逻辑）

## 接口列表

| Handler | HTTP | 说明 |
|---|---|---|
| ValidateStudent | GET /api/validate-student?idNumber= | 直接调 validator，返回 StudentValidationResult |
| GetCardByIDNumber | GET /api/cards?idNumber= | 调 repo.FindCurrentCardByIDNumber，返回 CardDetail |
| IssueCard | POST /api/cards | 读 idNumber + preDeposit，调 cardSvc.IssueCard |
| GetCard | GET /api/cards/:cardNo | 按 16 位卡号查卡，返回 CardDetail |
| Deposit | POST /api/cards/:cardNo/deposits | 充值，返回收据 |
| CreateTransaction | POST /api/cards/:cardNo/transactions | 消费扣款 |
| ReportLoss | PUT /api/cards/:cardNo/loss-report | 挂失 |
| CancelLossReport | DELETE /api/cards/:cardNo/loss-report | 取消挂失 |
| CancelCard | POST /api/cards/:cardNo/cancellation | 注销 |

## 错误码→HTTP 状态码映射

| 错误码 | HTTP 状态 |
|---|---|
| CARD_NOT_FOUND / WINDOW_NOT_FOUND / STUDENT_NOT_FOUND | 404 |
| CARD_ALREADY_ACTIVE / CARD_NOT_ACTIVE / CARD_NOT_LOST / CARD_CANCELLED / CARD_LOST / CARD_ALREADY_CANCELLED / INSUFFICIENT_BALANCE | 409 |
| INVALID_AMOUNT / VALIDATION_ERROR | 400 |
| STUDENT_SERVICE_ERROR | 502 |
| 其他 | 500 |

## 响应字段规范（cardNo 替换旧 id）
- cardToJSON 输出中使用 `cardNo`（16 位字符串），不再暴露数据库自增 `id`
- Deposit 响应字段 `cardNo`（旧为 `cardId`）
- Transaction 响应字段 `cardNo`（旧为 `cardId`）
- IssueCard refund 字段使用 `oldCardNo`（旧为 `oldCardId`）

## 关键实现点
- GetCardByIDNumber 直接注入 CardRepository，绕过 service，因为该查询无业务规则，只是数据读取
- ValidateStudent 直接注入 StudentValidator，可在发卡前独立调用，供前端二次确认展示
- 路径参数统一为 `:cardNo`，不再使用 `:id`
- 日志策略：各写操作入口打 Info（方法/路径/关键参数）；`handleError` 集中打日志——BizError 打 Warn（含 code），系统错误打 Error；service 层关键操作成功后打 Info（含 cardNo、金额等字段）。日志不分散在每个错误分支，统一在 handleError 输出
