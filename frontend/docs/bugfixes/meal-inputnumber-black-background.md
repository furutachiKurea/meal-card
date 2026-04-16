# 就餐页消费金额输入框背景色为黑色、文字不可见

## 现象
- 就餐页（MealPage.jsx）步骤 1 的消费金额输入框（InputNumber）背景色显示为黑色，输入的文字几乎不可见

## 触发条件
- 进入就餐页，刷卡验证通过后，进入"输入金额"步骤时出现
- 页面整体背景为深色主题（`#0a1628`），卡片背景为 `#0d1f3c`

## 影响范围
- `frontend/src/pages/MealPage.jsx` 中 InputNumber 组件
- 仅影响就餐页消费金额输入框，不影响其他页面

## 根因
- Ant Design v5 的 `InputNumber` 是复合组件，外层是 wrapper div，内部才是真正的 `<input>` 元素
- 在 `style` prop 上设置 `background`/`color` 只作用于外层 wrapper，内部 `<input>` 元素不继承
- 内部 `<input>` 元素继承了页面深色背景（`#0a1628`），导致背景变黑、文字不可见

## 修复思路（第一次尝试 → 无效）
- 首先尝试使用 Ant Design v5 的 `styles` prop（`styles={{ input: { background, color } }}`）直接设置内部 input 样式
- 背景色可以生效，但 `color` 字段在 Ant Design 内部样式优先级下被覆盖，文字仍为黑色
- `styles.input.color` 方案对 `InputNumber` 无效，放弃

## 修复思路（最终方案）
- 给 InputNumber 加 `className="meal-amount-input"`
- 在组件 JSX 顶部注入内联 `<style>` 标签，用 CSS 类选择器 + `!important` 强制覆盖 Ant Design 内部样式：
  ```css
  .meal-amount-input .ant-input-number-input {
    color: #e8f4fd !important;
    background: #061020 !important;
  }
  ```
- `!important` 可绕过 Ant Design 内部样式优先级，确保颜色正确应用

## 实际改动
- `frontend/src/pages/MealPage.jsx`：
  - `return` 后第一个元素前插入 `<style>` 标签注入 `.meal-amount-input` 样式规则
  - InputNumber 加 `className="meal-amount-input"`
  - 删除原先无效的 `styles={{ input: { ... } }}` prop

## 修复结果
- 消费金额输入框背景色正确显示为深蓝色（`#061020`），文字颜色正确显示为浅蓝白（`#e8f4fd`），清晰可见

## 遗留问题 / 待确认
- 无
