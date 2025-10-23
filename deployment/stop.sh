#!/bin/bash

# 停止服务脚本

set -e

echo "=== 停止 Power Supply System ==="

if [ ! -f "power-supply-sys.pid" ]; then
    echo "PID 文件不存在，服务可能未运行"
    exit 1
fi

PID=$(cat power-supply-sys.pid)

if ps -p $PID > /dev/null; then
    echo "停止进程 PID: $PID"
    kill $PID
    
    # 等待进程结束
    for i in {1..10}; do
        if ! ps -p $PID > /dev/null; then
            echo "进程已停止"
            rm -f power-supply-sys.pid
            exit 0
        fi
        echo "等待进程结束... ($i/10)"
        sleep 1
    done
    
    # 强制结束
    echo "强制结束进程"
    kill -9 $PID
    rm -f power-supply-sys.pid
else
    echo "进程不存在"
    rm -f power-supply-sys.pid
fi

echo "服务已停止"

