package handler

import (
	"backend/service"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

// StatisticsHandler 统计相关 HTTP 处理
type StatisticsHandler struct {
	statsSvc *service.StatisticsService
}

// NewStatisticsHandler 创建 StatisticsHandler 实例
func NewStatisticsHandler(statsSvc *service.StatisticsService) *StatisticsHandler {
	return &StatisticsHandler{statsSvc: statsSvc}
}

// parseTimeParam 解析查询参数中的 ISO 8601 时间字符串
func parseTimeParam(c echo.Context, name string, required bool) (*time.Time, error) {
	val := c.QueryParam(name)
	if val == "" {
		if required {
			return nil, &service.BizError{Code: service.ErrCodeValidationError, Message: "缺少必填参数: " + name}
		}
		return nil, nil
	}
	t, err := time.Parse(time.RFC3339, val)
	if err != nil {
		return nil, &service.BizError{Code: service.ErrCodeValidationError, Message: "时间格式无效，应为 ISO 8601: " + name}
	}
	return &t, nil
}

// GetMealRevenue GET /api/statistics/meal-revenue 本餐售饭总收入
func (h *StatisticsHandler) GetMealRevenue(c echo.Context) error {
	start, err := parseTimeParam(c, "startTime", true)
	if err != nil {
		return handleError(c, err)
	}
	end, err := parseTimeParam(c, "endTime", true)
	if err != nil {
		return handleError(c, err)
	}

	result, err := h.statsSvc.GetMealRevenue(*start, *end)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, map[string]any{
		"totalRevenue": result.TotalRevenue,
	})
}

// GetWindowRevenue GET /api/statistics/window-revenue 各窗口收入
func (h *StatisticsHandler) GetWindowRevenue(c echo.Context) error {
	start, err := parseTimeParam(c, "startTime", true)
	if err != nil {
		return handleError(c, err)
	}
	end, err := parseTimeParam(c, "endTime", true)
	if err != nil {
		return handleError(c, err)
	}

	result, err := h.statsSvc.GetWindowRevenue(*start, *end)
	if err != nil {
		return handleError(c, err)
	}

	windows := make([]map[string]any, 0, len(result.Windows))
	for _, w := range result.Windows {
		windows = append(windows, map[string]any{
			"windowId":   w.WindowID,
			"windowName": w.WindowName,
			"revenue":    w.Revenue,
		})
	}
	return c.JSON(http.StatusOK, map[string]any{"windows": windows})
}

// GetDepositDetails GET /api/statistics/deposit-details 各持卡人存款明细（支持分页）
func (h *StatisticsHandler) GetDepositDetails(c echo.Context) error {
	start, err := parseTimeParam(c, "startTime", false)
	if err != nil {
		return handleError(c, err)
	}
	end, err := parseTimeParam(c, "endTime", false)
	if err != nil {
		return handleError(c, err)
	}

	page := 1
	pageSize := 10
	if p := c.QueryParam("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v > 0 {
			page = v
		}
	}
	if ps := c.QueryParam("pageSize"); ps != "" {
		if v, err := strconv.Atoi(ps); err == nil && v > 0 {
			pageSize = v
		}
	}

	result, err := h.statsSvc.GetDepositDetails(start, end, page, pageSize)
	if err != nil {
		return handleError(c, err)
	}

	holders := make([]map[string]any, 0, len(result.Holders))
	for _, h := range result.Holders {
		deposits := make([]map[string]any, 0, len(h.Deposits))
		for _, d := range h.Deposits {
			deposits = append(deposits, map[string]any{
				"id":        d.ID,
				"cardNo":    d.CardNo,
				"amount":    d.Amount,
				"createdAt": d.CreatedAt.Format(time.RFC3339),
			})
		}
		holders = append(holders, map[string]any{
			"holderId":    h.HolderID,
			"holderName":  h.HolderName,
			"idNumber":    h.IDNumber,
			"deposits":    deposits,
			"totalAmount": h.TotalAmount,
		})
	}
	return c.JSON(http.StatusOK, map[string]any{
		"holders":  holders,
		"total":    result.Total,
		"page":     result.Page,
		"pageSize": result.PageSize,
	})
}

