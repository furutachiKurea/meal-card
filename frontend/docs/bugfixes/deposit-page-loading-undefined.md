# DepositPage loading 变量未定义

## 现象
- 存款页面查询按钮和存款按钮的 `loading` prop 引用了未定义的变量，导致按钮 loading 状态永远不生效

## 触发条件
- 打开存款页面，点击查询或存款按钮时，按钮不会显示加载中状态
- 在严格模式或构建工具静态分析下可能直接报错

## 影响范围
- 仅影响 DepositPage 的用户体验（无 loading 反馈），不影响功能正确性

## 根因
- 之前将各页面共享的单一 `loading` 状态拆分为独立状态时，DepositPage 定义了 `queryLoading` 和 `depositLoading` 两个变量，但 JSX 中的 `<Button loading={loading}>` 没有同步更新，仍然引用旧的 `loading` 变量名

## 修复思路
- 查询按钮使用 `queryLoading`，存款按钮使用 `depositLoading`，与 state 定义对齐

## 实际改动
- `frontend/src/pages/DepositPage.jsx`：查询按钮 `loading={loading}` → `loading={queryLoading}`，存款按钮 `loading={loading}` → `loading={depositLoading}`

## 修复结果
- `pnpm build` 通过，无未定义变量警告
- 两个按钮各自独立显示 loading 状态，互不干扰

## 遗留问题 / 待确认
- 无
