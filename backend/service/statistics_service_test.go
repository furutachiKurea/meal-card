package service_test

import (
	"backend/db"
	"backend/model"
	"backend/repository"
	"backend/service"
	"testing"
	"time"
)

// setupStatisticsService 初始化内存数据库并返回 StatisticsService 和辅助用 CardService
func setupStatisticsService(t *testing.T) (*service.StatisticsService, *service.CardService, uint) {
	t.Helper()
	gormDB, err := db.Init(":memory:")
	if err != nil {
		t.Fatalf("初始化内存数据库失败: %v", err)
	}
	cardRepo := repository.NewCardRepository(gormDB)
	windowRepo := repository.NewWindowRepository(gormDB)

	// 创建一个测试窗口
	w := &model.Window{Name: "测试窗口"}
	if err := windowRepo.CreateWindow(w); err != nil {
		t.Fatalf("创建窗口失败: %v", err)
	}

	v := newFakeValidator("ID001", "ID002", "ID003")
	cardSvc := service.NewCardService(cardRepo, windowRepo, v)
	statsSvc := service.NewStatisticsService(cardRepo)
	return statsSvc, cardSvc, w.ID
}

func TestGetMealRevenue_EmptyRange(t *testing.T) {
	statsSvc, _, _ := setupStatisticsService(t)

	start := time.Now().Add(-1 * time.Hour)
	end := time.Now().Add(1 * time.Hour)

	result, err := statsSvc.GetMealRevenue(start, end)
	if err != nil {
		t.Fatalf("GetMealRevenue 失败: %v", err)
	}
	if result.TotalRevenue != 0 {
		t.Errorf("无消费时总收入应为 0，得到 %d", result.TotalRevenue)
	}
}

func TestGetMealRevenue_WithTransactions(t *testing.T) {
	statsSvc, cardSvc, windowID := setupStatisticsService(t)

	// 发卡并消费
	r, err := cardSvc.IssueCard("ID001", 5000)
	if err != nil {
		t.Fatalf("发卡失败: %v", err)
	}
	cardNo := r.Card.CardNo

	if _, err := cardSvc.CreateTransaction(cardNo, windowID, 300); err != nil {
		t.Fatalf("消费失败: %v", err)
	}
	if _, err := cardSvc.CreateTransaction(cardNo, windowID, 200); err != nil {
		t.Fatalf("消费失败: %v", err)
	}

	start := time.Now().Add(-1 * time.Hour)
	end := time.Now().Add(1 * time.Hour)

	result, err := statsSvc.GetMealRevenue(start, end)
	if err != nil {
		t.Fatalf("GetMealRevenue 失败: %v", err)
	}
	if result.TotalRevenue != 500 {
		t.Errorf("总收入期望 500，得到 %d", result.TotalRevenue)
	}
}

func TestGetActiveBalance(t *testing.T) {
	statsSvc, cardSvc, windowID := setupStatisticsService(t)

	// 发两张卡，余额分别为 500、300
	r1, _ := cardSvc.IssueCard("ID001", 500)
	r2, _ := cardSvc.IssueCard("ID002", 300)

	result, err := statsSvc.GetActiveBalance()
	if err != nil {
		t.Fatalf("GetActiveBalance 失败: %v", err)
	}
	if result.TotalBalance != 800 {
		t.Errorf("流动资金期望 800，得到 %d", result.TotalBalance)
	}

	// 注销一张卡后，流动资金应减少
	if _, err := cardSvc.CancelCard(r1.Card.CardNo); err != nil {
		t.Fatalf("注销失败: %v", err)
	}
	result2, err := statsSvc.GetActiveBalance()
	if err != nil {
		t.Fatalf("GetActiveBalance 失败: %v", err)
	}
	if result2.TotalBalance != 300 {
		t.Errorf("注销后流动资金期望 300，得到 %d", result2.TotalBalance)
	}

	// 消费后流动资金减少
	if _, err := cardSvc.CreateTransaction(r2.Card.CardNo, windowID, 100); err != nil {
		t.Fatalf("消费失败: %v", err)
	}
	result3, err := statsSvc.GetActiveBalance()
	if err != nil {
		t.Fatalf("GetActiveBalance 失败: %v", err)
	}
	if result3.TotalBalance != 200 {
		t.Errorf("消费后流动资金期望 200，得到 %d", result3.TotalBalance)
	}
}

