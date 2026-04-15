package service_test

import (
	"backend/db"
	"backend/model"
	"backend/repository"
	"backend/service"
	"errors"
	"testing"
)

// setupCardService 初始化内存数据库并返回 CardService
func setupCardService(t *testing.T) *service.CardService {
	t.Helper()
	gormDB, err := db.Init(":memory:")
	if err != nil {
		t.Fatalf("初始化内存数据库失败: %v", err)
	}
	cardRepo := repository.NewCardRepository(gormDB)
	windowRepo := repository.NewWindowRepository(gormDB)
	return service.NewCardService(cardRepo, windowRepo)
}

// setupCardServiceWithWindow 初始化内存数据库、创建一个窗口并返回 CardService 和窗口 ID
func setupCardServiceWithWindow(t *testing.T) (*service.CardService, uint) {
	t.Helper()
	gormDB, err := db.Init(":memory:")
	if err != nil {
		t.Fatalf("初始化内存数据库失败: %v", err)
	}
	cardRepo := repository.NewCardRepository(gormDB)
	windowRepo := repository.NewWindowRepository(gormDB)

	w := &model.Window{Name: "一号窗口"}
	if err := windowRepo.CreateWindow(w); err != nil {
		t.Fatalf("创建窗口失败: %v", err)
	}
	return service.NewCardService(cardRepo, windowRepo), w.ID
}

// asBizError 断言 err 是 BizError 并返回，否则 t.Fatal
func asBizError(t *testing.T, err error) *service.BizError {
	t.Helper()
	var bizErr *service.BizError
	if !errors.As(err, &bizErr) {
		t.Fatalf("期望 BizError，得到: %v (%T)", err, err)
	}
	return bizErr
}

func TestIssueCard(t *testing.T) {
	tests := []struct {
		name      string
		req       service.IssueCardRequest
		wantErr   bool
		errCode   string
	}{
		{
			name: "正常发卡",
			req: service.IssueCardRequest{
				Name: "张三", IDNumber: "ID001", Deposit: 1000, PreDeposit: 500,
			},
		},
		{
			name: "押金为零应报错",
			req: service.IssueCardRequest{
				Name: "李四", IDNumber: "ID002", Deposit: 0, PreDeposit: 0,
			},
			wantErr: true,
			errCode: service.ErrCodeInvalidAmount,
		},
		{
			name: "押金为负应报错",
			req: service.IssueCardRequest{
				Name: "王五", IDNumber: "ID003", Deposit: -100, PreDeposit: 0,
			},
			wantErr: true,
			errCode: service.ErrCodeInvalidAmount,
		},
		{
			name: "预存款为负应报错",
			req: service.IssueCardRequest{
				Name: "赵六", IDNumber: "ID004", Deposit: 1000, PreDeposit: -1,
			},
			wantErr: true,
			errCode: service.ErrCodeInvalidAmount,
		},
		{
			name: "姓名为空应报错",
			req: service.IssueCardRequest{
				Name: "", IDNumber: "ID005", Deposit: 1000, PreDeposit: 0,
			},
			wantErr: true,
			errCode: service.ErrCodeValidationError,
		},
		{
			name: "证件号为空应报错",
			req: service.IssueCardRequest{
				Name: "张三", IDNumber: "", Deposit: 1000, PreDeposit: 0,
			},
			wantErr: true,
			errCode: service.ErrCodeValidationError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := setupCardService(t)
			result, err := svc.IssueCard(tt.req)
			if tt.wantErr {
				if err == nil {
					t.Fatal("期望错误，但未返回错误")
				}
				bizErr := asBizError(t, err)
				if bizErr.Code != tt.errCode {
					t.Errorf("错误码期望 %s，得到 %s", tt.errCode, bizErr.Code)
				}
				return
			}
			if err != nil {
				t.Fatalf("不期望错误，得到: %v", err)
			}
			if result.Card == nil {
				t.Fatal("Card 不应为 nil")
			}
			if result.CardHolder == nil {
				t.Fatal("CardHolder 不应为 nil")
			}
			if result.Card.Deposit != tt.req.Deposit {
				t.Errorf("押金期望 %d，得到 %d", tt.req.Deposit, result.Card.Deposit)
			}
			if result.Card.Balance != tt.req.PreDeposit {
				t.Errorf("余额期望 %d，得到 %d", tt.req.PreDeposit, result.Card.Balance)
			}
		})
	}
}

