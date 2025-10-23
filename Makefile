.PHONY: run build clean test docker-build docker-run

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
	go test -v ./...

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

