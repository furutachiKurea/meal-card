// Package service 提供业务逻辑层
package service

import (
	"backend/model"
	"backend/repository"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// 业务错误码常量
const (
	ErrCodeCardNotFound         = "CARD_NOT_FOUND"
	ErrCodeCardAlreadyActive    = "CARD_ALREADY_ACTIVE"
	ErrCodeCardNotActive        = "CARD_NOT_ACTIVE"
	ErrCodeCardNotLost          = "CARD_NOT_LOST"
	ErrCodeCardCancelled        = "CARD_CANCELLED"
	ErrCodeCardLost             = "CARD_LOST"
	ErrCodeCardAlreadyCancelled = "CARD_ALREADY_CANCELLED"
	ErrCodeInsufficientBalance  = "INSUFFICIENT_BALANCE"
	ErrCodeInvalidAmount        = "INVALID_AMOUNT"
	ErrCodeWindowNotFound       = "WINDOW_NOT_FOUND"
	ErrCodeValidationError      = "VALIDATION_ERROR"
	ErrCodeStudentNotFound      = "STUDENT_NOT_FOUND"
	ErrCodeStudentServiceError  = "STUDENT_SERVICE_ERROR"
	ErrCodeExceedSingleLimit    = "EXCEED_SINGLE_LIMIT"
	ErrCodeExceedDailyLimit     = "EXCEED_DAILY_LIMIT"
)

// 消费限额（单位：分）
const (
	MaxSingleTransaction = 20000 // 单笔消费上限 200 元
	MaxDailyTransaction  = 50000 // 单日累计消费上限 500 元
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

// IssueCardResult 发卡结果
type IssueCardResult struct {
	Card       *model.Card
	CardHolder *model.CardHolder
	// Refund 旧卡退款信息，仅在有 lost 卡被自动注销时不为 nil
	Refund *OldCardRefund
}

// OldCardRefund 旧卡退款详情
type OldCardRefund struct {
	OldCardNo string
	Deposit   int64
	Balance   int64
	Total     int64
}

// DepositResult 存款结果（收据信息）
type DepositResult struct {
	ID         uint
	CardNo     string
	HolderName string
	Amount     int64
	NewBalance int64
	CreatedAt  time.Time
}

// TransactionResult 消费结果
type TransactionResult struct {
	ID         uint
	CardNo     string
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
	validator  StudentValidator
}

// NewCardService 创建 CardService 实例
func NewCardService(cardRepo *repository.CardRepository, windowRepo *repository.WindowRepository, validator StudentValidator) *CardService {
	return &CardService{
		cardRepo:   cardRepo,
		windowRepo: windowRepo,
		validator:  validator,
	}
}

// generateCardNo 生成 16 位随机数字卡号
func generateCardNo() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%016d", r.Int63n(10000000000000000))
}

// IssueCard 发卡业务
// 参数：idNumber 证件号，preDeposit 预存款（分）
// 规则：
//  1. 调用 validator 验证证件号，失败返回对应错误码
//  2. 同证件号已有 active 卡 → 返回 409 CARD_ALREADY_ACTIVE
//  3. 同证件号有 lost 卡 → 自动注销旧卡，创建新卡，返回退款信息
//  4. 无任何卡 → 直接创建，押金固定 2000 分
func (s *CardService) IssueCard(idNumber string, preDeposit int64) (*IssueCardResult, error) {
	if preDeposit < 0 {
		log.Error().Str("code", ErrCodeInvalidAmount).Str("msg", "预存款不能为负数").Msg("业务错误")
		return nil, newBizError(ErrCodeInvalidAmount, "预存款不能为负数")
	}

	// 验证证件号（外部服务调用，不在事务内）
	studentInfo, err := s.validator.Validate(idNumber)
	if err != nil {
		return nil, err
	}

	var result *IssueCardResult
	err = s.cardRepo.DB().Transaction(func(tx *gorm.DB) error {
		txRepo := s.cardRepo.WithTx(tx)

		// 查找或创建持卡人
		holder, err := txRepo.FindCardHolderByIDNumber(idNumber)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error().Err(err).Msg("查询持卡人失败")
			return err
		}

		if holder == nil {
			holder = &model.CardHolder{
				Name:     studentInfo.Name,
				IDNumber: idNumber,
			}
			if err := txRepo.CreateCardHolder(holder); err != nil {
				log.Error().Err(err).Msg("创建持卡人失败")
				return err
			}
		}

		// 检查是否已有 active 卡
		activeCard, err := txRepo.FindActiveCardByHolderID(holder.ID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error().Err(err).Msg("查询 active 卡失败")
			return err
		}
		if activeCard != nil {
			log.Error().Str("code", ErrCodeCardAlreadyActive).Str("idNumber", idNumber).Msg("业务错误")
			return newBizError(ErrCodeCardAlreadyActive, "该证件号已有正在使用的卡")
		}

		// 检查是否有 lost 卡，若有则自动注销
		var refund *OldCardRefund
		lostCard, err := txRepo.FindLostCardByHolderID(holder.ID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error().Err(err).Msg("查询 lost 卡失败")
			return err
		}
		if lostCard != nil {
			oldDeposit := lostCard.Deposit
			oldBalance := lostCard.Balance
			oldCardNo := lostCard.CardNo
			lostCard.Status = model.CardStatusCancelled
			lostCard.Balance = 0
			if err := txRepo.UpdateCard(lostCard); err != nil {
				log.Error().Err(err).Str("oldCardNo", oldCardNo).Msg("自动注销旧卡失败")
				return err
			}
			refund = &OldCardRefund{
				OldCardNo: oldCardNo,
				Deposit:   oldDeposit,
				Balance:   oldBalance,
				Total:     oldDeposit + oldBalance,
			}
		}

		// 生成新卡号，押金固定 2000 分
		newCard := &model.Card{
			CardNo:       generateCardNo(),
			CardHolderID: holder.ID,
			Deposit:      2000,
			Balance:      preDeposit,
			Status:       model.CardStatusActive,
		}
		if err := txRepo.CreateCard(newCard); err != nil {
			log.Error().Err(err).Msg("创建卡失败")
			return err
		}

		if preDeposit > 0 {
			record := &model.DepositRecord{
				CardID: newCard.ID,
				Amount: preDeposit,
			}
			if err := txRepo.CreateDepositRecord(record); err != nil {
				log.Error().Err(err).Str("cardNo", newCard.CardNo).Msg("创建预存款记录失败")
				return err
			}
		}

		result = &IssueCardResult{
			Card:       newCard,
			CardHolder: holder,
			Refund:     refund,
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetCard 根据 16 位卡号查询卡片详情
func (s *CardService) GetCard(cardNo string) (*model.Card, error) {
	card, err := s.cardRepo.FindCardByCardNo(cardNo)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error().Str("code", ErrCodeCardNotFound).Str("cardNo", cardNo).Msg("业务错误")
			return nil, newBizError(ErrCodeCardNotFound, "卡号不存在")
		}
		log.Error().Err(err).Str("cardNo", cardNo).Msg("查询卡失败")
		return nil, err
	}
	return card, nil
}

// GetCurrentCardByIDNumber 根据证件号查询当前有效卡（active 或 lost）
func (s *CardService) GetCurrentCardByIDNumber(idNumber string) (*model.Card, error) {
	card, err := s.cardRepo.FindCurrentCardByIDNumber(idNumber)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error().Str("code", ErrCodeCardNotFound).Str("idNumber", idNumber).Msg("业务错误")
			return nil, newBizError(ErrCodeCardNotFound, "该证件号无有效卡")
		}
		log.Error().Err(err).Str("idNumber", idNumber).Msg("查询卡失败")
		return nil, err
	}
	return card, nil
}

