package power

import (
	"time"

	"gorm.io/gorm"
)

// PowerSupply 电源模型
type PowerSupply struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	Name        string    `gorm:"size:100;not null" json:"name"`
	Brand       string    `gorm:"size:50" json:"brand"`
	Model       string    `gorm:"size:50" json:"model"`
	Power       int       `gorm:"comment:功率(W)" json:"power"`
	Efficiency  string    `gorm:"size:20;comment:能效等级" json:"efficiency"`
	Modular     bool      `gorm:"comment:是否模组化" json:"modular"`
	Price       float64   `gorm:"type:decimal(10,2)" json:"price"`
	Stock       int       `gorm:"default:0;comment:库存数量" json:"stock"`
	Description string    `gorm:"type:text" json:"description"`
	Status      int       `gorm:"default:1;comment:状态 1-上架 0-下架" json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TableName 指定表名
func (PowerSupply) TableName() string {
	return "power_supplies"
}

// AutoMigrate 自动迁移电源表结构
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&PowerSupply{})
}

// PowerSupplyCreateRequest 创建电源请求
type PowerSupplyCreateRequest struct {
	Name        string  `json:"name" binding:"required,min=1,max=100"`
	Brand       string  `json:"brand" binding:"omitempty,max=50"`
	Model       string  `json:"model" binding:"omitempty,max=50"`
	Power       int     `json:"power" binding:"required,min=0"`
	Efficiency  string  `json:"efficiency" binding:"omitempty,max=20"`
	Modular     bool    `json:"modular"`
	Price       float64 `json:"price" binding:"required,min=0"`
	Stock       int     `json:"stock" binding:"omitempty,min=0"`
	Description string  `json:"description" binding:"omitempty"`
}

// PowerSupplyUpdateRequest 更新电源请求
type PowerSupplyUpdateRequest struct {
	Name        string   `json:"name" binding:"omitempty,min=1,max=100"`
	Brand       string   `json:"brand" binding:"omitempty,max=50"`
	Model       string   `json:"model" binding:"omitempty,max=50"`
	Power       *int     `json:"power" binding:"omitempty,min=0"`
	Efficiency  string   `json:"efficiency" binding:"omitempty,max=20"`
	Modular     *bool    `json:"modular"`
	Price       *float64 `json:"price" binding:"omitempty,min=0"`
	Stock       *int     `json:"stock" binding:"omitempty,min=0"`
	Description string   `json:"description" binding:"omitempty"`
	Status      *int     `json:"status" binding:"omitempty,oneof=0 1"`
}

// PowerSupplyQueryRequest 查询电源请求
type PowerSupplyQueryRequest struct {
	Page       int      `form:"page" binding:"omitempty,min=1"`
	PageSize   int      `form:"page_size" binding:"omitempty,min=1,max=100"`
	Name       string   `form:"name" binding:"omitempty"`
	Brand      string   `form:"brand" binding:"omitempty"`
	MinPower   *int     `form:"min_power" binding:"omitempty,min=0"`
	MaxPower   *int     `form:"max_power" binding:"omitempty,min=0"`
	MinPrice   *float64 `form:"min_price" binding:"omitempty,min=0"`
	MaxPrice   *float64 `form:"max_price" binding:"omitempty,min=0"`
	Efficiency string   `form:"efficiency" binding:"omitempty"`
	Status     *int     `form:"status" binding:"omitempty,oneof=0 1"`
}
