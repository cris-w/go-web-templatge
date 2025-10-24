# 测试文档

本项目包含全面的单元测试和集成测试，确保代码质量和功能正确性。

## 运行测试

### 运行所有测试

```bash
# 运行所有测试
go test ./pkg/... -v

# 运行测试并显示覆盖率
go test ./pkg/... -cover

# 生成覆盖率报告
go test ./pkg/... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

### 运行特定包的测试

```bash
# 运行auth包测试
go test ./pkg/auth/... -v

# 运行user包测试
go test ./pkg/user/... -v

# 运行power包测试
go test ./pkg/power/... -v

# 运行common包测试
go test ./pkg/common/... -v
```

### 运行单个测试

```bash
# 运行特定测试函数
go test ./pkg/auth -run TestGenerateToken -v

# 运行匹配模式的测试
go test ./pkg/user -run TestService_ -v
```

### 测试选项

```bash
# 显示详细输出
go test ./pkg/... -v

# 禁用测试缓存（强制重新运行）
go test ./pkg/... -count=1

# 运行基准测试
go test ./pkg/... -bench=. -benchmem

# 并行运行测试
go test ./pkg/... -parallel 4

# 设置超时时间
go test ./pkg/... -timeout 30s
```

## 测试依赖

项目使用以下测试库：

- **testify** - 断言和 Mock 框架

  ```go
  github.com/stretchr/testify/assert
  github.com/stretchr/testify/mock
  github.com/stretchr/testify/require
  ```

- **sqlmock** - 数据库 Mock (未来可选)

  ```go
  github.com/DATA-DOG/go-sqlmock
  ```

- **sqlite** - 测试数据库
  ```go
  gorm.io/driver/sqlite
  ```

## 测试最佳实践

### 1. 测试命名

- 测试文件以 `_test.go` 结尾
- 测试函数以 `Test` 开头
- 使用表驱动测试（Table-Driven Tests）

```go
func TestService_Create(t *testing.T) {
    tests := []struct {
        name    string
        input   interface{}
        wantErr bool
    }{
        // test cases
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // test logic
        })
    }
}
```

### 2. Mock 使用

```go
// 创建Mock对象
mockRepo := new(MockRepository)

// 设置期望
mockRepo.On("FindByID", ctx, uint(1)).Return(expectedUser, nil)

// 执行测试
result, err := service.GetByID(ctx, 1)

// 验证
mockRepo.AssertExpectations(t)
```

### 3. 测试隔离

- 每个测试使用独立的数据库实例
- 使用 `setup` 和 `cleanup` 函数
- 避免测试之间的数据污染

```go
func TestXXX(t *testing.T) {
    repo, cleanup := setupTestRepo(t)
    defer cleanup()

    // test logic
}
```

### 4. 断言风格

```go
// 基本断言
assert.NoError(t, err)
assert.Equal(t, expected, actual)
assert.NotNil(t, result)

// Require (失败时立即停止)
require.NoError(t, err)
require.NotNil(t, result)
```

## 持续集成

测试可以轻松集成到 CI/CD 流程中：

```yaml
# GitHub Actions 示例
- name: Run Tests
  run: go test ./pkg/... -v -cover

- name: Generate Coverage
  run: |
    go test ./pkg/... -coverprofile=coverage.out
    go tool cover -func=coverage.out
```

## 未来改进

- [ ] 提高测试覆盖率到 80%以上
- [ ] 添加 handler 层的集成测试
- [ ] 添加性能基准测试
- [ ] 添加并发测试
- [ ] 集成测试数据工厂
- [ ] 添加 E2E 测试

## 注意事项

1. **数据库测试**：集成测试使用 SQLite 内存数据库，实际生产环境使用 MySQL
2. **异步操作**：测试中的异步操作需要适当的同步机制
3. **外部依赖**：尽量使用 Mock 避免依赖外部服务
4. **测试数据**：测试数据应该具有代表性，覆盖边界情况

## 故障排查

### 常见问题

1. **测试超时**

   ```bash
   go test ./pkg/... -timeout 60s
   ```

2. **并发问题**

   ```bash
   go test ./pkg/... -race
   ```

3. **内存问题**

   ```bash
   go test ./pkg/... -memprofile=mem.out
   ```

4. **清理缓存**
   ```bash
   go clean -testcache
   go test ./pkg/... -count=1
   ```