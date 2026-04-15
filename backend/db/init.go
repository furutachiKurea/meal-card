// Package db 负责数据库连接与初始化
package db

import (
	"backend/model"

	"github.com/glebarez/sqlite"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// Init 初始化 SQLite 数据库连接并自动迁移表结构。
// dsn 为数据库文件路径，传入 ":memory:" 可使用内存数据库（测试用）。
func Init(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(
		&model.CardHolder{},
		&model.Window{},
		&model.Card{},
		&model.DepositRecord{},
		&model.Transaction{},
	)
	if err != nil {
		return nil, err
	}

	log.Info().Str("dsn", dsn).Msg("数据库初始化完成")
	return db, nil
}
