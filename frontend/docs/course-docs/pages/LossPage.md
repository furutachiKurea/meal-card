# LossPage（挂失管理页）

## 作用
- 管理员查询卡状态，对正常卡申请挂失，或对已挂失卡取消挂失

## 结构
- 单 Card 布局，最大宽 520px，居中
- 卡号查询区：Input + 查询按钮（Space.Compact）
- 卡片信息区（查询后显示）：卡号 / 持卡人 / 证件号 / 状态 Tag / 余额（Descriptions）
- 操作按钮区（在信息卡片内，根据状态条件渲染）
- 反馈 Alert（成功绿色 / 失败红色）

## 状态 Tag 映射

| status | Tag color | 显示文字 |
|--------|-----------|---------|
| active | success | 正常 |
| lost | warning | 已挂失 |
| cancelled | default | 已注销 |

## 操作按钮逻辑
- status === 'active' → 红色「申请挂失」按钮
- status === 'lost' → 蓝色「取消挂失」按钮
- status === 'cancelled' → 灰色文字「该卡已注销，无法操作」（无按钮）

## 用户操作
1. 输入卡号 → 查询（回车或点击按钮）
2. 查看卡片状态
3. 若正常：点击「申请挂失」→ 卡状态更新为 lost，成功提示
4. 若已挂失：点击「取消挂失」→ 卡状态更新为 active，成功提示

## 状态变化
```
cardId 变化 → 清空 cardInfo / successMsg / error
查询成功 → cardInfo 填入
挂失/取消挂失成功 → cardInfo 更新为接口返回值；successMsg 填入
操作失败 → error 填入
```

## 接口依赖
- `GET /api/cards/{id}`（getCard）—— 查询卡信息
- `PUT /api/cards/{id}/loss-report`（reportLoss）—— 挂失
- `DELETE /api/cards/{id}/loss-report`（cancelLossReport）—— 取消挂失

## 备注 / 待确认
- 挂失/取消挂失后，cardInfo 直接用接口返回值替换，不需要重新查询
