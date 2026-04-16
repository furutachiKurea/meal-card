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
- 在 `style` prop 上设置 `background` 只作用于外层 wrapper，内部 `<input>` 元素不继承该背景色
- 内部 `<input>` 元素继承了页面深色背景（`#0a1628`），导致背景变黑、文字（`color: '#e8f4fd'`）不可见
- 对比：普通 `Input` 组件的 `style.background` 可以直接作用于 input 元素本身，所以卡号输入框不受影响

## 修复思路
- 使用 Ant Design v5 提供的 `styles` prop，通过 `styles.input` 直接设置内部 input 元素的样式，使背景色和文字颜色正确应用

## 实际改动
- `frontend/src/pages/MealPage.jsx` 第 266-281 行：为 `InputNumber` 新增 `styles` prop
  ```jsx
  styles={{
    input: {
      background: '#061020',
      color: '#e8f4fd',
    },
  }}
  ```

## 修复结果
- 消费金额输入框背景色正确显示为深蓝色（`#061020`），文字清晰可见

## 遗留问题 / 待确认
- 无
