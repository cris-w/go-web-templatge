# Power Supply System

一个基于 Go、Gin 和 GORM 的WEB模板。

## 技术栈

- **Go 1.24**
- **Gin** - Web 框架
- **GORM** - ORM 框架
- **Viper** - 配置管理
- **MySQL** - 数据库

## 项目结构

```
power-supply-sys/
├── cmd/                    # 应用程序入口
│   └── main.go            # 主程序（使用 App 结构体）
├── configs/               # 配置文件
│   ├── config_dev.yml     # 开发环境配置
│   ├── config_test.yml    # 测试环境配置
│   └── config_prod.yml    # 生产环境配置
├── internal/              # 内部应用代码
│   └── app/              # 应用核心
│       ├── app.go        # App 结构体（核心）
│       ├── config.go     # 配置管理
│       ├── database.go   # 数据库初始化
│       ├── response.go   # 统一响应
│       └── middleware/   # 中间件
│           ├── cors.go   # CORS 中间件
│           ├── logger.go # 日志中间件
│           └── recovery.go # 错误恢复中间件
├── pkg/                   # 业务模块
│   ├── user/             # 用户模块
│   │   ├── model.go      # 用户模型
│   │   ├── service.go    # 用户服务（依赖注入）
│   │   └── handler.go    # 用户处理器（依赖注入）
│   └── power/            # 电源模块
│       ├── model.go      # 电源模型
│       ├── service.go    # 电源服务（依赖注入）
│       └── handler.go    # 电源处理器（依赖注入）
├── deployment/           # 部署相关
├── logs/                 # 日志目录
├── Dockerfile           # Docker 配置
├── go.mod              # Go 模块依赖
└── README.md           # 项目说明
```

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
