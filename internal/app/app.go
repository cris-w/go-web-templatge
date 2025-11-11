package app

import (
	"context"
	"fmt"
	"net/http"
	"power-supply-sys/internal/domain/power"
	"power-supply-sys/internal/domain/user"
	"power-supply-sys/internal/infra/db"
	"power-supply-sys/internal/service"
	httputil "power-supply-sys/internal/transport/http"
	httphandler "power-supply-sys/internal/transport/http/handler"
	httpmiddleware "power-supply-sys/internal/transport/http/middleware"
	"power-supply-sys/pkg/auth"
	"power-supply-sys/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// App 应用结构体
type App struct {
	config *Config
	db     *gorm.DB
	router *gin.Engine
	server *http.Server
}

// New 创建新的应用实例
func New() (*App, error) {
	app := &App{}

	// 1. 加载配置
	config, err := LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("加载配置失败: %w", err)
	}
	app.config = config

	// 2. 初始化日志系统
	if err := app.initLogger(); err != nil {
		return nil, fmt.Errorf("初始化日志系统失败: %w", err)
	}
	logger.Info("Logger initialized successfully")

	// 3. 初始化数据库
	dbConfig := &db.Config{
		DSN:             config.DB.DSN,
		MaxIdleConns:    config.DB.MaxIdleConns,
		MaxOpenConns:    config.DB.MaxOpenConns,
		ConnMaxLifetime: config.DB.GetConnMaxLifetime(),
		Debug:           config.Debug,
	}
	database, err := db.InitDatabase(dbConfig)
	if err != nil {
		return nil, fmt.Errorf("初始化数据库失败: %w", err)
	}
	app.db = database
	logger.Info("Database initialized successfully")

	// 4. 执行数据库迁移
	if err := app.migrate(); err != nil {
		return nil, fmt.Errorf("数据库迁移失败: %w", err)
	}
	logger.Info("Database migration completed")

	// 5. 设置路由
	app.setupRouter()
	logger.Info("Router setup completed")

	return app, nil
}

// initLogger 初始化日志系统
func (a *App) initLogger() error {
	logConfig := &logger.Config{
		Level:      a.config.Log.Level,
		FilePath:   a.config.Log.FilePath,
		MaxSize:    100,
		MaxBackups: 7,
		MaxAge:     30,
		Compress:   true,
		Debug:      a.config.Debug,
	}
	return logger.Init(logConfig)
}

// migrate 执行数据库迁移
func (a *App) migrate() error {
	// 迁移用户表
	if err := user.AutoMigrate(a.db); err != nil {
		return fmt.Errorf("用户表迁移失败: %w", err)
	}

	// 迁移电源表
	if err := power.AutoMigrate(a.db); err != nil {
		return fmt.Errorf("电源表迁移失败: %w", err)
	}

	return nil
}

// healthCheckHandler 健康检查处理函数
func (a *App) healthCheckHandler(c *gin.Context) {
	health := gin.H{
		"status":   "ok",
		"database": "disconnected",
	}

	// 检查数据库连接
	if sqlDB, err := a.db.DB(); err == nil {
		if err := sqlDB.Ping(); err == nil {
			health["database"] = "connected"
		}
	}

	httputil.SuccessResponse(c, health)
}

// setupRouter 设置路由
func (a *App) setupRouter() {
	if !a.config.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	// 使用中间件
	r.Use(httpmiddleware.Logger())
	r.Use(httpmiddleware.Recovery())
	r.Use(httpmiddleware.CORS())

	// 健康检查
	r.GET("/health", a.healthCheckHandler)

	// 创建 JWT 管理器（依赖注入）
	jwtManager := auth.NewJWTManager(a.config.JWT.Secret, a.config.JWT.ExpireHours)

	// 初始化 Services（整个应用共享）
	userService := service.NewUserService(a.db)
	powerService := service.NewPowerService(a.db)

	// 初始化 Handlers（整个应用共享）
	userHandler := httphandler.NewUserHandler(userService, jwtManager)
	powerHandler := httphandler.NewPowerHandler(powerService)

	// 注册 API 路由
	a.registerAPIRoutes(r, userHandler, powerHandler, jwtManager)

	a.router = r
}

