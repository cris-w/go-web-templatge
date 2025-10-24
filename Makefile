.PHONY: run build clean test test-cover test-coverage test-coverage-html test-nocache test-race docker-build docker-run

# 运行项目
run:
	go run cmd/main.go

# 编译项目
build:
	go build -o bin/power-supply-sys cmd/main.go

# 清理编译产物
clean:
	rm -rf bin/
	rm -rf logs/

# 运行测试
test:
	go test -v ./pkg/...

# 运行测试（显示覆盖率）
test-cover:
	go test -v -cover ./pkg/...

# 生成覆盖率报告
test-coverage:
	go test ./pkg/... -coverprofile=coverage.out
	go tool cover -func=coverage.out

# 生成HTML覆盖率报告
test-coverage-html:
	go test ./pkg/... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# 运行测试（禁用缓存）
test-nocache:
	go test -v -count=1 ./pkg/...

# 运行竞态检测
test-race:
	go test -race ./pkg/...

# 下载依赖
deps:
	go mod download
	go mod tidy

# 构建 Docker 镜像
docker-build:
	docker build -t power-supply-sys:latest .

# 运行 Docker 容器
docker-run:
	docker run -d -p 8080:8080 --name power-supply-sys power-supply-sys:latest

# 运行 Docker compose
docker-compose:
	docker compose -f deployment/docker-compose.yml up -d

# 运行 Docker compose
docker-compose-build:
	docker compose -f deployment/docker-compose.yml up -d --build app

# 开发环境运行
dev:
	APP_ENV=dev go run cmd/main.go

# 生产环境运行
prod:
	APP_ENV=prod go run cmd/main.go

# 格式化代码
fmt:
	go fmt ./...

# 代码检查
lint:
	golangci-lint run ./...

