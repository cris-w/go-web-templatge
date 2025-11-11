package repo

import (
	"context"
	"power-supply-sys/internal/domain/user"
	"power-supply-sys/pkg/common"

	"gorm.io/gorm"
)

// UserRepository 用户数据访问层接口
type UserRepository interface {
	// 基础 CRUD 方法
	Create(ctx context.Context, u *user.User) error
	FindByID(ctx context.Context, id uint) (*user.User, error)
	FindOne(ctx context.Context, opts ...common.QueryOption) (*user.User, error)
	Update(ctx context.Context, u *user.User, updates map[string]any) error
	UpdateByID(ctx context.Context, id uint, updates map[string]any) error
	Delete(ctx context.Context, id uint) error
	Exists(ctx context.Context, opts ...common.QueryOption) (bool, error)

	// 自定义方法
	FindByUsername(ctx context.Context, username string) (*user.User, error)
	List(ctx context.Context, query *user.QueryOptions) ([]*user.User, error)
	Count(ctx context.Context, query *user.QueryOptions) (int64, error)
}

// userRepository 用户数据访问层实现
type userRepository struct {
	*common.BaseRepository[user.User]
}

// NewUserRepository 创建用户仓储
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		BaseRepository: common.NewBaseRepository[user.User](db),
	}
}

// FindByUsername 根据用户名查询用户
func (r *userRepository) FindByUsername(ctx context.Context, username string) (*user.User, error) {
	return r.FindOne(ctx, common.Where("username", username))
}

// Count 统计用户数量
func (r *userRepository) Count(ctx context.Context, query *user.QueryOptions) (int64, error) {
	if query == nil {
		return r.BaseRepository.Count(ctx)
	}

	return r.BaseRepository.Count(ctx,
		common.WhereLike("username", query.Username),
		common.WhereLike("email", query.Email),
		common.WhereIfNotNil("status", query.Status),
	)
}

// List 查询用户列表
func (r *userRepository) List(ctx context.Context, query *user.QueryOptions) ([]*user.User, error) {
	if query == nil {
		return r.BaseRepository.List(ctx, common.OrderByDesc("id"))
	}

	return r.BaseRepository.List(ctx,
		common.WhereLike("username", query.Username),
		common.WhereLike("email", query.Email),
		common.WhereIfNotNil("status", query.Status),
		common.OrderByDesc("id"),
		common.Paginate(query.Page, query.PageSize),
	)
}