// GetHolderDeposits GET /api/statistics/holder-deposits 指定持卡人的存款明细（支持分页）
func (h *StatisticsHandler) GetHolderDeposits(c echo.Context) error {
	holderIDStr := c.QueryParam("holderId")
	if holderIDStr == "" {
		return c.JSON(http.StatusBadRequest, errorResponse{Code: "VALIDATION_ERROR", Message: "缺少必填参数: holderId"})
	}
	holderIDInt, err := strconv.ParseUint(holderIDStr, 10, 64)
	if err != nil || holderIDInt == 0 {
		return c.JSON(http.StatusBadRequest, errorResponse{Code: "VALIDATION_ERROR", Message: "holderId 参数无效"})
	}

	start, err := parseTimeParam(c, "startTime", false)
	if err != nil {
		return handleError(c, err)
	}
	end, err := parseTimeParam(c, "endTime", false)
	if err != nil {
		return handleError(c, err)
	}

	page := 1
	pageSize := 10
	if p := c.QueryParam("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v > 0 {
			page = v
		}
	}
	if ps := c.QueryParam("pageSize"); ps != "" {
		if v, err := strconv.Atoi(ps); err == nil && v > 0 {
			pageSize = v
		}
	}

	result, err := h.statsSvc.GetHolderDeposits(uint(holderIDInt), start, end, page, pageSize)
	if err != nil {
		return handleError(c, err)
	}

	deposits := make([]map[string]any, 0, len(result.Deposits))
	for _, d := range result.Deposits {
		deposits = append(deposits, map[string]any{
			"id":        d.ID,
			"cardNo":    d.CardNo,
			"amount":    d.Amount,
			"createdAt": d.CreatedAt.Format(time.RFC3339),
		})
	}
	return c.JSON(http.StatusOK, map[string]any{
		"holderId":   result.HolderID,
		"holderName": result.HolderName,
		"idNumber":   result.IDNumber,
		"deposits":   deposits,
		"total":      result.Total,
		"page":       result.Page,
		"pageSize":   result.PageSize,
	})
}

// GetDepositSummary GET /api/statistics/deposit-summary 本日/本月存款金额
func (h *StatisticsHandler) GetDepositSummary(c echo.Context) error {
	result, err := h.statsSvc.GetDepositSummary()
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, map[string]any{
		"todayTotal": result.TodayTotal,
		"monthTotal": result.MonthTotal,
	})
}

// GetActiveBalance GET /api/statistics/active-balance 卡中流动资金总额
func (h *StatisticsHandler) GetActiveBalance(c echo.Context) error {
	result, err := h.statsSvc.GetActiveBalance()
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, map[string]any{
		"totalBalance": result.TotalBalance,
	})
}

// GetDailyReport GET /api/statistics/daily-report 日餐报表
func (h *StatisticsHandler) GetDailyReport(c echo.Context) error {
	date := c.QueryParam("date")
	if date == "" {
		return c.JSON(http.StatusBadRequest, errorResponse{Code: "VALIDATION_ERROR", Message: "缺少必填参数: date"})
	}

	result, err := h.statsSvc.GetDailyReport(date)
	if err != nil {
		return handleError(c, err)
	}

	windows := make([]map[string]any, 0, len(result.Windows))
	for _, w := range result.Windows {
		windows = append(windows, map[string]any{
			"windowId":         w.WindowID,
			"windowName":       w.WindowName,
			"revenue":          w.Revenue,
			"transactionCount": w.TransactionCount,
		})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"date":             result.Date,
		"totalRevenue":     result.TotalRevenue,
		"transactionCount": result.TransactionCount,
		"windows":          windows,
	})
}

// GetYearlyReport GET /api/statistics/yearly-report 年餐报表
func (h *StatisticsHandler) GetYearlyReport(c echo.Context) error {
	yearStr := c.QueryParam("year")
	if yearStr == "" {
		return c.JSON(http.StatusBadRequest, errorResponse{Code: "VALIDATION_ERROR", Message: "缺少必填参数: year"})
	}
	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		return c.JSON(http.StatusBadRequest, errorResponse{Code: "VALIDATION_ERROR", Message: "year 参数无效"})
	}

	result, err := h.statsSvc.GetYearlyReport(year)
	if err != nil {
		return handleError(c, err)
	}

	months := make([]map[string]any, 0, len(result.Months))
	for _, m := range result.Months {
		months = append(months, map[string]any{
			"month":            m.Month,
			"revenue":          m.Revenue,
			"transactionCount": m.TransactionCount,
		})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"year":             result.Year,
		"totalRevenue":     result.TotalRevenue,
		"transactionCount": result.TransactionCount,
		"months":           months,
	})
}
