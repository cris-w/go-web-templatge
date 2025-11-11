package repo

import (
	"context"
	"power-supply-sys/internal/domain/power"
	"power-supply-sys/pkg/common"

	"gorm.io/gorm"
)

// PowerRepository 电源数据访问层接口
type PowerRepository interface {
	// 基础 CRUD 方法
	Create(ctx context.Context, ps *power.PowerSupply) error
	FindByID(ctx context.Context, id uint) (*power.PowerSupply, error)
	FindOne(ctx context.Context, opts ...common.QueryOption) (*power.PowerSupply, error)
	Update(ctx context.Context, ps *power.PowerSupply, updates map[string]any) error
	UpdateByID(ctx context.Context, id uint, updates map[string]any) error
	Delete(ctx context.Context, id uint) error
	Exists(ctx context.Context, opts ...common.QueryOption) (bool, error)

	// 自定义方法
	List(ctx context.Context, query *power.QueryOptions) ([]*power.PowerSupply, error)
	Count(ctx context.Context, query *power.QueryOptions) (int64, error)
}

// powerRepository 电源数据访问层实现
type powerRepository struct {
	*common.BaseRepository[power.PowerSupply]
}

// NewPowerRepository 创建电源仓储
func NewPowerRepository(db *gorm.DB) PowerRepository {
	return &powerRepository{
		BaseRepository: common.NewBaseRepository[power.PowerSupply](db),
	}
}

// Count 统计电源数量
func (r *powerRepository) Count(ctx context.Context, query *power.QueryOptions) (int64, error) {
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
func (r *powerRepository) List(ctx context.Context, query *power.QueryOptions) ([]*power.PowerSupply, error) {
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
