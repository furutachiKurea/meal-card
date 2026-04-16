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

// FindCardHolderByID 根据数据库自增 ID 查找持卡人，不存在时返回 gorm.ErrRecordNotFound
func (r *CardRepository) FindCardHolderByID(id uint) (*model.CardHolder, error) {
	var holder model.CardHolder
	err := r.db.First(&holder, id).Error
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

// FindCardByID 根据数据库自增 ID 查询卡片（含持卡人信息）
func (r *CardRepository) FindCardByID(id uint) (*model.Card, error) {
	var card model.Card
	err := r.db.Preload("CardHolder").First(&card, id).Error
	if err != nil {
		return nil, err
	}
	return &card, nil
}

// FindCardByCardNo 根据 16 位业务卡号查询卡片（含持卡人信息）
func (r *CardRepository) FindCardByCardNo(cardNo string) (*model.Card, error) {
	var card model.Card
	err := r.db.Preload("CardHolder").Where("card_no = ?", cardNo).First(&card).Error
	if err != nil {
		return nil, err
	}
	return &card, nil
}

// FindCurrentCardByIDNumber 根据证件号查询该持卡人当前有效卡（active 或 lost），不返回 cancelled 卡
func (r *CardRepository) FindCurrentCardByIDNumber(idNumber string) (*model.Card, error) {
	var card model.Card
	err := r.db.Preload("CardHolder").
		Joins("JOIN card_holders ON card_holders.id = cards.card_holder_id").
		Where("card_holders.id_number = ? AND cards.status != ?", idNumber, model.CardStatusCancelled).
		Order("cards.created_at DESC").
		First(&card).Error
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
	CardNo    string
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

// GetDepositDetails 获取各持卡人存款明细，支持可选时间范围和分页（以持卡人为单位）。
// page 从 1 开始，pageSize 为每页持卡人数量。
// 返回当前页持卡人列表和满足条件的持卡人总数。
func (r *CardRepository) GetDepositDetails(start, end *time.Time, page, pageSize int) (holders []HolderDepositDetail, total int64, err error) {
	// 构建基础过滤条件（关联 deposit_records → cards → card_holders）
	baseQuery := r.db.Model(&model.DepositRecord{}).
		Joins("LEFT JOIN cards ON cards.id = deposit_records.card_id").
		Joins("LEFT JOIN card_holders ON card_holders.id = cards.card_holder_id")
	if start != nil {
		baseQuery = baseQuery.Where("deposit_records.created_at >= ?", *start)
	}
	if end != nil {
		baseQuery = baseQuery.Where("deposit_records.created_at <= ?", *end)
	}

	// 统计满足条件的不重复持卡人总数
	err = baseQuery.Distinct("card_holders.id").Count(&total).Error
	if err != nil {
		return
	}

	// 分页查询持卡人 ID（按 card_holders.id 排序保证稳定分页）
	offset := (page - 1) * pageSize
	type holderIDRow struct {
		HolderID uint
	}
	var holderIDs []holderIDRow
	err = baseQuery.
		Select("card_holders.id as holder_id").
		Group("card_holders.id").
		Order("card_holders.id ASC").
		Limit(pageSize).Offset(offset).
		Scan(&holderIDs).Error
	if err != nil {
		return
	}
	if len(holderIDs) == 0 {
		holders = []HolderDepositDetail{}
		return
	}

	// 收集本页持卡人 ID
	ids := make([]uint, 0, len(holderIDs))
	for _, h := range holderIDs {
		ids = append(ids, h.HolderID)
	}

	// 查询本页持卡人的全部存款明细
	type rawRow struct {
		ID         uint
		CardNo     string
		Amount     int64
		CreatedAt  time.Time
		HolderID   uint
		HolderName string
		IDNumber   string
	}
	detailQuery := r.db.Model(&model.DepositRecord{}).
		Select("deposit_records.id, cards.card_no, deposit_records.amount, deposit_records.created_at, card_holders.id as holder_id, card_holders.name as holder_name, card_holders.id_number").
		Joins("LEFT JOIN cards ON cards.id = deposit_records.card_id").
		Joins("LEFT JOIN card_holders ON card_holders.id = cards.card_holder_id").
		Where("card_holders.id IN ?", ids)
	if start != nil {
		detailQuery = detailQuery.Where("deposit_records.created_at >= ?", *start)
	}
	if end != nil {
		detailQuery = detailQuery.Where("deposit_records.created_at <= ?", *end)
	}

	var rows []rawRow
	if err = detailQuery.Scan(&rows).Error; err != nil {
		return
	}

	// 按持卡人分组，保持 holderIDs 返回的顺序
	holderMap := make(map[uint]*HolderDepositDetail, len(ids))
	for _, row := range rows {
		if _, ok := holderMap[row.HolderID]; !ok {
			holderMap[row.HolderID] = &HolderDepositDetail{
				HolderID:   row.HolderID,
				HolderName: row.HolderName,
				IDNumber:   row.IDNumber,
				Deposits:   []DepositDetailItem{},
			}
		}
		h := holderMap[row.HolderID]
		h.Deposits = append(h.Deposits, DepositDetailItem{
			ID:        row.ID,
			CardNo:    row.CardNo,
			Amount:    row.Amount,
			CreatedAt: row.CreatedAt,
		})
		h.TotalAmount += row.Amount
	}

	holders = make([]HolderDepositDetail, 0, len(ids))
	for _, id := range ids {
		if h, ok := holderMap[id]; ok {
			holders = append(holders, *h)
		}
	}
	return
}

// GetHolderDeposits 获取指定持卡人的存款明细，支持可选时间范围和分页。
// page 从 1 开始，pageSize 为每页记录数。
// 返回当前页存款记录列表和满足条件的总记录数。
func (r *CardRepository) GetHolderDeposits(holderID uint, start, end *time.Time, page, pageSize int) (deposits []DepositDetailItem, total int64, err error) {
	baseQuery := r.db.Model(&model.DepositRecord{}).
		Joins("LEFT JOIN cards ON cards.id = deposit_records.card_id").
		Joins("LEFT JOIN card_holders ON card_holders.id = cards.card_holder_id").
		Where("card_holders.id = ?", holderID)
	if start != nil {
		baseQuery = baseQuery.Where("deposit_records.created_at >= ?", *start)
	}
	if end != nil {
		baseQuery = baseQuery.Where("deposit_records.created_at <= ?", *end)
	}

	err = baseQuery.Count(&total).Error
	if err != nil {
		return
	}

	offset := (page - 1) * pageSize
	type rawRow struct {
		ID        uint
		CardNo    string
		Amount    int64
		CreatedAt time.Time
	}
	var rows []rawRow
	err = baseQuery.
		Select("deposit_records.id, cards.card_no, deposit_records.amount, deposit_records.created_at").
		Order("deposit_records.created_at DESC").
		Limit(pageSize).Offset(offset).
		Scan(&rows).Error
	if err != nil {
		return
	}

	deposits = make([]DepositDetailItem, 0, len(rows))
	for _, row := range rows {
		deposits = append(deposits, DepositDetailItem{
			ID:        row.ID,
			CardNo:    row.CardNo,
			Amount:    row.Amount,
			CreatedAt: row.CreatedAt,
		})
	}
	return
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
