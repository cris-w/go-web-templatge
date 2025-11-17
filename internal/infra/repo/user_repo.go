package repo

import (
	"context"
	"power-supply-sys/internal/domain/user"
	"power-supply-sys/pkg/common"

	"gorm.io/gorm"
)

// userRepository 用户数据访问层实现（实现 domain 层的 Repository 接口）
type userRepository struct {
	*common.BaseRepository[user.User]
}

// NewUserRepository 创建用户仓储
func NewUserRepository(db *gorm.DB) user.Repository {
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