// Deposit 存款（充值）
func (s *CardService) Deposit(cardNo string, amount int64) (*DepositResult, error) {
	if amount <= 0 {
		log.Error().Str("code", ErrCodeInvalidAmount).Str("cardNo", cardNo).Msg("业务错误")
		return nil, newBizError(ErrCodeInvalidAmount, "充值金额必须大于 0")
	}

	var result *DepositResult
	err := s.cardRepo.DB().Transaction(func(tx *gorm.DB) error {
		txRepo := s.cardRepo.WithTx(tx)

		card, err := txRepo.FindCardByCardNo(cardNo)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				log.Error().Str("code", ErrCodeCardNotFound).Str("cardNo", cardNo).Msg("业务错误")
				return newBizError(ErrCodeCardNotFound, "卡号不存在")
			}
			log.Error().Err(err).Str("cardNo", cardNo).Msg("查询卡失败")
			return err
		}

		if card.Status != model.CardStatusActive {
			if card.Status == model.CardStatusLost {
				log.Error().Str("code", ErrCodeCardNotActive).Str("cardNo", cardNo).Msg("业务错误")
				return newBizError(ErrCodeCardNotActive, "该卡已挂失，无法充值")
			}
			log.Error().Str("code", ErrCodeCardNotActive).Str("cardNo", cardNo).Msg("业务错误")
			return newBizError(ErrCodeCardNotActive, "该卡已注销，无法充值")
		}

		card.Balance += amount
		if err := txRepo.UpdateCard(card); err != nil {
			log.Error().Err(err).Str("cardNo", cardNo).Msg("更新卡余额失败")
			return err
		}

		record := &model.DepositRecord{
			CardID: card.ID,
			Amount: amount,
		}
		if err := txRepo.CreateDepositRecord(record); err != nil {
			log.Error().Err(err).Str("cardNo", cardNo).Msg("创建存款记录失败")
			return err
		}

		result = &DepositResult{
			ID:         record.ID,
			CardNo:     card.CardNo,
			HolderName: card.CardHolder.Name,
			Amount:     amount,
			NewBalance: card.Balance,
			CreatedAt:  record.CreatedAt,
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// CreateTransaction 就餐消费结算
// 三重校验：卡存在 → 非 cancelled → 非 lost → 余额充足
// 限额校验：单笔不超 200 元、当日累计不超 500 元
func (s *CardService) CreateTransaction(cardNo string, windowID uint, amount int64) (*TransactionResult, error) {
	if amount <= 0 {
		log.Error().Str("code", ErrCodeInvalidAmount).Str("cardNo", cardNo).Msg("业务错误")
		return nil, newBizError(ErrCodeInvalidAmount, "消费金额必须大于 0")
	}
	if amount > MaxSingleTransaction {
		log.Error().Str("code", ErrCodeExceedSingleLimit).Str("cardNo", cardNo).Int64("amount", amount).Msg("业务错误")
		return nil, newBizError(ErrCodeExceedSingleLimit, fmt.Sprintf("单笔消费不能超过 %.0f 元", float64(MaxSingleTransaction)/100))
	}

	// 先在事务外验证窗口存在（只读查询，无需事务保护）
	_, err := s.windowRepo.FindWindowByID(windowID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error().Str("code", ErrCodeWindowNotFound).Uint("windowID", windowID).Msg("业务错误")
			return nil, newBizError(ErrCodeWindowNotFound, "窗口不存在")
		}
		log.Error().Err(err).Uint("windowID", windowID).Msg("查询窗口失败")
		return nil, err
	}

	var result *TransactionResult
	err = s.cardRepo.DB().Transaction(func(tx *gorm.DB) error {
		txRepo := s.cardRepo.WithTx(tx)

		card, err := txRepo.FindCardByCardNo(cardNo)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				log.Error().Str("code", ErrCodeCardNotFound).Str("cardNo", cardNo).Msg("业务错误")
				return newBizError(ErrCodeCardNotFound, "此卡非本单位所发")
			}
			log.Error().Err(err).Str("cardNo", cardNo).Msg("查询卡失败")
			return err
		}

		if card.Status == model.CardStatusCancelled {
			log.Error().Str("code", ErrCodeCardCancelled).Str("cardNo", cardNo).Msg("业务错误")
			return newBizError(ErrCodeCardCancelled, "此卡已注销")
		}
		if card.Status == model.CardStatusLost {
			log.Error().Str("code", ErrCodeCardLost).Str("cardNo", cardNo).Msg("业务错误")
			return newBizError(ErrCodeCardLost, "此卡已挂失")
		}

		// 日消费限额校验
		now := time.Now()
		dayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		dayEnd := dayStart.Add(24 * time.Hour)
		todaySpent, err := txRepo.SumCardTransactionsByTimeRange(card.ID, dayStart, dayEnd)
		if err != nil {
			log.Error().Err(err).Str("cardNo", cardNo).Msg("查询日消费额失败")
			return err
		}
		if todaySpent+amount > MaxDailyTransaction {
			log.Error().Str("code", ErrCodeExceedDailyLimit).Str("cardNo", cardNo).Int64("todaySpent", todaySpent).Int64("amount", amount).Msg("业务错误")
			return newBizError(ErrCodeExceedDailyLimit, fmt.Sprintf("今日累计消费将超过 %.0f 元限额", float64(MaxDailyTransaction)/100))
		}

		if card.Balance < amount {
			log.Error().Str("code", ErrCodeInsufficientBalance).Str("cardNo", cardNo).Int64("balance", card.Balance).Int64("amount", amount).Msg("业务错误")
			return newBizError(ErrCodeInsufficientBalance, "余额不足")
		}

		card.Balance -= amount
		if err := txRepo.UpdateCard(card); err != nil {
			log.Error().Err(err).Str("cardNo", cardNo).Msg("更新卡余额失败")
			return err
		}

		txn := &model.Transaction{
			CardID:   card.ID,
			WindowID: windowID,
			Amount:   amount,
		}
		if err := txRepo.CreateTransaction(txn); err != nil {
			log.Error().Err(err).Str("cardNo", cardNo).Msg("创建消费记录失败")
			return err
		}

		result = &TransactionResult{
			ID:         txn.ID,
			CardNo:     card.CardNo,
			WindowID:   windowID,
			Amount:     amount,
			NewBalance: card.Balance,
			CreatedAt:  txn.CreatedAt,
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// CardTransactionsResult 消费历史查询结果
type CardTransactionsResult struct {
	Records  []repository.TransactionRecord
	Total    int64
	Page     int
	PageSize int
}

// GetCardTransactions 查询指定卡的消费历史
func (s *CardService) GetCardTransactions(cardNo string, page, pageSize int) (*CardTransactionsResult, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	card, err := s.cardRepo.FindCardByCardNo(cardNo)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, newBizError(ErrCodeCardNotFound, "卡号不存在")
		}
		return nil, err
	}

	records, total, err := s.cardRepo.GetCardTransactions(card.ID, page, pageSize)
	if err != nil {
		return nil, err
	}

	return &CardTransactionsResult{
		Records:  records,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

// CardDepositsResult 存款历史查询结果
type CardDepositsResult struct {
	Records  []repository.DepositDetailItem
	Total    int64
	Page     int
	PageSize int
}

// GetCardDeposits 查询指定卡的存款历史
func (s *CardService) GetCardDeposits(cardNo string, page, pageSize int) (*CardDepositsResult, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	card, err := s.cardRepo.FindCardByCardNo(cardNo)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, newBizError(ErrCodeCardNotFound, "卡号不存在")
		}
		return nil, err
	}

	// 复用已有的 GetHolderDeposits，但按 card 的 holder 来查
	deposits, total, err := s.cardRepo.GetHolderDeposits(card.CardHolderID, nil, nil, page, pageSize)
	if err != nil {
		return nil, err
	}

	return &CardDepositsResult{
		Records:  deposits,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}
func (s *CardService) ReportLoss(cardNo string) (*model.Card, error) {
	card, err := s.cardRepo.FindCardByCardNo(cardNo)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error().Str("code", ErrCodeCardNotFound).Str("cardNo", cardNo).Msg("业务错误")
			return nil, newBizError(ErrCodeCardNotFound, "卡号不存在")
		}
		log.Error().Err(err).Str("cardNo", cardNo).Msg("查询卡失败")
		return nil, err
	}

	if card.Status != model.CardStatusActive {
		log.Error().Str("code", ErrCodeCardNotActive).Str("cardNo", cardNo).Msg("业务错误")
		return nil, newBizError(ErrCodeCardNotActive, "只有正常使用中的卡可以挂失")
	}

	card.Status = model.CardStatusLost
	if err := s.cardRepo.UpdateCard(card); err != nil {
		log.Error().Err(err).Str("cardNo", cardNo).Msg("更新卡状态失败")
		return nil, err
	}

	return card, nil
}

