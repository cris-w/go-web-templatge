package user

import (
	"context"
	"errors"
	"power-supply-sys/pkg/common"

	"gorm.io/gorm"
)

// Repository 用户数据访问层接口
type Repository interface {
	// Create 创建用户
	Create(ctx context.Context, user *User) error
	// FindByID 根据ID查询用户
	FindByID(ctx context.Context, id uint) (*User, error)
	// FindByUsername 根据用户名查询用户
	FindByUsername(ctx context.Context, username string) (*User, error)
	// Update 更新用户
	Update(ctx context.Context, user *User, updates map[string]interface{}) error
	// Delete 删除用户
	Delete(ctx context.Context, id uint) error
	// Count 统计用户数量（根据条件）
	Count(ctx context.Context, query *QueryOptions) (int64, error)
	// List 查询用户列表
	List(ctx context.Context, query *QueryOptions) ([]*User, error)
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
	db *gorm.DB
}

// NewRepository 创建用户仓储
func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

// Create 创建用户
func (r *repository) Create(ctx context.Context, user *User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return common.ErrDatabase(err)
	}
	return nil
}

// FindByID 根据ID查询用户
func (r *repository) FindByID(ctx context.Context, id uint) (*User, error) {
	var user User
	err := r.db.WithContext(ctx).First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrNotFound("用户")
		}
		return nil, common.ErrDatabase(err)
	}
	return &user, nil
}

// FindByUsername 根据用户名查询用户
func (r *repository) FindByUsername(ctx context.Context, username string) (*User, error) {
	var user User
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrNotFound("用户")
		}
		return nil, common.ErrDatabase(err)
	}
	return &user, nil
}

// Update 更新用户
func (r *repository) Update(ctx context.Context, user *User, updates map[string]interface{}) error {
	if err := r.db.WithContext(ctx).Model(user).Updates(updates).Error; err != nil {
		return common.ErrDatabase(err)
	}
	return nil
}

// Delete 删除用户
func (r *repository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&User{}, id)
	if result.Error != nil {
		return common.ErrDatabase(result.Error)
	}
	if result.RowsAffected == 0 {
		return common.ErrNotFound("用户")
	}
	return nil
}

// Count 统计用户数量
func (r *repository) Count(ctx context.Context, query *QueryOptions) (int64, error) {
	var count int64
	db := r.db.WithContext(ctx).Model(&User{})

	// 应用查询条件
	db = r.applyQueryConditions(db, query)

	if err := db.Count(&count).Error; err != nil {
		return 0, common.ErrDatabase(err)
	}
	return count, nil
}

// List 查询用户列表
func (r *repository) List(ctx context.Context, query *QueryOptions) ([]*User, error) {
	var users []*User
	db := r.db.WithContext(ctx).Model(&User{})

	// 应用查询条件
	db = r.applyQueryConditions(db, query)

	// 分页
	if query.Page > 0 && query.PageSize > 0 {
		offset := (query.Page - 1) * query.PageSize
		db = db.Offset(offset).Limit(query.PageSize)
	}

	// 排序
	db = db.Order("id DESC")

	if err := db.Find(&users).Error; err != nil {
		return nil, common.ErrDatabase(err)
	}

	return users, nil
}

// applyQueryConditions 应用查询条件
func (r *repository) applyQueryConditions(db *gorm.DB, query *QueryOptions) *gorm.DB {
	if query == nil {
		return db
	}

	if query.Username != "" {
		db = db.Where("username LIKE ?", "%"+query.Username+"%")
	}
	if query.Email != "" {
		db = db.Where("email LIKE ?", "%"+query.Email+"%")
	}
	if query.Status != nil {
		db = db.Where("status = ?", *query.Status)
	}

	return db
}
