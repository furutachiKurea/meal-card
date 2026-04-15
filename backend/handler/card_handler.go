// Package handler 提供 HTTP 接口处理层
package handler

import (
	"backend/model"
	"backend/service"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

// CardHandler 饭卡相关 HTTP 处理
type CardHandler struct {
	cardSvc *service.CardService
}

// NewCardHandler 创建 CardHandler 实例
func NewCardHandler(cardSvc *service.CardService) *CardHandler {
	return &CardHandler{cardSvc: cardSvc}
}

// errorResponse 统一错误响应结构
type errorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// bizErr 将 BizError 转换为对应 HTTP 状态码
func bizErrStatus(code string) int {
	switch code {
	case service.ErrCodeCardNotFound, service.ErrCodeWindowNotFound:
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
	default:
		return http.StatusInternalServerError
	}
}

// handleError 统一错误处理
func handleError(c echo.Context, err error) error {
	var bizErr *service.BizError
	if errors.As(err, &bizErr) {
		return c.JSON(bizErrStatus(bizErr.Code), errorResponse{
			Code:    bizErr.Code,
			Message: bizErr.Message,
		})
	}
	return c.JSON(http.StatusInternalServerError, errorResponse{
		Code:    "INTERNAL_ERROR",
		Message: "服务端内部错误",
	})
}

// parseCardID 解析路径参数中的卡号
func parseCardID(c echo.Context) (uint, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}

// cardToJSON 将 Card 模型转换为 JSON 响应格式
func cardToJSON(card *model.Card) map[string]any {
	return map[string]any{
		"id":           card.ID,
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

// IssueCard POST /api/cards 发卡
func (h *CardHandler) IssueCard(c echo.Context) error {
	var req struct {
		Name       string `json:"name"`
		IDNumber   string `json:"idNumber"`
		Deposit    int64  `json:"deposit"`
		PreDeposit int64  `json:"preDeposit"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse{Code: "VALIDATION_ERROR", Message: "请求体格式错误"})
	}

	result, err := h.cardSvc.IssueCard(service.IssueCardRequest{
		Name:       req.Name,
		IDNumber:   req.IDNumber,
		Deposit:    req.Deposit,
		PreDeposit: req.PreDeposit,
	})
	if err != nil {
		return handleError(c, err)
	}

	resp := map[string]any{
		"card":       cardToJSON(result.Card),
		"cardHolder": holderToJSON(result.CardHolder),
	}
	if result.Refund != nil {
		resp["refund"] = map[string]any{
			"oldCardId": result.Refund.OldCardID,
			"deposit":   result.Refund.Deposit,
			"balance":   result.Refund.Balance,
			"total":     result.Refund.Total,
		}
	}
	return c.JSON(http.StatusOK, resp)
}

// GetCard GET /api/cards/:id 查询卡信息
func (h *CardHandler) GetCard(c echo.Context) error {
	cardID, err := parseCardID(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse{Code: "VALIDATION_ERROR", Message: "卡号格式无效"})
	}

	card, err := h.cardSvc.GetCard(cardID)
	if err != nil {
		return handleError(c, err)
	}

	resp := cardToJSON(card)
	resp["cardHolder"] = holderToJSON(&card.CardHolder)
	return c.JSON(http.StatusOK, resp)
}

// Deposit POST /api/cards/:id/deposits 存款
func (h *CardHandler) Deposit(c echo.Context) error {
	cardID, err := parseCardID(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse{Code: "VALIDATION_ERROR", Message: "卡号格式无效"})
	}

	var req struct {
		Amount int64 `json:"amount"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse{Code: "VALIDATION_ERROR", Message: "请求体格式错误"})
	}

	result, err := h.cardSvc.Deposit(cardID, req.Amount)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, map[string]any{
		"id":         result.ID,
		"cardId":     result.CardID,
		"holderName": result.HolderName,
		"amount":     result.Amount,
		"newBalance": result.NewBalance,
		"createdAt":  result.CreatedAt.Format(time.RFC3339),
	})
}

// CreateTransaction POST /api/cards/:id/transactions 就餐消费
func (h *CardHandler) CreateTransaction(c echo.Context) error {
	cardID, err := parseCardID(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse{Code: "VALIDATION_ERROR", Message: "卡号格式无效"})
	}

	var req struct {
		WindowID int64 `json:"windowId"`
		Amount   int64 `json:"amount"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse{Code: "VALIDATION_ERROR", Message: "请求体格式错误"})
	}

	result, err := h.cardSvc.CreateTransaction(cardID, uint(req.WindowID), req.Amount)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, map[string]any{
		"id":         result.ID,
		"cardId":     result.CardID,
		"windowId":   result.WindowID,
		"amount":     result.Amount,
		"newBalance": result.NewBalance,
		"createdAt":  result.CreatedAt.Format(time.RFC3339),
	})
}

// ReportLoss PUT /api/cards/:id/loss-report 挂失
func (h *CardHandler) ReportLoss(c echo.Context) error {
	cardID, err := parseCardID(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse{Code: "VALIDATION_ERROR", Message: "卡号格式无效"})
	}

	card, err := h.cardSvc.ReportLoss(cardID)
	if err != nil {
		return handleError(c, err)
	}

	resp := cardToJSON(card)
	resp["cardHolder"] = holderToJSON(&card.CardHolder)
	return c.JSON(http.StatusOK, resp)
}

// CancelLossReport DELETE /api/cards/:id/loss-report 取消挂失
func (h *CardHandler) CancelLossReport(c echo.Context) error {
	cardID, err := parseCardID(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse{Code: "VALIDATION_ERROR", Message: "卡号格式无效"})
	}

	card, err := h.cardSvc.CancelLossReport(cardID)
	if err != nil {
		return handleError(c, err)
	}

	resp := cardToJSON(card)
	resp["cardHolder"] = holderToJSON(&card.CardHolder)
	return c.JSON(http.StatusOK, resp)
}

// CancelCard POST /api/cards/:id/cancellation 注销
func (h *CardHandler) CancelCard(c echo.Context) error {
	cardID, err := parseCardID(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse{Code: "VALIDATION_ERROR", Message: "卡号格式无效"})
	}

	result, err := h.cardSvc.CancelCard(cardID)
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
