package service_test

import (
	"backend/db"
	"backend/model"
	"backend/repository"
	"backend/service"
	"errors"
	"testing"
)

// fakeStudentValidator 测试用的假学籍验证器
type fakeStudentValidator struct {
	// validIDs 中的证件号验证通过，其余返回 STUDENT_NOT_FOUND
	validIDs map[string]*service.StudentInfo
	// forceError 为 true 时模拟服务不可用
	forceError bool
}

func newFakeValidator(ids ...string) *fakeStudentValidator {
	v := &fakeStudentValidator{validIDs: make(map[string]*service.StudentInfo)}
	for _, id := range ids {
		v.validIDs[id] = &service.StudentInfo{IDNumber: id, Name: "测试用户", Type: "student"}
	}
	return v
}

func (f *fakeStudentValidator) Validate(idNumber string) (*service.StudentInfo, error) {
	if f.forceError {
		return nil, &service.BizError{Code: service.ErrCodeStudentServiceError, Message: "学籍验证服务暂时不可用"}
	}
	if info, ok := f.validIDs[idNumber]; ok {
		return info, nil
	}
	return nil, &service.BizError{Code: service.ErrCodeStudentNotFound, Message: "证件号不存在，非本校学生/教职工"}
}

// setupCardService 初始化内存数据库并返回 CardService（传入验证器）
func setupCardService(t *testing.T, validator service.StudentValidator) *service.CardService {
	t.Helper()
	gormDB, err := db.Init(":memory:")
	if err != nil {
		t.Fatalf("初始化内存数据库失败: %v", err)
	}
	cardRepo := repository.NewCardRepository(gormDB)
	windowRepo := repository.NewWindowRepository(gormDB)
	return service.NewCardService(cardRepo, windowRepo, validator)
}

// setupCardServiceWithWindow 初始化内存数据库、创建一个窗口并返回 CardService 和窗口 ID
func setupCardServiceWithWindow(t *testing.T, validator service.StudentValidator) (*service.CardService, uint) {
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
	return service.NewCardService(cardRepo, windowRepo, validator), w.ID
}

// setupWithRepo 初始化内存数据库，同时返回 CardService 和 CardRepository（用于 FindCardByCardNo 辅助查询）
func setupWithRepo(t *testing.T, validator service.StudentValidator) (*service.CardService, *repository.CardRepository) {
	t.Helper()
	gormDB, err := db.Init(":memory:")
	if err != nil {
		t.Fatalf("初始化内存数据库失败: %v", err)
	}
	cardRepo := repository.NewCardRepository(gormDB)
	windowRepo := repository.NewWindowRepository(gormDB)
	svc := service.NewCardService(cardRepo, windowRepo, validator)
	return svc, cardRepo
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
		name       string
		idNumber   string
		preDeposit int64
		wantErr    bool
		errCode    string
	}{
		{
			name:       "正常发卡",
			idNumber:   "ID001",
			preDeposit: 500,
		},
		{
			name:       "预存款为零也可以发卡",
			idNumber:   "ID002",
			preDeposit: 0,
		},
		{
			name:       "预存款为负应报错",
			idNumber:   "ID003",
			preDeposit: -1,
			wantErr:    true,
			errCode:    service.ErrCodeInvalidAmount,
		},
		{
			name:       "证件号不在学籍库返回 STUDENT_NOT_FOUND",
			idNumber:   "UNKNOWN",
			preDeposit: 0,
			wantErr:    true,
			errCode:    service.ErrCodeStudentNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 只有 ID001/ID002/ID003 在学籍库中
			v := newFakeValidator("ID001", "ID002", "ID003")
			svc := setupCardService(t, v)
			result, err := svc.IssueCard(tt.idNumber, tt.preDeposit)
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
			// 押金固定 2000 分
			if result.Card.Deposit != 2000 {
				t.Errorf("押金期望 2000，得到 %d", result.Card.Deposit)
			}
			if result.Card.Balance != tt.preDeposit {
				t.Errorf("余额期望 %d，得到 %d", tt.preDeposit, result.Card.Balance)
			}
			if result.Card.CardNo == "" {
				t.Error("CardNo 不应为空")
			}
		})
	}
}

func TestIssueCard_StudentServiceError(t *testing.T) {
	v := &fakeStudentValidator{validIDs: make(map[string]*service.StudentInfo), forceError: true}
	svc := setupCardService(t, v)

	_, err := svc.IssueCard("ID001", 0)
	if err == nil {
		t.Fatal("期望 STUDENT_SERVICE_ERROR 错误")
	}
	bizErr := asBizError(t, err)
	if bizErr.Code != service.ErrCodeStudentServiceError {
		t.Errorf("错误码期望 %s，得到 %s", service.ErrCodeStudentServiceError, bizErr.Code)
	}
}

