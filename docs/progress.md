# 迭代进度

## 当前迭代目标

前后端联调与功能验收。

## 已完成

- 前端项目骨架（React + Vite + react-router-dom），全部 7 个页面实现完成，构建无报错
- PRD 文档（`docs/prd.md`），覆盖 6 项功能的完整行为和约束
- 数据库表设计：5 张表的 GORM 模型（`backend/model/`）及文档（`backend/docs/database/`）
- 架构文档（`docs/architecture.md`）
- 修复 Card 模型：CardHolderID 改为必填，不再重用旧卡号
- AGENTS.md 加入文档维护规范（PRD + architecture + progress 每次对话必读）
- OpenAPI 契约（`docs/api/openapi.yaml`），覆盖 16 个接口（6 项核心业务 + 7 项统计 + 窗口管理）

## 进行中

（无）

## 待办

- 前后端联调

## 阻塞

（无）

## 下一步

前后端联调，启动后端服务后通过前端页面进行功能验收。

## 变更记录

### 2026-04-15 第 1 轮：项目初始化与数据库设计

新增文件：
- `backend/model/card_holder.go` — CardHolder 模型
- `backend/model/card.go` — Card 模型（含 CardStatus 枚举）
- `backend/model/deposit_record.go` — DepositRecord 模型
- `backend/model/transaction.go` — Transaction 模型
- `backend/model/window.go` — Window 模型
- `backend/docs/database/card_holders.md`
- `backend/docs/database/cards.md`
- `backend/docs/database/deposit_records.md`
- `backend/docs/database/transactions.md`
- `backend/docs/database/windows.md`
- `docs/prd.md` — 产品需求文档
- `docs/architecture.md` — 架构与关键约束
- `docs/progress.md` — 本文件

修改文件：
- `AGENTS.md` — 新增「文档维护」和「重要信息」小节
- `backend/model/card.go` — CardHolderID 从 `*uint`(nullable) 改为 `uint`(not null)，不再支持卡片重用
- `backend/docs/database/cards.md` — 同步去掉卡片重用设计，状态流转图移除 cancelled→active
- `docs/architecture.md` — 同步去掉卡片重用、补全三重校验和完整数据流

关键决策：
- 每次发卡创建新记录新编号，注销的卡保留为历史，不重用
- 金额统一用 int64 存分
- 就餐消费三重校验：卡号存在 → 未注销 → 未挂失

### 2026-04-15 第 2 轮：OpenAPI 契约设计

新增文件：
- `docs/api/openapi.yaml` — 完整 OpenAPI 3.0.3 契约

接口清单（16 个操作）：
- `POST /api/cards` — 发卡（含旧卡自动注销退款）
- `GET /api/cards/{id}` — 查卡（就餐刷卡第一步）
- `POST /api/cards/{id}/deposits` — 存款（返回收据）
- `POST /api/cards/{id}/transactions` — 就餐消费结算
- `PUT /api/cards/{id}/loss-report` — 挂失
- `DELETE /api/cards/{id}/loss-report` — 取消挂失
- `POST /api/cards/{id}/cancellation` — 注销（返回退款明细）
- `GET /api/statistics/meal-revenue` — 售饭总收入
- `GET /api/statistics/window-revenue` — 各窗口收入
- `GET /api/statistics/deposit-details` — 各持卡人存款明细
- `GET /api/statistics/deposit-summary` — 本日/月存款金额
- `GET /api/statistics/active-balance` — 流动资金总额
- `GET /api/statistics/daily-report` — 日餐报表
- `GET /api/statistics/yearly-report` — 年餐报表
- `GET /api/windows` — 窗口列表
- `POST /api/windows` — 创建窗口

关键决策：
- 统计接口用多个子路径（每个端点固定响应结构），不用单接口 + type 参数
- 就餐消费分两步：GET 查卡展示余额 → POST 结算扣款
- 成功统一返回 200（不用 201），透明响应不加包装层
- 业务错误返回 4xx + {code, message}，其他意外错误统一 500 + INTERNAL_ERROR
- JSON 字段名与 Go 模型 json tag 完全一致（camelCase）

### 2026-04-15 第 4 轮：后端三层架构实现与单元测试

新增文件：
- `backend/db/init.go` — GORM + SQLite 数据库初始化，AutoMigrate 5 张表
- `backend/repository/card_repository.go` — Card/CardHolder/DepositRecord/Transaction CRUD 及统计查询
- `backend/repository/window_repository.go` — Window CRUD
- `backend/service/card_service.go` — 发卡、存款、就餐消费、挂失/取消挂失、注销核心业务逻辑
- `backend/service/statistics_service.go` — 7 项统计业务逻辑
- `backend/service/window_service.go` — 窗口管理业务逻辑
- `backend/handler/card_handler.go` — 7 个饭卡相关 HTTP handler
- `backend/handler/statistics_handler.go` — 7 个统计 HTTP handler
- `backend/handler/window_handler.go` — 2 个窗口 HTTP handler
- `backend/router/router.go` — 路由注册，CORS 允许所有来源
- `backend/main.go` — 程序入口，监听 :8080
- `backend/service/card_service_test.go` — CardService 表驱动单元测试
- `backend/service/statistics_service_test.go` — StatisticsService 单元测试

关键决策：
- 业务错误统一用 BizError 结构体，handler 层通过 bizErrStatus 映射到对应 HTTP 状态码
- 金额全部 int64（分），绝无 float
- 就餐三重校验顺序：卡存在 → 非 cancelled → 非 lost → 余额充足
- 测试使用 SQLite 内存数据库（":memory:"），不依赖文件

新增文件：
- `frontend/package.json` — Vite + React + react-router-dom 项目配置
- `frontend/vite.config.js` — Vite 配置
- `frontend/index.html` — 入口 HTML
- `frontend/src/main.jsx` — React 入口
- `frontend/src/api.js` — 封装全部 16 个 API（base URL: http://localhost:8080）
- `frontend/src/App.jsx` — 顶部导航 + React Router 路由配置
- `frontend/src/pages/IssuePage.jsx` — 发卡页（POST /api/cards）
- `frontend/src/pages/DepositPage.jsx` — 存款页（含收据展示）
- `frontend/src/pages/MealPage.jsx` — 就餐消费页（三重校验 + 报警 + 结算）
- `frontend/src/pages/LossPage.jsx` — 挂失管理页（挂失/取消挂失）
- `frontend/src/pages/CancelPage.jsx` — 注销页（二次确认 + 退款明细）
- `frontend/src/pages/StatisticsPage.jsx` — 汇总统计页（全部 7 个统计区块）
- `frontend/src/pages/WindowsPage.jsx` — 窗口管理页（列表 + 新建）

关键决策：
- 金额在前端全部转换：显示时 ÷100 保留 2 位小数，提交时 Math.round(value * 100)
- 就餐页卡片异常（cancelled/lost/404）显示红色醒目报警提示
- 注销操作设置二次确认防止误操作
- 统计页各区块独立发请求，避免一次加载过慢
