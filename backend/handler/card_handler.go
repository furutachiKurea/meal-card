// Package handler 提供 HTTP 接口处理层
package handler

import (
	"backend/model"
	"backend/repository"
	"backend/service"
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// CardHandler 饭卡相关 HTTP 处理
type CardHandler struct {
	cardSvc   *service.CardService
	validator service.StudentValidator
	cardRepo  *repository.CardRepository
}

// NewCardHandler 创建 CardHandler 实例
func NewCardHandler(cardSvc *service.CardService, validator service.StudentValidator, cardRepo *repository.CardRepository) *CardHandler {
	return &CardHandler{
		cardSvc:   cardSvc,
		validator: validator,
		cardRepo:  cardRepo,
	}
}

// errorResponse 统一错误响应结构
type errorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// bizErrStatus 将 BizError 转换为对应 HTTP 状态码
func bizErrStatus(code string) int {
	switch code {
	case service.ErrCodeCardNotFound, service.ErrCodeWindowNotFound, service.ErrCodeStudentNotFound:
		return http.StatusNotFound
	case service.ErrCodeCardAlreadyActive,
		service.ErrCodeCardNotActive,
		service.ErrCodeCardNotLost,
		service.ErrCodeCardCancelled,
		service.ErrCodeCardLost,
		service.ErrCodeCardAlreadyCancelled,
		service.ErrCodeInsufficientBalance:
		return http.StatusConflict
	case service.ErrCodeInvalidAmount, service.ErrCodeValidationError:
		return http.StatusBadRequest
	case service.ErrCodeStudentServiceError:
		return http.StatusBadGateway
	default:
		return http.StatusInternalServerError
	}
}

// handleError 统一错误处理
func handleError(c echo.Context, err error) error {
	var bizErr *service.BizError
	if errors.As(err, &bizErr) {
		log.Warn().Str("path", c.Request().URL.Path).Str("method", c.Request().Method).Str("code", bizErr.Code).Msg(bizErr.Message)
		return c.JSON(bizErrStatus(bizErr.Code), errorResponse{
			Code:    bizErr.Code,
			Message: bizErr.Message,
		})
	}
	log.Error().Str("path", c.Request().URL.Path).Str("method", c.Request().Method).Err(err).Msg("内部错误")
	return c.JSON(http.StatusInternalServerError, errorResponse{
		Code:    "INTERNAL_ERROR",
		Message: "服务端内部错误",
	})
}

// cardToJSON 将 Card 模型转换为 JSON 响应格式
func cardToJSON(card *model.Card) map[string]any {
	return map[string]any{
		"cardNo":       card.CardNo,
		"cardHolderId": card.CardHolderID,
		"deposit":      card.Deposit,
		"balance":      card.Balance,
		"status":       card.Status,
		"createdAt":    card.CreatedAt.Format(time.RFC3339),
		"updatedAt":    card.UpdatedAt.Format(time.RFC3339),
	}
}

// holderToJSON 将 CardHolder 模型转换为 JSON 响应格式
func holderToJSON(holder *model.CardHolder) map[string]any {
	return map[string]any{
		"id":        holder.ID,
		"name":      holder.Name,
		"idNumber":  holder.IDNumber,
		"createdAt": holder.CreatedAt.Format(time.RFC3339),
	}
}

// ValidateStudent GET /api/validate-student 验证证件号
func (h *CardHandler) ValidateStudent(c echo.Context) error {
	idNumber := c.QueryParam("idNumber")
	if idNumber == "" {
		return c.JSON(http.StatusBadRequest, errorResponse{Code: "VALIDATION_ERROR", Message: "idNumber 不能为空"})
	}

	info, err := h.validator.Validate(idNumber)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, map[string]any{
		"valid":    true,
		"idNumber": info.IDNumber,
		"name":     info.Name,
		"type":     info.Type,
	})
}

// GetCardByIDNumber GET /api/cards?idNumber=xxx 按证件号查询当前有效卡
func (h *CardHandler) GetCardByIDNumber(c echo.Context) error {
	idNumber := c.QueryParam("idNumber")
	if idNumber == "" {
		return c.JSON(http.StatusBadRequest, errorResponse{Code: "VALIDATION_ERROR", Message: "idNumber 不能为空"})
	}

	card, err := h.cardRepo.FindCurrentCardByIDNumber(idNumber)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, errorResponse{Code: "CARD_NOT_FOUND", Message: "该证件号无有效卡"})
		}
		return c.JSON(http.StatusInternalServerError, errorResponse{Code: "INTERNAL_ERROR", Message: "服务端内部错误"})
	}

	resp := cardToJSON(card)
	resp["cardHolder"] = holderToJSON(&card.CardHolder)
	return c.JSON(http.StatusOK, resp)
}