func TestGetDepositSummary(t *testing.T) {
	statsSvc, cardSvc, _ := setupStatisticsService(t)

	// 发卡时没有存款记录（预存款不产生 DepositRecord），执行存款操作
	r, _ := cardSvc.IssueCard("ID001", 0)
	cardNo := r.Card.CardNo

	if _, err := cardSvc.Deposit(cardNo, 600); err != nil {
		t.Fatalf("充值失败: %v", err)
	}
	if _, err := cardSvc.Deposit(cardNo, 400); err != nil {
		t.Fatalf("充值失败: %v", err)
	}

	result, err := statsSvc.GetDepositSummary()
	if err != nil {
		t.Fatalf("GetDepositSummary 失败: %v", err)
	}
	if result.TodayTotal != 1000 {
		t.Errorf("今日存款期望 1000，得到 %d", result.TodayTotal)
	}
	if result.MonthTotal != 1000 {
		t.Errorf("本月存款期望 1000，得到 %d", result.MonthTotal)
	}
}

func TestGetDailyReport(t *testing.T) {
	statsSvc, cardSvc, windowID := setupStatisticsService(t)

	r, _ := cardSvc.IssueCard("ID001", 5000)
	cardNo := r.Card.CardNo

	if _, err := cardSvc.CreateTransaction(cardNo, windowID, 150); err != nil {
		t.Fatalf("消费失败: %v", err)
	}
	if _, err := cardSvc.CreateTransaction(cardNo, windowID, 250); err != nil {
		t.Fatalf("消费失败: %v", err)
	}

	today := time.Now().Format("2006-01-02")
	result, err := statsSvc.GetDailyReport(today)
	if err != nil {
		t.Fatalf("GetDailyReport 失败: %v", err)
	}
	if result.TotalRevenue != 400 {
		t.Errorf("日报总收入期望 400，得到 %d", result.TotalRevenue)
	}
	if result.TransactionCount != 2 {
		t.Errorf("日报笔数期望 2，得到 %d", result.TransactionCount)
	}
	if result.Date != today {
		t.Errorf("日期期望 %s，得到 %s", today, result.Date)
	}
}

func TestGetDailyReport_InvalidDate(t *testing.T) {
	statsSvc, _, _ := setupStatisticsService(t)

	_, err := statsSvc.GetDailyReport("invalid-date")
	if err == nil {
		t.Fatal("期望日期格式错误")
	}
	bizErr := asBizError(t, err)
	if bizErr.Code != service.ErrCodeValidationError {
		t.Errorf("错误码期望 %s，得到 %s", service.ErrCodeValidationError, bizErr.Code)
	}
}

func TestGetYearlyReport(t *testing.T) {
	statsSvc, cardSvc, windowID := setupStatisticsService(t)

	r, _ := cardSvc.IssueCard("ID001", 5000)
	cardNo := r.Card.CardNo

	if _, err := cardSvc.CreateTransaction(cardNo, windowID, 100); err != nil {
		t.Fatalf("消费失败: %v", err)
	}
	if _, err := cardSvc.CreateTransaction(cardNo, windowID, 200); err != nil {
		t.Fatalf("消费失败: %v", err)
	}

	year := time.Now().Year()
	result, err := statsSvc.GetYearlyReport(year)
	if err != nil {
		t.Fatalf("GetYearlyReport 失败: %v", err)
	}
	if result.TotalRevenue != 300 {
		t.Errorf("年报总收入期望 300，得到 %d", result.TotalRevenue)
	}
	if result.TransactionCount != 2 {
		t.Errorf("年报笔数期望 2，得到 %d", result.TransactionCount)
	}
	if result.Year != year {
		t.Errorf("年份期望 %d，得到 %d", year, result.Year)
	}
}

func TestGetWindowRevenue(t *testing.T) {
	statsSvc, cardSvc, windowID := setupStatisticsService(t)

	r, _ := cardSvc.IssueCard("ID001", 5000)
	cardNo := r.Card.CardNo

	if _, err := cardSvc.CreateTransaction(cardNo, windowID, 400); err != nil {
		t.Fatalf("消费失败: %v", err)
	}
	if _, err := cardSvc.CreateTransaction(cardNo, windowID, 100); err != nil {
		t.Fatalf("消费失败: %v", err)
	}

	start := time.Now().Add(-1 * time.Hour)
	end := time.Now().Add(1 * time.Hour)
	result, err := statsSvc.GetWindowRevenue(start, end)
	if err != nil {
		t.Fatalf("GetWindowRevenue 失败: %v", err)
	}
	if len(result.Windows) != 1 {
		t.Fatalf("窗口数量期望 1，得到 %d", len(result.Windows))
	}
	if result.Windows[0].Revenue != 500 {
		t.Errorf("窗口收入期望 500，得到 %d", result.Windows[0].Revenue)
	}
}

