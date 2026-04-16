# 存款明细卡号列不显示

## 现象
- 汇总统计页"各持卡人存款明细"中，每笔存款的"卡号"列为空白，无任何内容显示

## 触发条件
- 在 `/admin/statistics` 页面点击"各持卡人存款明细"的查询按钮，有存款数据时必现

## 影响范围
- `StatisticsPage.jsx` 存款明细表格卡号列

## 根因
- 后端全量适配 `cardNo` 时，`statistics_handler.go` 中存款明细的字段名由 `cardId` 改为 `cardNo`
- 但 `StatisticsPage.jsx` 中 `depositDetailColumns` 的 `dataIndex` 和 `key` 仍为旧值 `cardId`，未同步更新
- Ant Design Table 找不到 `cardId` 字段，渲染为空

## 修复思路
- 将 `depositDetailColumns` 中卡号列的 `dataIndex` 和 `key` 从 `cardId` 改为 `cardNo`

## 实际改动
- `frontend/src/pages/StatisticsPage.jsx`：卡号列 `dataIndex: 'cardId', key: 'cardId'` → `dataIndex: 'cardNo', key: 'cardNo'`

## 修复结果
- 存款明细卡号列正常显示 16 位卡号

## 遗留问题 / 待确认
- 无
