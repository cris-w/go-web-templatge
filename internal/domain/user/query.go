package user

// QueryOptions 查询选项（保留在领域层，属于领域概念）
type QueryOptions struct {
	Username string
	Email    string
	Status   *int
	Page     int
	PageSize int
}

