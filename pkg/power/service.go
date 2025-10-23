package power

import (
	"context"
	"power-supply-sys/pkg/common"

	"gorm.io/gorm"
)

// Service 电源服务接口
type Service interface {
	Create(ctx context.Context, req *PowerSupplyCreateRequest) (*PowerSupply, error)
	GetByID(ctx context.Context, id uint) (*PowerSupply, error)
	Update(ctx context.Context, id uint, req *PowerSupplyUpdateRequest) (*PowerSupply, error)
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, req *PowerSupplyQueryRequest) ([]*PowerSupply, int64, error)
}

// service 电源服务实现
type service struct {
	repo Repository
}

var _ Service = &service{}

// NewService 创建电源服务
func NewService(db *gorm.DB) (Service, error) {
	// 执行数据库迁移
	if err := AutoMigrate(db); err != nil {
		return nil, err
	}

	return &service{
		repo: NewRepository(db),
	}, nil
}

// Create 创建电源
func (s *service) Create(ctx context.Context, req *PowerSupplyCreateRequest) (*PowerSupply, error) {
	powerSupply := &PowerSupply{
		Name:        req.Name,
		Brand:       req.Brand,
		Model:       req.Model,
		Power:       req.Power,
		Efficiency:  req.Efficiency,
		Modular:     req.Modular,
		Price:       req.Price,
		Stock:       req.Stock,
		Description: req.Description,
		Status:      1,
	}

	if err := s.repo.Create(ctx, powerSupply); err != nil {
		return nil, err
	}

	return powerSupply, nil
}

// GetByID 根据ID获取电源
func (s *service) GetByID(ctx context.Context, id uint) (*PowerSupply, error) {
	return s.repo.FindByID(ctx, id)
}

// Update 更新电源
func (s *service) Update(ctx context.Context, id uint, req *PowerSupplyUpdateRequest) (*PowerSupply, error) {
	powerSupply, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Brand != "" {
		updates["brand"] = req.Brand
	}
	if req.Model != "" {
		updates["model"] = req.Model
	}
	if req.Power != nil {
		updates["power"] = *req.Power
	}
	if req.Efficiency != "" {
		updates["efficiency"] = req.Efficiency
	}
	if req.Modular != nil {
		updates["modular"] = *req.Modular
	}
	if req.Price != nil {
		updates["price"] = *req.Price
	}
	if req.Stock != nil {
		updates["stock"] = *req.Stock
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	if len(updates) == 0 {
		return powerSupply, nil
	}

	if err := s.repo.Update(ctx, powerSupply, updates); err != nil {
		return nil, err
	}

	// 重新查询更新后的电源信息
	return s.repo.FindByID(ctx, id)
}

// Delete 删除电源
func (s *service) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

// List 获取电源列表
func (s *service) List(ctx context.Context, req *PowerSupplyQueryRequest) ([]*PowerSupply, int64, error) {
	// 构建查询选项
	page, pageSize := common.GetPageInfo(req.Page, req.PageSize)

	queryOpts := &QueryOptions{
		Name:       req.Name,
		Brand:      req.Brand,
		MinPower:   req.MinPower,
		MaxPower:   req.MaxPower,
		MinPrice:   req.MinPrice,
		MaxPrice:   req.MaxPrice,
		Efficiency: req.Efficiency,
		Status:     req.Status,
		Page:       page,
		PageSize:   pageSize,
	}

	// 获取总数
	total, err := s.repo.Count(ctx, queryOpts)
	if err != nil {
		return nil, 0, err
	}

	// 获取列表
	powerSupplies, err := s.repo.List(ctx, queryOpts)
	if err != nil {
		return nil, 0, err
	}

	return powerSupplies, total, nil
}
