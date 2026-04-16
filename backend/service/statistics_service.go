package service

import (
	"backend/repository"
	"time"

	"github.com/rs/zerolog/log"
)

// StatisticsService 统计业务逻辑
type StatisticsService struct {
	cardRepo *repository.CardRepository
}

// NewStatisticsService 创建 StatisticsService 实例
func NewStatisticsService(cardRepo *repository.CardRepository) *StatisticsService {
	return &StatisticsService{cardRepo: cardRepo}
}

// MealRevenueResult 本餐售饭总收入结果
type MealRevenueResult struct {
	TotalRevenue int64
}

// GetMealRevenue 统计指定时间范围内的消费总收入
func (s *StatisticsService) GetMealRevenue(start, end time.Time) (*MealRevenueResult, error) {
	total, err := s.cardRepo.SumTransactionsByTimeRange(start, end)
	if err != nil {
		return nil, err
	}
	return &MealRevenueResult{TotalRevenue: total}, nil
}

// WindowRevenueItem 单个窗口收入
type WindowRevenueItem struct {
	WindowID   int64
	WindowName string
	Revenue    int64
}

// WindowRevenueResult 各窗口收入结果
type WindowRevenueResult struct {
	Windows []WindowRevenueItem
}

// GetWindowRevenue 统计指定时间范围内各窗口的消费收入
func (s *StatisticsService) GetWindowRevenue(start, end time.Time) (*WindowRevenueResult, error) {
	rows, err := s.cardRepo.SumTransactionsByWindowAndTimeRange(start, end)
	if err != nil {
		return nil, err
	}

	items := make([]WindowRevenueItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, WindowRevenueItem{
			WindowID:   int64(row.WindowID),
			WindowName: row.WindowName,
			Revenue:    row.Revenue,
		})
	}
	return &WindowRevenueResult{Windows: items}, nil
}

// DepositDetailHolder 持卡人存款明细
type DepositDetailHolder struct {
	HolderID    int64
	HolderName  string
	IDNumber    string
	Deposits    []DepositDetailEntry
	TotalAmount int64
}

// DepositDetailEntry 单条存款记录
type DepositDetailEntry struct {
	ID        int64
	CardNo    string
	Amount    int64
	CreatedAt time.Time
}

// DepositDetailsResult 各持卡人存款明细结果（含分页信息）
type DepositDetailsResult struct {
	Holders  []DepositDetailHolder
	Total    int64
	Page     int
	PageSize int
}

