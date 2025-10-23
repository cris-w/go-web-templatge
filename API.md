# Power Supply System API 文档

## 基础信息

- **Base URL**: `http://localhost:9090`
- **认证方式**: JWT Bearer Token

## 认证相关 API

### 1. 用户注册

**POST** `/api/v1/auth/register`

**请求体:**

```json
{
  "username": "testuser",
  "password": "123456",
  "email": "test@example.com",
  "phone": "13800138000",
  "nickname": "测试用户",
  "avatar": "https://example.com/avatar.jpg"
}
```

**响应:**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com",
    "phone": "13800138000",
    "nickname": "测试用户",
    "avatar": "https://example.com/avatar.jpg",
    "status": 1,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

### 2. 用户登录

**POST** `/api/v1/auth/login`

**请求体:**

```json
{
  "username": "testuser",
  "password": "123456"
}
```

**响应:**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "testuser",
      "email": "test@example.com",
      "phone": "13800138000",
      "nickname": "测试用户",
      "avatar": "https://example.com/avatar.jpg",
      "status": 1,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  }
}
```

---

## 用户管理 API（需要认证）

**认证头:** `Authorization: Bearer <token>`

### 3. 获取用户列表

**GET** `/api/v1/users?page=1&page_size=10&username=test&status=1`

**查询参数:**

- `page`: 页码（默认 1）
- `page_size`: 每页数量（默认 10，最大 100）
- `username`: 用户名模糊查询（可选）
- `email`: 邮箱模糊查询（可选）
- `status`: 状态（0-禁用，1-正常）（可选）

**响应:**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "username": "testuser",
        "email": "test@example.com",
        "phone": "13800138000",
        "nickname": "测试用户",
        "avatar": "https://example.com/avatar.jpg",
        "status": 1,
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 100,
    "page": 1,
    "size": 10
  }
}
```

### 4. 获取用户详情

**GET** `/api/v1/users/:id`

**响应:**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com",
    "phone": "13800138000",
    "nickname": "测试用户",
    "avatar": "https://example.com/avatar.jpg",
    "status": 1,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

### 5. 更新用户

**PUT** `/api/v1/users/:id`

**请求体:**

```json
{
  "email": "newemail@example.com",
  "phone": "13900139000",
  "nickname": "新昵称",
  "avatar": "https://example.com/new-avatar.jpg",
  "status": 1
}
```

**响应:**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "username": "testuser",
    "email": "newemail@example.com",
    "phone": "13900139000",
    "nickname": "新昵称",
    "avatar": "https://example.com/new-avatar.jpg",
    "status": 1,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:01Z"
  }
}
```

### 6. 删除用户

**DELETE** `/api/v1/users/:id`

**响应:**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "message": "删除成功"
  }
}
```

---

## 电源管理 API（需要认证）

**认证头:** `Authorization: Bearer <token>`

### 7. 获取电源列表

**GET** `/api/v1/powers?page=1&page_size=10&name=海盗船&min_power=500&max_power=1000`

**查询参数:**

- `page`: 页码（默认 1）
- `page_size`: 每页数量（默认 10，最大 100）
- `name`: 名称模糊查询（可选）
- `brand`: 品牌模糊查询（可选）
- `min_power`: 最小功率（可选）
- `max_power`: 最大功率（可选）
- `min_price`: 最小价格（可选）
- `max_price`: 最大价格（可选）
- `efficiency`: 能效等级（可选）
- `status`: 状态（0-下架，1-上架）（可选）

**响应:**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "name": "海盗船 RM850x",
        "brand": "海盗船",
        "model": "RM850x",
        "power": 850,
        "efficiency": "80Plus金牌",
        "modular": true,
        "price": 899.0,
        "stock": 100,
        "description": "全模组电源",
        "status": 1,
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 50,
    "page": 1,
    "size": 10
  }
}
```

### 8. 获取电源详情

**GET** `/api/v1/powers/:id`

**响应:**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "name": "海盗船 RM850x",
    "brand": "海盗船",
    "model": "RM850x",
    "power": 850,
    "efficiency": "80Plus金牌",
    "modular": true,
    "price": 899.0,
    "stock": 100,
    "description": "全模组电源",
    "status": 1,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

### 9. 创建电源

**POST** `/api/v1/powers`

**请求体:**

```json
{
  "name": "海盗船 RM850x",
  "brand": "海盗船",
  "model": "RM850x",
  "power": 850,
  "efficiency": "80Plus金牌",
  "modular": true,
  "price": 899.0,
  "stock": 100,
  "description": "全模组电源"
}
```

**响应:**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "name": "海盗船 RM850x",
    "brand": "海盗船",
    "model": "RM850x",
    "power": 850,
    "efficiency": "80Plus金牌",
    "modular": true,
    "price": 899.0,
    "stock": 100,
    "description": "全模组电源",
    "status": 1,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

### 10. 更新电源

**PUT** `/api/v1/powers/:id`

**请求体:**

```json
{
  "price": 799.0,
  "stock": 150,
  "status": 1
}
```

**响应:**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "name": "海盗船 RM850x",
    "brand": "海盗船",
    "model": "RM850x",
    "power": 850,
    "efficiency": "80Plus金牌",
    "modular": true,
    "price": 799.0,
    "stock": 150,
    "description": "全模组电源",
    "status": 1,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:01Z"
  }
}
```

### 11. 删除电源

**DELETE** `/api/v1/powers/:id`

**响应:**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "message": "删除成功"
  }
}
```

---

## 健康检查

### 12. 健康检查（无需认证）

**GET** `/health`

**响应:**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "status": "ok",
    "database": "connected"
  }
}
```

---

## 错误码说明

| 错误码 | 说明             |
| ------ | ---------------- |
| 0      | 成功             |
| 1001   | 参数错误         |
| 1002   | 未授权，请先登录 |
| 1003   | 禁止访问         |
| 1004   | 资源不存在       |
| 1005   | 资源已存在       |
| 1006   | Token 无效       |
| 1007   | Token 过期       |
| 1008   | 请求格式错误     |
| 5000   | 服务器内部错误   |
| 5001   | 数据库操作失败   |
| 5002   | 缓存操作失败     |
| 5003   | 服务调用失败     |

---

## 使用示例

### 1. 注册并登录

```bash
# 注册用户
curl -X POST http://localhost:9090/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "123456",
    "email": "test@example.com",
    "nickname": "测试用户"
  }'

# 登录获取 token
curl -X POST http://localhost:9090/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "123456"
  }'
```

### 2. 使用 Token 访问受保护的 API

```bash
# 将登录返回的 token 替换到下面的命令中
export TOKEN="your-jwt-token-here"

# 获取用户列表
curl -X GET http://localhost:9090/api/v1/users \
  -H "Authorization: Bearer $TOKEN"

# 创建电源
curl -X POST http://localhost:9090/api/v1/powers \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "海盗船 RM850x",
    "brand": "海盗船",
    "model": "RM850x",
    "power": 850,
    "efficiency": "80Plus金牌",
    "modular": true,
    "price": 899.00,
    "stock": 100,
    "description": "全模组电源"
  }'
```

---

## 注意事项

1. 除了 `/health`、`/api/v1/auth/register`、`/api/v1/auth/login` 外，所有 API 都需要 JWT 认证
2. JWT Token 默认有效期为 72 小时（可在配置文件中修改）
3. Token 需要放在 HTTP Header 中：`Authorization: Bearer <token>`
4. 分页查询默认页码为 1，默认每页 10 条，最大 100 条
5. 所有时间格式均为 ISO8601 格式
6. 价格字段使用 decimal(10,2) 格式
