# 迭代进度

## 当前迭代目标

重新对齐设计，完成后端重构与前后端联调。

## 已完成

- 前端项目骨架（React + Vite + react-router-dom），全部 7 个页面实现完成，构建无报错
- PRD 文档（`docs/prd.md`），覆盖 6 项功能的完整行为和约束
- 数据库表设计：5 张表的 GORM 模型（`backend/model/`）及文档（`backend/docs/database/`）
- 架构文档（`docs/architecture.md`）
- 修复 Card 模型：CardHolderID 改为必填，不再重用旧卡号
- AGENTS.md 加入文档维护规范（PRD + architecture + progress 每次对话必读）
- OpenAPI 契约（`docs/api/openapi.yaml`），覆盖 16 个接口（6 项核心业务 + 7 项统计 + 窗口管理）

## 进行中

- 后端重构：card_no、StudentValidator 接口、mock 服务、repo/service/handler/router 全量适配（Task #1 进行中）

## 待办

- 后端：repository 新增 FindCardByCardNo、FindCurrentCardByIDNumber 方法
- 后端：CardService 重构（注入 StudentValidator，IssueCard 改签名，所有方法 cardID→cardNo）
- 后端：CardHandler 重构（路径参数 :id→:cardNo，新增 ValidateStudent、GetCardByIDNumber handler）
- 后端：router 更新路由
- 后端：main.go 接入 HttpStudentValidator
- 后端：单元测试更新（IssueCardRequest 去掉 Name/Deposit，改用 FakeStudentValidator）
- 前后端联调，启动后端服务后通过前端页面进行功能验收

## 阻塞

（无）

## 下一步

前后端联调，启动后端服务后通过前端页面进行功能验收。

## 变更记录

### 2026-04-16 第 9 轮：前端 cardNo 全量适配

修改文件：
- `frontend/src/api.js` — 新增 `validateStudent`、`getCardByIDNumber`；将所有路径参数从 `id` 改为 `cardNo`；`issueCard` 请求体去掉 `name`/`deposit`，只传 `idNumber`/`preDeposit`；`createTransaction` 改名签名对齐
- `frontend/src/pages/IssuePage.jsx` — 完全重写为两步流程（Step 1 验证证件号调 `validateStudent`；Step 2 录入预存款调 `issueCard`）；移除姓名/押金输入框
- `frontend/src/pages/LossPage.jsx` — 查询入口由卡号改为证件号，调 `getCardByIDNumber`；操作时使用 `cardInfo.cardNo` 作路径参数
- `frontend/src/pages/CancelPage.jsx` — 查询入口由卡号改为证件号，调 `getCardByIDNumber`；退款收据卡号字段改为 `result.card.cardNo`
- `frontend/src/pages/DepositPage.jsx` — 查询入口改为 16 位卡号提示；收据卡号字段从 `receipt.cardId` 改为 `receipt.cardNo`
- `frontend/src/pages/MealPage.jsx` — 变量名 `cardId`→`cardNo`；卡信息展示改为 `cardInfo.cardNo`

构建验证：`pnpm build` 通过，无报错。

关键决策：
- 前端所有涉及卡主键的地方统一使用 `cardNo`（16位字符串），不再使用数据库自增 id
- IssuePage 发卡流程拆分两步，第一步通过学籍验证获得姓名，第二步才提交发卡，避免操作员录错信息
- 挂失/注销页面统一改用证件号查询，符合 PRD 流程（持卡人到管理处时报证件号，不一定知道卡号）

### 2026-04-16 第 8 轮：后端重构启动（已完成部分）

新增文件：
- `mock-services/student-service/main.go` — 学籍验证 Mock 服务，硬编码5条学生/教职工记录，GET /validate
- `backend/service/student_validator.go` — StudentValidator 接口 + StudentInfo 结构体定义
- `backend/client/student_client.go` — HttpStudentValidator，HTTP 调用 Mock 服务（未完成，缺少错误码定义）

修改文件：
- `backend/model/card.go` — 新增 CardNo 字段（string，uniqueIndex，size:16），ID json tag 改为 `-` 不对外暴露

待完成（本轮中断）：
- card_service.go 错误码新增 STUDENT_NOT_FOUND / STUDENT_SERVICE_ERROR
- repository/service/handler/router/main.go 全量适配
- 测试更新
- 前端全量适配



