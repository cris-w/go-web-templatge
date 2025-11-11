package user

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Username  string    `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Password  string    `gorm:"size:255;not null" json:"-"`
	Email     string    `gorm:"uniqueIndex;size:100" json:"email"`
	Phone     string    `gorm:"size:20" json:"phone"`
	Nickname  string    `gorm:"size:50" json:"nickname"`
	Avatar    string    `gorm:"size:255" json:"avatar"`
	Status    int       `gorm:"default:1;comment:状态 1-正常 0-禁用" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// AutoMigrate 自动迁移用户表结构
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&User{})
}

