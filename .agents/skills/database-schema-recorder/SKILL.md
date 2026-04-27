---
name: database-schema-recorder
description: 当新增表、修改表结构、补充字段含义、调整索引、增加约束、明确表关系或记录迁移说明时使用。只处理数据库设计内容，并把文档写到 backend/docs/course-docs 下。这个是为我们课设服务的，你必须及时的调用这个 skill
---

# Purpose

记录表结构和数据约束，避免后面只能靠看 migration、DDL 和代码猜数据库设计。

# When to use

在这些情况下使用：

1. 新增表
2. 新增字段
3. 修改字段类型或含义
4. 新增或修改索引
5. 新增约束
6. 明确表之间关系
7. 有结构迁移或数据迁移
8. 某个设计理由值得记下来

# Output location

统一写到：

- `backend/docs/course-docs/database/`

如果目录不存在，先创建。

# What to record

- 这张表是干什么的
- 为什么需要这张表或这次变更
- 字段有哪些
- 字段分别是什么意思
- 主键 / 唯一键 / 索引是什么
- 有哪些约束
- 和其他表什么关系
- 迁移时要注意什么

# Workflow

1. 先确认这次变化是否真的影响表结构或数据约束。
2. 不要只贴 DDL，要解释设计意图。
3. 如果是字段改动，要写清改动前后含义。
4. 如果有迁移风险，要明确写出来。
5. 不写接口或测试内容。

# File suggestion

- `backend/docs/course-docs/database/<table-name>.md`
- `backend/docs/course-docs/database/migrations/<migration-name>.md`

# Output template

```md
# <表名或迁移名>

## 作用
- ...

## 设计原因
- ...

## 字段
- 字段名:
  - 含义:
  - 类型:
  - 是否必填:
  - 默认值:
  - 备注:

## 主键 / 索引 / 约束
- ...

## 表关系
- ...

## 变更内容
- ...

## 迁移注意事项
- ...

## 待确认
- ...
```

# Rules

1. 不写整理型字段。
2. 重点记录字段语义和约束。
3. 不要只贴 SQL。
4. 如果没有实际结构变化，不要硬记一条。
