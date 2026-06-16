# CustomerScreen 顾客屏

## 作用
- 学生面对的刷卡入口 + 结果展示屏
- 学生在此输入卡号完成刷卡验证，然后等待阿姨结算，最后看到消费结果

## 路由
- `/window/customer` — 默认频道
- `/window/customer?id=1` — 绑定到窗口 1 的频道

## 结构
- 全屏深色背景，居中显示单一状态卡片
- 四种状态视图：idle（刷卡输入）、waiting（等待结算）、settled（结算完成）、alarm（报警）

## 用户操作
- 学生输入 16 位卡号，点确认或回车
- 系统验证卡片状态，通过后显示余额并通知操作端
- 阿姨结算后自动切换到结算结果
- 阿姨点"下一位"后自动回到刷卡输入

## 状态变化
```
idle → 学生输入卡号点确认
  验证通过 → waiting（显示余额，等待结算）
  验证失败 → alarm（显示报警）
waiting → 收到操作端 settled 消息 → settled
settled → 收到操作端 reset 消息 → idle
alarm → 收到操作端 reset 消息 → idle
```

## 通信机制（BroadcastChannel 双向）
- 频道名：`window-{id}` 或 `window-default`
- **发送**：
  - 刷卡成功 → `{ type: 'card_read', cardNo, holderName, balance }` 通知操作端
  - 刷卡异常 → `{ type: 'alarm', message, cardNo }` 通知操作端
- **接收**：
  - `settled` → 切换到结算完成视图
  - `reset` → 回到刷卡输入

## 接口依赖
- `GET /api/cards/{cardNo}`（getCard）— 顾客屏直接调用后端验证卡片

## 关键实现点
- 顾客屏负责调用 getCard API 做卡片验证（不是操作端调）
- 验证结果通过 BroadcastChannel 推送给操作端，操作端被动接收
- 操作端结算完成后通过 BroadcastChannel 回传结果给顾客屏
- 两端必须在同一浏览器同源下才能通信（BroadcastChannel 限制）
