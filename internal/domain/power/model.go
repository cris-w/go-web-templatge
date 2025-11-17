package power

import (
	"time"
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

