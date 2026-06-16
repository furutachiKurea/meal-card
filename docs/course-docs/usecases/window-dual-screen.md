# 窗口机双屏协同

## 当前行为
- 窗口机拆分为两个独立页面：操作端（`/window`）和顾客屏（`/window/customer`）
- 顾客屏是学生的刷卡入口，学生在此输入卡号完成验证
- 操作端被动接收顾客屏推送的卡信息，阿姨只负责输入金额和结算
- 两端通过浏览器 BroadcastChannel API 双向通信

## 前置条件
- 操作端和顾客屏必须在同一浏览器实例的不同标签页中打开（同源限制）
- 两端的 URL 参数 `?id=N` 必须一致，才能连接到同一频道

## 成功路径
1. 学生在顾客屏输入卡号 → 顾客屏调 API 验证 → 显示余额
2. 顾客屏通过 BroadcastChannel 通知操作端 → 操作端自动显示卡信息和金额输入框
3. 阿姨输入金额确认结算 → 操作端调 API 扣款
4. 操作端通过 BroadcastChannel 通知顾客屏 → 顾客屏显示消费结果
5. 阿姨点"下一位" → 广播 reset → 顾客屏回到刷卡输入

## 异常 / 边界情况
- 刷卡验证失败（挂失/注销/非本单位卡）→ 顾客屏显示报警 + 通知操作端显示警报
- 操作端未打开 → 顾客屏正常验证和显示余额，但无人结算
- 顾客屏未打开 → 操作端停留在"等待学生刷卡"，无法进入下一步

## 规则 / 限制
- 频道命名：有 `?id=N` 时为 `window-{N}`，无参数时为 `window-default`
- 消息类型和方向：
  - 顾客屏 → 操作端：`card_read`（含 cardNo/holderName/balance）、`alarm`（含 message/cardNo）
  - 操作端 → 顾客屏：`settled`（含 amount/newBalance）、`reset`
- 顾客屏负责调用 `GET /api/cards/{cardNo}` 验证卡片
- 操作端负责调用 `POST /api/cards/{cardNo}/transactions` 结算
- 职责分离：验证在学生端，结算在工作人员端

## 技术实现
- 使用浏览器原生 BroadcastChannel API，无需 WebSocket 或后端中转
- 两端各自 new BroadcastChannel(channelName)，channelName 相同即可通信
- 同源限制：必须是同一 origin（协议+域名+端口）下的不同标签页
- 消息为 JSON 对象，通过 postMessage 发送，onmessage 接收
