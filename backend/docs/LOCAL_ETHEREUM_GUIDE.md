# 🌐 本地以太坊网络启动指南

## 🚀 **快速启动**

### **方法1：使用一键启动脚本**
```bash
# 直接运行启动脚本
start-local-ethereum.bat
```

### **方法2：手动启动**

#### **步骤1：启动Hardhat网络**
```bash
# 打开新的终端窗口
cd contracts
npx hardhat node
```

您将看到类似输出：
```
Started HTTP and WebSocket JSON-RPC server at http://127.0.0.1:8545/

Accounts
========

WARNING: These accounts, and their private keys, are publicly known.
Any funds sent to them on Mainnet or any other live network WILL BE LOST.

Account #0: 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266 (10000 ETH)
Private Key: 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80

Account #1: 0x70997970C51812dc3A010C7d01b50e0d17dc79C8 (10000 ETH)
Private Key: 0x59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d

... (更多账户)
```

#### **步骤2：部署合约**
```bash
# 在另一个终端中
cd contracts
npx hardhat run scripts/deploy.js --network localhost
```

成功部署后您会看到：
```
开始部署NFT市场合约...
部署账户: 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266
账户余额: 10000.0
NFTMarketplace合约部署到: 0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512
手续费收取地址: 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266
验证合约部署状态...
平台手续费率: 250 基点
部署完成！
合约地址已保存到 contract-addresses.json
```

## ⚙️ **配置后端**

### **创建.env文件**
在 `backend/` 目录下创建 `.env` 文件：

```env
# 数据库配置
DATABASE_URL=root:root123@tcp(localhost:3306)/nft_market?charset=utf8mb4&parseTime=True&loc=Local

# 服务器配置
PORT=8080
ENVIRONMENT=development

# 以太坊网络配置
ETHEREUM_RPC=http://localhost:8545
CONTRACT_ADDRESS=0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512
PRIVATE_KEY=0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80

# JWT密钥
JWT_SECRET=your_jwt_secret_key

# CORS配置
ALLOWED_ORIGINS=http://localhost:3000,http://localhost:3001
```

**⚠️ 重要**：
- `CONTRACT_ADDRESS`: 使用部署输出中的实际合约地址
- `PRIVATE_KEY`: 使用Account #0的私钥（仅用于本地测试）

## 🔧 **验证网络状态**

### **检查网络连接**
```bash
# 使用curl检查RPC连接
curl -X POST -H "Content-Type: application/json" \
  --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' \
  http://localhost:8545
```

### **检查账户余额**
```bash
curl -X POST -H "Content-Type: application/json" \
  --data '{"jsonrpc":"2.0","method":"eth_getBalance","params":["0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266","latest"],"id":1}' \
  http://localhost:8545
```

## 🏃‍♂️ **启动完整系统**

### **1. 启动MySQL数据库**
```bash
# 如果使用Docker
docker-compose up mysql -d

# 或使用批处理文件
start-mysql.bat
```

### **2. 启动后端服务**
```bash
cd backend
go run main.go
```

### **3. 启动前端应用**
```bash
cd frontend
npm start
```

## 📊 **网络信息总览**

| 服务 | 地址 | 说明 |
|------|------|------|
| Hardhat网络 | http://localhost:8545 | 本地以太坊RPC |
| 后端API | http://localhost:8080 | Go后端服务 |
| 前端应用 | http://localhost:3000 | React前端 |
| MySQL数据库 | localhost:3306 | 数据库服务 |

## 🔍 **测试账户**

Hardhat自动提供20个测试账户，每个都有10000 ETH：

```
Account #0: 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266
Private Key: 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80

Account #1: 0x70997970C51812dc3A010C7d01b50e0d17dc79C8
Private Key: 0x59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d
```

## 🎯 **测试流程**

1. **创建订单**：使用前端界面创建NFT订单
2. **查看日志**：观察后端日志中的区块链交互
3. **检查交易**：在Hardhat终端中查看交易详情
4. **验证同步**：确认数据库和区块链数据一致

## 🛠️ **常见问题**

### **Q: Hardhat网络无法启动**
```bash
# 检查端口是否被占用
netstat -ano | findstr :8545

# 杀死占用端口的进程
taskkill /PID <PID> /F
```

### **Q: 合约部署失败**
- 确保Hardhat网络正在运行
- 检查网络连接：http://localhost:8545
- 重新编译合约：`npx hardhat compile`

### **Q: 后端连接区块链失败**
- 验证.env文件中的ETHEREUM_RPC地址
- 确保CONTRACT_ADDRESS是正确的合约地址
- 检查PRIVATE_KEY格式（需要0x前缀）

### **Q: 前端无法连接钱包**
- 在MetaMask中添加本地网络：
  - 网络名称: Hardhat Local
  - RPC URL: http://localhost:8545
  - 链ID: 31337
  - 货币符号: ETH
- 导入测试账户私钥到MetaMask

## 🔄 **重置网络**

如需重置本地网络：
```bash
# 停止Hardhat网络 (Ctrl+C)
# 重新启动
cd contracts
npx hardhat node

# 重新部署合约
npx hardhat run scripts/deploy.js --network localhost
```

现在您可以开始使用完整的本地以太坊环境进行NFT市场的开发和测试了！🎉
