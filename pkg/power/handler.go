package power

import (
	"power-supply-sys/pkg/common"
	"power-supply-sys/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Handler 电源处理器
type Handler struct {
	service Service
}

// NewHandler 创建电源处理器
func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

// Create 创建电源
func (h *Handler) Create(c *gin.Context) {
	ctx := c.Request.Context()
	var req PowerSupplyCreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("Invalid request parameters", zap.Error(err))
		common.HandleError(c, common.ErrInvalidParam("参数格式错误"))
		return
	}

	powerSupply, err := h.service.Create(ctx, &req)
	if err != nil {
		logger.Error("Failed to create power supply", zap.Error(err), zap.String("name", req.Name))
		common.HandleError(c, err)
		return
	}

	logger.Info("Power supply created successfully", zap.Uint("power_supply_id", powerSupply.ID), zap.String("name", powerSupply.Name))
	common.HandleSuccess(c, powerSupply)
}

// Get 获取电源详情
func (h *Handler) Get(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := common.ParseUintParam(c, "id")
	if err != nil {
		common.HandleError(c, err)
		return
	}

	powerSupply, err := h.service.GetByID(ctx, id)
	if err != nil {
		logger.Warn("Power supply not found", zap.Uint("power_supply_id", id), zap.Error(err))
		common.HandleError(c, err)
		return
	}

	common.HandleSuccess(c, powerSupply)
}

// Update 更新电源
func (h *Handler) Update(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := common.ParseUintParam(c, "id")
	if err != nil {
		common.HandleError(c, err)
		return
	}

	var req PowerSupplyUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("Invalid request parameters", zap.Error(err))
		common.HandleError(c, common.ErrInvalidParam("参数格式错误"))
		return
	}

	powerSupply, err := h.service.Update(ctx, id, &req)
	if err != nil {
		logger.Error("Failed to update power supply", zap.Uint("power_supply_id", id), zap.Error(err))
		common.HandleError(c, err)
		return
	}

	logger.Info("Power supply updated successfully", zap.Uint("power_supply_id", id))
	common.HandleSuccess(c, powerSupply)
}

// Delete 删除电源
func (h *Handler) Delete(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := common.ParseUintParam(c, "id")
	if err != nil {
		common.HandleError(c, err)
		return
	}

	if err := h.service.Delete(ctx, id); err != nil {
		logger.Error("Failed to delete power supply", zap.Uint("power_supply_id", id), zap.Error(err))
		common.HandleError(c, err)
		return
	}

	logger.Info("Power supply deleted successfully", zap.Uint("power_supply_id", id))
	common.HandleSuccess(c, gin.H{"message": "删除成功"})
}

// List 获取电源列表
func (h *Handler) List(c *gin.Context) {
	ctx := c.Request.Context()
	var req PowerSupplyQueryRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		logger.Warn("Invalid query parameters", zap.Error(err))
		common.HandleError(c, common.ErrInvalidParam("查询参数错误"))
		return
	}

	powerSupplies, total, err := h.service.List(ctx, &req)
	if err != nil {
		logger.Error("Failed to list power supplies", zap.Error(err))
		common.HandleError(c, err)
		return
	}

	page, pageSize := common.GetPageInfo(req.Page, req.PageSize)
	common.HandlePageSuccess(c, powerSupplies, total, page, pageSize)
}
