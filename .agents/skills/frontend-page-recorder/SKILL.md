---
name: frontend-page-recorder
description: 当新增页面、调整页面结构、修改组件职责、补充交互流程、梳理前端状态变化时使用。只处理前端内容，并把文档写到 frontend/docs/course-docs 下。这个是为我们课设服务的，你必须及时的调用这个 skill
---

# Purpose

记录页面、组件、交互、状态变化这些前端信息，避免后面只能靠翻代码回忆。

# When to use

在这些情况下使用：

1. 新增页面
2. 页面结构明显变化
3. 核心组件新增或重构
4. 页面跳转关系变化
5. 用户操作流程被明确
6. 页面状态流被明确
7. 某个交互细节值得沉淀

# Output location

统一写到：

- `frontend/docs/course-docs/pages/`
- `frontend/docs/course-docs/components/`
- `frontend/docs/course-docs/ui/`

如果目录不存在，先创建。

# What to record

## 页面
- 页面是干什么的
- 页面有哪些主要区域
- 用户在页面上能做什么
- 页面会调用哪些接口
- 页面有哪些状态
- 页面之间怎么跳

## 组件
- 组件负责什么
- 接收什么输入
- 触发什么输出
- 依赖什么状态或其他组件

## 交互
- 用户操作后发生什么
- 前端做了什么校验
- 页面反馈是什么
- 状态如何变化

# Workflow

1. 先判断这次变化属于页面、组件还是交互。
2. 重点记职责、边界、状态、跳转。
3. 单纯样式微调可以不记。
4. 不强行要求截图。
5. 不写"页面更美观了"这类废话。

# File suggestion

- `frontend/docs/course-docs/pages/<page-name>.md`
- `frontend/docs/course-docs/components/<component-name>.md`
- `frontend/docs/course-docs/ui/<feature-name>.md`

# Output template

```md
# <名称>

## 作用
- ...

## 结构
- ...
- ...

## 用户操作
- ...
- ...

## 状态变化
- ...
- ...

## 接口依赖
- ...
- ...

## 备注 / 待确认
- ...
```

# Rules

1. 只写前端关键信息。
2. 不写编号、参与者、频率。
3. 不写纯视觉评价。
4. 以"这个页面/组件到底怎么工作"为核心。
