package power

// Service 层使用的请求类型（从 DTO 转换而来）
// 这些类型用于 Service 层接口，保持 Service 层与传输层解耦

// PowerSupplyCreateRequest Service 层创建电源请求
type PowerSupplyCreateRequest struct {
	Name        string
	Brand       string
	Model       string
	Power       int
	Efficiency  string
	Modular     bool
	Price       float64
	Stock       int
	Description string
}

// PowerSupplyUpdateRequest Service 层更新电源请求
type PowerSupplyUpdateRequest struct {
	Name        string
	Brand       string
	Model       string
	Power       *int
	Efficiency  string
	Modular     *bool
	Price       *float64
	Stock       *int
	Description string
	Status      *int
}

// PowerSupplyQueryRequest Service 层查询电源请求
type PowerSupplyQueryRequest struct {
	Page       int
	PageSize   int
	Name       string
	Brand      string
	MinPower   *int
	MaxPower   *int
	MinPrice   *float64
	MaxPrice   *float64
	Efficiency string
	Status     *int
}

