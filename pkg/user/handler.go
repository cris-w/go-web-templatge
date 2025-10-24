package user

import (
	"power-supply-sys/pkg/auth"
	"power-supply-sys/pkg/common"
	"power-supply-sys/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Handler 用户处理器
type Handler struct {
	service    Service
	jwtManager *auth.JWTManager
}

// NewHandler 创建用户处理器
func NewHandler(service Service, jwtManager *auth.JWTManager) *Handler {
	return &Handler{
		service:    service,
		jwtManager: jwtManager,
	}
}

// Create 创建用户
func (h *Handler) Create(c *gin.Context) {
	ctx := c.Request.Context()
	var req UserCreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("Invalid request parameters", zap.Error(err))
		common.HandleError(c, common.ErrInvalidParam("参数格式错误"))
		return
	}

	user, err := h.service.Create(ctx, &req)
	if err != nil {
		logger.Error("Failed to create user", zap.Error(err), zap.String("username", req.Username))
		common.HandleError(c, err)
		return
	}

	logger.Info("User created successfully", zap.Uint("user_id", user.ID), zap.String("username", user.Username))
	common.HandleSuccess(c, user)
}

// Get 获取用户详情
func (h *Handler) Get(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := common.ParseUintParam(c, "id")
	if err != nil {
		common.HandleError(c, err)
		return
	}

	user, err := h.service.GetByID(ctx, id)
	if err != nil {
		logger.Warn("User not found", zap.Uint("user_id", id), zap.Error(err))
		common.HandleError(c, err)
		return
	}

	common.HandleSuccess(c, user)
}

// Update 更新用户
func (h *Handler) Update(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := common.ParseUintParam(c, "id")
	if err != nil {
		common.HandleError(c, err)
		return
	}

	var req UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("Invalid request parameters", zap.Error(err))
		common.HandleError(c, common.ErrInvalidParam("参数格式错误"))
		return
	}

	user, err := h.service.Update(ctx, id, &req)
	if err != nil {
		logger.Error("Failed to update user", zap.Uint("user_id", id), zap.Error(err))
		common.HandleError(c, err)
		return
	}

	logger.Info("User updated successfully", zap.Uint("user_id", id))
	common.HandleSuccess(c, user)
}

// Delete 删除用户
func (h *Handler) Delete(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := common.ParseUintParam(c, "id")
	if err != nil {
		common.HandleError(c, err)
		return
	}

	if err := h.service.Delete(ctx, id); err != nil {
		logger.Error("Failed to delete user", zap.Uint("user_id", id), zap.Error(err))
		common.HandleError(c, err)
		return
	}

	logger.Info("User deleted successfully", zap.Uint("user_id", id))
	common.HandleSuccess(c, gin.H{"message": "删除成功"})
}

// List 获取用户列表
func (h *Handler) List(c *gin.Context) {
	ctx := c.Request.Context()
	var req UserQueryRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		logger.Warn("Invalid query parameters", zap.Error(err))
		common.HandleError(c, common.ErrInvalidParam("查询参数错误"))
		return
	}

	users, total, err := h.service.List(ctx, &req)
	if err != nil {
		logger.Error("Failed to list users", zap.Error(err))
		common.HandleError(c, err)
		return
	}

	page, pageSize := common.GetPageInfo(req.Page, req.PageSize)
	common.HandlePageSuccess(c, users, total, page, pageSize)
}

// Login 用户登录
func (h *Handler) Login(c *gin.Context) {
	ctx := c.Request.Context()
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("Invalid login request", zap.Error(err))
		common.HandleError(c, common.ErrInvalidParam("参数格式错误"))
		return
	}

	// 验证用户名和密码
	user, err := h.service.Login(ctx, &req)
	if err != nil {
		logger.Warn("Login failed", zap.String("username", req.Username), zap.Error(err))
		common.HandleError(c, err)
		return
	}

	// 生成 JWT token
	if h.jwtManager == nil {
		logger.Error("JWT manager not initialized")
		common.HandleError(c, common.ErrInternal(nil))
		return
	}

	token, err := h.jwtManager.GenerateToken(user.ID, user.Username)
	if err != nil {
		logger.Error("Failed to generate token", zap.Error(err))
		common.HandleError(c, common.ErrInternal(err))
		return
	}

	logger.Info("User logged in successfully", zap.Uint("user_id", user.ID), zap.String("username", user.Username))

	common.HandleSuccess(c, LoginResponse{
		Token: token,
		User:  user,
	})
}
