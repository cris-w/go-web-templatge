package service

import (
	"context"
	"power-supply-sys/internal/domain/power"
	"power-supply-sys/pkg/common"
)

// PowerService 电源服务接口
type PowerService interface {
	Create(ctx context.Context, req *power.PowerSupplyCreateRequest) (*power.PowerSupply, error)
	GetByID(ctx context.Context, id uint) (*power.PowerSupply, error)
	Update(ctx context.Context, id uint, req *power.PowerSupplyUpdateRequest) (*power.PowerSupply, error)
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, req *power.PowerSupplyQueryRequest) ([]*power.PowerSupply, int64, error)
}

// powerService 电源服务实现
type powerService struct {
	repo power.Repository
}

var _ PowerService = &powerService{}

// NewPowerService 创建电源服务（接收 Repository 接口而非 GORM）
func NewPowerService(repo power.Repository) PowerService {
	return &powerService{
		repo: repo,
	}
}

// Create 创建电源
func (s *powerService) Create(ctx context.Context, req *power.PowerSupplyCreateRequest) (*power.PowerSupply, error) {
	ps := &power.PowerSupply{
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

	if err := s.repo.Create(ctx, ps); err != nil {
		return nil, err
	}

	return ps, nil
}

// GetByID 根据ID获取电源
func (s *powerService) GetByID(ctx context.Context, id uint) (*power.PowerSupply, error) {
	return s.repo.FindByID(ctx, id)
}

// Update 更新电源
func (s *powerService) Update(ctx context.Context, id uint, req *power.PowerSupplyUpdateRequest) (*power.PowerSupply, error) {
	ps, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	updates := make(map[string]any)
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
		return ps, nil
	}

	if err := s.repo.Update(ctx, ps, updates); err != nil {
		return nil, err
	}

	// 重新查询更新后的电源信息
	return s.repo.FindByID(ctx, id)
}

// Delete 删除电源
func (s *powerService) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

// List 获取电源列表
func (s *powerService) List(ctx context.Context, req *power.PowerSupplyQueryRequest) ([]*power.PowerSupply, int64, error) {
	// 构建查询选项
	page, pageSize := common.GetPageInfo(req.Page, req.PageSize)

	queryOpts := &power.QueryOptions{
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
