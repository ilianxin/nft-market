#!/bin/bash

# NFT市场应用开发环境启动脚本

echo "🚀 启动NFT市场开发环境..."

# 检查是否在正确的目录
if [ ! -f "setup.sh" ]; then
    echo "❌ 请在项目根目录运行此脚本"
    exit 1
fi

# 启动Hardhat本地网络
echo "🔗 启动本地区块链网络..."
cd contracts
npx hardhat node &
HARDHAT_PID=$!
echo "✅ Hardhat网络已启动 (PID: $HARDHAT_PID)"

# 等待网络启动
sleep 5

# 部署合约
echo "📄 部署智能合约..."
npx hardhat run scripts/deploy.js --network localhost
cd ..

# 启动后端服务
echo "⚙️ 启动后端服务..."
cd backend
go run main.go &
BACKEND_PID=$!
echo "✅ 后端服务已启动 (PID: $BACKEND_PID)"
cd ..

# 等待后端启动
sleep 3

# 启动前端应用
echo "🎨 启动前端应用..."
cd frontend
npm start &
FRONTEND_PID=$!
echo "✅ 前端应用已启动 (PID: $FRONTEND_PID)"
cd ..

echo ""
echo "🎉 NFT市场开发环境启动完成！"
echo ""
echo "🌐 服务地址："
echo "- 前端应用: http://localhost:3000"
echo "- 后端API: http://localhost:8080"
echo "- 区块链RPC: http://localhost:8545"
echo ""
echo "💡 提示："
echo "- 使用 Ctrl+C 停止此脚本"
echo "- 合约地址信息保存在 contracts/contract-addresses.json"
echo "- 日志文件位置："
echo "  - 后端日志: backend/logs/"
echo "  - 前端日志: 浏览器控制台"
echo ""

# 等待用户中断
trap "echo '⏹️ 正在停止服务...'; kill $HARDHAT_PID $BACKEND_PID $FRONTEND_PID 2>/dev/null; exit 0" INT

# 保持脚本运行
wait
