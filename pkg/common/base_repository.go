package common

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

// BaseRepository 通用仓储基类
type BaseRepository[T any] struct {
	db *gorm.DB
}

// NewBaseRepository 创建通用仓储
func NewBaseRepository[T any](db *gorm.DB) *BaseRepository[T] {
	return &BaseRepository[T]{db: db}
}

// Create 创建记录
func (r *BaseRepository[T]) Create(ctx context.Context, entity *T) error {
	if err := r.db.WithContext(ctx).Create(entity).Error; err != nil {
		return ErrDatabase(err)
	}
	return nil
}

// FindByID 根据ID查询记录
func (r *BaseRepository[T]) FindByID(ctx context.Context, id uint) (*T, error) {
	var entity T
	err := r.db.WithContext(ctx).First(&entity, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound("记录")
		}
		return nil, ErrDatabase(err)
	}
	return &entity, nil
}

// FindOne 根据条件查询单条记录
func (r *BaseRepository[T]) FindOne(ctx context.Context, opts ...QueryOption) (*T, error) {
	var entity T
	db := r.db.WithContext(ctx).Model(new(T))
	db = ApplyQuery(db, opts...)

	err := db.First(&entity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound("记录")
		}
		return nil, ErrDatabase(err)
	}
	return &entity, nil
}

// Update 更新记录
func (r *BaseRepository[T]) Update(ctx context.Context, entity *T, updates map[string]interface{}) error {
	if err := r.db.WithContext(ctx).Model(entity).Updates(updates).Error; err != nil {
		return ErrDatabase(err)
	}
	return nil
}

// UpdateByID 根据ID更新记录
func (r *BaseRepository[T]) UpdateByID(ctx context.Context, id uint, updates map[string]interface{}) error {
	result := r.db.WithContext(ctx).Model(new(T)).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return ErrDatabase(result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrNotFound("记录")
	}
	return nil
}

// Delete 删除记录
func (r *BaseRepository[T]) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(new(T), id)
	if result.Error != nil {
		return ErrDatabase(result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrNotFound("记录")
	}
	return nil
}

// DeleteByCondition 根据条件删除记录
func (r *BaseRepository[T]) DeleteByCondition(ctx context.Context, opts ...QueryOption) error {
	db := r.db.WithContext(ctx).Model(new(T))
	db = ApplyQuery(db, opts...)

	result := db.Delete(new(T))
	if result.Error != nil {
		return ErrDatabase(result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrNotFound("记录")
	}
	return nil
}

// List 查询记录列表
func (r *BaseRepository[T]) List(ctx context.Context, opts ...QueryOption) ([]*T, error) {
	var entities []*T
	db := r.db.WithContext(ctx).Model(new(T))
	db = ApplyQuery(db, opts...)

	if err := db.Find(&entities).Error; err != nil {
		return nil, ErrDatabase(err)
	}
	return entities, nil
}

// Count 统计记录数量
func (r *BaseRepository[T]) Count(ctx context.Context, opts ...QueryOption) (int64, error) {
	var count int64
	db := r.db.WithContext(ctx).Model(new(T))
	db = ApplyQuery(db, opts...)

	if err := db.Count(&count).Error; err != nil {
		return 0, ErrDatabase(err)
	}
	return count, nil
}

// Exists 检查记录是否存在
func (r *BaseRepository[T]) Exists(ctx context.Context, opts ...QueryOption) (bool, error) {
	count, err := r.Count(ctx, opts...)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// First 查询第一条记录
func (r *BaseRepository[T]) First(ctx context.Context, opts ...QueryOption) (*T, error) {
	var entity T
	db := r.db.WithContext(ctx).Model(new(T))
	db = ApplyQuery(db, opts...)

	err := db.First(&entity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound("记录")
		}
		return nil, ErrDatabase(err)
	}
	return &entity, nil
}

// BatchCreate 批量创建记录
func (r *BaseRepository[T]) BatchCreate(ctx context.Context, entities []*T) error {
	if len(entities) == 0 {
		return nil
	}
	if err := r.db.WithContext(ctx).Create(&entities).Error; err != nil {
		return ErrDatabase(err)
	}
	return nil
}

// BatchUpdate 批量更新记录
func (r *BaseRepository[T]) BatchUpdate(ctx context.Context, updates map[string]interface{}, opts ...QueryOption) (int64, error) {
	db := r.db.WithContext(ctx).Model(new(T))
	db = ApplyQuery(db, opts...)

	result := db.Updates(updates)
	if result.Error != nil {
		return 0, ErrDatabase(result.Error)
	}
	return result.RowsAffected, nil
}

// Transaction 执行事务
func (r *BaseRepository[T]) Transaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return r.db.WithContext(ctx).Transaction(fn)
}

// GetDB 获取数据库连接(用于复杂查询)
func (r *BaseRepository[T]) GetDB(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx)
}
