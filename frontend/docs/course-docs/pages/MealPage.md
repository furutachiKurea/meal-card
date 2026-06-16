# MealPage（窗口机就餐页）

## 作用
- 模拟食堂窗口机：刷卡验证 → 显示余额 → 输入消费金额 → 结算。专为操作员站立使用场景设计，采用深色大字体界面。

## 结构
- 全屏深色背景（#0a1628）
- 顶栏：标题「就餐窗口机」/ 窗口选择 Select / 返回首页按钮
- 步骤条（Steps，3 步）：刷卡验证 → 输入金额 → 结算完成
- 主内容区（居中，最大宽 560px）：根据当前步骤条件渲染不同 Card

## 步骤内容

### 步骤 0：刷卡输入
- 大号 Input（fontSize 22，height 56）自动聚焦
- 「刷卡」按钮或回车触发查询

### 警报区（覆盖步骤 0）
- 触发条件：卡已注销 / 已挂失 / 非本单位卡（404）
- 红色边框 Card（#2d0a0a），48px WarningFilled 图标，22px 红色警报文字
- 「重新刷卡」按钮调用 handleReset

### 步骤 1：余额 + 金额输入
- 余额大显示：64px 蓝色（#4fc3f7）字体
- 金额 InputNumber（fontSize 28，height 64）自动聚焦
- 「确认结算」按钮（block，height 56，fontSize 20）

### 步骤 2：结算完成
- 绿色边框 Card（#061a0a）
- 本次消费（36px 橙色）/ 扣款后余额（36px 绿色）并排展示
- 「下一位」按钮调用 handleReset，回到步骤 0

## 用户操作
1. 进入页面时自动加载窗口列表
2. 若 URL 带 `?id=N` 参数，自动锁定到该窗口（隐藏选择器）；否则默认选中第一个窗口，可手动切换
3. 输入卡号 → 刷卡 → 验证
4. 验证通过 → 余额显示 → 输入消费金额 → 确认结算
5. 结算成功 → 展示本次消费 + 新余额 → 「下一位」重置

## 状态变化
```
初始: step=0
刷卡成功: cardInfo 填入 → step=1
刷卡失败(挂失/注销/404): alarm 填入（警报覆盖）
结算成功: txResult 填入 → step=2；cardInfo.balance 本地更新
handleReset: 所有状态清空 → step=0
```

currentStep 计算规则：`txResult ? 2 : cardInfo ? 1 : 0`

## 接口依赖
- `GET /api/windows`（listWindows）—— 页面加载时获取窗口列表
- `GET /api/cards/{cardNo}`（getCard）—— 刷卡验证（路径参数为 16 位 cardNo）
- `POST /api/cards/{cardNo}/transactions`（createTransaction）
  - 入参：windowId, amount（分）
  - 金额处理：显示元，提交前 × 100 转分

## 备注 / 待确认
- 警报文字：
  - 已注销：「警报：此卡已注销，禁止就餐！」
  - 已挂失：「警报：此卡已挂失，禁止就餐！」
  - 非本单位卡（HTTP 404）：「警报：此卡非本单位所发，禁止就餐！」
- 窗口机页面不使用 AdminLayout，是独立全屏路由 `/window`
- 深色主题设计意图：模拟实体窗口机终端，便于站立操作
- 窗口绑定：URL 参数 `/window?id=1` 锁定窗口，真实部署时每台机器配固定 URL
