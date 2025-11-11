package handler

import (
	"power-supply-sys/internal/domain/user"
	"power-supply-sys/internal/service"
	"power-supply-sys/pkg/auth"
	"power-supply-sys/pkg/common"
	"power-supply-sys/pkg/logger"
	httputil "power-supply-sys/internal/transport/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// UserHandler 用户处理器
type UserHandler struct {
	service    service.UserService
	jwtManager *auth.JWTManager
}

// NewUserHandler 创建用户处理器
func NewUserHandler(userService service.UserService, jwtManager *auth.JWTManager) *UserHandler {
	return &UserHandler{
		service:    userService,
		jwtManager: jwtManager,
	}
}

// Create 创建用户
func (h *UserHandler) Create(c *gin.Context) {
	ctx := c.Request.Context()
	var req user.UserCreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("Invalid request parameters", zap.Error(err))
		httputil.HandleError(c, common.ErrInvalidParam("参数格式错误"))
		return
	}

	u, err := h.service.Create(ctx, &req)
	if err != nil {
		logger.Error("Failed to create user", zap.Error(err), zap.String("username", req.Username))
		httputil.HandleError(c, err)
		return
	}

	logger.Info("User created successfully", zap.Uint("user_id", u.ID), zap.String("username", u.Username))
	httputil.HandleSuccess(c, u)
}

// Get 获取用户详情
func (h *UserHandler) Get(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := common.ParseUintParam(c, "id")
	if err != nil {
		httputil.HandleError(c, err)
		return
	}

	u, err := h.service.GetByID(ctx, id)
	if err != nil {
		logger.Warn("User not found", zap.Uint("user_id", id), zap.Error(err))
		httputil.HandleError(c, err)
		return
	}

	httputil.HandleSuccess(c, u)
}

// Update 更新用户
func (h *UserHandler) Update(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := common.ParseUintParam(c, "id")
	if err != nil {
		httputil.HandleError(c, err)
		return
	}

	var req user.UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("Invalid request parameters", zap.Error(err))
		httputil.HandleError(c, common.ErrInvalidParam("参数格式错误"))
		return
	}

	u, err := h.service.Update(ctx, id, &req)
	if err != nil {
		logger.Error("Failed to update user", zap.Uint("user_id", id), zap.Error(err))
		httputil.HandleError(c, err)
		return
	}

	logger.Info("User updated successfully", zap.Uint("user_id", id))
	httputil.HandleSuccess(c, u)
}

// Delete 删除用户
func (h *UserHandler) Delete(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := common.ParseUintParam(c, "id")
	if err != nil {
		httputil.HandleError(c, err)
		return
	}

	if err := h.service.Delete(ctx, id); err != nil {
		logger.Error("Failed to delete user", zap.Uint("user_id", id), zap.Error(err))
		httputil.HandleError(c, err)
		return
	}

	logger.Info("User deleted successfully", zap.Uint("user_id", id))
	httputil.HandleSuccess(c, gin.H{"message": "删除成功"})
}

// List 获取用户列表
func (h *UserHandler) List(c *gin.Context) {
	ctx := c.Request.Context()
	var req user.UserQueryRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		logger.Warn("Invalid query parameters", zap.Error(err))
		httputil.HandleError(c, common.ErrInvalidParam("查询参数错误"))
		return
	}

	users, total, err := h.service.List(ctx, &req)
	if err != nil {
		logger.Error("Failed to list users", zap.Error(err))
		httputil.HandleError(c, err)
		return
	}

	page, pageSize := common.GetPageInfo(req.Page, req.PageSize)
	httputil.HandlePageSuccess(c, users, total, page, pageSize)
}

// Login 用户登录
func (h *UserHandler) Login(c *gin.Context) {
	ctx := c.Request.Context()
	var req user.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("Invalid login request", zap.Error(err))
		httputil.HandleError(c, common.ErrInvalidParam("参数格式错误"))
		return
	}

	// 验证用户名和密码
	u, err := h.service.Login(ctx, &req)
	if err != nil {
		logger.Warn("Login failed", zap.String("username", req.Username), zap.Error(err))
		httputil.HandleError(c, err)
		return
	}

	// 生成 JWT token
	if h.jwtManager == nil {
		logger.Error("JWT manager not initialized")
		httputil.HandleError(c, common.ErrInternal(nil))
		return
	}

	token, err := h.jwtManager.GenerateToken(u.ID, u.Username)
	if err != nil {
		logger.Error("Failed to generate token", zap.Error(err))
		httputil.HandleError(c, common.ErrInternal(err))
		return
	}

	logger.Info("User logged in successfully", zap.Uint("user_id", u.ID), zap.String("username", u.Username))

	httputil.HandleSuccess(c, user.LoginResponse{
		Token: token,
		User:  u,
	})
}

