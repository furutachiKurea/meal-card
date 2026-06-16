# 迭代进度

## 当前迭代目标

功能拓展：窗口机双屏分离、种子数据、系统闭环完善。

## 已完成

- 前端项目骨架（React + Vite + react-router-dom），全部 7 个页面实现完成，构建无报错
- PRD 文档（`docs/prd.md`），覆盖 6 项功能的完整行为和约束
- 数据库表设计：5 张表的 GORM 模型（`backend/model/`）及文档（`backend/docs/database/`）
- 架构文档（`docs/architecture.md`）
- OpenAPI 契约（`docs/api/openapi.yaml`），覆盖全部接口
- 后端全量适配 cardNo + StudentValidator
- 前端全量适配 cardNo
- 后端事务化重构（IssueCard/Deposit/CreateTransaction 使用数据库事务）
- GetCurrentCardByIDNumber 下沉到 service 层
- 前端 DepositPage loading 状态独立（queryLoading / depositLoading）
- 窗口机支持 URL 参数绑定窗口（`/window?id=1`）
- 404 页面
- 前端 api.js 改为相对路径 + Vite 代理
- 全流程集成测试通过（发卡→存款→消费→挂失→取消挂失→注销→统计）
- 窗口机拆分为操作端（/window）和顾客屏（/window/customer），BroadcastChannel 实时同步
- 后端启动时自动初始化 5 个默认窗口（种子数据）
- 首页三入口卡片（管理端/窗口操作端/顾客屏）

- 消费限额：单笔 200 元、日累计 500 元，超限拒绝结算
- 历史查询 API：GET /api/cards/:cardNo/transactions 和 GET /api/cards/:cardNo/deposits

## 进行中

（无）

## 待办

（无）

## 阻塞

（无）

## 下一步

系统已形成完整闭环：6 项基本功能 + 消费限额 + 历史查询 + 双屏窗口机 + 种子数据。

## 变更记录

### 2026-06-16 第 20 轮：消费限额 + 历史查询 API

修改文件：
- `backend/service/card_service.go` — CreateTransaction 增加单笔限额（200元）和日累计限额（500元）校验；新增 GetCardTransactions / GetCardDeposits 方法
- `backend/repository/card_repository.go` — 新增 SumCardTransactionsByTimeRange（按卡+时间范围统计消费）、GetCardTransactions（消费记录分页）
- `backend/handler/card_handler.go` — 新增 GetCardTransactions / GetCardDeposits handler；bizErrStatus 增加 EXCEED_SINGLE_LIMIT / EXCEED_DAILY_LIMIT → 403
- `backend/router/router.go` — 注册 GET /api/cards/:cardNo/transactions 和 GET /api/cards/:cardNo/deposits
- `backend/service/card_service_test.go` — 修正余额不足测试用例金额；新增 TestCreateTransaction_ExceedSingleLimit / ExceedDailyLimit

新增文件：
- `docs/course-docs/usecases/transaction-limit.md` — 消费限额 use case 文档

测试结果：go test 全部通过；集成测试验证限额拦截和历史查询正确

关键决策：
- 单笔 200 元 / 日累计 500 元为系统常量，不可配置（课设不需要运行时配置）
- 日累计计算在事务内执行，避免并发消费绕过限额
- 历史查询复用已有 repository 方法，不引入新表

### 2026-06-16 第 19 轮：窗口机双屏拆分 + 种子数据

修改文件：
- `frontend/src/pages/MealPage.jsx` — 重构为"窗口操作端"，新增 BroadcastChannel 广播能力
- `frontend/src/App.jsx` — 首页三入口卡片，新增 `/window/customer` 路由
- `backend/main.go` — 启动时调用 db.Seed 初始化种子数据