// registerAPIRoutes 注册 API 路由
func (a *App) registerAPIRoutes(r *gin.Engine, userHandler *httphandler.UserHandler, powerHandler *httphandler.PowerHandler, jwtManager *auth.JWTManager) {
	// API v1 路由组
	v1 := r.Group("/api/v1")
	{
		// 认证相关路由（无需JWT验证）
		a.registerAuthRoutes(v1, userHandler)

		// 需要JWT认证的路由
		authorized := v1.Group("")
		authorized.Use(httpmiddleware.JWTAuth(jwtManager))
		{
			a.registerUserRoutes(authorized, userHandler)
			a.registerPowerRoutes(authorized, powerHandler)
		}
	}
}

// registerAuthRoutes 注册认证路由
func (a *App) registerAuthRoutes(rg *gin.RouterGroup, handler *httphandler.UserHandler) {
	authGroup := rg.Group("/auth")
	{
		authGroup.POST("/login", handler.Login)
		authGroup.POST("/register", handler.Create)
	}
}

// registerUserRoutes 注册用户路由
func (a *App) registerUserRoutes(rg *gin.RouterGroup, handler *httphandler.UserHandler) {
	userGroup := rg.Group("/users")
	{
		userGroup.GET("", handler.List)
		userGroup.GET("/:id", handler.Get)
		userGroup.PUT("/:id", handler.Update)
		userGroup.DELETE("/:id", handler.Delete)
	}
}

// registerPowerRoutes 注册电源路由
func (a *App) registerPowerRoutes(rg *gin.RouterGroup, handler *httphandler.PowerHandler) {
	powerGroup := rg.Group("/powers")
	{
		powerGroup.GET("", handler.List)
		powerGroup.GET("/:id", handler.Get)
		powerGroup.POST("", handler.Create)
		powerGroup.PUT("/:id", handler.Update)
		powerGroup.DELETE("/:id", handler.Delete)
	}
}

// Run 启动应用
func (a *App) Run() error {
	a.server = &http.Server{
		Addr:         a.config.Addr,
		Handler:      a.router,
		ReadTimeout:  a.config.Server.GetReadTimeout(),
		WriteTimeout: a.config.Server.GetWriteTimeout(),
		IdleTimeout:  a.config.Server.GetIdleTimeout(),
	}

	logger.Info("Starting server",
		zap.String("addr", a.config.Addr),
		zap.Duration("read_timeout", a.config.Server.GetReadTimeout()),
		zap.Duration("write_timeout", a.config.Server.GetWriteTimeout()),
		zap.Duration("idle_timeout", a.config.Server.GetIdleTimeout()),
	)

	// ListenAndServe 会阻塞直到出现错误或调用 Shutdown
	if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("Server start failed", zap.Error(err))
		return fmt.Errorf("服务器启动失败: %w", err)
	}

	return nil
}

// Shutdown 优雅关闭应用
func (a *App) Shutdown(ctx context.Context) error {
	logger.Info("Shutting down server...")

	// 关闭 HTTP 服务器
	if a.server != nil {
		if err := a.server.Shutdown(ctx); err != nil {
			logger.Error("Failed to shutdown server", zap.Error(err))
			return fmt.Errorf("服务器关闭失败: %w", err)
		}
		logger.Info("HTTP server shutdown successfully")
	}

	// 关闭数据库连接
	if a.db != nil {
		sqlDB, err := a.db.DB()
		if err != nil {
			logger.Error("Failed to get DB instance", zap.Error(err))
			return fmt.Errorf("获取数据库连接失败: %w", err)
		}
		if err := sqlDB.Close(); err != nil {
			logger.Error("Failed to close database connection", zap.Error(err))
			return fmt.Errorf("关闭数据库连接失败: %w", err)
		}
		logger.Info("Database connection closed successfully")
	}

	// 同步日志
	logger.Sync()

	logger.Info("Application shutdown successfully")
	return nil
}
