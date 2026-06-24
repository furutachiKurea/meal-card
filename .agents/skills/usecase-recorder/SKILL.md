---
name: usecase-recorder
description: 当新增功能、修改功能行为、明确边界条件、补充业务规则、梳理状态约束时使用。用于记录当前已经确定的功能行为和规则，不记录编号、参与者之类的整理性信息。这个是为我们课设服务的，你必须及时的调用这个 skill
---

# Purpose

记录功能本身怎么工作、有哪些限制、有哪些异常或边界情况。

# When to use

在这些情况下使用：

1. 新增一个功能
2. 修改一个已有功能的行为
3. 某个功能的前置条件、成功条件、失败条件被明确
4. 新增业务规则
5. 新增状态约束
6. 某个边界条件或异常流程被讨论清楚

# Output location

- `docs/course-docs/usecases/`
- `frontend` 单侧内容可写到 `frontend/docs/course-docs/usecases/`
- `backend` 单侧内容可写到 `backend/docs/course-docs/usecases/`

如果目录不存在，先创建。

# What to record

只记录这些关键信息：

- 这是哪个功能
- 它现在应该怎么工作
- 成功路径是什么
- 失败或异常路径是什么
- 需要满足哪些前置条件
- 有哪些明确规则或限制
- 哪些地方还没定

# Workflow

1. 只提取当前已经确定的事实。
2. 不补充编号、参与者、频率这类整理性信息。
3. 如果只是修改已有功能，直接更新对应记录，不要重复造一份。
4. 如果规则跨前后端，写清规则本身，不写协作流程。
5. 不生成截图占位。

# File suggestion

一个功能一个文件，例如：

- `docs/course-docs/usecases/login.md`
- `frontend/docs/course-docs/usecases/feed.md`
- `backend/docs/course-docs/usecases/order-submit.md`

# Output template

```md
# <功能名>

## 当前行为
- ...

## 前置条件
- ...

## 成功路径
1. ...
2. ...

## 异常 / 边界情况
- ...
- ...

## 规则 / 限制
- ...
- ...

## 待确认
- ...
```

# Rules

1. 用中文。
2. 不写空话，不写模板字段。
3. 只写对理解功能有用的信息。
4. 不把实现细节写进这里，除非它已经影响功能行为。
