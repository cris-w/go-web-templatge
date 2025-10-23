#!/bin/bash

# Power Supply System 部署脚本

set -e

echo "=== Power Supply System 部署脚本 ==="

# 检查环境变量
if [ -z "$APP_ENV" ]; then
    echo "请设置 APP_ENV 环境变量 (dev/test/prod)"
    exit 1
fi

echo "部署环境: $APP_ENV"

# 拉取最新代码（如果需要）
if [ "$1" == "pull" ]; then
    echo "拉取最新代码..."
    git pull
fi

# 安装依赖
echo "安装依赖..."
go mod download
go mod tidy

# 编译项目
echo "编译项目..."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/power-supply-sys cmd/main.go

# 创建必要的目录
mkdir -p logs

# 停止旧进程（如果存在）
if [ -f "power-supply-sys.pid" ]; then
    echo "停止旧进程..."
    PID=$(cat power-supply-sys.pid)
    if ps -p $PID > /dev/null; then
        kill $PID
        sleep 2
    fi
    rm -f power-supply-sys.pid
fi

# 启动新进程
echo "启动新进程..."
nohup ./bin/power-supply-sys > logs/app.log 2>&1 &
echo $! > power-supply-sys.pid

echo "部署完成！"
echo "PID: $(cat power-supply-sys.pid)"
echo "日志: logs/app.log"

# 检查服务状态
sleep 3
if ps -p $(cat power-supply-sys.pid) > /dev/null; then
    echo "服务运行正常"
    
    # 检查健康检查接口
    if command -v curl &> /dev/null; then
        echo "测试健康检查接口..."
        curl -s http://localhost:8080/health || echo "健康检查失败，请手动检查"
    fi
else
    echo "服务启动失败，请查看日志"
    exit 1
fi

