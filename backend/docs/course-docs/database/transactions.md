# transactions（消费记录）

## 作用
- 记录每次就餐扣款，用于收入统计和报表

## 设计原因
- 每次就餐结算产生一条不可变记录
- 同时关联饭卡和窗口，支持"各窗口收入"、"本餐售饭总收入"等多维度统计
- CreatedAt 即消费时间，按时间范围聚合即可生成日餐/年餐报表

## 字段
- ID:
  - 含义: 记录主键
  - 类型: uint
  - 是否必填: 自动生成
  - 默认值: 自增
- CardID:
  - 含义: 消费的饭卡
  - 类型: uint
  - 是否必填: 是
  - 默认值: 无
- WindowID:
  - 含义: 消费窗口
  - 类型: uint
  - 是否必填: 是
  - 默认值: 无
- Amount:
  - 含义: 消费金额
  - 类型: int64
  - 是否必填: 是
  - 默认值: 无
  - 备注: 单位：分，必须为正数（业务层校验）
- CreatedAt:
  - 含义: 消费时间
  - 类型: time.Time
  - 是否必填: GORM 自动填充
  - 默认值: 当前时间

## 主键 / 索引 / 约束
- 主键: ID
- 普通索引: CardID（按卡查询消费记录）、WindowID（按窗口统计收入）
- NOT NULL: CardID, WindowID, Amount
- 外键: CardID → cards.ID, WindowID → windows.ID

## 表关系
- 属于 cards（多对一，通过 CardID）
- 属于 windows（多对一，通过 WindowID）
