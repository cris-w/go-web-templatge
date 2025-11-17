package power

import (
	"context"
	"power-supply-sys/pkg/common"
)

// Reader 读取操作接口（接口隔离原则）
type Reader interface {
	FindByID(ctx context.Context, id uint) (*PowerSupply, error)
	FindOne(ctx context.Context, opts ...common.QueryOption) (*PowerSupply, error)
	List(ctx context.Context, query *QueryOptions) ([]*PowerSupply, error)
	Count(ctx context.Context, query *QueryOptions) (int64, error)
	Exists(ctx context.Context, opts ...common.QueryOption) (bool, error)
}

// Writer 写入操作接口（接口隔离原则）
type Writer interface {
	Create(ctx context.Context, ps *PowerSupply) error
	Update(ctx context.Context, ps *PowerSupply, updates map[string]any) error
	UpdateByID(ctx context.Context, id uint, updates map[string]any) error
	Delete(ctx context.Context, id uint) error
}

// Repository 电源仓储接口（组合 Reader 和 Writer）
type Repository interface {
	Reader
	Writer
}