// GetDepositDetails 获取各持卡人存款明细，可选时间范围，支持分页
func (s *StatisticsService) GetDepositDetails(start, end *time.Time, page, pageSize int) (*DepositDetailsResult, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	details, total, err := s.cardRepo.GetDepositDetails(start, end, page, pageSize)
	if err != nil {
		return nil, err
	}

	holders := make([]DepositDetailHolder, 0, len(details))
	for _, d := range details {
		entries := make([]DepositDetailEntry, 0, len(d.Deposits))
		for _, dep := range d.Deposits {
			entries = append(entries, DepositDetailEntry{
				ID:        int64(dep.ID),
				CardNo:    dep.CardNo,
				Amount:    dep.Amount,
				CreatedAt: dep.CreatedAt,
			})
		}
		holders = append(holders, DepositDetailHolder{
			HolderID:    int64(d.HolderID),
			HolderName:  d.HolderName,
			IDNumber:    d.IDNumber,
			Deposits:    entries,
			TotalAmount: d.TotalAmount,
		})
	}
	return &DepositDetailsResult{
		Holders:  holders,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

// DepositSummaryResult 本日/本月存款汇总
type DepositSummaryResult struct {
	TodayTotal int64
	MonthTotal int64
}

// GetDepositSummary 获取今日和本月的存款总额
func (s *StatisticsService) GetDepositSummary() (*DepositSummaryResult, error) {
	now := time.Now()

	// 今日范围
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	todayEnd := todayStart.Add(24 * time.Hour)

	todayTotal, err := s.cardRepo.SumDepositsByTimeRange(todayStart, todayEnd)
	if err != nil {
		return nil, err
	}

	// 本月范围
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	// 下个月第一天即本月结束
	monthEnd := time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, now.Location())

	monthTotal, err := s.cardRepo.SumDepositsByTimeRange(monthStart, monthEnd)
	if err != nil {
		return nil, err
	}

	return &DepositSummaryResult{
		TodayTotal: todayTotal,
		MonthTotal: monthTotal,
	}, nil
}

// ActiveBalanceResult 流动资金总额结果
type ActiveBalanceResult struct {
	TotalBalance int64
}

// GetActiveBalance 获取所有 active 卡的余额总和
func (s *StatisticsService) GetActiveBalance() (*ActiveBalanceResult, error) {
	total, err := s.cardRepo.SumActiveBalance()
	if err != nil {
		return nil, err
	}
	return &ActiveBalanceResult{TotalBalance: total}, nil
}

// DailyReportWindowItem 日报中窗口明细
type DailyReportWindowItem struct {
	WindowID         int64
	WindowName       string
	Revenue          int64
	TransactionCount int
}

// DailyReportResult 日餐报表结果
type DailyReportResult struct {
	Date             string
	TotalRevenue     int64
	TransactionCount int
	Windows          []DailyReportWindowItem
}

// GetDailyReport 获取指定日期的日餐报表
func (s *StatisticsService) GetDailyReport(date string) (*DailyReportResult, error) {
	// 解析日期
	day, err := time.Parse("2006-01-02", date)
	if err != nil {
		log.Error().Str("code", ErrCodeValidationError).Str("date", date).Msg("业务错误")
		return nil, newBizError(ErrCodeValidationError, "日期格式无效，应为 YYYY-MM-DD")
	}

	start := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, time.Local)
	end := start.Add(24 * time.Hour)

	totalRevenue, transactionCount, windows, err := s.cardRepo.GetDailyReport(start, end)
	if err != nil {
		log.Error().Err(err).Str("date", date).Msg("获取日餐报表失败")
		return nil, err
	}

	items := make([]DailyReportWindowItem, 0, len(windows))
	for _, w := range windows {
		items = append(items, DailyReportWindowItem{
			WindowID:         int64(w.WindowID),
			WindowName:       w.WindowName,
			Revenue:          w.Revenue,
			TransactionCount: w.TransactionCount,
		})
	}

	return &DailyReportResult{
		Date:             date,
		TotalRevenue:     totalRevenue,
		TransactionCount: transactionCount,
		Windows:          items,
	}, nil
}

// MonthlyReportItem 年报月度明细
type MonthlyReportItem struct {
	Month            int
	Revenue          int64
	TransactionCount int
}

// YearlyReportResult 年餐报表结果
type YearlyReportResult struct {
	Year             int
	TotalRevenue     int64
	TransactionCount int
	Months           []MonthlyReportItem
}

// GetYearlyReport 获取指定年份的年餐报表
func (s *StatisticsService) GetYearlyReport(year int) (*YearlyReportResult, error) {
	if year <= 0 {
		log.Error().Str("code", ErrCodeValidationError).Int("year", year).Msg("业务错误")
		return nil, newBizError(ErrCodeValidationError, "年份无效")
	}

	totalRevenue, transactionCount, months, err := s.cardRepo.GetYearlyReport(year)
	if err != nil {
		log.Error().Err(err).Int("year", year).Msg("获取年餐报表失败")
		return nil, err
	}

	items := make([]MonthlyReportItem, 0, len(months))
	for _, m := range months {
		items = append(items, MonthlyReportItem{
			Month:            m.Month,
			Revenue:          m.Revenue,
			TransactionCount: m.TransactionCount,
		})
	}

	return &YearlyReportResult{
		Year:             year,
		TotalRevenue:     totalRevenue,
		TransactionCount: transactionCount,
		Months:           items,
	}, nil
}
