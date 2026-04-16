# AdminLayout（管理端布局）

## 作用
- 管理端所有页面的外层容器，提供固定侧边导航和顶部标题栏

## 结构
- Sider（宽 180px，固定定位，深色主题）
  - 顶部品牌区："食堂饭卡 / 管理端"
  - Menu（inline 模式）：6 个导航项
  - 底部「返回首页」按钮 → navigate('/')
- Header（白色，高 64px）：显示当前页面名称（根据路径自动匹配）
- Content（浅灰背景 #f5f5f5）：渲染 `<Outlet />`

## 导航菜单项

| key | 图标 | 标签 |
|-----|------|------|
| /admin/issue | CreditCardOutlined | 发卡 |
| /admin/deposit | WalletOutlined | 存款 |
| /admin/loss | WarningOutlined | 挂失管理 |
| /admin/cancel | StopOutlined | 注销 |
| /admin/statistics | BarChartOutlined | 汇总统计 |
| /admin/windows | AppstoreOutlined | 窗口管理 |

## 用户操作
- 点击菜单项 → navigate(key)，activeKey 通过 useLocation 路径前缀匹配
- 点击「返回首页」→ navigate('/')

## 状态变化
- 无本地 state，selectedKey 由 location.pathname 推导

## 接口依赖
- 无

## 备注 / 待确认
- Sider 使用 position: fixed，Content 区域用 marginLeft: 180 补偿
