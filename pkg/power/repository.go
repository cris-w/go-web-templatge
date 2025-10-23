package power

import (
	"context"
	"errors"
	"power-supply-sys/pkg/common"

	"gorm.io/gorm"
)

// Repository 电源数据访问层接口
type Repository interface {
	// Create 创建电源
	Create(ctx context.Context, powerSupply *PowerSupply) error
	// FindByID 根据ID查询电源
	FindByID(ctx context.Context, id uint) (*PowerSupply, error)
	// Update 更新电源
	Update(ctx context.Context, powerSupply *PowerSupply, updates map[string]interface{}) error
	// Delete 删除电源
	Delete(ctx context.Context, id uint) error
	// Count 统计电源数量（根据条件）
	Count(ctx context.Context, query *QueryOptions) (int64, error)
	// List 查询电源列表
	List(ctx context.Context, query *QueryOptions) ([]*PowerSupply, error)
}

// QueryOptions 查询选项
type QueryOptions struct {
	Name       string
	Brand      string
	MinPower   *int
	MaxPower   *int
	MinPrice   *float64
	MaxPrice   *float64
	Efficiency string
	Status     *int
	Page       int
	PageSize   int
}

// repository 电源数据访问层实现
type repository struct {
	db *gorm.DB
}

// NewRepository 创建电源仓储
func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

// Create 创建电源
func (r *repository) Create(ctx context.Context, powerSupply *PowerSupply) error {
	if err := r.db.WithContext(ctx).Create(powerSupply).Error; err != nil {
		return common.ErrDatabase(err)
	}
	return nil
}

// FindByID 根据ID查询电源
func (r *repository) FindByID(ctx context.Context, id uint) (*PowerSupply, error) {
	var powerSupply PowerSupply
	err := r.db.WithContext(ctx).First(&powerSupply, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrNotFound("电源")
		}
		return nil, common.ErrDatabase(err)
	}
	return &powerSupply, nil
}

// Update 更新电源
func (r *repository) Update(ctx context.Context, powerSupply *PowerSupply, updates map[string]interface{}) error {
	if err := r.db.WithContext(ctx).Model(powerSupply).Updates(updates).Error; err != nil {
		return common.ErrDatabase(err)
	}
	return nil
}

// Delete 删除电源
func (r *repository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&PowerSupply{}, id)
	if result.Error != nil {
		return common.ErrDatabase(result.Error)
	}
	if result.RowsAffected == 0 {
		return common.ErrNotFound("电源")
	}
	return nil
}

// Count 统计电源数量
func (r *repository) Count(ctx context.Context, query *QueryOptions) (int64, error) {
	var count int64
	db := r.db.WithContext(ctx).Model(&PowerSupply{})

	// 应用查询条件
	db = r.applyQueryConditions(db, query)

	if err := db.Count(&count).Error; err != nil {
		return 0, common.ErrDatabase(err)
	}
	return count, nil
}

// List 查询电源列表
func (r *repository) List(ctx context.Context, query *QueryOptions) ([]*PowerSupply, error) {
	var powerSupplies []*PowerSupply
	db := r.db.WithContext(ctx).Model(&PowerSupply{})

	// 应用查询条件
	db = r.applyQueryConditions(db, query)

	// 分页
	if query.Page > 0 && query.PageSize > 0 {
		offset := (query.Page - 1) * query.PageSize
		db = db.Offset(offset).Limit(query.PageSize)
	}

	// 排序
	db = db.Order("id DESC")

	if err := db.Find(&powerSupplies).Error; err != nil {
		return nil, common.ErrDatabase(err)
	}

	return powerSupplies, nil
}

// applyQueryConditions 应用查询条件
func (r *repository) applyQueryConditions(db *gorm.DB, query *QueryOptions) *gorm.DB {
	if query == nil {
		return db
	}

	if query.Name != "" {
		db = db.Where("name LIKE ?", "%"+query.Name+"%")
	}
	if query.Brand != "" {
		db = db.Where("brand LIKE ?", "%"+query.Brand+"%")
	}
	if query.MinPower != nil {
		db = db.Where("power >= ?", *query.MinPower)
	}
	if query.MaxPower != nil {
		db = db.Where("power <= ?", *query.MaxPower)
	}
	if query.MinPrice != nil {
		db = db.Where("price >= ?", *query.MinPrice)
	}
	if query.MaxPrice != nil {
		db = db.Where("price <= ?", *query.MaxPrice)
	}
	if query.Efficiency != "" {
		db = db.Where("efficiency = ?", query.Efficiency)
	}
	if query.Status != nil {
		db = db.Where("status = ?", *query.Status)
	}

	return db
}
