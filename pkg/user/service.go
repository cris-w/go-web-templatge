package user

import (
	"context"
	"power-supply-sys/pkg/common"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Service 用户服务接口
type Service interface {
	Create(ctx context.Context, req *UserCreateRequest) (*User, error)
	GetByID(ctx context.Context, id uint) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	Update(ctx context.Context, id uint, req *UserUpdateRequest) (*User, error)
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, req *UserQueryRequest) ([]*User, int64, error)
	Login(ctx context.Context, req *LoginRequest) (*User, error)
	VerifyPassword(hashedPassword, password string) error
}

// service 用户服务实现
type service struct {
	repo Repository
}

var _ Service = &service{}

// NewService 创建用户服务
func NewService(db *gorm.DB) (Service, error) {
	// 执行数据库迁移
	if err := AutoMigrate(db); err != nil {
		return nil, err
	}

	return &service{
		repo: NewRepository(db),
	}, nil
}

// Create 创建用户
func (s *service) Create(ctx context.Context, req *UserCreateRequest) (*User, error) {
	// 检查用户名是否已存在
	existUser, err := s.repo.FindByUsername(ctx, req.Username)
	if err != nil && !common.IsAppError(err) {
		return nil, err
	}
	if existUser != nil {
		return nil, common.ErrAlreadyExists("用户名")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, common.ErrInternal(err)
	}

	user := &User{
		Username: req.Username,
		Password: string(hashedPassword),
		Email:    req.Email,
		Phone:    req.Phone,
		Nickname: req.Nickname,
		Avatar:   req.Avatar,
		Status:   1,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// GetByID 根据ID获取用户
func (s *service) GetByID(ctx context.Context, id uint) (*User, error) {
	return s.repo.FindByID(ctx, id)
}

// GetByUsername 根据用户名获取用户
func (s *service) GetByUsername(ctx context.Context, username string) (*User, error) {
	return s.repo.FindByUsername(ctx, username)
}

// Update 更新用户
func (s *service) Update(ctx context.Context, id uint, req *UserUpdateRequest) (*User, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	updates := make(map[string]interface{})
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}
	if req.Nickname != "" {
		updates["nickname"] = req.Nickname
	}
	if req.Avatar != "" {
		updates["avatar"] = req.Avatar
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	if len(updates) == 0 {
		return user, nil
	}

	if err := s.repo.Update(ctx, user, updates); err != nil {
		return nil, err
	}

	// 重新查询更新后的用户信息
	return s.repo.FindByID(ctx, id)
}

// Delete 删除用户
func (s *service) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

// List 获取用户列表
func (s *service) List(ctx context.Context, req *UserQueryRequest) ([]*User, int64, error) {
	// 构建查询选项
	page, pageSize := common.GetPageInfo(req.Page, req.PageSize)

	queryOpts := &QueryOptions{
		Username: req.Username,
		Email:    req.Email,
		Status:   req.Status,
		Page:     page,
		PageSize: pageSize,
	}

	// 获取总数
	total, err := s.repo.Count(ctx, queryOpts)
	if err != nil {
		return nil, 0, err
	}

	// 获取列表
	users, err := s.repo.List(ctx, queryOpts)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// Login 用户登录
func (s *service) Login(ctx context.Context, req *LoginRequest) (*User, error) {
	// 根据用户名查询用户
	user, err := s.repo.FindByUsername(ctx, req.Username)
	if err != nil {
		if common.IsAppError(err) {
			return nil, common.ErrUnauthorized("用户名或密码错误")
		}
		return nil, err
	}

	// 验证密码
	if err := s.VerifyPassword(user.Password, req.Password); err != nil {
		return nil, common.ErrUnauthorized("用户名或密码错误")
	}

	// 检查用户状态
	if user.Status != 1 {
		return nil, common.ErrForbidden("用户已被禁用")
	}

	return user, nil
}

// VerifyPassword 验证密码
func (s *service) VerifyPassword(hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return err
	}
	return nil
}
