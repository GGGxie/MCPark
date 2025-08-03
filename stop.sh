#!/bin/bash

# MCP Info 停止脚本
# 简单停止后端服务

# 停止可能运行的Go服务
pkill -f "go run cmd/main.go" 2>/dev/null && echo "✅ 后端服务已停止" || echo "ℹ️  没有运行的服务"

# 释放8080端口
if lsof -ti:8080 >/dev/null 2>&1; then
    echo "🔧 释放端口8080..."
    lsof -ti:8080 | xargs kill -9 2>/dev/null || true
fi

echo "🛑 停止完成"