// CancelLossReport 取消挂失：lost → active
func (s *CardService) CancelLossReport(cardNo string) (*model.Card, error) {
	card, err := s.cardRepo.FindCardByCardNo(cardNo)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error().Str("code", ErrCodeCardNotFound).Str("cardNo", cardNo).Msg("业务错误")
			return nil, newBizError(ErrCodeCardNotFound, "卡号不存在")
		}
		log.Error().Err(err).Str("cardNo", cardNo).Msg("查询卡失败")
		return nil, err
	}

	if card.Status != model.CardStatusLost {
		log.Error().Str("code", ErrCodeCardNotLost).Str("cardNo", cardNo).Msg("业务错误")
		return nil, newBizError(ErrCodeCardNotLost, "只有已挂失的卡可以取消挂失")
	}

	card.Status = model.CardStatusActive
	if err := s.cardRepo.UpdateCard(card); err != nil {
		log.Error().Err(err).Str("cardNo", cardNo).Msg("更新卡状态失败")
		return nil, err
	}

	return card, nil
}

// CancelCard 注销卡片，退还押金和余额
func (s *CardService) CancelCard(cardNo string) (*CancellationResult, error) {
	card, err := s.cardRepo.FindCardByCardNo(cardNo)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error().Str("code", ErrCodeCardNotFound).Str("cardNo", cardNo).Msg("业务错误")
			return nil, newBizError(ErrCodeCardNotFound, "卡号不存在")
		}
		log.Error().Err(err).Str("cardNo", cardNo).Msg("查询卡失败")
		return nil, err
	}

	if card.Status == model.CardStatusCancelled {
		log.Error().Str("code", ErrCodeCardAlreadyCancelled).Str("cardNo", cardNo).Msg("业务错误")
		return nil, newBizError(ErrCodeCardAlreadyCancelled, "该卡已注销")
	}

	deposit := card.Deposit
	balance := card.Balance
	card.Status = model.CardStatusCancelled
	card.Balance = 0
	if err := s.cardRepo.UpdateCard(card); err != nil {
		log.Error().Err(err).Str("cardNo", cardNo).Msg("注销卡失败")
		return nil, err
	}

	return &CancellationResult{
		Card:    card,
		Deposit: deposit,
		Balance: balance,
		Total:   deposit + balance,
	}, nil
}