新增文件：
- `frontend/src/pages/CustomerScreen.jsx` — 顾客屏（学生面对的只读大字界面）
- `backend/db/seed.go` — 种子数据逻辑（5 个默认窗口）
- `frontend/docs/course-docs/pages/CustomerScreen.md` — 顾客屏页面文档
- `docs/course-docs/usecases/window-dual-screen.md` — 双屏协同 use case 文档

测试结果：后端 go test 全通过；前端 pnpm build 通过；集成测试全流程 API 验证通过

关键决策：
- 操作端与顾客屏通过 BroadcastChannel 通信，频道名与窗口 ID 绑定
- 种子数据仅在 windows 表为空时插入，不影响已有数据
- 首页从两入口扩展为三入口（管理端/操作端/顾客屏）

### 2026-06-16 第 18 轮：系统闭环完善

修改文件：
- `backend/service/card_service.go` — IssueCard/Deposit/CreateTransaction 事务化重构
- `backend/repository/card_repository.go` — 新增 WithTx 方法
- `backend/handler/card_handler.go` — GetCardByIDNumber 改为调用 cardSvc.GetCurrentCardByIDNumber
- `backend/main.go` — 使用 MockStudentValidator（内嵌，无需外部服务）
- `frontend/src/pages/DepositPage.jsx` — 修复未定义 loading 变量 bug（拆分为 queryLoading/depositLoading）
- `frontend/src/pages/MealPage.jsx` — 支持 URL 参数绑定窗口（`/window?id=1`）
- `frontend/src/pages/NotFoundPage.jsx` — 新增 404 页面
- `frontend/src/App.jsx` — 注册 404 路由
- `frontend/vite.config.js` — 添加 Vite 代理（/api → localhost:8080）
- `frontend/src/api.js` — BASE_URL 改为空字符串（使用相对路径）
- `docs/architecture.md` — 更新学籍验证实现说明、新增事务策略
- `docs/progress.md` — 更新进度
- `README.md` — 重写，精简内容，补充项目结构和窗口机说明
- `backend/docs/course-docs/modules/card-service.md` — 补充事务化策略和 GetCurrentCardByIDNumber 说明

集成测试：
- 启动后端服务，通过 curl 验证全部 API 接口（20 个请求覆盖全流程）
- 发卡→存款→消费→余额不足→挂失→挂失后消费拒绝→取消挂失→注销→注销后消费拒绝→再次发卡
- 统计接口全部返回正确（售饭收入、窗口收入、存款汇总、流动资金、日报、年报、持卡人存款明细）
- 后端单元测试全部通过（29 个），前端构建通过

关键决策：
- 学籍验证改为内嵌 MockStudentValidator，不再依赖外部 mock-services 进程
- 前端使用 Vite 代理替代硬编码 localhost:8080，开发和构建均可正常工作
- 窗口机通过 URL query param 绑定窗口，真实部署时每台机器配不同 URL

修改文件：
- `backend/repository/card_repository.go` — 新增 `GetHolderDeposits`，支持 holderID、可选时间范围、LIMIT/OFFSET 分页
- `backend/service/statistics_service.go` — 新增 `HolderDepositsResult` 和 `GetHolderDeposits` 方法
- `backend/handler/statistics_handler.go` — 新增 `GetHolderDeposits` handler，解析 holderId/page/pageSize/startTime/endTime
- `backend/router/router.go` — 注册 `GET /api/statistics/holder-deposits`
- `backend/service/statistics_service_test.go` — 新增 GetHolderDeposits 单测
- `frontend/src/api.js` — 新增 `getHolderDeposits`
- `frontend/src/pages/StatisticsPage.jsx` — 内层存款改为服务端分页；外层持卡人改用 Ant Design Pagination；所有分页加 `showSizeChanger`（选项 5/10/20/50），用户可在页面实时切换每页条数；内层每个持卡人独立保存自己的 `pageSize`
- `docs/course-docs/usecases/statistics.md` — 新增汇总统计全部功能的 use case 文档
- `backend/docs/course-docs/modules/statistics-service.md` — 补充第 8 项 GetHolderDeposits 说明

