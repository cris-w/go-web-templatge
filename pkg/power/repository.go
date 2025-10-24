package power

import (
	"context"
	"power-supply-sys/pkg/common"

	"gorm.io/gorm"
)

// Repository 电源数据访问层接口
type Repository interface {
	// 基础 CRUD 方法
	Create(ctx context.Context, powerSupply *PowerSupply) error
	FindByID(ctx context.Context, id uint) (*PowerSupply, error)
	FindOne(ctx context.Context, opts ...common.QueryOption) (*PowerSupply, error)
	Update(ctx context.Context, powerSupply *PowerSupply, updates map[string]interface{}) error
	UpdateByID(ctx context.Context, id uint, updates map[string]interface{}) error
	Delete(ctx context.Context, id uint) error
	Exists(ctx context.Context, opts ...common.QueryOption) (bool, error)

	// 自定义方法
	List(ctx context.Context, query *QueryOptions) ([]*PowerSupply, error)
	Count(ctx context.Context, query *QueryOptions) (int64, error)
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
	*common.BaseRepository[PowerSupply]
}

// NewRepository 创建电源仓储
func NewRepository(db *gorm.DB) Repository {
	return &repository{
		BaseRepository: common.NewBaseRepository[PowerSupply](db),
	}
}

// Count 统计电源数量
func (r *repository) Count(ctx context.Context, query *QueryOptions) (int64, error) {
	if query == nil {
		return r.BaseRepository.Count(ctx)
	}

	return r.BaseRepository.Count(ctx,
		common.WhereLike("name", query.Name),
		common.WhereLike("brand", query.Brand),
		common.WhereGTEIfNotNil("power", query.MinPower),
		common.WhereLTEIfNotNil("power", query.MaxPower),
		common.WhereGTEIfNotNil("price", query.MinPrice),
		common.WhereLTEIfNotNil("price", query.MaxPrice),
		common.WhereIf(query.Efficiency != "", "efficiency", query.Efficiency),
		common.WhereIfNotNil("status", query.Status),
	)
}

// List 查询电源列表
func (r *repository) List(ctx context.Context, query *QueryOptions) ([]*PowerSupply, error) {
	if query == nil {
		return r.BaseRepository.List(ctx, common.OrderByDesc("id"))
	}

	return r.BaseRepository.List(ctx,
		common.WhereLike("name", query.Name),
		common.WhereLike("brand", query.Brand),
		common.WhereGTEIfNotNil("power", query.MinPower),
		common.WhereLTEIfNotNil("power", query.MaxPower),
		common.WhereGTEIfNotNil("price", query.MinPrice),
		common.WhereLTEIfNotNil("price", query.MaxPrice),
		common.WhereIf(query.Efficiency != "", "efficiency", query.Efficiency),
		common.WhereIfNotNil("status", query.Status),
		common.OrderByDesc("id"),
		common.Paginate(query.Page, query.PageSize),
	)
}
