# CancelPage（注销页）

## 作用
- 管理员查询持卡人当前有效卡信息，二次确认后注销该卡，展示退款明细（押金 + 余额合计）

## 结构
- 单 Card 布局，最大宽 520px，居中
- 证件号查询区：Input（placeholder：请输入证件号） + 查询按钮（Space.Compact）
- 卡片信息区（查询后显示）：卡号 / 持卡人 / 证件号 / 状态 Tag / 余额 / 押金 / 预计退款（Descriptions）
- 操作按钮区（在信息卡片内）
- 退款明细 Card（注销成功后显示，绿色背景）

## 变更说明（相对旧版）
- **变更**：查询入口从"输入卡号"改为"输入证件号（学号/工号）"
- **变更**：查询接口从 `GET /api/cards/{id}` 改为 `GET /api/cards?idNumber=xxx`
- 其余展示和操作逻辑不变

## 预计退款计算
- 前端本地计算：`(cardInfo.balance + cardInfo.deposit) / 100` 元，仅作预览
- 实际退款以后端 response 中 `refund` 字段为准

## 操作按钮逻辑
- status !== 'cancelled' → 红色实心「申请注销」按钮
- status === 'cancelled' → 灰色文字「该卡已注销」（无按钮）

## 用户操作
1. 输入证件号 → 查询
2. 查看卡片信息（含预计退款金额）
3. 点击「申请注销」→ 弹出 Modal.confirm 二次确认
   - 确认文案：「确认要注销该卡吗？此操作不可撤销！」
   - okText：「确认注销」（danger 类型）
4. 确认后执行注销 → 展示退款明细（押金 / 余额 / 合计）

## 状态变化
```
输入值变化 → 清空 cardInfo / result / error
查询成功 → cardInfo 填入（包含 cardNo、持卡人信息、状态）
注销成功 → result 填入；cardInfo 清空（设为 null）
注销失败 → error 填入
```

## 接口依赖
- `GET /api/cards?idNumber=xxx`（getCardByIDNumber）—— 按证件号查询当前有效卡
- `POST /api/cards/{cardNo}/cancellation`（cancelCard）—— 注销（路径参数为 16 位 cardNo）

## 备注 / 待确认
- 使用 Modal.confirm 实现二次确认，防止误操作
- 注销成功后清空 cardInfo，避免用户对已注销卡重复操作
- 退款明细中「应退合计」用绿色加粗字体突出显示
