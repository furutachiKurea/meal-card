package main

import (
	"backend/db"
	"backend/handler"
	"backend/repository"
	"backend/router"
	"backend/service"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// 初始化日志：
	// APP_ENV=production  → JSON 结构化输出，Info 级别
	// 其他（开发环境）    → 彩色控制台输出，Info 级别
	// LOG_LEVEL=debug     → 无论哪个环境，强制开启 Debug 级别（调试开关）
	level := zerolog.InfoLevel
	if os.Getenv("LOG_LEVEL") == "debug" {
		level = zerolog.DebugLevel
	}
	zerolog.SetGlobalLevel(level)

	if os.Getenv("APP_ENV") == "production" {
		log.Logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
	} else {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	// 初始化数据库
	gormDB, err := db.Init("meal_card.db")
	if err != nil {
		log.Fatal().Err(err).Msg("数据库初始化失败")
	}

	// 初始化 repository
	cardRepo := repository.NewCardRepository(gormDB)
	windowRepo := repository.NewWindowRepository(gormDB)

	// 初始化 service
	cardSvc := service.NewCardService(cardRepo, windowRepo)
	statsSvc := service.NewStatisticsService(cardRepo)
	windowSvc := service.NewWindowService(windowRepo)

	// 初始化 handler
	cardH := handler.NewCardHandler(cardSvc)
	statsH := handler.NewStatisticsHandler(statsSvc)
	windowH := handler.NewWindowHandler(windowSvc)

	// 初始化 echo 并注册路由
	e := echo.New()
	e.HideBanner = true
	router.Register(e, cardH, statsH, windowH)

	log.Info().Msg("服务启动，监听 :8080")
	if err := e.Start(":8080"); err != nil {
		log.Fatal().Err(err).Msg("服务启动失败")
	}
}
