package handler

import (
	"power-supply-sys/internal/domain/power"
	"power-supply-sys/internal/service"
	"power-supply-sys/internal/transport/http/dto"
	"power-supply-sys/pkg/common"
	"power-supply-sys/pkg/logger"
	httputil "power-supply-sys/internal/transport/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// PowerHandler 电源处理器
type PowerHandler struct {
	service service.PowerService
}

// NewPowerHandler 创建电源处理器
func NewPowerHandler(powerService service.PowerService) *PowerHandler {
	return &PowerHandler{
		service: powerService,
	}
}

// Create 创建电源
func (h *PowerHandler) Create(c *gin.Context) {
	ctx := c.Request.Context()
	var req dto.PowerSupplyCreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("Invalid request parameters", zap.Error(err))
		c.Error(common.ErrInvalidParam("参数格式错误"))
		return
	}

	// 转换 DTO 为 Service 层需要的格式
	serviceReq := &power.PowerSupplyCreateRequest{
		Name:        req.Name,
		Brand:       req.Brand,
		Model:       req.Model,
		Power:       req.Power,
		Efficiency:  req.Efficiency,
		Modular:     req.Modular,
		Price:       req.Price,
		Stock:       req.Stock,
		Description: req.Description,
	}
	ps, err := h.service.Create(ctx, serviceReq)
	if err != nil {
		logger.Error("Failed to create power supply", zap.Error(err), zap.String("name", req.Name))
		c.Error(err)
		return
	}

	logger.Info("Power supply created successfully", zap.Uint("power_supply_id", ps.ID), zap.String("name", ps.Name))
	httputil.HandleSuccess(c, ps)
}

// Get 获取电源详情
func (h *PowerHandler) Get(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := common.ParseUintParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}

	ps, err := h.service.GetByID(ctx, id)
	if err != nil {
		logger.Warn("Power supply not found", zap.Uint("power_supply_id", id), zap.Error(err))
		c.Error(err)
		return
	}

	httputil.HandleSuccess(c, ps)
}

// Update 更新电源
func (h *PowerHandler) Update(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := common.ParseUintParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}

	var req dto.PowerSupplyUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("Invalid request parameters", zap.Error(err))
		c.Error(common.ErrInvalidParam("参数格式错误"))
		return
	}

	// 转换 DTO 为 Service 层需要的格式
	serviceReq := &power.PowerSupplyUpdateRequest{
		Name:        req.Name,
		Brand:       req.Brand,
		Model:       req.Model,
		Power:       req.Power,
		Efficiency:  req.Efficiency,
		Modular:     req.Modular,
		Price:       req.Price,
		Stock:       req.Stock,
		Description: req.Description,
		Status:      req.Status,
	}
	ps, err := h.service.Update(ctx, id, serviceReq)
	if err != nil {
		logger.Error("Failed to update power supply", zap.Uint("power_supply_id", id), zap.Error(err))
		c.Error(err)
		return
	}

	logger.Info("Power supply updated successfully", zap.Uint("power_supply_id", id))
	httputil.HandleSuccess(c, ps)
}

// Delete 删除电源
func (h *PowerHandler) Delete(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := common.ParseUintParam(c, "id")
	if err != nil {
		c.Error(err)
		return
	}

	if err := h.service.Delete(ctx, id); err != nil {
		logger.Error("Failed to delete power supply", zap.Uint("power_supply_id", id), zap.Error(err))
		c.Error(err)
		return
	}

	logger.Info("Power supply deleted successfully", zap.Uint("power_supply_id", id))
	httputil.HandleSuccess(c, gin.H{"message": "删除成功"})
}

// List 获取电源列表
func (h *PowerHandler) List(c *gin.Context) {
	ctx := c.Request.Context()
	var req dto.PowerSupplyQueryRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		logger.Warn("Invalid query parameters", zap.Error(err))
		c.Error(common.ErrInvalidParam("查询参数错误"))
		return
	}

	// 转换 DTO 为 Service 层需要的格式
	serviceReq := &power.PowerSupplyQueryRequest{
		Page:       req.Page,
		PageSize:   req.PageSize,
		Name:       req.Name,
		Brand:      req.Brand,
		MinPower:   req.MinPower,
		MaxPower:   req.MaxPower,
		MinPrice:   req.MinPrice,
		MaxPrice:   req.MaxPrice,
		Efficiency: req.Efficiency,
		Status:     req.Status,
	}
	powerSupplies, total, err := h.service.List(ctx, serviceReq)
	if err != nil {
		logger.Error("Failed to list power supplies", zap.Error(err))
		c.Error(err)
		return
	}

	page, pageSize := common.GetPageInfo(req.Page, req.PageSize)
	httputil.HandlePageSuccess(c, powerSupplies, total, page, pageSize)
}

