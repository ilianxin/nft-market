#!/bin/bash

# 启动MySQL数据库服务脚本

echo "🚀 启动MySQL数据库服务..."

# 检查Docker是否已安装
if ! command -v docker &> /dev/null; then
    echo "❌ Docker 未安装，请先安装 Docker"
    echo "或者手动安装MySQL 8.0并创建数据库 'nft_market'"
    exit 1
fi

# 检查Docker Compose是否已安装
if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose 未安装，请先安装 Docker Compose"
    exit 1
fi

# 启动MySQL服务
echo "📦 启动MySQL容器..."
docker-compose up -d mysql

# 等待MySQL启动
echo "⏳ 等待MySQL服务启动..."
sleep 10

# 检查MySQL是否已启动
if docker-compose ps mysql | grep -q "Up"; then
    echo "✅ MySQL服务启动成功！"
    echo ""
    echo "📊 数据库连接信息："
    echo "- 主机: localhost"
    echo "- 端口: 3306"
    echo "- 数据库: nft_market"
    echo "- 用户名: nft_user"
    echo "- 密码: nft_password"
    echo "- Root密码: rootpassword"
    echo ""
    echo "🌐 phpMyAdmin 管理界面: http://localhost:8081"
    echo "   用户名: root"
    echo "   密码: rootpassword"
    echo ""
    echo "💡 更新后端 .env 文件中的数据库连接字符串："
    echo "DATABASE_URL=nft_user:nft_password@tcp(localhost:3306)/nft_market?charset=utf8mb4&parseTime=True&loc=Local"
    echo ""
    echo "🛑 停止服务: docker-compose down"
else
    echo "❌ MySQL服务启动失败"
    docker-compose logs mysql
    exit 1
fi
