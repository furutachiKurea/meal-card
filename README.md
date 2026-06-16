# 食堂饭卡管理系统

课设项目，前后端分离架构，实现发卡、存款、就餐消费、挂失、注销、汇总统计 6 项功能。

[交互式代码讲解](https://furutachikurea.github.io/meal-card/)

## 环境准备

- Go 1.26+
- Node.js 18+、pnpm 8+

## 快速启动

```bash
# 终端 1：后端（监听 :8080，首次启动自动创建默认窗口）
make dev-backend

# 终端 2：前端（监听 :5173，自动代理 /api 到后端）
make dev-frontend

# 运行后端单测
make test
```

前端通过 Vite 代理将 `/api` 请求转发到后端 `:8080`，无需配置 CORS。

## 系统角色与入口

| 角色 | 入口 | 说明 |
|------|------|------|
| 管理处操作员 | `/admin/*` | 发卡、存款、挂失、注销、统计、窗口管理 |
| 窗口工作人员 | `/window?id=N` | 刷卡验证、输入金额、确认结算 |
| 学生（顾客屏） | `/window/customer?id=N` | 只读查看余额和消费结果 |

窗口操作端与顾客屏通过 BroadcastChannel 实时同步状态，部署时同一窗口机的两块屏幕各打开对应 URL。

## 学籍验证（Mock 外部服务）

系统设计上依赖外部学籍验证服务接口（`StudentValidator`），当前实现使用内嵌 stub 数据 mock 该外部调用。证件号统一为 12 位数字：

| 证件号 | 姓名 | 类型 |
|------|------|------|
| `202100010001` | 张三 | 学生 |
| `202100010002` | 李四 | 学生 |
| `202100010003` | 王五 | 学生 |
| `202100010004` | 赵六 | 学生 |
| `202100010005` | 陈七 | 学生 |
| `200900010001` | 刘老师 | 教职工 |
| `200900010002` | 孙老师 | 教职工 |
| `200900010003` | 周主任 | 教职工 |

如需增加，编辑 `backend/client/student_client.go` 中的 `studentDB` map。如需对接真实学籍系统，实现 `service.StudentValidator` 接口并替换注入即可。

## 项目结构

```
meal-card/
├── backend/           Go 后端（Echo + GORM + SQLite）
│   ├── handler/       HTTP 处理层
│   ├── service/       业务逻辑层
│   ├── repository/    数据访问层
│   ├── model/         GORM 模型
│   ├── client/        学籍验证实现（mock）
│   ├── db/            数据库初始化 + 种子数据
│   └── docs/          后端文档
├── frontend/          React 前端（Vite + Ant Design）
│   └── src/pages/     8 个业务页面 + 404
├── docs/              全局文档
│   ├── api/           OpenAPI 契约
│   ├── prd.md         产品需求
│   ├── architecture.md 架构设计
│   └── progress.md    迭代进度
└── Makefile
```
