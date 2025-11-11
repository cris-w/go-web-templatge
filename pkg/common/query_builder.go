package common

import (
	"fmt"
	"reflect"

	"gorm.io/gorm"
)

// QueryOption 查询选项函数类型 - 函数式选项模式
type QueryOption func(*gorm.DB) *gorm.DB

// ApplyQuery 应用查询选项
func ApplyQuery(db *gorm.DB, opts ...QueryOption) *gorm.DB {
	for _, opt := range opts {
		db = opt(db)
	}
	return db
}

// ============ 基础查询条件 ============

// Where 等值查询
func Where(field string, value any) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("%s = ?", field), value)
	}
}

// WhereIn IN 查询
func WhereIn(field string, values any) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("%s IN ?", field), values)
	}
}

// WhereLike 模糊查询
func WhereLike(field, value string) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		if value == "" {
			return db
		}
		return db.Where(fmt.Sprintf("%s LIKE ?", field), "%"+value+"%")
	}
}

// WhereBetween 区间查询
func WhereBetween(field string, min, max any) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("%s BETWEEN ? AND ?", field), min, max)
	}
}

// WhereGTE 大于等于
func WhereGTE(field string, value any) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("%s >= ?", field), value)
	}
}

// WhereLTE 小于等于
func WhereLTE(field string, value any) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("%s <= ?", field), value)
	}
}

// WhereGT 大于
func WhereGT(field string, value any) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("%s > ?", field), value)
	}
}

// WhereLT 小于
func WhereLT(field string, value any) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("%s < ?", field), value)
	}
}

// WhereNotNull NOT NULL 条件
func WhereNotNull(field string) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("%s IS NOT NULL", field))
	}
}

// WhereNull IS NULL 条件
func WhereNull(field string) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("%s IS NULL", field))
	}
}

// ============ 条件查询(智能过滤) ============

// isNil 正确检查值是否为nil（包括接口中包含nil指针的情况）
func isNil(v any) bool {
	if v == nil {
		return true
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return rv.IsNil()
	}
	return false
}

// WhereIf 条件为真时才添加查询
func WhereIf(condition bool, field string, value any) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		if condition {
			return db.Where(fmt.Sprintf("%s = ?", field), value)
		}
		return db
	}
}

// WhereLikeIf 条件为真时才添加模糊查询
func WhereLikeIf(condition bool, field, value string) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		if condition && value != "" {
			return db.Where(fmt.Sprintf("%s LIKE ?", field), "%"+value+"%")
		}
		return db
	}
}

// WhereIfNotNil 值不为nil时才添加查询
func WhereIfNotNil(field string, value any) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		if !isNil(value) {
			return db.Where(fmt.Sprintf("%s = ?", field), value)
		}
		return db
	}
}

// WhereGTEIfNotNil 值不为nil时才添加大于等于查询
func WhereGTEIfNotNil(field string, value any) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		if !isNil(value) {
			return db.Where(fmt.Sprintf("%s >= ?", field), value)
		}
		return db
	}
}

// WhereLTEIfNotNil 值不为nil时才添加小于等于查询
func WhereLTEIfNotNil(field string, value any) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		if !isNil(value) {
			return db.Where(fmt.Sprintf("%s <= ?", field), value)
		}
		return db
	}
}

// ============ 排序 ============

// OrderBy 升序排序
func OrderBy(field string) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Order(field)
	}
}

// OrderByDesc 降序排序
func OrderByDesc(field string) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Order(fmt.Sprintf("%s DESC", field))
	}
}

// OrderByMulti 多字段排序
func OrderByMulti(fields ...string) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		for _, field := range fields {
			db = db.Order(field)
		}
		return db
	}
}

// ============ 分页 ============

// Paginate 分页
func Paginate(page, pageSize int) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 || pageSize <= 0 {
			return db
		}
		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

// Limit 限制数量
func Limit(limit int) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Limit(limit)
	}
}

// Offset 偏移量
func Offset(offset int) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(offset)
	}
}

// ============ 关联查询 ============

// Preload 预加载关联
func Preload(query string, args ...any) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Preload(query, args...)
	}
}

// Joins 连接查询
func Joins(query string, args ...any) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Joins(query, args...)
	}
}

// ============ 分组与聚合 ============

// GroupBy 分组
func GroupBy(field string) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Group(field)
	}
}

// Having HAVING 子句
func Having(query string, args ...any) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Having(query, args...)
	}
}

// ============ 其他 ============

// Distinct 去重
func Distinct(fields ...string) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		if len(fields) == 0 {
			return db.Distinct()
		}
		return db.Distinct(fields)
	}
}

// Select 选择字段
func Select(fields ...string) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Select(fields)
	}
}

// Omit 忽略字段
func Omit(fields ...string) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Omit(fields...)
	}
}

// ============ 组合查询选项 ============

// Combine 组合多个查询选项
func Combine(opts ...QueryOption) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		for _, opt := range opts {
			db = opt(db)
		}
		return db
	}
}
