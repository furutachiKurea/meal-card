# statistics_service 统计模块

## 作用
- 提供 7 项统计查询，支撑管理端汇总分析需求

## 职责边界
- 负责：统计时间范围计算、调用 Repository 聚合查询、数据结构组装
- 不负责：卡片状态变更、HTTP 参数解析

---

## 7 项统计接口

### 1. GetMealRevenue（本餐售饭总收入）
- 输入：startTime、endTime（ISO 8601，必填）
- 输出：totalRevenue（int64，单位分）
- 口径：transactions 表中 created_at 在 [start, end] 区间内的 amount 求和

### 2. GetWindowRevenue（各窗口收入）
- 输入：startTime、endTime（ISO 8601，必填）
- 输出：windows 列表，每项含 windowId、windowName、revenue
- 口径：transactions 按 window_id 分组求和，LEFT JOIN windows 表获取名称

### 3. GetDepositDetails（各持卡人存款明细）
- 输入：startTime、endTime（ISO 8601，可选）
- 输出：holders 列表，每项含 holderId、holderName、idNumber、deposits 明细列表、totalAmount
- 口径：deposit_records LEFT JOIN cards LEFT JOIN card_holders，按时间范围过滤后按持卡人分组
- 注：时间参数均可选，不传则返回全量数据

### 4. GetDepositSummary（本日/本月存款金额）
- 输入：无
- 输出：todayTotal、monthTotal
- 口径：
  - 今日：当天 00:00:00 ~ 次日 00:00:00
  - 本月：当月 1 日 00:00:00 ~ 下月 1 日 00:00:00
  - 均对 deposit_records.amount 求和

### 5. GetActiveBalance（卡中流动资金总额）
- 输入：无
- 输出：totalBalance
- 口径：cards 表中 status=active 的所有记录的 balance 求和

### 6. GetDailyReport（日餐报表）
- 输入：date（YYYY-MM-DD，必填）
- 输出：date、totalRevenue、transactionCount、windows 明细（含每窗口 revenue 和 transactionCount）
- 口径：指定日期 00:00:00 ~ 次日 00:00:00，transactions 表求总额+笔数，再按 window_id 分组求明细

### 7. GetYearlyReport（年餐报表）
- 输入：year（正整数，必填）
- 输出：year、totalRevenue、transactionCount、months 列表（1~12 月，含 revenue 和 transactionCount）
- 口径：SQLite `strftime('%Y', created_at) = year`，按月份（`strftime('%m')`）分组，只返回有数据的月份

## 异常处理
- date 格式非 YYYY-MM-DD：VALIDATION_ERROR（400）
- year <= 0：VALIDATION_ERROR（400）
- startTime/endTime 格式非 ISO 8601：VALIDATION_ERROR（400）
- startTime/endTime 为必填时缺少参数：VALIDATION_ERROR（400）

## 关键实现点
- GetDepositSummary 的时间范围在 service 层用 time.Now() 计算，不接受外部参数
- GetYearlyReport 月份聚合仅返回有数据的月份，不补零
- 所有金额单位为分（int64）
