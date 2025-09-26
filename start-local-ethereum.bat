@echo off
echo 启动本地以太坊网络...
echo.

echo 1. 启动Hardhat网络
cd contracts
start "Hardhat Network" npx hardhat node

echo.
echo 等待网络启动...
timeout /t 5 /nobreak > nul

echo.
echo 2. 部署合约
npx hardhat run scripts/deploy.js --network localhost

echo.
echo ✅ 本地以太坊网络已启动！
echo.
echo 📍 网络信息:
echo    RPC地址: http://localhost:8545
echo    链ID: 31337
echo    合约地址: 见上方输出
echo.
echo 📝 请将合约地址复制到 backend/.env 文件中的 CONTRACT_ADDRESS
echo.
echo 按任意键继续...
pause > nul
