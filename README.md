# Power Supply System

一个基于 Go、Gin 和 GORM 的 WEB 模板。

## 技术栈

- **Go 1.24**
- **Gin** - Web 框架
- **GORM** - ORM 框架
- **Viper** - 配置管理
- **MySQL** - 数据库
- **Zap** - 高性能日志库
- **JWT** - 身份认证

## 项目结构

```
power-supply-sys/
├── cmd/                    # 应用程序入口
│   └── main.go            # 主程序入口
├── configs/               # 配置文件
│   ├── config_dev.yml     # 开发环境配置
│   ├── config_test.yml    # 测试环境配置
│   └── config_prod.yml    # 生产环境配置
├── internal/              # 内部应用代码（不对外暴露）
│   ├── app/               # 应用装配与启动
│   │   ├── app.go         # App 结构体（核心）
│   │   └── config.go      # 配置管理
│   ├── domain/            # 领域模型（无基础设施依赖）
│   │   ├── user/
│   │   │   ├── model.go   # 用户模型
│   │   │   └── dto.go     # 数据传输对象
│   │   └── power/
│   │       ├── model.go   # 电源模型
│   │       └── dto.go     # 数据传输对象
│   ├── service/           # 领域服务（依赖接口）
│   │   ├── user_service.go
│   │   ├── user_service_test.go
│   │   ├── power_service.go
│   │   └── power_service_test.go
│   ├── infra/             # 基础设施实现
│   │   ├── db/
│   │   │   └── database.go # 数据库初始化
│   │   └── repo/
│   │       ├── user_repo.go
│   │       ├── user_repo_test.go
│   │       ├── power_repo.go
│   │       └── power_repo_test.go
│   └── transport/         # 传输层
│       └── http/
│           ├── handler/    # HTTP 处理器
│           │   ├── user_handler.go
│           │   └── power_handler.go
│           ├── middleware/  # HTTP 中间件
│           │   ├── auth.go      # JWT 认证
│           │   ├── cors.go      # CORS 跨域
│           │   ├── logger.go    # 日志中间件
│           │   └── recovery.go  # 错误恢复
│           └── response.go      # 统一响应
├── pkg/                   # 可被外部导入的库
│   ├── auth/              # JWT 认证库
│   ├── logger/            # 日志库
│   └── common/            # 通用工具
│       ├── errors.go      # 错误处理
│       ├── utils.go       # 工具函数
│       ├── base_repository.go # 基础仓储
│       └── query_builder.go   # 查询构建器
├── deployment/            # 部署相关
├── logs/                 # 日志目录（已加入 .gitignore）
├── Dockerfile            # Docker 配置
├── go.mod                # Go 模块依赖
└── README.md             # 项目说明
```

## 架构设计

### 分层架构

项目采用清晰的分层架构，符合 Go 社区最佳实践：

1. **Domain Layer（领域层）**：`internal/domain/`

   - 包含领域模型和 DTO
   - 无基础设施依赖，保持纯净

2. **Service Layer（服务层）**：`internal/service/`

   - 实现业务逻辑
   - 依赖仓储接口，不直接依赖数据库

3. **Infrastructure Layer（基础设施层）**：`internal/infra/`

   - 数据库初始化：`internal/infra/db/`
   - 仓储实现：`internal/infra/repo/`

4. **Transport Layer（传输层）**：`internal/transport/http/`

   - HTTP 处理器、中间件、响应封装
   - 与业务逻辑解耦

5. **Application Layer（应用层）**：`internal/app/`
   - 应用装配、配置管理、生命周期管理

### 设计原则

- **依赖倒置**：服务层依赖接口，不依赖具体实现
- **单一职责**：每个包职责明确
- **接口隔离**：接口设计精简，避免臃肿
- **测试友好**：支持单元测试和集成测试

## 功能特性

### 核心功能

- ✅ 多环境配置支持（dev/test/prod）
- ✅ RESTful API 设计
- ✅ 用户管理（CRUD）
- ✅ 电源供应管理（CRUD）
- ✅ 分页查询支持
- ✅ CORS 跨域支持

### 架构特性

- ✅ **分层架构**：Handler → Service → Repository → Model
- ✅ **依赖注入**：松耦合，易测试
- ✅ **上下文传递**：支持超时和取消
- ✅ **生命周期管理**：优雅启动和关闭
- ✅ **领域驱动设计**：清晰的领域边界

### 认证与安全

- ✅ **JWT 认证**：Token 生成和验证
- ✅ **密码加密**：bcrypt 加密存储
- ✅ **权限控制**：路由级别认证
- ✅ **安全响应**：不泄露敏感信息

### 日志与监控

- ✅ **结构化日志**：zap 高性能日志库
- ✅ **日志轮转**：基于大小和时间
- ✅ **请求追踪**：记录每个请求的详细信息
- ✅ **健康检查**：包含数据库状态

### 错误处理

- ✅ **统一错误码**：标准错误码体系
- ✅ **错误分类**：客户端错误 vs 服务端错误
- ✅ **错误链**：支持错误包装和追踪
- ✅ **自动恢复**：Panic 自动恢复

## 快速开始

### 1. 安装依赖

```bash
go mod download
```

### 2. 配置环境

复制并修改配置文件：

```bash
# 开发环境使用 config_dev.yml
# 生产环境使用 config_prod.yml
# 测试环境使用 config_test.yml
```

修改配置文件中的数据库连接信息。

### 3. 设置环境变量

```bash
# 设置运行环境（dev/test/prod），默认为 dev
export APP_ENV=dev
```

### 4. 运行项目

```bash
# 开发环境
go run cmd/main.go

# 或者编译后运行
go build -o power-supply-sys cmd/main.go
./power-supply-sys
```

### 5. 访问 API

服务默认运行在 `http://localhost:9090`（开发环境）

健康检查：

```bash
curl http://localhost:9090/health
```

## 测试

### 运行测试

```bash
# 运行所有测试
go test ./...

# 运行特定包的测试
go test ./internal/service/...
go test ./internal/infra/repo/...

# 运行测试并查看覆盖率
go test -cover ./...

# 生成覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### 测试覆盖

项目包含完整的单元测试：

- ✅ Repository 层测试：覆盖所有 CRUD 操作和查询方法
- ✅ Service 层测试：覆盖业务逻辑和边界情况
- ✅ 使用内存数据库（SQLite）进行快速测试

## 部署

### Docker 部署

```bash
# 构建镜像
docker build -t power-supply-sys .

# 运行容器
docker run -d -p 9090:9090 \
  -e APP_ENV=prod \
  -e DB_DSN="user:password@tcp(db:3306)/database" \
  power-supply-sys
```

### Docker Compose

```bash
cd deployment
# docker compose up -d --build app   重新构建app
docker-compose up -d
```

## 开发指南

### 添加新功能模块

1. 在 `internal/domain/` 创建领域模型和 DTO
2. 在 `internal/infra/repo/` 实现仓储接口
3. 在 `internal/service/` 实现业务逻辑
4. 在 `internal/transport/http/handler/` 创建处理器
5. 在 `internal/app/app.go` 注册路由

### 代码规范

- 遵循 Go 官方代码规范
- 使用 `gofmt` 格式化代码
- 使用 `golint` 检查代码质量
- 编写单元测试，保持测试覆盖率 > 80%

## 许可证

MIT License
