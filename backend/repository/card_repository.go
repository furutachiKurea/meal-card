// Package repository 提供数据访问层，封装 GORM 操作
package repository

import (
	"backend/model"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// CardRepository 饭卡及相关数据的增删改查
type CardRepository struct {
	db *gorm.DB
}

// NewCardRepository 创建 CardRepository 实例
func NewCardRepository(db *gorm.DB) *CardRepository {
	return &CardRepository{db: db}
}

// FindCardHolderByIDNumber 根据证件号查找持卡人，不存在时返回 gorm.ErrRecordNotFound
func (r *CardRepository) FindCardHolderByIDNumber(idNumber string) (*model.CardHolder, error) {
	var holder model.CardHolder
	err := r.db.Where("id_number = ?", idNumber).First(&holder).Error
	if err != nil {
		return nil, err
	}
	return &holder, nil
}

// CreateCardHolder 创建持卡人记录
func (r *CardRepository) CreateCardHolder(holder *model.CardHolder) error {
	return r.db.Create(holder).Error
}

// FindActiveCardByHolderID 查找持卡人名下 active 状态的卡
func (r *CardRepository) FindActiveCardByHolderID(holderID uint) (*model.Card, error) {
	var card model.Card
	err := r.db.Where("card_holder_id = ? AND status = ?", holderID, model.CardStatusActive).First(&card).Error
	if err != nil {
		return nil, err
	}
	return &card, nil
}

// FindLostCardByHolderID 查找持卡人名下 lost 状态的卡
func (r *CardRepository) FindLostCardByHolderID(holderID uint) (*model.Card, error) {
	var card model.Card
	err := r.db.Where("card_holder_id = ? AND status = ?", holderID, model.CardStatusLost).First(&card).Error
	if err != nil {
		return nil, err
	}
	return &card, nil
}

// CreateCard 创建饭卡
func (r *CardRepository) CreateCard(card *model.Card) error {
	return r.db.Create(card).Error
}

// FindCardByID 根据卡号查询卡片（含持卡人信息）
func (r *CardRepository) FindCardByID(id uint) (*model.Card, error) {
	var card model.Card
	err := r.db.Preload("CardHolder").First(&card, id).Error
	if err != nil {
		return nil, err
	}
	return &card, nil
}

// UpdateCard 更新卡片信息
func (r *CardRepository) UpdateCard(card *model.Card) error {
	return r.db.Save(card).Error
}

// CreateDepositRecord 创建存款记录
func (r *CardRepository) CreateDepositRecord(record *model.DepositRecord) error {
	return r.db.Create(record).Error
}

// CreateTransaction 创建消费记录
func (r *CardRepository) CreateTransaction(tx *model.Transaction) error {
	return r.db.Create(tx).Error
}

// SumTransactionsByTimeRange 统计指定时间范围内的消费总额
func (r *CardRepository) SumTransactionsByTimeRange(start, end time.Time) (int64, error) {
	var total int64
	err := r.db.Model(&model.Transaction{}).
		Where("created_at >= ? AND created_at <= ?", start, end).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&total).Error
	return total, err
}

// WindowRevenue 单个窗口的收入统计
type WindowRevenue struct {
	WindowID   uint
	WindowName string
	Revenue    int64
}

// SumTransactionsByWindowAndTimeRange 按窗口统计指定时间范围内的消费收入
func (r *CardRepository) SumTransactionsByWindowAndTimeRange(start, end time.Time) ([]WindowRevenue, error) {
	var results []WindowRevenue
	err := r.db.Model(&model.Transaction{}).
		Select("transactions.window_id, windows.name as window_name, COALESCE(SUM(transactions.amount), 0) as revenue").
		Joins("LEFT JOIN windows ON windows.id = transactions.window_id").
		Where("transactions.created_at >= ? AND transactions.created_at <= ?", start, end).
		Group("transactions.window_id").
		Scan(&results).Error
	return results, err
}

// DepositDetailItem 存款明细中单条存款信息
type DepositDetailItem struct {
	ID        uint
	CardID    uint
	Amount    int64
	CreatedAt time.Time
}

// HolderDepositDetail 持卡人存款明细
type HolderDepositDetail struct {
	HolderID    uint
	HolderName  string
	IDNumber    string
	Deposits    []DepositDetailItem
	TotalAmount int64
}