测试结果：go test ./... 全部通过；pnpm build 无报错

关键决策：
- 内层存款记录 pageSize 按持卡人独立存储（depositRecordsData[holderId].pageSize），切换某人分页大小不影响其他人
- 简单表格（窗口收入、日/年报明细）用客户端分页 + showSizeChanger，无需额外状态

### 2026-04-16 第 16 轮：zerolog 日志策略调整 + 前端 bug 修复

修改文件：
- `backend/service/card_service.go` — 删除 log.Info() 成功日志；在 newBizError() 调用点统一加 log.Warn()，在数据库等底层错误处加 log.Error()；日志在错误源头打印一次
- `backend/handler/card_handler.go` — 删除所有 log.Info()；handleError 不打日志（service 层已打）；去掉 zerolog import
- `backend/handler/statistics_handler.go` — 删除所有 log.Info()，去掉 zerolog import
- `backend/handler/window_handler.go` — 删除 log.Info()，去掉 zerolog import
- `frontend/src/pages/MealPage.jsx` — 消费金额 InputNumber 改用内联 `<style>` + `className` 方案覆盖 Ant Design CSS 变量（`styles.input.color` 因 CSS 变量优先级问题无效）
- `frontend/src/pages/StatisticsPage.jsx` — 外层持卡人分页由 Button.Group 替换为 Ant Design `<Pagination>` 组件，total > 0 时常驻显示
- `frontend/docs/bugfixes/meal-inputnumber-black-background.md` — 补充说明第一次方案无效原因及最终方案

测试结果：go build ./... + go test ./... 全部通过；pnpm build 无报错

关键决策：
- 日志打在错误最初产生的 service 层，handler 层只做 HTTP 响应映射，避免同一错误重复打印
- Ant Design v5 InputNumber 文字颜色受 CSS 变量控制，`styles.input.color` 优先级不足，需用 CSS className + `!important` 覆盖



修改文件：
- `backend/repository/card_repository.go` — GetDepositDetails 新增 page/pageSize 参数，改为返回 (holders, total, err)；先 DISTINCT COUNT 统计持卡人总数，再 GROUP BY + LIMIT/OFFSET 分页取持卡人 ID，最后 IN 查询本页全量存款明细
- `backend/service/statistics_service.go` — GetDepositDetails 新增 page/pageSize 参数，page/pageSize < 1 时自动修正；DepositDetailsResult 新增 Total/Page/PageSize 字段
- `backend/handler/statistics_handler.go` — 解析 page/pageSize query 参数（非法值使用默认值）；响应体新增 total/page/pageSize
- `backend/service/statistics_service_test.go` — 新增 TestGetDepositDetails_Pagination，覆盖：全量查询、分页第 1 页、第 2 页、超出范围页

测试结果：go test -count=1 ./... 全部通过

关键决策：
- 分页单位为「持卡人」（不是存款记录），与前端展示逻辑一致
- 排序以 card_holders.id ASC 保证跨页稳定
- 分两次查询（先 COUNT+分页持卡人，再查明细），避免大结果集一次性加载

### 2026-04-16 第 14 轮：修复发卡预存款不生成流水记录

修改文件：
- `backend/service/card_service.go` — IssueCard 在 CreateCard 后，若 preDeposit > 0 补建 DepositRecord
- `backend/service/card_service_test.go` — 新增 TestIssueCard_PreDepositCreatesDepositRecord，覆盖 preDeposit > 0 和 == 0 两个分支

新增文件：
- `backend/docs/bugfixes/issue-card-pre-deposit-no-record.md`

测试结果：go test ./... 全部通过

### 2026-04-16 第 13 轮：统计页表格加分页

修改文件：
- `frontend/src/pages/StatisticsPage.jsx` — 各窗口收入、日餐/年餐报表明细表每页 10 条；各持卡人存款明细外层持卡人每页 5 条（手动分页），内层存款记录每页 10 条

