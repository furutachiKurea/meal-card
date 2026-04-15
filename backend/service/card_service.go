// Package service 提供业务逻辑层
package service

import (
	"backend/model"
	"backend/repository"
	"errors"
	"time"

	"gorm.io/gorm"
)

// 业务错误码常量
const (
	ErrCodeCardNotFound        = "CARD_NOT_FOUND"
	ErrCodeCardAlreadyActive   = "CARD_ALREADY_ACTIVE"
	ErrCodeCardNotActive       = "CARD_NOT_ACTIVE"
	ErrCodeCardNotLost         = "CARD_NOT_LOST"
	ErrCodeCardCancelled       = "CARD_CANCELLED"
	ErrCodeCardLost            = "CARD_LOST"
	ErrCodeCardAlreadyCancelled = "CARD_ALREADY_CANCELLED"
	ErrCodeInsufficientBalance = "INSUFFICIENT_BALANCE"
	ErrCodeInvalidAmount       = "INVALID_AMOUNT"
	ErrCodeWindowNotFound      = "WINDOW_NOT_FOUND"
	ErrCodeValidationError     = "VALIDATION_ERROR"
)

// BizError 业务错误，携带错误码和人类可读消息
type BizError struct {
	Code    string
	Message string
}

func (e *BizError) Error() string {
	return e.Message
}

// newBizError 创建业务错误
func newBizError(code, message string) *BizError {
	return &BizError{Code: code, Message: message}
}

// IssueCardRequest 发卡请求参数
type IssueCardRequest struct {
	Name       string
	IDNumber   string
	Deposit    int64
	PreDeposit int64
}

// IssueCardResult 发卡结果
type IssueCardResult struct {
	Card       *model.Card
	CardHolder *model.CardHolder
	// Refund 旧卡退款信息，仅在有 lost 卡被自动注销时不为 nil
	Refund *OldCardRefund
}

// OldCardRefund 旧卡退款详情
type OldCardRefund struct {
	OldCardID uint
	Deposit   int64
	Balance   int64
	Total     int64
}

// DepositResult 存款结果（收据信息）
type DepositResult struct {
	ID         uint
	CardID     uint
	HolderName string
	Amount     int64
	NewBalance int64
	CreatedAt  time.Time
}

// TransactionResult 消费结果
type TransactionResult struct {
	ID         uint
	CardID     uint
	WindowID   uint
	Amount     int64
	NewBalance int64
	CreatedAt  time.Time
}

// CancellationResult 注销结果
type CancellationResult struct {
	Card    *model.Card
	Deposit int64
	Balance int64
	Total   int64
}

// CardService 饭卡业务逻辑
type CardService struct {
	cardRepo   *repository.CardRepository
	windowRepo *repository.WindowRepository
}

// NewCardService 创建 CardService 实例
func NewCardService(cardRepo *repository.CardRepository, windowRepo *repository.WindowRepository) *CardService {
	return &CardService{
		cardRepo:   cardRepo,
		windowRepo: windowRepo,
	}
}

// IssueCard 发卡业务
// 规则：
//  1. 同证件号已有 active 卡 → 返回 409 CARD_ALREADY_ACTIVE
//  2. 同证件号有 lost 卡 → 自动注销旧卡，创建新卡，返回退款信息
//  3. 无任何卡 → 直接创建
func (s *CardService) IssueCard(req IssueCardRequest) (*IssueCardResult, error) {
	if req.Deposit <= 0 {
		return nil, newBizError(ErrCodeInvalidAmount, "押金必须大于 0")
	}
	if req.PreDeposit < 0 {
		return nil, newBizError(ErrCodeInvalidAmount, "预存款不能为负数")
	}
	if req.Name == "" {
		return nil, newBizError(ErrCodeValidationError, "姓名不能为空")
	}
	if req.IDNumber == "" {
		return nil, newBizError(ErrCodeValidationError, "证件号不能为空")
	}

	// 查找或创建持卡人
	holder, err := s.cardRepo.FindCardHolderByIDNumber(req.IDNumber)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if holder == nil {
		holder = &model.CardHolder{
			Name:     req.Name,
			IDNumber: req.IDNumber,
		}
		if err := s.cardRepo.CreateCardHolder(holder); err != nil {
			return nil, err
		}
	}

	// 检查是否已有 active 卡
	activeCard, err := s.cardRepo.FindActiveCardByHolderID(holder.ID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if activeCard != nil {
		return nil, newBizError(ErrCodeCardAlreadyActive, "该证件号已有正在使用的卡")
	}

	// 检查是否有 lost 卡，若有则自动注销
	var refund *OldCardRefund
	lostCard, err := s.cardRepo.FindLostCardByHolderID(holder.ID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if lostCard != nil {
		// 自动注销旧卡
		oldDeposit := lostCard.Deposit
		oldBalance := lostCard.Balance
		lostCard.Status = model.CardStatusCancelled
		lostCard.Balance = 0
		if err := s.cardRepo.UpdateCard(lostCard); err != nil {
			return nil, err
		}
		refund = &OldCardRefund{
			OldCardID: lostCard.ID,
			Deposit:   oldDeposit,
			Balance:   oldBalance,
			Total:     oldDeposit + oldBalance,
		}
	}

	// 创建新卡
	newCard := &model.Card{
		CardHolderID: holder.ID,
		Deposit:      req.Deposit,
		Balance:      req.PreDeposit,
		Status:       model.CardStatusActive,
	}
	if err := s.cardRepo.CreateCard(newCard); err != nil {
		return nil, err
	}

	return &IssueCardResult{
		Card:       newCard,
		CardHolder: holder,
		Refund:     refund,
	}, nil
}

// GetCard 根据卡号查询卡片详情
func (s *CardService) GetCard(cardID uint) (*model.Card, error) {
	card, err := s.cardRepo.FindCardByID(cardID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, newBizError(ErrCodeCardNotFound, "卡号不存在")
		}
		return nil, err
	}
	return card, nil
}

// Deposit 存款（充值）
func (s *CardService) Deposit(cardID uint, amount int64) (*DepositResult, error) {
	if amount <= 0 {
		return nil, newBizError(ErrCodeInvalidAmount, "充值金额必须大于 0")
	}

	card, err := s.cardRepo.FindCardByID(cardID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, newBizError(ErrCodeCardNotFound, "卡号不存在")
		}
		return nil, err
	}

	if card.Status != model.CardStatusActive {
		if card.Status == model.CardStatusLost {
			return nil, newBizError(ErrCodeCardNotActive, "该卡已挂失，无法充值")
		}
		return nil, newBizError(ErrCodeCardNotActive, "该卡已注销，无法充值")
	}

	card.Balance += amount
	if err := s.cardRepo.UpdateCard(card); err != nil {
		return nil, err
	}

	record := &model.DepositRecord{
		CardID: card.ID,
		Amount: amount,
	}
	if err := s.cardRepo.CreateDepositRecord(record); err != nil {
		return nil, err
	}

	return &DepositResult{
		ID:         record.ID,
		CardID:     card.ID,
		HolderName: card.CardHolder.Name,
		Amount:     amount,
		NewBalance: card.Balance,
		CreatedAt:  record.CreatedAt,
	}, nil
}