func TestIssueCard_OnePersonOneCard(t *testing.T) {
	svc := setupCardService(t)

	req := service.IssueCardRequest{
		Name: "张三", IDNumber: "ID001", Deposit: 1000, PreDeposit: 500,
	}

	// 第一次发卡成功
	_, err := svc.IssueCard(req)
	if err != nil {
		t.Fatalf("第一次发卡失败: %v", err)
	}

	// 同一证件号再次发卡应返回 CARD_ALREADY_ACTIVE
	_, err = svc.IssueCard(req)
	if err == nil {
		t.Fatal("期望 CARD_ALREADY_ACTIVE 错误，但未返回错误")
	}
	bizErr := asBizError(t, err)
	if bizErr.Code != service.ErrCodeCardAlreadyActive {
		t.Errorf("错误码期望 %s，得到 %s", service.ErrCodeCardAlreadyActive, bizErr.Code)
	}
}

func TestIssueCard_AutoCancelLostCard(t *testing.T) {
	svc := setupCardService(t)

	// 先发一张卡
	result1, err := svc.IssueCard(service.IssueCardRequest{
		Name: "张三", IDNumber: "ID001", Deposit: 1000, PreDeposit: 500,
	})
	if err != nil {
		t.Fatalf("发卡失败: %v", err)
	}

	// 挂失
	_, err = svc.ReportLoss(result1.Card.ID)
	if err != nil {
		t.Fatalf("挂失失败: %v", err)
	}

	// 同证件号再次发卡，旧卡应被自动注销
	result2, err := svc.IssueCard(service.IssueCardRequest{
		Name: "张三", IDNumber: "ID001", Deposit: 2000, PreDeposit: 0,
	})
	if err != nil {
		t.Fatalf("二次发卡失败: %v", err)
	}
	if result2.Refund == nil {
		t.Fatal("期望包含旧卡退款信息，但 Refund 为 nil")
	}
	if result2.Refund.OldCardID != result1.Card.ID {
		t.Errorf("旧卡号期望 %d，得到 %d", result1.Card.ID, result2.Refund.OldCardID)
	}
	if result2.Refund.Deposit != 1000 {
		t.Errorf("退还押金期望 1000，得到 %d", result2.Refund.Deposit)
	}
	if result2.Refund.Balance != 500 {
		t.Errorf("退还余额期望 500，得到 %d", result2.Refund.Balance)
	}
	if result2.Refund.Total != 1500 {
		t.Errorf("退款总额期望 1500，得到 %d", result2.Refund.Total)
	}
	// 新卡号与旧卡号不同
	if result2.Card.ID == result1.Card.ID {
		t.Error("新卡号不应与旧卡号相同")
	}
}

