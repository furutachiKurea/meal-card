// Package router 负责注册所有 HTTP 路由
package router

import (
	"backend/handler"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
)

// zerologMiddleware 将每条 HTTP 请求通过 zerolog 输出
func zerologMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			err := next(c)
			req := c.Request()
			res := c.Response()
			log.Info().
				Str("method", req.Method).
				Str("path", req.URL.Path).
				Int("status", res.Status).
				Dur("latency", time.Since(start)).
				Msg("request")
			return err
		}
	}
}

// Register 注册所有路由，并配置 CORS 与日志中间件
func Register(e *echo.Echo, cardH *handler.CardHandler, statsH *handler.StatisticsHandler, windowH *handler.WindowHandler) {
	e.Use(zerologMiddleware())

	// 允许所有来源的跨域请求
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Content-Type", "Authorization"},
	}))

	api := e.Group("/api")

	// 学籍验证
	api.GET("/validate-student", cardH.ValidateStudent)

	// 饭卡
	api.POST("/cards", cardH.IssueCard)
	api.GET("/cards", cardH.GetCardByIDNumber)
	api.GET("/cards/:cardNo", cardH.GetCard)
	api.POST("/cards/:cardNo/deposits", cardH.Deposit)
	api.GET("/cards/:cardNo/deposits", cardH.GetCardDeposits)
	api.POST("/cards/:cardNo/transactions", cardH.CreateTransaction)
	api.GET("/cards/:cardNo/transactions", cardH.GetCardTransactions)
	api.PUT("/cards/:cardNo/loss-report", cardH.ReportLoss)
	api.DELETE("/cards/:cardNo/loss-report", cardH.CancelLossReport)
	api.POST("/cards/:cardNo/cancellation", cardH.CancelCard)

	// 统计
	api.GET("/statistics/meal-revenue", statsH.GetMealRevenue)
	api.GET("/statistics/window-revenue", statsH.GetWindowRevenue)
	api.GET("/statistics/deposit-details", statsH.GetDepositDetails)
	api.GET("/statistics/holder-deposits", statsH.GetHolderDeposits)
	api.GET("/statistics/deposit-summary", statsH.GetDepositSummary)
	api.GET("/statistics/active-balance", statsH.GetActiveBalance)
	api.GET("/statistics/daily-report", statsH.GetDailyReport)
	api.GET("/statistics/yearly-report", statsH.GetYearlyReport)

	// 窗口
	api.GET("/windows", windowH.ListWindows)
	api.POST("/windows", windowH.CreateWindow)
}
