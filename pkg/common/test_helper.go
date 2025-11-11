package common

import (
    "testing"

    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

// SetupTestDB 创建测试用的内存数据库
func SetupTestDB(t *testing.T) *gorm.DB {
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
    })
    if err != nil {
        t.Fatalf("Failed to connect to test database: %v", err)
    }
    return db
}

// TeardownTestDB 清理测试数据库
func TeardownTestDB(t *testing.T, db *gorm.DB) {
    sqlDB, err := db.DB()
    if err != nil {
        t.Logf("Failed to get database instance: %v", err)
        return
    }
    if err := sqlDB.Close(); err != nil {
        t.Logf("Failed to close database: %v", err)
    }
}