func TestIssueCard_OnePersonOneCard(t *testing.T) {
	v := newFakeValidator("ID001")
	svc := setupCardService(t, v)

	// 第一次发卡成功
	_, err := svc.IssueCard("ID001", 500)
	if err != nil {
		t.Fatalf("第一次发卡失败: %v", err)
	}

	// 同一证件号再次发卡应返回 CARD_ALREADY_ACTIVE
	_, err = svc.IssueCard("ID001", 0)
	if err == nil {
		t.Fatal("期望 CARD_ALREADY_ACTIVE 错误，但未返回错误")
	}
	bizErr := asBizError(t, err)
	if bizErr.Code != service.ErrCodeCardAlreadyActive {
		t.Errorf("错误码期望 %s，得到 %s", service.ErrCodeCardAlreadyActive, bizErr.Code)
	}
}

func TestIssueCard_AutoCancelLostCard(t *testing.T) {
	v := newFakeValidator("ID001")
	svc, cardRepo := setupWithRepo(t, v)

	// 先发一张卡
	result1, err := svc.IssueCard("ID001", 500)
	if err != nil {
		t.Fatalf("发卡失败: %v", err)
	}
	cardNo1 := result1.Card.CardNo

	// 挂失
	_, err = svc.ReportLoss(cardNo1)
	if err != nil {
		t.Fatalf("挂失失败: %v", err)
	}

	// 同证件号再次发卡，旧卡应被自动注销
	result2, err := svc.IssueCard("ID001", 0)
	if err != nil {
		t.Fatalf("二次发卡失败: %v", err)
	}
	if result2.Refund == nil {
		t.Fatal("期望包含旧卡退款信息，但 Refund 为 nil")
	}
	if result2.Refund.OldCardNo != cardNo1 {
		t.Errorf("旧卡号期望 %s，得到 %s", cardNo1, result2.Refund.OldCardNo)
	}
	if result2.Refund.Deposit != 2000 {
		t.Errorf("退还押金期望 2000，得到 %d", result2.Refund.Deposit)
	}
	if result2.Refund.Balance != 500 {
		t.Errorf("退还余额期望 500，得到 %d", result2.Refund.Balance)
	}
	if result2.Refund.Total != 2500 {
		t.Errorf("退款总额期望 2500，得到 %d", result2.Refund.Total)
	}
	// 新卡号与旧卡号不同
	if result2.Card.CardNo == cardNo1 {
		t.Error("新卡号不应与旧卡号相同")
	}

	// 通过 cardRepo 验证旧卡状态为 cancelled
	oldCard, err := cardRepo.FindCardByCardNo(cardNo1)
	if err != nil {
		t.Fatalf("查询旧卡失败: %v", err)
	}
	if oldCard.Status != model.CardStatusCancelled {
		t.Errorf("旧卡状态期望 cancelled，得到 %s", oldCard.Status)
	}
}

func TestDeposit(t *testing.T) {
	tests := []struct {
		name        string
		setupStatus model.CardStatus
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
			v := newFakeValidator("ID001")
			svc := setupCardService(t, v)

			// 先发一张卡
			result, err := svc.IssueCard("ID001", 200)
			if err != nil {
				t.Fatalf("发卡失败: %v", err)
			}
			cardNo := result.Card.CardNo

			// 按测试需要调整卡状态
			if tt.setupStatus == model.CardStatusLost {
				if _, err := svc.ReportLoss(cardNo); err != nil {
					t.Fatalf("挂失失败: %v", err)
				}
			} else if tt.setupStatus == model.CardStatusCancelled {
				if _, err := svc.CancelCard(cardNo); err != nil {
					t.Fatalf("注销失败: %v", err)
				}
			}

			depositResult, err := svc.Deposit(cardNo, tt.amount)
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
			if depositResult.CardNo != cardNo {
				t.Errorf("收据 CardNo 期望 %s，得到 %s", cardNo, depositResult.CardNo)
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
			v := newFakeValidator("ID001")
			svc, windowID := setupCardServiceWithWindow(t, v)

			// 发卡，余额 500
			result, err := svc.IssueCard("ID001", 500)
			if err != nil {
				t.Fatalf("发卡失败: %v", err)
			}
			cardNo := result.Card.CardNo

			// 按需调整卡状态
			if tt.setupStatus == model.CardStatusLost {
				if _, err := svc.ReportLoss(cardNo); err != nil {
					t.Fatalf("挂失失败: %v", err)
				}
			} else if tt.setupStatus == model.CardStatusCancelled {
				if _, err := svc.CancelCard(cardNo); err != nil {
					t.Fatalf("注销失败: %v", err)
				}
			}

			txResult, err := svc.CreateTransaction(cardNo, windowID, tt.amount)
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
			if txResult.CardNo != cardNo {
				t.Errorf("结果 CardNo 期望 %s，得到 %s", cardNo, txResult.CardNo)
			}
		})
	}
}

