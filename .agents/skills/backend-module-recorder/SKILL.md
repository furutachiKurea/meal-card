---
name: backend-module-recorder
description: 当新增后端模块、修改模块职责、调整核心流程、梳理状态流转、沉淀关键逻辑时使用。只处理后端内容，并把文档写到 backend/docs/course-docs 下。这个是为我们课设服务的，你必须及时的调用这个 skill
---

# Purpose

记录后端模块是怎么分工的、关键逻辑怎么走、状态怎么变化。

# When to use

在这些情况下使用：

1. 新增后端模块
2. 修改模块职责
3. 核心业务流程变化
4. 状态流转被明确
5. 某段关键逻辑值得单独说明
6. 某个实现决策已经稳定

# Output location

统一写到：

- `backend/docs/course-docs/modules/`
- `backend/docs/course-docs/flows/`

如果目录不存在，先创建。

# What to record

- 这个模块负责什么
- 它不负责什么
- 输入是什么
- 输出是什么
- 内部关键步骤是什么
- 状态如何变化
- 异常情况怎么处理
- 哪些点最容易出错

# Workflow

1. 先确认这次变化是不是模块级或流程级变化。
2. 不要只是复述代码 diff。
3. 要写清职责和边界。
4. 如果是重构，重点写"职责怎么变了"。
5. 如果只是小修小补且不影响理解，可以不记。

# File suggestion

- `backend/docs/course-docs/modules/<module-name>.md`
- `backend/docs/course-docs/flows/<flow-name>.md`

# Output template

```md
# <名称>

## 作用
- ...

## 职责边界
- 负责:
- 不负责:

## 输入
- ...

## 输出
- ...

## 核心流程
1. ...
2. ...
3. ...

## 状态 / 数据变化
- ...

## 异常处理
- ...

## 关键实现点
- ...

## 待确认
- ...
```

# Rules

1. 只写后端关键信息。
2. 不写整理型字段。
3. 不写测试内容。
4. 不写"优化了一下逻辑"这种空话，要写清具体变了什么。