修改文件：
- `docs/prd.md` — 发卡流程加入学籍验证+二次确认步骤；押金改为固定20元不可录入；挂失/注销入口改为证件号
- `docs/architecture.md` — 新增 mock-services/student-service 模块；新增卡号格式约束（16位字符串）；新增押金固定常量；新增学籍服务集成方式（环境变量 STUDENT_SERVICE_URL）
- `docs/api/openapi.yaml` — v2.0.0：路径参数 {id}(int) 全部替换为 {cardNo}(string)；新增 GET /api/validate-student；新增 GET /api/cards?idNumber=；发卡请求移除 deposit 字段；Card schema 的主标识改为 cardNo；新增错误码 STUDENT_NOT_FOUND/STUDENT_SERVICE_ERROR

关键决策：
- 卡号用 16 位随机数字字符串（card_no 字段），数据库自增 id 仅作内部 FK，不对外暴露
- 押金 2000 分（20元）系统常量，不可操作员录入，DB 仍存储以便退款
- 学籍验证服务（mock）作为独立 Go 进程，通过接口定义与后端解耦，硬编码学生/教职工名单
- 挂失/注销均以证件号为入口（前端先 GET /api/cards?idNumber= 查卡，再用 cardNo 调用操作接口）
- 存款和就餐消费依然以卡号（card_no）为入口（持卡到窗口/管理处）

新增文件：
- `frontend/src/layouts/AdminLayout.jsx` — 管理端侧边栏 Layout，使用 Ant Design Sider + Menu，包含返回首页按钮

修改文件：
- `frontend/src/App.jsx` — 重写路由结构：`/` 首页（双入口卡片）、`/admin/*`（管理端嵌套路由）、`/window`（窗口机独立路由）
- `frontend/src/pages/MealPage.jsx` — 改为深色全屏窗口机风格，余额大字展示（64px），按钮/输入框放大，支持 `useNavigate` 返回首页

路由映射：
- `/admin/issue` — 发卡
- `/admin/deposit` — 存款
- `/admin/loss` — 挂失管理
- `/admin/cancel` — 注销
- `/admin/statistics` — 汇总统计
- `/admin/windows` — 窗口管理
- `/window` — 窗口机就餐消费

关键决策：
- 管理端使用固定宽度（180px）侧边栏，Content 区 `marginLeft: 180` 避免遮挡
- 窗口机页面完全独立于管理导航，使用深色主题（`#0a1628` 背景）与大字体，方便站立操作
- 首页不含任何业务功能，仅作两个入口的跳板

### 2026-04-16 第 5 轮：前端全面迁移至 Ant Design

新增依赖：
- `antd` 6.3.5
- `@ant-design/icons` 6.1.1
- `dayjs` 1.11.20（antd DatePicker 所需）

修改文件：
- `frontend/src/main.jsx` — 引入 `antd/dist/reset.css`
- `frontend/src/App.jsx` — 改用 Layout + Header + Menu 导航，Menu 通过 react-router-dom navigate 切换路由
- `frontend/src/pages/IssuePage.jsx` — Form + InputNumber + Descriptions + Card 展示结果
- `frontend/src/pages/DepositPage.jsx` — Space.Compact 查询栏 + Descriptions 展示卡信息 + Descriptions 收据
- `frontend/src/pages/MealPage.jsx` — Select 切换窗口 + Steps 两步流程 + Alert 报警
- `frontend/src/pages/LossPage.jsx` — Tag 显示状态 + Button danger 申请挂失
- `frontend/src/pages/CancelPage.jsx` — Modal.confirm 二次确认 + Descriptions 退款明细
- `frontend/src/pages/StatisticsPage.jsx` — Row/Col 响应式布局 + Card 分区块 + RangePicker + DatePicker + Table
- `frontend/src/pages/WindowsPage.jsx` — Table 窗口列表 + Form inline 新建窗口

关键决策：
- 使用 Ant Design 默认商务风格，不引入深色主题或自定义视觉效果
- StatisticsPage 中 RangePicker 的值类型为 dayjs 对象，调用 .toISOString() 转 ISO 字符串传给 API
- CancelPage 改用 Modal.confirm 替代内联二次确认按钮，体验更清晰
- 所有页面业务逻辑（api.js 调用、金额换算规则）保持不变

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