// GetDepositDetails 获取各持卡人存款明细，可选时间范围
func (r *CardRepository) GetDepositDetails(start, end *time.Time) ([]HolderDepositDetail, error) {
	// 先查询满足条件的存款记录（关联 card 和 card_holder）
	type rawRow struct {
		ID         uint
		CardID     uint
		Amount     int64
		CreatedAt  time.Time
		HolderID   uint
		HolderName string
		IDNumber   string
	}

	query := r.db.Model(&model.DepositRecord{}).
		Select("deposit_records.id, deposit_records.card_id, deposit_records.amount, deposit_records.created_at, card_holders.id as holder_id, card_holders.name as holder_name, card_holders.id_number").
		Joins("LEFT JOIN cards ON cards.id = deposit_records.card_id").
		Joins("LEFT JOIN card_holders ON card_holders.id = cards.card_holder_id")

	if start != nil {
		query = query.Where("deposit_records.created_at >= ?", *start)
	}
	if end != nil {
		query = query.Where("deposit_records.created_at <= ?", *end)
	}

	var rows []rawRow
	if err := query.Scan(&rows).Error; err != nil {
		return nil, err
	}

	// 按持卡人分组
	holderMap := make(map[uint]*HolderDepositDetail)
	order := []uint{}
	for _, row := range rows {
		if _, ok := holderMap[row.HolderID]; !ok {
			holderMap[row.HolderID] = &HolderDepositDetail{
				HolderID:   row.HolderID,
				HolderName: row.HolderName,
				IDNumber:   row.IDNumber,
				Deposits:   []DepositDetailItem{},
			}
			order = append(order, row.HolderID)
		}
		h := holderMap[row.HolderID]
		h.Deposits = append(h.Deposits, DepositDetailItem{
			ID:        row.ID,
			CardID:    row.CardID,
			Amount:    row.Amount,
			CreatedAt: row.CreatedAt,
		})
		h.TotalAmount += row.Amount
	}

	result := make([]HolderDepositDetail, 0, len(order))
	for _, id := range order {
		result = append(result, *holderMap[id])
	}
	return result, nil
}

// SumDepositsByTimeRange 统计指定时间范围内的存款总额
func (r *CardRepository) SumDepositsByTimeRange(start, end time.Time) (int64, error) {
	var total int64
	err := r.db.Model(&model.DepositRecord{}).
		Where("created_at >= ? AND created_at <= ?", start, end).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&total).Error
	return total, err
}

// SumActiveBalance 统计所有 active 卡的余额总和
func (r *CardRepository) SumActiveBalance() (int64, error) {
	var total int64
	err := r.db.Model(&model.Card{}).
		Where("status = ?", model.CardStatusActive).
		Select("COALESCE(SUM(balance), 0)").
		Scan(&total).Error
	return total, err
}

// DailyWindowStat 日报窗口统计
type DailyWindowStat struct {
	WindowID         uint
	WindowName       string
	Revenue          int64
	TransactionCount int
}

// GetDailyReport 获取指定日期的消费汇总（含总额、笔数、各窗口明细）
func (r *CardRepository) GetDailyReport(start, end time.Time) (totalRevenue int64, transactionCount int, windows []DailyWindowStat, err error) {
	// 总额和笔数
	type summary struct {
		Total int64
		Count int
	}
	var s summary
	err = r.db.Model(&model.Transaction{}).
		Select("COALESCE(SUM(amount), 0) as total, COUNT(*) as count").
		Where("created_at >= ? AND created_at <= ?", start, end).
		Scan(&s).Error
	if err != nil {
		return
	}
	totalRevenue = s.Total
	transactionCount = s.Count

	// 按窗口分组
	err = r.db.Model(&model.Transaction{}).
		Select("transactions.window_id, windows.name as window_name, COALESCE(SUM(transactions.amount), 0) as revenue, COUNT(*) as transaction_count").
		Joins("LEFT JOIN windows ON windows.id = transactions.window_id").
		Where("transactions.created_at >= ? AND transactions.created_at <= ?", start, end).
		Group("transactions.window_id").
		Scan(&windows).Error
	return
}

// MonthlyTransactionStat 年报月度统计
type MonthlyTransactionStat struct {
	Month            int
	Revenue          int64
	TransactionCount int
}

// GetYearlyReport 获取指定年份按月汇总的消费数据
func (r *CardRepository) GetYearlyReport(year int) (totalRevenue int64, transactionCount int, months []MonthlyTransactionStat, err error) {
	type monthRow struct {
		Month            int
		Revenue          int64
		TransactionCount int
	}

	var rows []monthRow
	err = r.db.Model(&model.Transaction{}).
		Select("CAST(strftime('%m', created_at) AS INTEGER) as month, COALESCE(SUM(amount), 0) as revenue, COUNT(*) as transaction_count").
		Where("strftime('%Y', created_at) = ?", fmt.Sprintf("%04d", year)).
		Group("month").
		Order("month ASC").
		Scan(&rows).Error
	if err != nil {
		return
	}

	for _, row := range rows {
		totalRevenue += row.Revenue
		transactionCount += row.TransactionCount
		months = append(months, MonthlyTransactionStat{
			Month:            row.Month,
			Revenue:          row.Revenue,
			TransactionCount: row.TransactionCount,
		})
	}
	return
}
