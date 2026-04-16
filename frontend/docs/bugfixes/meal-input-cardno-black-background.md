# 就餐页刷卡号输入框背景色为黑色、文字不可见

## 现象
- 就餐页（MealPage.jsx）步骤 0 的卡号输入框（Input）背景色显示为黑色，输入的文字几乎不可见

## 触发条件
- 进入就餐页，初始刷卡输入步骤时出现
- 页面整体背景为深色主题（`#0a1628`），卡片背景为 `#0d1f3c`

## 影响范围
- `frontend/src/pages/MealPage.jsx` 中刷卡号 Input 组件
- 仅影响就餐页卡号输入框，不影响其他页面

## 根因
- 与 Task #8（InputNumber 黑色背景）根因相同：Ant Design v5 的 `Input` 也是复合组件，外层是 wrapper div，内部才是真正的 `<input>` 元素
- 在 `style` prop 上设置 `background` 只作用于外层 wrapper，内部 `<input>` 元素不继承该背景色
- 内部 `<input>` 元素继承了页面深色背景（`#0a1628`），导致背景变黑、文字（`color: '#e8f4fd'`）不可见

## 修复思路
- 使用 Ant Design v5 提供的 `styles` prop，通过 `styles.input` 直接设置内部 input 元素的样式

## 实际改动
- `frontend/src/pages/MealPage.jsx` 卡号 Input 组件新增 `styles` prop：
  ```jsx
  styles={{
    input: {
      background: '#061020',
      color: '#e8f4fd',
    },
  }}
  ```

## 修复结果
- 卡号输入框背景色正确显示为深蓝色（`#061020`），文字清晰可见，与消费金额 InputNumber 样式一致

## 遗留问题 / 待确认
- 无