func TestGetDepositDetails_Pagination(t *testing.T) {
	statsSvc, cardSvc, _ := setupStatisticsService(t)

	// 发 3 张卡，每人各充一笔款（ID001/ID002/ID003）
	for _, id := range []string{"ID001", "ID002", "ID003"} {
		r, err := cardSvc.IssueCard(id, 0)
		if err != nil {
			t.Fatalf("发卡失败 %s: %v", id, err)
		}
		if _, err := cardSvc.Deposit(r.Card.CardNo, 100); err != nil {
			t.Fatalf("充值失败 %s: %v", id, err)
		}
	}

	// 第 1 页 pageSize=2：应返回 2 条，total=3
	result, err := statsSvc.GetDepositDetails(nil, nil, 1, 2)
	if err != nil {
		t.Fatalf("GetDepositDetails 失败: %v", err)
	}
	if result.Total != 3 {
		t.Errorf("total 期望 3，得到 %d", result.Total)
	}
	if len(result.Holders) != 2 {
		t.Errorf("第 1 页 holders 数量期望 2，得到 %d", len(result.Holders))
	}
	if result.Page != 1 {
		t.Errorf("page 期望 1，得到 %d", result.Page)
	}
	if result.PageSize != 2 {
		t.Errorf("pageSize 期望 2，得到 %d", result.PageSize)
	}

	// 第 2 页 pageSize=2：应返回 1 条
	result2, err := statsSvc.GetDepositDetails(nil, nil, 2, 2)
	if err != nil {
		t.Fatalf("GetDepositDetails 第 2 页失败: %v", err)
	}
	if result2.Total != 3 {
		t.Errorf("第 2 页 total 期望 3，得到 %d", result2.Total)
	}
	if len(result2.Holders) != 1 {
		t.Errorf("第 2 页 holders 数量期望 1，得到 %d", len(result2.Holders))
	}

	// 超出范围的页：应返回空列表，total 不变
	result3, err := statsSvc.GetDepositDetails(nil, nil, 99, 2)
	if err != nil {
		t.Fatalf("GetDepositDetails 超出范围失败: %v", err)
	}
	if result3.Total != 3 {
		t.Errorf("超出范围 total 期望 3，得到 %d", result3.Total)
	}
	if len(result3.Holders) != 0 {
		t.Errorf("超出范围 holders 应为空，得到 %d 条", len(result3.Holders))
	}

	// 每条持卡人的存款记录应包含正确金额
	for _, h := range result.Holders {
		if h.TotalAmount != 100 {
			t.Errorf("持卡人 %s 存款总额期望 100，得到 %d", h.IDNumber, h.TotalAmount)
		}
	}
}

func TestGetHolderDeposits(t *testing.T) {
	statsSvc, cardSvc, _ := setupStatisticsService(t)

	// 发卡并存款
	r, err := cardSvc.IssueCard("ID001", 0)
	if err != nil {
		t.Fatalf("发卡失败: %v", err)
	}
	cardNo := r.Card.CardNo
	holderID := uint(r.CardHolder.ID)

	if _, err := cardSvc.Deposit(cardNo, 300); err != nil {
		t.Fatalf("充值失败: %v", err)
	}
	if _, err := cardSvc.Deposit(cardNo, 200); err != nil {
		t.Fatalf("充值失败: %v", err)
	}
	if _, err := cardSvc.Deposit(cardNo, 100); err != nil {
		t.Fatalf("充值失败: %v", err)
	}

	// 第 1 页 pageSize=2：应返回 2 条，total=3
	result, err := statsSvc.GetHolderDeposits(holderID, nil, nil, 1, 2)
	if err != nil {
		t.Fatalf("GetHolderDeposits 失败: %v", err)
	}
	if result.Total != 3 {
		t.Errorf("total 期望 3，得到 %d", result.Total)
	}
	if len(result.Deposits) != 2 {
		t.Errorf("第 1 页 deposits 数量期望 2，得到 %d", len(result.Deposits))
	}
	if result.Page != 1 {
		t.Errorf("page 期望 1，得到 %d", result.Page)
	}
	if result.PageSize != 2 {
		t.Errorf("pageSize 期望 2，得到 %d", result.PageSize)
	}
	if result.HolderName == "" {
		t.Error("holderName 不应为空")
	}

	// 第 2 页 pageSize=2：应返回 1 条
	result2, err := statsSvc.GetHolderDeposits(holderID, nil, nil, 2, 2)
	if err != nil {
		t.Fatalf("GetHolderDeposits 第 2 页失败: %v", err)
	}
	if result2.Total != 3 {
		t.Errorf("第 2 页 total 期望 3，得到 %d", result2.Total)
	}
	if len(result2.Deposits) != 1 {
		t.Errorf("第 2 页 deposits 数量期望 1，得到 %d", len(result2.Deposits))
	}
}

func TestGetHolderDeposits_HolderNotFound(t *testing.T) {
	statsSvc, _, _ := setupStatisticsService(t)

	// holderID 9999 不存在，应返回错误
	_, err := statsSvc.GetHolderDeposits(9999, nil, nil, 1, 10)
	if err == nil {
		t.Fatal("期望持卡人不存在错误，但未返回错误")
	}
}
