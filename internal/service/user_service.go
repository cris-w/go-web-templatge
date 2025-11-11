package service

import (
	"context"
	"power-supply-sys/internal/domain/user"
	"power-supply-sys/internal/infra/repo"
	"power-supply-sys/pkg/common"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserService 用户服务接口
type UserService interface {
	Create(ctx context.Context, req *user.UserCreateRequest) (*user.User, error)
	GetByID(ctx context.Context, id uint) (*user.User, error)
	GetByUsername(ctx context.Context, username string) (*user.User, error)
	Update(ctx context.Context, id uint, req *user.UserUpdateRequest) (*user.User, error)
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, req *user.UserQueryRequest) ([]*user.User, int64, error)
	Login(ctx context.Context, req *user.LoginRequest) (*user.User, error)
	VerifyPassword(hashedPassword, password string) error
}

// userService 用户服务实现
type userService struct {
	repo repo.UserRepository
}

var _ UserService = &userService{}

// NewUserService 创建用户服务
func NewUserService(db *gorm.DB) UserService {
	return &userService{
		repo: repo.NewUserRepository(db),
	}
}

// Create 创建用户
func (s *userService) Create(ctx context.Context, req *user.UserCreateRequest) (*user.User, error) {
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

	u := &user.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Email:    req.Email,
		Phone:    req.Phone,
		Nickname: req.Nickname,
		Avatar:   req.Avatar,
		Status:   1,
	}

	if err := s.repo.Create(ctx, u); err != nil {
		return nil, err
	}

	return u, nil
}

// GetByID 根据ID获取用户
func (s *userService) GetByID(ctx context.Context, id uint) (*user.User, error) {
	return s.repo.FindByID(ctx, id)
}

// GetByUsername 根据用户名获取用户
func (s *userService) GetByUsername(ctx context.Context, username string) (*user.User, error) {
	return s.repo.FindByUsername(ctx, username)
}

// Update 更新用户
func (s *userService) Update(ctx context.Context, id uint, req *user.UserUpdateRequest) (*user.User, error) {
	u, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	updates := make(map[string]any)
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
		return u, nil
	}

	if err := s.repo.Update(ctx, u, updates); err != nil {
		return nil, err
	}

	// 重新查询更新后的用户信息
	return s.repo.FindByID(ctx, id)
}

// Delete 删除用户
func (s *userService) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

// List 获取用户列表
func (s *userService) List(ctx context.Context, req *user.UserQueryRequest) ([]*user.User, int64, error) {
	// 构建查询选项
	page, pageSize := common.GetPageInfo(req.Page, req.PageSize)

	queryOpts := &user.QueryOptions{
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
func (s *userService) Login(ctx context.Context, req *user.LoginRequest) (*user.User, error) {
	// 根据用户名查询用户
	u, err := s.repo.FindByUsername(ctx, req.Username)
	if err != nil {
		if common.IsAppError(err) {
			return nil, common.ErrUnauthorized("用户名或密码错误")
		}
		return nil, err
	}

	// 验证密码
	if err := s.VerifyPassword(u.Password, req.Password); err != nil {
		return nil, common.ErrUnauthorized("用户名或密码错误")
	}

	// 检查用户状态
	if u.Status != 1 {
		return nil, common.ErrForbidden("用户已被禁用")
	}

	return u, nil
}

// VerifyPassword 验证密码
func (s *userService) VerifyPassword(hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return err
	}
	return nil
}
