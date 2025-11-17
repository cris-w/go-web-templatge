package db

import (
	"power-supply-sys/internal/domain/power"
	"power-supply-sys/internal/domain/user"

	"gorm.io/gorm"
)

// Migrate 执行数据库迁移
func Migrate(db *gorm.DB) error {
	// 迁移用户表
	if err := db.AutoMigrate(&user.User{}); err != nil {
		return err
	}

	// 迁移电源表
	if err := db.AutoMigrate(&power.PowerSupply{}); err != nil {
		return err
	}

	return nil
}