func TestCreateTransaction_CardNotFound(t *testing.T) {
	v := newFakeValidator("ID001")
	svc, windowID := setupCardServiceWithWindow(t, v)

	_, err := svc.CreateTransaction("9999999999999999", windowID, 100)
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
			v := newFakeValidator("ID001")
			svc := setupCardService(t, v)

			result, err := svc.IssueCard("ID001", 0)
			if err != nil {
				t.Fatalf("发卡失败: %v", err)
			}
			cardNo := result.Card.CardNo

			// 预先调整状态
			if tt.setupStatus == model.CardStatusLost {
				if _, err := svc.ReportLoss(cardNo); err != nil {
					t.Fatalf("预置挂失失败: %v", err)
				}
			} else if tt.setupStatus == model.CardStatusCancelled {
				if _, err := svc.CancelCard(cardNo); err != nil {
					t.Fatalf("预置注销失败: %v", err)
				}
			}

			_, err = svc.ReportLoss(cardNo)
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
			v := newFakeValidator("ID001")
			svc := setupCardService(t, v)

			result, err := svc.IssueCard("ID001", 0)
			if err != nil {
				t.Fatalf("发卡失败: %v", err)
			}
			cardNo := result.Card.CardNo

			// 预先调整状态
			if tt.setupStatus == model.CardStatusLost {
				if _, err := svc.ReportLoss(cardNo); err != nil {
					t.Fatalf("预置挂失失败: %v", err)
				}
			} else if tt.setupStatus == model.CardStatusCancelled {
				if _, err := svc.CancelCard(cardNo); err != nil {
					t.Fatalf("预置注销失败: %v", err)
				}
			}

			_, err = svc.CancelLossReport(cardNo)
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
		name         string
		preDeposit   int64
		extraDeposit int64 // 额外充值金额
		doLoss       bool  // 发卡后先挂失再注销
		wantErr      bool
		errCode      string
		wantDeposit  int64
		wantBalance  int64
	}{
		{
			name:        "注销 active 卡，退还押金和余额",
			preDeposit:  500,
			wantDeposit: 2000,
			wantBalance: 500,
		},
		{
			name:         "注销 active 卡，有额外充值",
			preDeposit:   500,
			extraDeposit: 300,
			wantDeposit:  2000,
			wantBalance:  800,
		},
		{
			name:        "注销 lost 卡",
			preDeposit:  200,
			doLoss:      true,
			wantDeposit: 2000,
			wantBalance: 200,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := newFakeValidator("ID001")
			svc := setupCardService(t, v)

			result, err := svc.IssueCard("ID001", tt.preDeposit)
			if err != nil {
				t.Fatalf("发卡失败: %v", err)
			}
			cardNo := result.Card.CardNo

			if tt.extraDeposit > 0 {
				if _, err := svc.Deposit(cardNo, tt.extraDeposit); err != nil {
					t.Fatalf("充值失败: %v", err)
				}
			}
			if tt.doLoss {
				if _, err := svc.ReportLoss(cardNo); err != nil {
					t.Fatalf("挂失失败: %v", err)
				}
			}

			cancelResult, err := svc.CancelCard(cardNo)
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
	v := newFakeValidator("ID001")
	svc := setupCardService(t, v)

	result, err := svc.IssueCard("ID001", 0)
	if err != nil {
		t.Fatalf("发卡失败: %v", err)
	}
	cardNo := result.Card.CardNo

	// 第一次注销
	if _, err := svc.CancelCard(cardNo); err != nil {
		t.Fatalf("第一次注销失败: %v", err)
	}

	// 第二次注销应报错
	_, err = svc.CancelCard(cardNo)
	if err == nil {
		t.Fatal("期望 CARD_ALREADY_CANCELLED 错误")
	}
	bizErr := asBizError(t, err)
	if bizErr.Code != service.ErrCodeCardAlreadyCancelled {
		t.Errorf("错误码期望 %s，得到 %s", service.ErrCodeCardAlreadyCancelled, bizErr.Code)
	}
}

func TestIssueCard_PreDepositCreatesDepositRecord(t *testing.T) {
	v := newFakeValidator("ID001", "ID002")
	svc, cardRepo := setupWithRepo(t, v)

	// preDeposit > 0 时，应生成 DepositRecord
	_, err := svc.IssueCard("ID001", 500)
	if err != nil {
		t.Fatalf("发卡失败: %v", err)
	}

	details, err := cardRepo.GetDepositDetails(nil, nil)
	if err != nil {
		t.Fatalf("查询存款明细失败: %v", err)
	}
	if len(details) == 0 {
		t.Fatal("期望有 DepositRecord，但结果为空")
	}
	total := int64(0)
	for _, d := range details {
		total += d.TotalAmount
	}
	if total != 500 {
		t.Errorf("存款总额期望 500，得到 %d", total)
	}

	// preDeposit == 0 时，不应生成 DepositRecord
	_, err = svc.IssueCard("ID002", 0)
	if err != nil {
		t.Fatalf("发卡失败: %v", err)
	}

	details2, err := cardRepo.GetDepositDetails(nil, nil)
	if err != nil {
		t.Fatalf("查询存款明细失败: %v", err)
	}
	// 仍然只有 ID001 那条记录，总金额不变
	total2 := int64(0)
	for _, d := range details2 {
		total2 += d.TotalAmount
	}
	if total2 != 500 {
		t.Errorf("preDeposit=0 不应新增记录，总金额期望 500，得到 %d", total2)
	}
}
