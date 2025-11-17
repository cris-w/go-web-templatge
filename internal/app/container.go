package app

import (
	"power-supply-sys/internal/domain/power"
	"power-supply-sys/internal/domain/user"
	"power-supply-sys/internal/infra/repo"
	"power-supply-sys/internal/service"
	"power-supply-sys/pkg/auth"

	"gorm.io/gorm"
)

// Container 依赖容器
type Container struct {
	// 数据库
	DB *gorm.DB

	// Repositories
	UserRepo  user.Repository
	PowerRepo power.Repository

	// Services
	UserService  service.UserService
	PowerService service.PowerService

	// Auth
	JWTManager *auth.JWTManager
}

// NewContainer 创建依赖容器
func NewContainer(cfg *Config, database *gorm.DB) *Container {
	// 创建 Repositories
	userRepo := repo.NewUserRepository(database)
	powerRepo := repo.NewPowerRepository(database)

	// 创建 Services
	userService := service.NewUserService(userRepo)
	powerService := service.NewPowerService(powerRepo)

	// 创建 JWT Manager
	jwtManager := auth.NewJWTManager(cfg.JWT.Secret, cfg.JWT.ExpireHours)

	return &Container{
		DB:           database,
		UserRepo:     userRepo,
		PowerRepo:    powerRepo,
		UserService:  userService,
		PowerService: powerService,
		JWTManager:   jwtManager,
	}
}

