package db

import (
	"fmt"
	"power-supply-sys/pkg/logger"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

// Config 数据库配置
type Config struct {
	DSN             string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
	Debug           bool
}

// InitDatabase 初始化数据库连接并返回数据库实例
func InitDatabase(config *Config) (*gorm.DB, error) {
	var loggerConfig gormLogger.Interface

	if config.Debug {
		loggerConfig = gormLogger.Default.LogMode(gormLogger.Info)
	} else {
		loggerConfig = gormLogger.Default.LogMode(gormLogger.Error)
	}

	database, err := gorm.Open(mysql.Open(config.DSN), &gorm.Config{
		Logger: loggerConfig,
	})
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %w", err)
	}

	sqlDB, err := database.DB()
	if err != nil {
		return nil, fmt.Errorf("获取数据库实例失败: %w", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(config.ConnMaxLifetime)

	logger.Info("Database connected successfully",
		zap.Int("max_idle_conns", config.MaxIdleConns),
		zap.Int("max_open_conns", config.MaxOpenConns),
	)
	return database, nil
}

