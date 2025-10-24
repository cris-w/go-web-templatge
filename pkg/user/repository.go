package user

import (
	"context"
	"power-supply-sys/pkg/common"

	"gorm.io/gorm"
)

// Repository 用户数据访问层接口
type Repository interface {
	// 基础 CRUD 方法
	Create(ctx context.Context, user *User) error
	FindByID(ctx context.Context, id uint) (*User, error)
	FindOne(ctx context.Context, opts ...common.QueryOption) (*User, error)
	Update(ctx context.Context, user *User, updates map[string]interface{}) error
	UpdateByID(ctx context.Context, id uint, updates map[string]interface{}) error
	Delete(ctx context.Context, id uint) error
	Exists(ctx context.Context, opts ...common.QueryOption) (bool, error)

	// 自定义方法
	FindByUsername(ctx context.Context, username string) (*User, error)
	List(ctx context.Context, query *QueryOptions) ([]*User, error)
	Count(ctx context.Context, query *QueryOptions) (int64, error)
}

// QueryOptions 查询选项
type QueryOptions struct {
	Username string
	Email    string
	Status   *int
	Page     int
	PageSize int
}

// repository 用户数据访问层实现
type repository struct {
	*common.BaseRepository[User]
}

// NewRepository 创建用户仓储
func NewRepository(db *gorm.DB) Repository {
	return &repository{
		BaseRepository: common.NewBaseRepository[User](db),
	}
}

// FindByUsername 根据用户名查询用户
func (r *repository) FindByUsername(ctx context.Context, username string) (*User, error) {
	return r.FindOne(ctx, common.Where("username", username))
}

// Count 统计用户数量
func (r *repository) Count(ctx context.Context, query *QueryOptions) (int64, error) {
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
func (r *repository) List(ctx context.Context, query *QueryOptions) ([]*User, error) {
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
