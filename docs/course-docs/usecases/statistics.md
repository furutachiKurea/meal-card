# 汇总统计

## 当前行为
- 统计页提供 7 项查询功能，所有分页常驻显示，不需要数据才显示分页控件
- 分页大小全部通过常量配置（默认每页 10 条/人）

---

## 本餐售饭总收入

### 前置条件
- startTime、endTime 必填，格式为 ISO 8601

### 成功路径
1. 传入时间范围，返回该范围内 transactions 表的 amount 总和（单位：分）

### 异常 / 边界情况
- startTime 或 endTime 缺失：返回 VALIDATION_ERROR
- 时间格式非 ISO 8601：返回 VALIDATION_ERROR

---

## 各窗口收入

### 前置条件
- startTime、endTime 必填，格式为 ISO 8601

### 成功路径
1. 传入时间范围，按 windowId 分组返回各窗口收入

### 异常 / 边界情况
- startTime 或 endTime 缺失：返回 VALIDATION_ERROR

---

## 各持卡人存款明细（含服务端分页）

### 前置条件
- startTime、endTime 可选
- page 默认 1，pageSize 默认 10，小于 1 时重置为默认值

### 成功路径
1. 查询 deposit_records，按持卡人分组
2. 返回当前页的持卡人列表，每人含存款明细列表和汇总金额
3. 返回满足条件的持卡人总数、page、pageSize

### 规则 / 限制
- 分页单位为「持卡人」，一个持卡人的所有存款记录不会拆分到不同页
- 内层每人存款记录同样服务端分页，分页大小独立于外层

---

## 按持卡人查询存款明细（单人，含服务端分页）

### 前置条件
- holderId 必填，须为正整数
- startTime、endTime 可选
- page 默认 1，pageSize 默认 10

### 成功路径
1. 先通过 holderId 查询持卡人信息（姓名、证件号）
2. 查询该持卡人在 deposit_records 的记录，可选时间范围过滤
3. 返回持卡人信息、当前页存款列表、total、page、pageSize

### 异常 / 边界情况
- holderId 对应持卡人不存在：返回 HOLDER_NOT_FOUND（404）
- holderId 缺失或非正整数：返回 VALIDATION_ERROR（400）
- 时间格式非 ISO 8601：返回 VALIDATION_ERROR（400）

### 规则 / 限制
- 分页单位为「存款记录」
- 此接口专用于前端点击某个持卡人后查看其完整存款明细

---

## 本日 / 本月存款金额

### 成功路径
1. 无需任何参数，服务端按当前时间自动计算今日和本月范围
2. 返回 todayTotal、monthTotal（单位：分）

### 规则 / 限制
- 今日范围：当天 00:00:00 ~ 次日 00:00:00
- 本月范围：当月 1 日 00:00:00 ~ 下月 1 日 00:00:00
- 时间范围在 service 层用 time.Now() 计算，不接受外部参数

---

## 卡中流动资金总额

### 成功路径
1. 无需参数，返回所有 status=active 的卡的 balance 总和（单位：分）

---

## 日餐报表

### 前置条件
- date 必填，格式 YYYY-MM-DD

### 成功路径
1. 统计指定日期 00:00:00 ~ 次日 00:00:00 的消费总额、笔数
2. 按窗口分组返回明细（每窗口含 revenue 和 transactionCount）

### 异常 / 边界情况
- date 缺失：返回 VALIDATION_ERROR
- date 格式非 YYYY-MM-DD：返回 VALIDATION_ERROR

---

## 年餐报表

### 前置条件
- year 必填，须为正整数

### 成功路径
1. 统计指定年份全年消费总额和笔数
2. 返回各月明细列表，只返回有数据的月份（不补零）

### 异常 / 边界情况
- year 缺失：返回 VALIDATION_ERROR
- year <= 0：返回 VALIDATION_ERROR
