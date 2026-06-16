# MealPage（窗口操作端）

## 作用
- 食堂工作人员（阿姨）使用的操作界面：看到学生刷卡信息 → 输入消费金额 → 确认结算
- 深色大字体界面，适合站立操作场景
- 被动接收顾客屏的刷卡信息，主动发送结算结果

## 路由
- `/window` — 手动选择窗口
- `/window?id=1` — 绑定到窗口 1（隐藏选择器）

## 结构
- 全屏深色背景（#0a1628）
- 顶栏：标题「窗口操作端」/ 窗口选择 / 返回首页
- 步骤条（3 步）：刷卡验证 → 输入金额 → 结算完成
- 主内容区：根据状态条件渲染

## 步骤内容

### 步骤 0：等待刷卡
- 显示"等待学生刷卡..."提示
- 无输入框，阿姨无法在此端输入卡号
- 顾客屏推送 `card_read` 消息后自动进入步骤 1

### 警报区
- 顾客屏推送 `alarm` 消息后显示红色警报
- 「重新刷卡」按钮广播 `reset` 让顾客屏回到刷卡界面

### 步骤 1：输入金额
- 显示学生姓名、卡号、余额
- 显示限额提示（单笔200/日累计500）
- 阿姨输入消费金额 → 确认结算

### 步骤 2：结算完成
- 显示消费金额和扣款后余额
- 「下一位」按钮广播 `reset`

## 通信机制（BroadcastChannel 双向）
- 频道名：`window-{id}` 或 `window-default`
- **接收**：
  - `card_read`（含 cardNo, holderName, balance）→ 自动填充卡信息，进入步骤 1
  - `alarm`（含 message, cardNo）→ 显示警报
- **发送**：
  - `settled`（含 amount, newBalance）→ 通知顾客屏结算完成
  - `reset` → 通知顾客屏回到刷卡界面

## 状态变化
```
初始: step=0（等待刷卡）
收到 card_read: cardInfo 填入 → step=1
收到 alarm: alarm 填入（警报覆盖）
结算成功: txResult 填入 → step=2 + broadcast settled
handleReset: 清空 → step=0 + broadcast reset
```

## 接口依赖
- `GET /api/windows`（listWindows）— 加载窗口列表
- `POST /api/cards/{cardNo}/transactions`（createTransaction）— 结算扣款

## 关键实现点
- 操作端**不调用** getCard，卡片验证由顾客屏完成
- 操作端只负责：接收卡信息 → 输入金额 → 调用结算 API → 广播结果
- 职责分离：学生刷卡在学生屏，阿姨结算在阿姨屏
