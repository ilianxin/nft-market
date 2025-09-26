@echo off
chcp 65001 >nul

echo 🚀 启动NFT市场开发环境...

:: 检查是否在正确的目录
if not exist "setup.bat" (
    echo ❌ 请在项目根目录运行此脚本
    pause
    exit /b 1
)

:: 启动Hardhat本地网络
echo 🔗 启动本地区块链网络...
cd contracts
start "Hardhat Network" cmd /k "npx hardhat node"
timeout /t 5 >nul

:: 部署合约
echo 📄 部署智能合约...
call npx hardhat run scripts/deploy.js --network localhost
cd ..

:: 启动后端服务
echo ⚙️ 启动后端服务...
cd backend
start "Backend API" cmd /k "go run main.go"
cd ..

:: 等待后端启动
timeout /t 3 >nul

:: 启动前端应用
echo 🎨 启动前端应用...
cd frontend
start "Frontend App" cmd /k "npm start"
cd ..

echo.
echo 🎉 NFT市场开发环境启动完成！
echo.
echo 🌐 服务地址：
echo - 前端应用: http://localhost:3000
echo - 后端API: http://localhost:8080
echo - 区块链RPC: http://localhost:8545
echo.
echo 💡 提示：
echo - 所有服务都在独立的命令行窗口中运行
echo - 关闭相应窗口即可停止对应服务
echo - 合约地址信息保存在 contracts\contract-addresses.json
echo.
pause
