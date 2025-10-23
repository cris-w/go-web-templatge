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

// UserCreateRequest 创建用户请求
type UserCreateRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6"`
	Email    string `json:"email" binding:"omitempty,email"`
	Phone    string `json:"phone" binding:"omitempty"`
	Nickname string `json:"nickname" binding:"omitempty,max=50"`
	Avatar   string `json:"avatar" binding:"omitempty"`
}

// UserUpdateRequest 更新用户请求
type UserUpdateRequest struct {
	Email    string `json:"email" binding:"omitempty,email"`
	Phone    string `json:"phone" binding:"omitempty"`
	Nickname string `json:"nickname" binding:"omitempty,max=50"`
	Avatar   string `json:"avatar" binding:"omitempty"`
	Status   *int   `json:"status" binding:"omitempty,oneof=0 1"`
}

// UserQueryRequest 查询用户请求
type UserQueryRequest struct {
	Page     int    `form:"page" binding:"omitempty,min=1"`
	PageSize int    `form:"page_size" binding:"omitempty,min=1,max=100"`
	Username string `form:"username" binding:"omitempty"`
	Email    string `form:"email" binding:"omitempty"`
	Status   *int   `form:"status" binding:"omitempty,oneof=0 1"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required,min=3"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}
