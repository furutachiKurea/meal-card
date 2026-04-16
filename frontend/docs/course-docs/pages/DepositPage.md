# DepositPage（存款页）

## 作用
- 管理员为持卡人办理充值存款，并打印（展示）收据

## 结构
- 单 Card 布局，最大宽 520px，居中
- 卡号查询区（顶部）：Input + 查询按钮（Space.Compact）
- 卡片信息区（查询后显示）：持卡人 / 卡号 / 当前余额（Descriptions bordered）
- 存款表单区（卡片信息出现后显示）：存款金额 InputNumber
- 收据区（存款成功后显示）：绿色背景 Card，展示充值明细

## 用户操作
1. 输入卡号，点击「查询」或回车
   - 若卡状态为 lost / cancelled → 显示错误提示，不展示存款表单
   - 若正常 → 展示卡片信息和存款表单
2. 输入存款金额，点击「确认存款」
   - 成功：展示收据；更新卡片信息中的余额（本地同步，不再重新查询）；表单重置
   - 失败：显示错误 Alert

## 状态变化
```
cardId 变化 → 清空 cardInfo / receipt / error
查询成功 → cardInfo 填入
存款成功 → receipt 填入；cardInfo.balance 本地更新为 res.newBalance
```

## 接口依赖
- `GET /api/cards/{id}`（getCard）—— 查询卡信息
- `POST /api/cards/{id}/deposits`（deposit）
  - 入参：amount（分）
  - 金额处理：显示元，提交前 × 100 转分

## 备注 / 待确认
- 收据中 createdAt 使用 `new Date(v).toLocaleString('zh-CN')` 格式化
- 存款表单仅在 cardInfo 存在时渲染（条件渲染控制）
