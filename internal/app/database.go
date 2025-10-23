package app

import (
	"fmt"
	"power-supply-sys/pkg/logger"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

var db *gorm.DB

// InitDatabase 初始化数据库连接并返回数据库实例
func InitDatabase(config *Config) (*gorm.DB, error) {
	var err error
	var loggerConfig gormLogger.Interface

	if config.Debug {
		loggerConfig = gormLogger.Default.LogMode(gormLogger.Info)
	} else {
		loggerConfig = gormLogger.Default.LogMode(gormLogger.Error)
	}

	db, err = gorm.Open(mysql.Open(config.DB.DSN), &gorm.Config{
		Logger: loggerConfig,
	})
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("获取数据库实例失败: %w", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(config.DB.MaxIdleConns)
	sqlDB.SetMaxOpenConns(config.DB.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(config.DB.GetConnMaxLifetime())

	logger.Info("Database connected successfully",
		zap.Int("max_idle_conns", config.DB.MaxIdleConns),
		zap.Int("max_open_conns", config.DB.MaxOpenConns),
	)
	return db, nil
}
