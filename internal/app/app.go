package app

import (
	"context"
	"fmt"
	"net/http"
	"power-supply-sys/internal/app/middleware"
	"power-supply-sys/pkg/auth"
	"power-supply-sys/pkg/common"
	"power-supply-sys/pkg/logger"
	"power-supply-sys/pkg/power"
	"power-supply-sys/pkg/user"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// App 应用结构体
type App struct {
	Config       *Config
	DB           *gorm.DB
	Router       *gin.Engine
	Server       *http.Server
	UserService  user.Service
	PowerService power.Service
	JWTManager   *auth.JWTManager
}

// New 创建新的应用实例
func New() (*App, error) {
	app := &App{}

	// 1. 加载配置
	config, err := LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("加载配置失败: %w", err)
	}
	app.Config = config

	// 2. 初始化日志系统
	if err := app.initLogger(); err != nil {
		return nil, fmt.Errorf("初始化日志系统失败: %w", err)
	}
	logger.Info("Logger initialized successfully")

	// 3. 初始化数据库
	db, err := app.initDatabase()
	if err != nil {
		return nil, fmt.Errorf("初始化数据库失败: %w", err)
	}
	app.DB = db
	logger.Info("Database initialized successfully")

	// 4. 初始化 JWT 管理器
	app.initJWT()
	logger.Info("JWT manager initialized successfully")

	// 5. 初始化服务(包含数据库迁移)
	if err := app.initServices(); err != nil {
		return nil, fmt.Errorf("初始化服务失败: %w", err)
	}
	logger.Info("Services initialized successfully")

	// 6. 设置路由
	app.setupRouter()
	logger.Info("Router setup completed")

	return app, nil
}

// initLogger 初始化日志系统
func (a *App) initLogger() error {
	logConfig := &logger.Config{
		Level:      a.Config.Log.Level,
		FilePath:   a.Config.Log.FilePath,
		MaxSize:    100,
		MaxBackups: 7,
		MaxAge:     30,
		Compress:   true,
		Debug:      a.Config.Debug,
	}
	return logger.Init(logConfig)
}

// initDatabase 初始化数据库连接
func (a *App) initDatabase() (*gorm.DB, error) {
	db, err := InitDatabase(a.Config)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// initJWT 初始化 JWT 管理器
func (a *App) initJWT() {
	a.JWTManager = auth.NewJWTManager(a.Config.JWT.Secret, a.Config.JWT.ExpireHours)
}

// initServices 初始化所有服务
func (a *App) initServices() error {
	// 初始化用户服务(包含数据库迁移)
	userService, err := user.NewService(a.DB)
	if err != nil {
		return fmt.Errorf("初始化用户服务失败: %w", err)
	}
	a.UserService = userService

	// 初始化电源服务(包含数据库迁移)
	powerService, err := power.NewService(a.DB)
	if err != nil {
		return fmt.Errorf("初始化电源服务失败: %w", err)
	}
	a.PowerService = powerService

	return nil
}

// setupRouter 设置路由
func (a *App) setupRouter() {
	if !a.Config.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	// 使用中间件
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())
	r.Use(middleware.CORS())

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		health := gin.H{
			"status":   "ok",
			"database": "disconnected",
		}

		// 检查数据库连接
		if sqlDB, err := a.DB.DB(); err == nil {
			if err := sqlDB.Ping(); err == nil {
				health["database"] = "connected"
			}
		}

		common.SuccessResponse(c, health)
	})

	// API v1 路由组
	v1 := r.Group("/api/v1")
	{
		// 认证相关路由（无需JWT验证）
		authGroup := v1.Group("/auth")
		{
			userHandler := user.NewHandler(a.UserService)
			userHandler.SetJWTManager(a.JWTManager)
			authGroup.POST("/login", userHandler.Login)
			authGroup.POST("/register", userHandler.Create)
		}

		// 需要JWT认证的路由
		authorized := v1.Group("")
		authorized.Use(middleware.JWTAuth(a.JWTManager))
		{
			// 用户相关路由
			userHandler := user.NewHandler(a.UserService)
			userGroup := authorized.Group("/users")
			{
				userGroup.GET("", userHandler.List)
				userGroup.GET("/:id", userHandler.Get)
				userGroup.PUT("/:id", userHandler.Update)
				userGroup.DELETE("/:id", userHandler.Delete)
			}

			// 电源相关路由
			powerHandler := power.NewHandler(a.PowerService)
			powerGroup := authorized.Group("/powers")
			{
				powerGroup.GET("", powerHandler.List)
				powerGroup.GET("/:id", powerHandler.Get)
				powerGroup.POST("", powerHandler.Create)
				powerGroup.PUT("/:id", powerHandler.Update)
				powerGroup.DELETE("/:id", powerHandler.Delete)
			}
		}
	}

	a.Router = r
}

// Run 启动应用
func (a *App) Run() error {
	a.Server = &http.Server{
		Addr:    a.Config.Addr,
		Handler: a.Router,
	}

	logger.Info("Starting server", zap.String("addr", a.Config.Addr))

	// ListenAndServe 会阻塞直到出现错误或调用 Shutdown
	if err := a.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("Server start failed", zap.Error(err))
		return fmt.Errorf("服务器启动失败: %w", err)
	}

	return nil
}

// Shutdown 优雅关闭应用
func (a *App) Shutdown(ctx context.Context) error {
	logger.Info("Shutting down server...")

	// 关闭 HTTP 服务器
	if a.Server != nil {
		if err := a.Server.Shutdown(ctx); err != nil {
			logger.Error("Failed to shutdown server", zap.Error(err))
			return fmt.Errorf("服务器关闭失败: %w", err)
		}
		logger.Info("HTTP server shutdown successfully")
	}

	// 关闭数据库连接
	if a.DB != nil {
		sqlDB, err := a.DB.DB()
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
