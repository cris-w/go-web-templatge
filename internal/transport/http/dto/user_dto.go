package dto

import "power-supply-sys/internal/domain/user"

// UserCreateRequest 创建用户请求（DTO 移至传输层）
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
	Token string      `json:"token"`
	User  *user.User  `json:"user"`
}

