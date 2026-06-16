# CustomerScreen 顾客屏

## 作用
- 学生面对的只读大字显示屏，放置在窗口机的学生侧
- 通过 BroadcastChannel 接收同窗口操作端（MealPage）推送的状态，无需任何用户交互
- 展示当前刷卡状态：余额、消费结果、报警信息

## 路由
- `/window/customer` — 默认频道
- `/window/customer?id=1` — 绑定到窗口 1 的频道

## 结构
- 全屏深色背景，居中显示单一状态卡片
- 四种状态视图：idle（等待刷卡）、card_read（显示余额）、settled（结算完成）、alarm（报警）

## 状态变化
- `idle` → 显示"请刷卡"提示图标
- `card_read` → 显示持卡人姓名 + 大字余额
- `settled` → 绿色成功卡片，显示本次消费和剩余余额
- `alarm` → 红色警报卡片，显示报警原因

## 通信机制
- 使用浏览器 BroadcastChannel API
- 频道名：`window-{id}` 或 `window-default`（无 id 参数时）
- 接收消息类型：`card_read`、`settled`、`alarm`、`reset`
- 顾客屏只接收不发送，完全只读

## 接口依赖
- 无（所有数据通过 BroadcastChannel 从操作端推送）

## 备注
- 操作端和顾客屏必须在同一浏览器的同源标签页中才能通信
- 实际部署时，窗口机的两块屏幕分别打开操作端和顾客屏 URL
