// Package db 负责数据库连接与初始化
package db

import (
	"backend/model"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// 默认窗口种子数据
var defaultWindows = []string{
	"一号窗口",
	"二号窗口",
	"三号窗口",
	"四号窗口",
	"五号窗口",
}

// Seed 初始化种子数据，仅在表为空时插入
func Seed(db *gorm.DB) {
	var count int64
	db.Model(&model.Window{}).Count(&count)
	if count > 0 {
		return
	}

	for _, name := range defaultWindows {
		db.Create(&model.Window{Name: name})
	}
	log.Info().Int("count", len(defaultWindows)).Msg("已初始化默认窗口")
}
