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
	// 初始化日志：APP_ENV=production 使用 JSON 结构化输出，否则使用彩色控制台输出
	if os.Getenv("APP_ENV") == "production" {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		log.Logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
	} else {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
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
