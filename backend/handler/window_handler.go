package handler

import (
	"backend/model"
	"backend/service"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

// WindowHandler 窗口管理 HTTP 处理
type WindowHandler struct {
	windowSvc *service.WindowService
}

// NewWindowHandler 创建 WindowHandler 实例
func NewWindowHandler(windowSvc *service.WindowService) *WindowHandler {
	return &WindowHandler{windowSvc: windowSvc}
}

// windowToJSON 将 Window 模型转换为 JSON 响应格式
func windowToJSON(w *model.Window) map[string]any {
	return map[string]any{
		"id":   w.ID,
		"name": w.Name,
	}
}

// ListWindows GET /api/windows 获取窗口列表
func (h *WindowHandler) ListWindows(c echo.Context) error {
	windows, err := h.windowSvc.ListWindows()
	if err != nil {
		return handleError(c, err)
	}

	items := make([]map[string]any, 0, len(windows))
	for i := range windows {
		items = append(items, windowToJSON(&windows[i]))
	}
	return c.JSON(http.StatusOK, map[string]any{"windows": items})
}

// CreateWindow POST /api/windows 创建窗口
func (h *WindowHandler) CreateWindow(c echo.Context) error {
	var req struct {
		Name string `json:"name"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse{Code: "VALIDATION_ERROR", Message: "请求体格式错误"})
	}
	log.Info().Str("path", "POST /api/windows").Str("name", req.Name).Msg("创建窗口请求")

	w, err := h.windowSvc.CreateWindow(req.Name)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, windowToJSON(w))
}
