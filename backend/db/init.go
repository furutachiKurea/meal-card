// Package db 负责数据库连接与初始化
package db

import (
	"backend/model"
	"context"
	"errors"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	gormlogger "gorm.io/gorm/logger"

	"gorm.io/gorm"
)

// zerologGormLogger 实现 gorm.io/gorm/logger.Interface，将 SQL 日志桥接到 zerolog。
type zerologGormLogger struct {
	// slowThreshold 超过此阈值的查询视为慢查询，输出 Warn 级日志
	slowThreshold time.Duration
}

// newZerologGormLogger 创建默认配置的 GORM zerolog logger。
func newZerologGormLogger() *zerologGormLogger {
	return &zerologGormLogger{
		slowThreshold: 200 * time.Millisecond,
	}
}

// LogMode 实现 logger.Interface，返回自身（日志级别由 zerolog 全局配置控制）。
func (l *zerologGormLogger) LogMode(_ gormlogger.LogLevel) gormlogger.Interface {
	return l
}

// Info 输出 Info 级别的 GORM 内部信息。
func (l *zerologGormLogger) Info(_ context.Context, msg string, args ...any) {
	log.Info().Msgf(msg, args...)
}

// Warn 输出 Warn 级别的 GORM 内部警告。
func (l *zerologGormLogger) Warn(_ context.Context, msg string, args ...any) {
	log.Warn().Msgf(msg, args...)
}

// Error 输出 Error 级别的 GORM 内部错误。
func (l *zerologGormLogger) Error(_ context.Context, msg string, args ...any) {
	log.Error().Msgf(msg, args...)
}

// Trace 记录每条 SQL 执行情况：
// - ErrRecordNotFound 降为 Debug，不视为错误；
// - 其他 SQL 错误输出 Error；
// - 慢查询（超过 slowThreshold）输出 Warn；
// - 正常查询输出 Debug。
func (l *zerologGormLogger) Trace(_ context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()

	// 构造基础 event
	baseEvent := func(level zerolog.Level) *zerolog.Event {
		return log.WithLevel(level).
			Str("sql", sql).
			Int64("rows", rows).
			Dur("elapsed", elapsed)
	}

	switch {
	case err != nil && !errors.Is(err, gorm.ErrRecordNotFound):
		baseEvent(zerolog.ErrorLevel).Err(err).Msg("gorm sql error")
	case elapsed >= l.slowThreshold:
		baseEvent(zerolog.WarnLevel).Msg("gorm slow query")
	default:
		baseEvent(zerolog.DebugLevel).Msg("gorm sql")
	}
}

// Init 初始化 SQLite 数据库连接并自动迁移表结构。
// dsn 为数据库文件路径，传入 ":memory:" 可使用内存数据库（测试用）。
func Init(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: newZerologGormLogger(),
	})
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
