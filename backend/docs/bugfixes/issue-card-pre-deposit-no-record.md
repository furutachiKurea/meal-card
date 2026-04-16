# 发卡时预存款不生成 DepositRecord

## 现象
- 通过发卡接口（POST /cards）传入 preDeposit > 0 成功发卡后，统计页面的存款明细中看不到该笔预存款记录
- 存款总额统计偏少，与实际流入金额不符

## 触发条件
- 调用 IssueCard 时 preDeposit 参数大于 0
- 即：每一次带预存款的正常发卡操作都会触发

## 影响范围
- 汇总统计：各持卡人存款明细缺少发卡时的预存款条目
- 日/月存款总额统计数据偏低
- 历史已发卡数据中所有 preDeposit > 0 的记录均受影响

## 根因
- `card_service.go` 的 `IssueCard` 方法在调用 `cardRepo.CreateCard` 后，只把 `preDeposit` 写入了卡的 `Balance` 字段，但没有同步调用 `cardRepo.CreateDepositRecord` 补建流水记录
- 对比 `Deposit`（充值）方法，充值时会先更新余额再写 DepositRecord，两处逻辑不一致

## 修复思路
- 在 `CreateCard` 成功后，判断 `preDeposit > 0`，若成立则构造一条 `DepositRecord` 并调用 `CreateDepositRecord` 写入

## 实际改动
- `backend/service/card_service.go`：`IssueCard` 方法末尾，`CreateCard` 成功后追加：
  ```go
  if preDeposit > 0 {
      record := &model.DepositRecord{CardID: newCard.ID, Amount: preDeposit}
      if err := s.cardRepo.CreateDepositRecord(record); err != nil {
          return nil, err
      }
  }
  ```
- `backend/service/card_service_test.go`：新增 `TestIssueCard_PreDepositCreatesDepositRecord` 测试，覆盖两个分支：
  - preDeposit > 0：验证 DepositRecord 确实存在，总金额正确
  - preDeposit == 0：验证不生成多余记录

## 修复结果
- `go test ./...` 全部通过，无回归

## 遗留问题 / 待确认
- 历史数据中已发卡但缺失的 DepositRecord 无法自动补齐（本系统不投入生产，忽略）
