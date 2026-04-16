# 路由结构

## 整体路由树

```
/                          → HomePage（模式选择）
├── /admin                 → AdminLayout（侧边导航外壳）
│   ├── index              → IssuePage（默认）
│   ├── /admin/issue       → IssuePage（发卡）
│   ├── /admin/deposit     → DepositPage（存款）
│   ├── /admin/loss        → LossPage（挂失管理）
│   ├── /admin/cancel      → CancelPage（注销）
│   ├── /admin/statistics  → StatisticsPage（汇总统计）
│   └── /admin/windows     → WindowsPage（窗口管理）
└── /window                → MealPage（窗口机，独立全屏）
```

## 跳转关系

```mermaid
graph TD
    Home[/ 首页] -->|点击管理端| Admin[/admin/issue 发卡]
    Home -->|点击窗口机| Window[/window 窗口机]
    Admin -->|侧边菜单| Deposit[/admin/deposit]
    Admin -->|侧边菜单| Loss[/admin/loss]
    Admin -->|侧边菜单| Cancel[/admin/cancel]
    Admin -->|侧边菜单| Stats[/admin/statistics]
    Admin -->|侧边菜单| Windows[/admin/windows]
    Deposit -->|返回首页按钮| Home
    Window -->|返回首页按钮| Home
```

## 关键设计决策
- 管理端所有子页面共用 AdminLayout，通过 React Router Outlet 渲染
- 窗口机（MealPage）不嵌套在 AdminLayout 内，是独立路由，独立全屏深色风格
- /admin 默认重定向到 /admin/issue（index route）
