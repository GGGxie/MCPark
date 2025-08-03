#!/bin/bash

# MCPark 简单启动脚本
# 构建前端 + 启动后端 + 打开浏览器

set -e

# 颜色定义
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${BLUE}🚀 启动 MCPark 系统...${NC}"

# 检查依赖
if ! command -v go &> /dev/null; then
    echo -e "${RED}❌ Go 未安装${NC}"
    exit 1
fi

if ! command -v npm &> /dev/null; then
    echo -e "${RED}❌ npm 未安装${NC}"
    exit 1
fi

# 构建前端
echo -e "${BLUE}📦 构建前端...${NC}"
cd frontend
if [ ! -d "node_modules" ]; then
    echo "安装前端依赖..."
    npm install
fi
npm run build
cd ..

# 复制前端构建产物到后端静态目录
mkdir -p backend/static
cp -r frontend/dist/* backend/static/

echo -e "${GREEN}✅ 前端构建完成${NC}"

# 启动后端
echo -e "${BLUE}🔧 启动后端服务...${NC}"
cd backend

# 检查go.mod
if [ ! -f "go.mod" ]; then
    go mod init MCPark
fi

# 下载依赖
go mod tidy

# 创建数据目录
mkdir -p ../data

# 启动后端
echo -e "${GREEN}🌟 正在启动后端服务...${NC}"

# 在后台启动后端服务
go run cmd/main.go &
BACKEND_PID=$!

# 等待服务启动
echo -e "${BLUE}等待服务启动...${NC}"
sleep 3

# 检查服务是否正常启动
if curl -s http://localhost:8080/api/v1/health > /dev/null 2>&1; then
    echo -e "${GREEN}✅ 服务启动成功！${NC}"
    echo -e "${BLUE}访问地址: http://localhost:8080${NC}"
    
    # 自动打开浏览器
    if command -v open &> /dev/null; then
        # macOS
        echo -e "${GREEN}📱 正在打开浏览器...${NC}"
        open http://localhost:8080
    elif command -v xdg-open &> /dev/null; then
        # Linux
        echo -e "${GREEN}📱 正在打开浏览器...${NC}"
        xdg-open http://localhost:8080
    elif command -v start &> /dev/null; then
        # Windows
        echo -e "${GREEN}📱 正在打开浏览器...${NC}"
        start http://localhost:8080
    fi
    
    echo -e "${GREEN}🎉 系统启动完成！按 Ctrl+C 停止服务${NC}"
    
    # 等待后端进程
    wait $BACKEND_PID
else
    echo -e "${RED}❌ 服务启动失败${NC}"
    kill $BACKEND_PID 2>/dev/null || true
    exit 1
fi