### 2026-04-16 第 12 轮：修复就餐页输入框黑色背景

修改文件：
- `frontend/src/pages/MealPage.jsx` — InputNumber 增加 `styles.input` 直接设置内部 input 背景色与文字色，修复深色主题下黑色背景导致文字不可见问题

新增文件：
- `frontend/docs/bugfixes/meal-inputnumber-black-background.md`

关键决策：
- Ant Design v5 InputNumber 的 `style` 只作用于外层 wrapper，需用 `styles={{ input: {...} }}` 才能穿透到内部 input 元素

### 2026-04-16 第 11 轮：修复两处前端 bug

修改文件：
- `frontend/src/pages/StatisticsPage.jsx` — 存款明细列 `dataIndex: 'cardId'` → `'cardNo'`，修复卡号列空白问题
- `frontend/src/pages/DepositPage.jsx` — 收据加"单号"（`receipt.id`）字段；加"打印收据"按钮（`window.print()`）及 `@media print` 样式隔离

新增文件：
- `frontend/docs/bugfixes/statistics-deposit-detail-cardno-missing.md`
- `frontend/docs/bugfixes/deposit-receipt-missing-id-and-print.md`

关键决策：
- 打印样式内嵌于组件，用 `.deposit-receipt-print` class 隔离打印区域，不影响全局样式

### 2026-04-16 第 10 轮：后端全量适配 cardNo + StudentValidator

修改文件：
- `backend/repository/card_repository.go` — 新增 `FindCardByCardNo`、`FindCurrentCardByIDNumber`；`DepositDetailItem.CardID` → `CardNo`
- `backend/service/card_service.go` — 新增错误码 `ErrCodeStudentNotFound`/`ErrCodeStudentServiceError`；`CardService` 注入 `StudentValidator`；`IssueCard` 签名改为 `(idNumber string, preDeposit int64)`；押金固定 2000 分由系统写入；Deposit/CreateTransaction/ReportLoss/CancelLossReport/CancelCard 参数由 `uint` → `cardNo string`；`OldCardRefund.OldCardID` → `OldCardNo`；`DepositResult/TransactionResult.CardID` → `CardNo`
- `backend/service/statistics_service.go` — `DepositDetailEntry.CardID` → `CardNo`
- `backend/handler/card_handler.go` — 新增 `ValidateStudent`/`GetCardByIDNumber` handler；路径参数 `:id` 全改为 `:cardNo`；响应字段 `cardId` → `cardNo`；`CardHandler` 注入 `StudentValidator` 和 `CardRepository`
- `backend/handler/statistics_handler.go` — 存款明细响应字段 `cardId` → `cardNo`
- `backend/router/router.go` — 路径全部改为 `:cardNo`；新增 `GET /api/validate-student` 和 `GET /api/cards`
- `backend/main.go` — 实例化 `client.NewHttpStudentValidator()`，注入 `NewCardService` 和 `NewCardHandler`
- `backend/service/card_service_test.go` — 全量重写：`FakeStudentValidator`；所有 `cardID uint` → `cardNo string`；新增 `STUDENT_NOT_FOUND`/`STUDENT_SERVICE_ERROR` 测试用例；新增 `setupWithRepo` 辅助函数
- `backend/service/statistics_service_test.go` — 全量更新：适配新 `IssueCard` 签名，`cardID` → `cardNo`

新增文件：
- `backend/docs/course-docs/modules/card-handler.md` — CardHandler 模块说明文档

测试结果：`go test ./...` 全部通过。

关键决策：
- 对外接口统一以 `cardNo`（16 位字符串）为主键，数据库自增 `id` 仅作内部 FK，不暴露给任何 HTTP 响应
- 押金 2000 分为系统常量，由后端写入，前端无需也无法传入
- `GetCardByIDNumber` handler 直接注入 `CardRepository`（无业务规则，纯读取），不走 service 层

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
