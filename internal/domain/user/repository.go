package user

import (
	"context"
	"power-supply-sys/pkg/common"
)

// Reader 读取操作接口（接口隔离原则）
type Reader interface {
	FindByID(ctx context.Context, id uint) (*User, error)
	FindOne(ctx context.Context, opts ...common.QueryOption) (*User, error)
	FindByUsername(ctx context.Context, username string) (*User, error)
	List(ctx context.Context, query *QueryOptions) ([]*User, error)
	Count(ctx context.Context, query *QueryOptions) (int64, error)
	Exists(ctx context.Context, opts ...common.QueryOption) (bool, error)
}

// Writer 写入操作接口（接口隔离原则）
type Writer interface {
	Create(ctx context.Context, u *User) error
	Update(ctx context.Context, u *User, updates map[string]any) error
	UpdateByID(ctx context.Context, id uint, updates map[string]any) error
	Delete(ctx context.Context, id uint) error
}

// Repository 用户仓储接口（组合 Reader 和 Writer）
type Repository interface {
	Reader
	Writer
}

