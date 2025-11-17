package dto

// PowerSupplyCreateRequest 创建电源请求（DTO 移至传输层）
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