func TestDeposit(t *testing.T) {
	tests := []struct {
		name        string
		setupStatus model.CardStatus // 在充值前将卡状态设为此值
		amount      int64
		wantErr     bool
		errCode     string
	}{
		{
			name:        "正常充值",
			setupStatus: model.CardStatusActive,
			amount:      500,
		},
		{
			name:        "充值金额为零应报错",
			setupStatus: model.CardStatusActive,
			amount:      0,
			wantErr:     true,
			errCode:     service.ErrCodeInvalidAmount,
		},
		{
			name:        "充值金额为负应报错",
			setupStatus: model.CardStatusActive,
			amount:      -100,
			wantErr:     true,
			errCode:     service.ErrCodeInvalidAmount,
		},
		{
			name:        "已挂失卡不能充值",
			setupStatus: model.CardStatusLost,
			amount:      500,
			wantErr:     true,
			errCode:     service.ErrCodeCardNotActive,
		},
		{
			name:        "已注销卡不能充值",
			setupStatus: model.CardStatusCancelled,
			amount:      500,
			wantErr:     true,
			errCode:     service.ErrCodeCardNotActive,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := setupCardService(t)

			// 先发一张卡
			result, err := svc.IssueCard(service.IssueCardRequest{
				Name: "张三", IDNumber: "ID001", Deposit: 1000, PreDeposit: 200,
			})
			if err != nil {
				t.Fatalf("发卡失败: %v", err)
			}
			cardID := result.Card.ID

			// 按测试需要调整卡状态
			if tt.setupStatus == model.CardStatusLost {
				if _, err := svc.ReportLoss(cardID); err != nil {
					t.Fatalf("挂失失败: %v", err)
				}
			} else if tt.setupStatus == model.CardStatusCancelled {
				if _, err := svc.CancelCard(cardID); err != nil {
					t.Fatalf("注销失败: %v", err)
				}
			}

			depositResult, err := svc.Deposit(cardID, tt.amount)
			if tt.wantErr {
				if err == nil {
					t.Fatal("期望错误，但未返回错误")
				}
				bizErr := asBizError(t, err)
				if bizErr.Code != tt.errCode {
					t.Errorf("错误码期望 %s，得到 %s", tt.errCode, bizErr.Code)
				}
				return
			}
			if err != nil {
				t.Fatalf("不期望错误，得到: %v", err)
			}
			if depositResult.NewBalance != 200+tt.amount {
				t.Errorf("充值后余额期望 %d，得到 %d", 200+tt.amount, depositResult.NewBalance)
			}
			if depositResult.Amount != tt.amount {
				t.Errorf("充值金额期望 %d，得到 %d", tt.amount, depositResult.Amount)
			}
		})
	}
}

func TestCreateTransaction_TripleValidation(t *testing.T) {
	tests := []struct {
		name        string
		setupStatus model.CardStatus
		amount      int64
		wantErr     bool
		errCode     string
	}{
		{
			name:        "正常消费",
			setupStatus: model.CardStatusActive,
			amount:      100,
		},
		{
			name:        "卡已注销不能消费",
			setupStatus: model.CardStatusCancelled,
			amount:      100,
			wantErr:     true,
			errCode:     service.ErrCodeCardCancelled,
		},
		{
			name:        "卡已挂失不能消费",
			setupStatus: model.CardStatusLost,
			amount:      100,
			wantErr:     true,
			errCode:     service.ErrCodeCardLost,
		},
		{
			name:        "余额不足",
			setupStatus: model.CardStatusActive,
			amount:      100000, // 远超余额
			wantErr:     true,
			errCode:     service.ErrCodeInsufficientBalance,
		},
		{
			name:        "消费金额为零应报错",
			setupStatus: model.CardStatusActive,
			amount:      0,
			wantErr:     true,
			errCode:     service.ErrCodeInvalidAmount,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, windowID := setupCardServiceWithWindow(t)

			// 发卡，余额 500
			result, err := svc.IssueCard(service.IssueCardRequest{
				Name: "张三", IDNumber: "ID001", Deposit: 1000, PreDeposit: 500,
			})
			if err != nil {
				t.Fatalf("发卡失败: %v", err)
			}
			cardID := result.Card.ID

			// 按需调整卡状态
			if tt.setupStatus == model.CardStatusLost {
				if _, err := svc.ReportLoss(cardID); err != nil {
					t.Fatalf("挂失失败: %v", err)
				}
			} else if tt.setupStatus == model.CardStatusCancelled {
				if _, err := svc.CancelCard(cardID); err != nil {
					t.Fatalf("注销失败: %v", err)
				}
			}

			txResult, err := svc.CreateTransaction(cardID, windowID, tt.amount)
			if tt.wantErr {
				if err == nil {
					t.Fatal("期望错误，但未返回错误")
				}
				bizErr := asBizError(t, err)
				if bizErr.Code != tt.errCode {
					t.Errorf("错误码期望 %s，得到 %s", tt.errCode, bizErr.Code)
				}
				return
			}
			if err != nil {
				t.Fatalf("不期望错误，得到: %v", err)
			}
			expectedBalance := int64(500) - tt.amount
			if txResult.NewBalance != expectedBalance {
				t.Errorf("消费后余额期望 %d，得到 %d", expectedBalance, txResult.NewBalance)
			}
		})
	}
}

