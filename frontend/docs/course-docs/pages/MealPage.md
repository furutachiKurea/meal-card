# MealPage（窗口操作端）

## 作用
- 食堂工作人员（阿姨）使用的操作界面：刷卡验证 → 显示余额 → 输入消费金额 → 结算
- 深色大字体界面，适合站立操作场景
- 通过 BroadcastChannel 向顾客屏（CustomerScreen）实时推送状态

## 路由
- `/window` — 手动选择窗口
- `/window?id=1` — 绑定到窗口 1（隐藏选择器）

## 结构
- 全屏深色背景（#0a1628）
- 顶栏：标题「窗口操作端」/ 窗口选择 Select / 返回首页按钮
- 步骤条（Steps，3 步）：刷卡验证 → 输入金额 → 结算完成
- 主内容区（居中，最大宽 560px）：根据当前步骤条件渲染不同 Card

## 步骤内容

### 步骤 0：刷卡输入
- 大号 Input（fontSize 22，height 56）自动聚焦
- 「刷卡」按钮或回车触发查询

### 警报区（覆盖步骤 0）
- 触发条件：卡已注销 / 已挂失 / 非本单位卡（404）
- 同时通过 BroadcastChannel 广播 `{ type: 'alarm', message }` 到顾客屏
- 「重新刷卡」按钮调用 handleReset

### 步骤 1：余额 + 金额输入
- 余额大显示：64px 蓝色字体
- 刷卡成功时广播 `{ type: 'card_read', holderName, balance }` 到顾客屏
- 金额 InputNumber 自动聚焦
- 「确认结算」按钮

### 步骤 2：结算完成
- 本次消费 / 扣款后余额并排展示
- 结算成功时广播 `{ type: 'settled', amount, newBalance }` 到顾客屏
- 「下一位」按钮调用 handleReset

## 通信机制（BroadcastChannel）
- 频道名：`window-{id}` 或 `window-default`（无 id 参数时）
- 广播时机：
  - 刷卡成功 → `card_read`
  - 刷卡异常 → `alarm`
  - 结算成功 → `settled`
  - 重置/下一位 → `reset`
- 只发送不接收

## 状态变化
```
初始: step=0
刷卡成功: cardInfo 填入 → step=1 + broadcast card_read
刷卡失败: alarm 填入 + broadcast alarm
结算成功: txResult 填入 → step=2 + broadcast settled
handleReset: 所有状态清空 → step=0 + broadcast reset
```

## 接口依赖
- `GET /api/windows`（listWindows）—— 页面加载时获取窗口列表
- `GET /api/cards/{cardNo}`（getCard）—— 刷卡验证
- `POST /api/cards/{cardNo}/transactions`（createTransaction）—— 结算

## 备注
- 操作端与顾客屏通过同源 BroadcastChannel 通信，必须在同一浏览器实例中
- 实际部署：窗口机的两块屏幕分别打开 `/window?id=N` 和 `/window/customer?id=N`
