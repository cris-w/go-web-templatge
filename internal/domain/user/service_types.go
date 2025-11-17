package user

// Service 层使用的请求类型（从 DTO 转换而来）
// 这些类型用于 Service 层接口，保持 Service 层与传输层解耦

// UserCreateRequest Service 层创建用户请求
type UserCreateRequest struct {
	Username string
	Password string
	Email    string
	Phone    string
	Nickname string
	Avatar   string
}

// UserUpdateRequest Service 层更新用户请求
type UserUpdateRequest struct {
	Email    string
	Phone    string
	Nickname string
	Avatar   string
	Status   *int
}

// UserQueryRequest Service 层查询用户请求
type UserQueryRequest struct {
	Page     int
	PageSize int
	Username string
	Email    string
	Status   *int
}

// LoginRequest Service 层登录请求
type LoginRequest struct {
	Username string
	Password string
}

