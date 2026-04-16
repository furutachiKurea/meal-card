# HomePage（首页）

## 作用
- 系统入口页，供用户选择进入「管理端」或「窗口机」两种使用模式

## 结构
- 居中布局，渐变蓝色背景
- 标题：食堂饭卡管理系统
- 两张可点击 Card：管理端（蓝色调）/ 窗口机（绿色调）

## 用户操作
- 点击「管理端」Card → navigate('/admin/issue')
- 点击「窗口机」Card → navigate('/window')

## 状态变化
- 无本地状态，纯导航跳转

## 接口依赖
- 无

## 备注 / 待确认
- 路由 /admin 默认重定向到 /admin/issue（App.jsx 中 index route 指向 IssuePage）
