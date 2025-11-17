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
│   │   ├── config.go      # 配置管理
│   │   └── container.go   # 依赖注入容器
│   ├── domain/            # 领域层（无基础设施依赖）
│   │   ├── user/
│   │   │   ├── model.go        # 用户领域模型
│   │   │   ├── repository.go   # Repository 接口定义
│   │   │   ├── query.go        # 查询选项
│   │   │   └── service_types.go # Service 层类型
│   │   └── power/
│   │       ├── model.go        # 电源领域模型
│   │       ├── repository.go   # Repository 接口定义
│   │       ├── query.go        # 查询选项
│   │       └── service_types.go # Service 层类型
│   ├── service/           # 服务层（依赖接口）
│   │   ├── user_service.go
│   │   ├── user_service_test.go
│   │   ├── power_service.go
│   │   └── power_service_test.go
│   ├── infra/             # 基础设施层
│   │   ├── db/
│   │   │   ├── database.go  # 数据库初始化
│   │   │   └── migrations.go # 数据库迁移
│   │   └── repo/
│   │       ├── user_repo.go      # Repository 实现
│   │       ├── user_repo_test.go
│   │       ├── power_repo.go     # Repository 实现
│   │       └── power_repo_test.go
│   └── transport/         # 传输层
│       └── http/
│           ├── dto/              # 数据传输对象（DTO）
│           │   ├── user_dto.go
│           │   └── power_dto.go
│           ├── handler/          # HTTP 处理器
│           │   ├── user_handler.go
│           │   └── power_handler.go
│           ├── middleware/       # HTTP 中间件
│           │   ├── auth.go          # JWT 认证
│           │   ├── cors.go          # CORS 跨域
│           │   ├── logger.go        # 日志中间件
│           │   ├── recovery.go      # 错误恢复
│           │   └── error_handler.go # 统一错误处理
│           └── response.go          # 统一响应
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

项目采用清晰的分层架构，符合 Go 社区最佳实践和 SOLID 原则：

1. **Domain Layer（领域层）**：`internal/domain/`

   - 包含领域模型、Repository 接口定义、查询选项
   - **无基础设施依赖**，保持领域层纯净
   - Repository 接口定义在领域层，实现依赖倒置原则

2. **Service Layer（服务层）**：`internal/service/`

   - 实现业务逻辑
   - **依赖 Repository 接口**，不直接依赖数据库或具体实现
   - 接收领域类型，与传输层解耦

3. **Infrastructure Layer（基础设施层）**：`internal/infra/`

   - 数据库初始化：`internal/infra/db/`
   - 数据库迁移：`internal/infra/db/migrations.go`
   - Repository 实现：`internal/infra/repo/`（实现 Domain 层定义的接口）

4. **Transport Layer（传输层）**：`internal/transport/http/`

   - DTO 定义：`internal/transport/http/dto/`（请求/响应对象）
   - HTTP 处理器、中间件、响应封装
   - 负责 DTO 与领域对象的转换

5. **Application Layer（应用层）**：`internal/app/`
   - 应用装配、配置管理、生命周期管理
   - 依赖注入容器：统一管理所有依赖

### 设计原则

- **依赖倒置（DIP）**：Service 层依赖 Repository 接口（定义在 Domain 层），不依赖具体实现
- **单一职责（SRP）**：每个包职责明确，Domain 层无基础设施依赖
- **接口隔离（ISP）**：Repository 接口拆分为 Reader 和 Writer，符合接口隔离原则
- **开闭原则（OCP）**：通过接口扩展，对修改关闭，对扩展开放
- **测试友好**：所有依赖可轻松 mock，支持单元测试和集成测试

### 架构特点

- ✅ **Repository 接口在 Domain 层**：符合依赖倒置原则
- ✅ **DTO 在 Transport 层**：传输层关注点分离
- ✅ **依赖注入容器**：统一管理依赖，易于测试和维护
- ✅ **接口隔离**：Reader/Writer 接口分离，职责清晰
- ✅ **无全局状态**：所有依赖通过容器注入

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
- ✅ **依赖注入**：通过 Container 统一管理，松耦合，易测试
- ✅ **依赖倒置**：Service 层依赖接口，Repository 接口定义在 Domain 层
- ✅ **接口隔离**：Repository 接口拆分为 Reader 和 Writer
- ✅ **上下文传递**：支持超时和取消
- ✅ **生命周期管理**：优雅启动和关闭
- ✅ **领域驱动设计**：清晰的领域边界，Domain 层无基础设施依赖
- ✅ **统一错误处理**：错误处理中间件，统一错误响应格式

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
- ✅ **错误处理中间件**：统一错误处理，自动转换错误响应

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

遵循分层架构和依赖倒置原则：

1. **Domain 层**（`internal/domain/{module}/`）

   - 创建 `model.go`：定义领域模型
   - 创建 `repository.go`：定义 Repository 接口（Reader/Writer）
   - 创建 `query.go`：定义查询选项
   - 创建 `service_types.go`：定义 Service 层使用的请求类型

2. **Infrastructure 层**（`internal/infra/`）

   - 在 `db/migrations.go` 中添加模型迁移
   - 在 `repo/` 中实现 Repository 接口

3. **Service 层**（`internal/service/`）

   - 实现业务逻辑，依赖 Repository 接口（而非具体实现）
   - Service 构造函数接收 Repository 接口

4. **Transport 层**（`internal/transport/http/`）

   - 在 `dto/` 中定义请求/响应 DTO
   - 在 `handler/` 中实现 HTTP 处理器
   - Handler 负责 DTO 与领域对象的转换

5. **Application 层**（`internal/app/`）
   - 在 `container.go` 中注册新依赖
   - 在 `app.go` 中注册路由

### 架构最佳实践

- **Domain 层**：保持纯净，无基础设施依赖，只包含领域概念
- **Repository 接口**：定义在 Domain 层，实现依赖倒置
- **DTO**：定义在 Transport 层，负责传输层的数据格式
- **依赖注入**：通过 Container 统一管理，避免全局状态
- **接口隔离**：Repository 接口拆分为 Reader/Writer，按需实现

### 代码规范

- 遵循 Go 官方代码规范
- 使用 `gofmt` 格式化代码
- 使用 `golint` 检查代码质量
- 编写单元测试，保持测试覆盖率 > 80%

## 许可证

MIT License