func TestCreateTransaction_CardNotFound(t *testing.T) {
	svc, windowID := setupCardServiceWithWindow(t)

	_, err := svc.CreateTransaction(99999, windowID, 100)
	if err == nil {
		t.Fatal("期望卡号不存在错误")
	}
	bizErr := asBizError(t, err)
	if bizErr.Code != service.ErrCodeCardNotFound {
		t.Errorf("错误码期望 %s，得到 %s", service.ErrCodeCardNotFound, bizErr.Code)
	}
}

func TestReportLoss(t *testing.T) {
	tests := []struct {
		name        string
		setupStatus model.CardStatus
		wantErr     bool
		errCode     string
	}{
		{
			name:        "active 卡可以挂失",
			setupStatus: model.CardStatusActive,
		},
		{
			name:        "lost 卡不能再次挂失",
			setupStatus: model.CardStatusLost,
			wantErr:     true,
			errCode:     service.ErrCodeCardNotActive,
		},
		{
			name:        "cancelled 卡不能挂失",
			setupStatus: model.CardStatusCancelled,
			wantErr:     true,
			errCode:     service.ErrCodeCardNotActive,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := setupCardService(t)

			result, err := svc.IssueCard(service.IssueCardRequest{
				Name: "张三", IDNumber: "ID001", Deposit: 1000, PreDeposit: 0,
			})
			if err != nil {
				t.Fatalf("发卡失败: %v", err)
			}
			cardID := result.Card.ID

			// 预先调整状态
			if tt.setupStatus == model.CardStatusLost {
				if _, err := svc.ReportLoss(cardID); err != nil {
					t.Fatalf("预置挂失失败: %v", err)
				}
			} else if tt.setupStatus == model.CardStatusCancelled {
				if _, err := svc.CancelCard(cardID); err != nil {
					t.Fatalf("预置注销失败: %v", err)
				}
			}

			_, err = svc.ReportLoss(cardID)
			if tt.wantErr {
				if err == nil {
					t.Fatal("期望错误，但未返回错误")
				}
				bizErr := asBizError(t, err)
				if bizErr.Code != tt.errCode {
					t.Errorf("错误码期望 %s，得到 %s", tt.errCode, bizErr.Code)
				}
				return
			}
			if err != nil {
				t.Fatalf("不期望错误，得到: %v", err)
			}
		})
	}
}

func TestCancelLossReport(t *testing.T) {
	tests := []struct {
		name        string
		setupStatus model.CardStatus
		wantErr     bool
		errCode     string
	}{
		{
			name:        "lost 卡可以取消挂失",
			setupStatus: model.CardStatusLost,
		},
		{
			name:        "active 卡不能取消挂失",
			setupStatus: model.CardStatusActive,
			wantErr:     true,
			errCode:     service.ErrCodeCardNotLost,
		},
		{
			name:        "cancelled 卡不能取消挂失",
			setupStatus: model.CardStatusCancelled,
			wantErr:     true,
			errCode:     service.ErrCodeCardNotLost,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := setupCardService(t)

			result, err := svc.IssueCard(service.IssueCardRequest{
				Name: "张三", IDNumber: "ID001", Deposit: 1000, PreDeposit: 0,
			})
			if err != nil {
				t.Fatalf("发卡失败: %v", err)
			}
			cardID := result.Card.ID

			// 预先调整状态
			if tt.setupStatus == model.CardStatusLost {
				if _, err := svc.ReportLoss(cardID); err != nil {
					t.Fatalf("预置挂失失败: %v", err)
				}
			} else if tt.setupStatus == model.CardStatusCancelled {
				if _, err := svc.CancelCard(cardID); err != nil {
					t.Fatalf("预置注销失败: %v", err)
				}
			}

			_, err = svc.CancelLossReport(cardID)
			if tt.wantErr {
				if err == nil {
					t.Fatal("期望错误，但未返回错误")
				}
				bizErr := asBizError(t, err)
				if bizErr.Code != tt.errCode {
					t.Errorf("错误码期望 %s，得到 %s", tt.errCode, bizErr.Code)
				}
				return
			}
			if err != nil {
				t.Fatalf("不期望错误，得到: %v", err)
			}
		})
	}
}

