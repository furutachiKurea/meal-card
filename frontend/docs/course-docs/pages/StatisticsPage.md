# StatisticsPage（汇总统计页）

## 作用
- 管理端统计中心，提供 7 种统计查询，每种独立触发、独立展示结果

## 结构
- Row/Col 栅格布局，最大宽 860px，居中
- 7 个独立 Card，每个 Card 内有自己的查询控件 + 结果区

## 统计模块列表

| 模块 | 查询控件 | 结果展示 |
|------|---------|---------|
| 本餐售饭总收入 | RangePicker（带时间）| 总收入文字 |
| 各窗口收入 | RangePicker（带时间）| Table（窗口/收入） |
| 各持卡人存款明细 | RangePicker（可选）| 按持卡人分组的嵌套 Card + Table |
| 本日/本月存款金额 | 无（直接查询）| Descriptions（今日/本月） |
| 卡中流动资金总额 | 无（直接查询）| 总额文字 |
| 日餐报表 | DatePicker | Descriptions（日期/笔数/总收入）+ Table（窗口明细） |
| 年餐报表 | InputNumber（年份）| Descriptions（年份/笔数/总收入）+ Table（按月明细） |

## 用户操作
- 各模块独立操作，互不影响
- 有时间范围的：先选时间，再点查询（未选时查询按钮 disabled）
- 本日/本月存款 和 流动资金：直接点查询按钮
- 存款明细：时间范围可选，不填则查全部

## 状态变化
- 每个模块有独立的 result state 和 error state
- loading/errors 用对象 key 区分（`wrap(key, fn)` 统一处理）
- 查询中：对应 key 的 loading 为 true，按钮 loading 状态
- 查询失败：对应 key 的 error 显示 Alert

## 接口依赖
- `GET /api/statistics/meal-revenue?startTime&endTime`
- `GET /api/statistics/window-revenue?startTime&endTime`
- `GET /api/statistics/deposit-details[?startTime&endTime]`
- `GET /api/statistics/deposit-summary`
- `GET /api/statistics/active-balance`
- `GET /api/statistics/daily-report?date=YYYY-MM-DD`
- `GET /api/statistics/yearly-report?year=YYYY`

## 金额格式化
- 统一用 `formatYuan(fen)` 函数：`(fen / 100).toFixed(2) + ' 元'`

## 备注 / 待确认
- 日餐报表的 Table 按窗口展示收入和笔数
- 年餐报表的 Table 按月份（1-12）展示收入和笔数
- 存款明细按持卡人分组，每组展示 holderName / idNumber / totalAmount，组内 Table 列：卡号/金额/时间
- 时间控件使用 dayjs（配合 Ant Design DatePicker）