// IssueCard POST /api/cards 发卡
func (h *CardHandler) IssueCard(c echo.Context) error {
	var req struct {
		IDNumber   string `json:"idNumber"`
		PreDeposit int64  `json:"preDeposit"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse{Code: "VALIDATION_ERROR", Message: "请求体格式错误"})
	}
	log.Info().Str("path", "POST /api/cards").Str("idNumber", req.IDNumber).Int64("preDeposit", req.PreDeposit).Msg("发卡请求")

	result, err := h.cardSvc.IssueCard(req.IDNumber, req.PreDeposit)
	if err != nil {
		return handleError(c, err)
	}

	resp := map[string]any{
		"card":       cardToJSON(result.Card),
		"cardHolder": holderToJSON(result.CardHolder),
	}
	if result.Refund != nil {
		resp["refund"] = map[string]any{
			"oldCardNo": result.Refund.OldCardNo,
			"deposit":   result.Refund.Deposit,
			"balance":   result.Refund.Balance,
			"total":     result.Refund.Total,
		}
	}
	return c.JSON(http.StatusOK, resp)
}

// GetCard GET /api/cards/:cardNo 按卡号查询卡信息
func (h *CardHandler) GetCard(c echo.Context) error {
	cardNo := c.Param("cardNo")

	card, err := h.cardSvc.GetCard(cardNo)
	if err != nil {
		return handleError(c, err)
	}

	resp := cardToJSON(card)
	resp["cardHolder"] = holderToJSON(&card.CardHolder)
	return c.JSON(http.StatusOK, resp)
}

// Deposit POST /api/cards/:cardNo/deposits 存款
func (h *CardHandler) Deposit(c echo.Context) error {
	cardNo := c.Param("cardNo")

	var req struct {
		Amount int64 `json:"amount"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse{Code: "VALIDATION_ERROR", Message: "请求体格式错误"})
	}
	log.Info().Str("path", "POST /api/cards/:cardNo/deposits").Str("cardNo", cardNo).Int64("amount", req.Amount).Msg("存款请求")

	result, err := h.cardSvc.Deposit(cardNo, req.Amount)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, map[string]any{
		"id":         result.ID,
		"cardNo":     result.CardNo,
		"holderName": result.HolderName,
		"amount":     result.Amount,
		"newBalance": result.NewBalance,
		"createdAt":  result.CreatedAt.Format(time.RFC3339),
	})
}

// CreateTransaction POST /api/cards/:cardNo/transactions 就餐消费
func (h *CardHandler) CreateTransaction(c echo.Context) error {
	cardNo := c.Param("cardNo")

	var req struct {
		WindowID int64 `json:"windowId"`
		Amount   int64 `json:"amount"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse{Code: "VALIDATION_ERROR", Message: "请求体格式错误"})
	}
	log.Info().Str("path", "POST /api/cards/:cardNo/transactions").Str("cardNo", cardNo).Int64("windowId", req.WindowID).Int64("amount", req.Amount).Msg("消费请求")

	result, err := h.cardSvc.CreateTransaction(cardNo, uint(req.WindowID), req.Amount)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, map[string]any{
		"id":         result.ID,
		"cardNo":     result.CardNo,
		"windowId":   result.WindowID,
		"amount":     result.Amount,
		"newBalance": result.NewBalance,
		"createdAt":  result.CreatedAt.Format(time.RFC3339),
	})
}

// ReportLoss PUT /api/cards/:cardNo/loss-report 挂失
func (h *CardHandler) ReportLoss(c echo.Context) error {
	cardNo := c.Param("cardNo")
	log.Info().Str("path", "PUT /api/cards/:cardNo/loss-report").Str("cardNo", cardNo).Msg("挂失请求")

	card, err := h.cardSvc.ReportLoss(cardNo)
	if err != nil {
		return handleError(c, err)
	}

	resp := cardToJSON(card)
	resp["cardHolder"] = holderToJSON(&card.CardHolder)
	return c.JSON(http.StatusOK, resp)
}

// CancelLossReport DELETE /api/cards/:cardNo/loss-report 取消挂失
func (h *CardHandler) CancelLossReport(c echo.Context) error {
	cardNo := c.Param("cardNo")
	log.Info().Str("path", "DELETE /api/cards/:cardNo/loss-report").Str("cardNo", cardNo).Msg("取消挂失请求")

	card, err := h.cardSvc.CancelLossReport(cardNo)
	if err != nil {
		return handleError(c, err)
	}

	resp := cardToJSON(card)
	resp["cardHolder"] = holderToJSON(&card.CardHolder)
	return c.JSON(http.StatusOK, resp)
}

// CancelCard POST /api/cards/:cardNo/cancellation 注销
func (h *CardHandler) CancelCard(c echo.Context) error {
	cardNo := c.Param("cardNo")
	log.Info().Str("path", "POST /api/cards/:cardNo/cancellation").Str("cardNo", cardNo).Msg("注销请求")

	result, err := h.cardSvc.CancelCard(cardNo)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, map[string]any{
		"card": cardToJSON(result.Card),
		"refund": map[string]any{
			"deposit": result.Deposit,
			"balance": result.Balance,
			"total":   result.Total,
		},
	})
}