func TestCancelCard(t *testing.T) {
	tests := []struct {
		name           string
		deposit        int64
		preDeposit     int64
		extraDeposit   int64 // 额外充值金额
		doLoss         bool  // 发卡后先挂失再注销
		wantErr        bool
		errCode        string
		wantDeposit    int64
		wantBalance    int64
	}{
		{
			name:        "注销 active 卡，退还押金和余额",
			deposit:     1000,
			preDeposit:  500,
			wantDeposit: 1000,
			wantBalance: 500,
		},
		{
			name:        "注销 active 卡，有额外充值",
			deposit:     1000,
			preDeposit:  500,
			extraDeposit: 300,
			wantDeposit: 1000,
			wantBalance: 800,
		},
		{
			name:        "注销 lost 卡",
			deposit:     1000,
			preDeposit:  200,
			doLoss:      true,
			wantDeposit: 1000,
			wantBalance: 200,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := setupCardService(t)

			result, err := svc.IssueCard(service.IssueCardRequest{
				Name: "张三", IDNumber: "ID001", Deposit: tt.deposit, PreDeposit: tt.preDeposit,
			})
			if err != nil {
				t.Fatalf("发卡失败: %v", err)
			}
			cardID := result.Card.ID

			if tt.extraDeposit > 0 {
				if _, err := svc.Deposit(cardID, tt.extraDeposit); err != nil {
					t.Fatalf("充值失败: %v", err)
				}
			}
			if tt.doLoss {
				if _, err := svc.ReportLoss(cardID); err != nil {
					t.Fatalf("挂失失败: %v", err)
				}
			}

			cancelResult, err := svc.CancelCard(cardID)
			if tt.wantErr {
				if err == nil {
					t.Fatal("期望错误，但未返回错误")
				}
				bizErr := asBizError(t, err)
				if bizErr.Code != tt.errCode {
					t.Errorf("错误码期望 %s，得到 %s", tt.errCode, bizErr.Code)
				}
				return
			}
			if err != nil {
				t.Fatalf("不期望错误，得到: %v", err)
			}
			if cancelResult.Deposit != tt.wantDeposit {
				t.Errorf("退还押金期望 %d，得到 %d", tt.wantDeposit, cancelResult.Deposit)
			}
			if cancelResult.Balance != tt.wantBalance {
				t.Errorf("退还余额期望 %d，得到 %d", tt.wantBalance, cancelResult.Balance)
			}
			if cancelResult.Total != tt.wantDeposit+tt.wantBalance {
				t.Errorf("退款总额期望 %d，得到 %d", tt.wantDeposit+tt.wantBalance, cancelResult.Total)
			}
			if cancelResult.Card.Status != model.CardStatusCancelled {
				t.Errorf("卡状态期望 cancelled，得到 %s", cancelResult.Card.Status)
			}
			if cancelResult.Card.Balance != 0 {
				t.Errorf("注销后余额应为 0，得到 %d", cancelResult.Card.Balance)
			}
		})
	}
}

func TestCancelCard_AlreadyCancelled(t *testing.T) {
	svc := setupCardService(t)

	result, err := svc.IssueCard(service.IssueCardRequest{
		Name: "张三", IDNumber: "ID001", Deposit: 1000, PreDeposit: 0,
	})
	if err != nil {
		t.Fatalf("发卡失败: %v", err)
	}
	cardID := result.Card.ID

	// 第一次注销
	if _, err := svc.CancelCard(cardID); err != nil {
		t.Fatalf("第一次注销失败: %v", err)
	}

	// 第二次注销应报错
	_, err = svc.CancelCard(cardID)
	if err == nil {
		t.Fatal("期望 CARD_ALREADY_CANCELLED 错误")
	}
	bizErr := asBizError(t, err)
	if bizErr.Code != service.ErrCodeCardAlreadyCancelled {
		t.Errorf("错误码期望 %s，得到 %s", service.ErrCodeCardAlreadyCancelled, bizErr.Code)
	}
}
