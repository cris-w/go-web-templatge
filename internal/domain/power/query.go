package power

// QueryOptions 查询选项（保留在领域层，属于领域概念）
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