// CreateTransaction 就餐消费结算
// 三重校验：卡存在 → 非 cancelled → 非 lost → 余额充足
func (s *CardService) CreateTransaction(cardID uint, windowID uint, amount int64) (*TransactionResult, error) {
	if amount <= 0 {
		return nil, newBizError(ErrCodeInvalidAmount, "消费金额必须大于 0")
	}

	card, err := s.cardRepo.FindCardByID(cardID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, newBizError(ErrCodeCardNotFound, "此卡非本单位所发")
		}
		return nil, err
	}

	if card.Status == model.CardStatusCancelled {
		return nil, newBizError(ErrCodeCardCancelled, "此卡已注销")
	}
	if card.Status == model.CardStatusLost {
		return nil, newBizError(ErrCodeCardLost, "此卡已挂失")
	}

	if card.Balance < amount {
		return nil, newBizError(ErrCodeInsufficientBalance, "余额不足")
	}

	// 验证窗口存在
	_, err = s.windowRepo.FindWindowByID(windowID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, newBizError(ErrCodeWindowNotFound, "窗口不存在")
		}
		return nil, err
	}

	card.Balance -= amount
	if err := s.cardRepo.UpdateCard(card); err != nil {
		return nil, err
	}

	tx := &model.Transaction{
		CardID:   card.ID,
		WindowID: windowID,
		Amount:   amount,
	}
	if err := s.cardRepo.CreateTransaction(tx); err != nil {
		return nil, err
	}

	return &TransactionResult{
		ID:         tx.ID,
		CardID:     card.ID,
		WindowID:   windowID,
		Amount:     amount,
		NewBalance: card.Balance,
		CreatedAt:  tx.CreatedAt,
	}, nil
}

// ReportLoss 挂失：active → lost
func (s *CardService) ReportLoss(cardID uint) (*model.Card, error) {
	card, err := s.cardRepo.FindCardByID(cardID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, newBizError(ErrCodeCardNotFound, "卡号不存在")
		}
		return nil, err
	}

	if card.Status != model.CardStatusActive {
		return nil, newBizError(ErrCodeCardNotActive, "只有正常使用中的卡可以挂失")
	}

	card.Status = model.CardStatusLost
	if err := s.cardRepo.UpdateCard(card); err != nil {
		return nil, err
	}
	return card, nil
}

// CancelLossReport 取消挂失：lost → active
func (s *CardService) CancelLossReport(cardID uint) (*model.Card, error) {
	card, err := s.cardRepo.FindCardByID(cardID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, newBizError(ErrCodeCardNotFound, "卡号不存在")
		}
		return nil, err
	}

	if card.Status != model.CardStatusLost {
		return nil, newBizError(ErrCodeCardNotLost, "只有已挂失的卡可以取消挂失")
	}

	card.Status = model.CardStatusActive
	if err := s.cardRepo.UpdateCard(card); err != nil {
		return nil, err
	}
	return card, nil
}

// CancelCard 注销卡片，退还押金和余额
func (s *CardService) CancelCard(cardID uint) (*CancellationResult, error) {
	card, err := s.cardRepo.FindCardByID(cardID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, newBizError(ErrCodeCardNotFound, "卡号不存在")
		}
		return nil, err
	}

	if card.Status == model.CardStatusCancelled {
		return nil, newBizError(ErrCodeCardAlreadyCancelled, "该卡已注销")
	}

	deposit := card.Deposit
	balance := card.Balance
	card.Status = model.CardStatusCancelled
	card.Balance = 0
	if err := s.cardRepo.UpdateCard(card); err != nil {
		return nil, err
	}

	return &CancellationResult{
		Card:    card,
		Deposit: deposit,
		Balance: balance,
		Total:   deposit + balance,
	}, nil
}
