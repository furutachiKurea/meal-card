---
name: bugfix-recorder
description: 当开发过程中发现 bug、定位根因、修改修复方案、补充兼容处理、记录回归结果时使用。用于沉淀 bug 的现象、触发条件、根因、修复方式和影响范围。前端内容写到 frontend/docs/bugfixes，下后端内容写到 backend/docs/bugfixes。这个是为我们课设服务的，你必须及时的调用这个 skill
---

# Purpose

记录开发过程中遇到的 bug，以及这个 bug 是怎么出现的、怎么定位的、怎么修掉的。

# When to use

在这些情况下使用：

1. 发现一个新 bug
2. bug 的触发条件被确认
3. bug 的根因被定位出来
4. 修复方案已经明确
5. 修复过程中增加了兼容逻辑、兜底逻辑或约束
6. 修完后需要补一条简短回归结果
7. 同一个 bug 出现了新的变体，值得补充说明

# Output location

- 前端 bug：`frontend/docs/bugfixes/`
- 后端 bug：`backend/docs/bugfixes/`
- 如果确实是跨前后端的问题，优先写到根因所在一侧；如果暂时不能判断，就先写到 `docs/bugfixes/` 或项目公共目录，后面再挪

如果目录不存在，先创建。

# What to record

只记录这些关键信息：

- bug 现象是什么
- 在什么条件下会触发
- 影响范围是什么
- 根因是什么
- 修复思路是什么
- 实际改了哪些点
- 修完后结果如何
- 还有没有遗留问题

# Workflow

1. 先确认这是值得记录的 bug，而不是临时代码没写完导致的正常中间状态。
2. 先写现象和触发条件，再写根因。
3. 根因没确定时，不要硬写；可以先留空或标记待确认。
4. 修复后补上"实际改动"和"结果"。
5. 如果同一个 bug 改了多次，直接更新原文件，不要散成很多份。
6. 不要把它写成测试报告，也不要写成流水账。

# File suggestion

建议一个 bug 一个文件，例如：

- `frontend/docs/bugfixes/login-form-submit-empty.md`
- `backend/docs/bugfixes/order-status-race-condition.md`

# Output template

```md
# <bug 名称>

## 现象
- ...

## 触发条件
- ...
- ...

## 影响范围
- ...

## 根因
- ...

## 修复思路
- ...

## 实际改动
- ...
- ...

## 修复结果
- ...

## 遗留问题 / 待确认
- ...
```
