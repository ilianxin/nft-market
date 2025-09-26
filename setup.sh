#!/bin/bash

# NFT市场应用安装脚本

echo "🚀 开始安装NFT市场应用..."

# 检查前置条件
echo "📋 检查前置条件..."

# 检查Node.js
if ! command -v node &> /dev/null; then
    echo "❌ Node.js 未安装，请先安装 Node.js 16+ 版本"
    exit 1
fi

# 检查Go
if ! command -v go &> /dev/null; then
    echo "❌ Go 未安装，请先安装 Go 1.19+ 版本"
    exit 1
fi

# 检查MySQL
if ! command -v mysql &> /dev/null; then
    echo "⚠️ MySQL 未安装，请确保已安装并运行 MySQL"
fi

echo "✅ 前置条件检查完成"

# 安装智能合约依赖
echo "📦 安装智能合约依赖..."
cd contracts
npm install
echo "✅ 智能合约依赖安装完成"

# 编译智能合约
echo "🔨 编译智能合约..."
npx hardhat compile
echo "✅ 智能合约编译完成"

# 安装后端依赖
echo "📦 安装后端依赖..."
cd ../backend
go mod tidy
echo "✅ 后端依赖安装完成"

# 安装前端依赖
echo "📦 安装前端依赖..."
cd ../frontend
npm install
echo "✅ 前端依赖安装完成"

# 复制环境变量文件
echo "📝 设置环境变量..."
cd ..
if [ ! -f backend/.env ]; then
    cp backend/env.example backend/.env
    echo "✅ 后端环境变量文件已创建，请编辑 backend/.env 文件"
fi

if [ ! -f frontend/.env ]; then
    cp frontend/env.example frontend/.env
    echo "✅ 前端环境变量文件已创建，请编辑 frontend/.env 文件"
fi

echo ""
echo "🎉 NFT市场应用安装完成！"
echo ""
echo "📖 下一步操作："
echo "1. 编辑 backend/.env 文件，配置数据库连接和区块链信息"
echo "2. 编辑 frontend/.env 文件，配置API地址和合约地址"
echo "3. 启动本地区块链网络（如Hardhat网络）"
echo "4. 部署智能合约：cd contracts && npm run deploy:localhost"
echo "5. 启动后端服务：cd backend && go run main.go"
echo "6. 启动前端应用：cd frontend && npm start"
echo ""
echo "🌐 应用地址："
echo "- 前端: http://localhost:3000"
echo "- 后端API: http://localhost:8080"
echo ""
