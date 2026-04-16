# WindowsPage（窗口管理页）

## 作用
- 管理员查看当前已创建的窗口列表，并可新建窗口

## 结构
- 单 Card 布局，最大宽 520px，居中
- 窗口列表区：Table（ID / 窗口名称），loading 态
- 新建窗口区：inline Form（窗口名称 Input + 创建按钮）

## 用户操作
1. 页面加载 → 自动调用 fetchWindows 获取列表
2. 输入窗口名称，点击「创建」→ 创建成功后表单重置，重新加载列表
3. 创建失败 → 底部 Alert 显示错误

## 状态变化
```
页面挂载(useEffect) → fetchWindows → windows 填入 / error 填入
提交创建 → creating=true → 成功: form 重置，重新 fetchWindows
                          → 失败: error 填入
```

## 接口依赖
- `GET /api/windows`（listWindows）—— 加载列表 & 创建后刷新
- `POST /api/windows`（createWindow）
  - 入参：name（字符串）

## 备注 / 待确认
- loading 和 creating 是两个独立 state：loading 控制 Table 加载态，creating 控制创建按钮 loading 态
- 窗口名称为必填项，无其他格式校验
- 此页面创建的窗口会被 MealPage 的窗口选择 Select 使用
