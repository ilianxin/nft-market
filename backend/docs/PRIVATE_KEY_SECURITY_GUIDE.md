# 🔐 私钥安全使用指南

## ⚠️ **重要安全提醒**

### **本地开发环境（推荐）**
```env
# ✅ 安全 - 使用Hardhat测试私钥
PRIVATE_KEY=0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80
```

### **生产环境（危险）**
```env
# ❌ 危险 - 永远不要在代码中暴露真实私钥！
PRIVATE_KEY=0x你的真实私钥  # 这样做会丢失资金！
```

## 🎯 **不同环境的私钥选择**

### **1. 本地开发（Hardhat网络）**
- **使用**: Hardhat提供的测试私钥
- **原因**: 安全、有测试ETH、无风险
- **配置**: `PRIVATE_KEY=0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80`

### **2. 测试网络（Goerli/Sepolia）**
- **使用**: 专门的测试账户私钥
- **原因**: 只有测试币，风险较低
- **获取测试币**: 通过水龙头获取

### **3. 生产环境（主网）**
- **使用**: 硬件钱包或密钥管理服务
- **原因**: 最高安全级别
- **绝对不要**: 将私钥写在代码或配置文件中

## 🔧 **当前项目配置建议**

### **推荐配置（安全）**
```env
# 本地开发 - 使用Hardhat测试账户
ETHEREUM_RPC=http://localhost:8545
CONTRACT_ADDRESS=0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512
PRIVATE_KEY=0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80
```

### **Hardhat测试账户列表**
```
Account #0: 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266
Private Key: 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80

Account #1: 0x70997970C51812dc3A010C7d01b50e0d17dc79C8
Private Key: 0x59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d

Account #2: 0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC
Private Key: 0x5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a
```

## 🛡️ **安全最佳实践**

### **DO（应该做的）**
✅ 本地开发使用公开的测试私钥  
✅ 生产环境使用环境变量  
✅ 使用硬件钱包签名重要交易  
✅ 定期轮换API密钥  
✅ 将.env文件添加到.gitignore  

### **DON'T（绝对不要做的）**
❌ 将真实私钥提交到Git仓库  
❌ 在代码中硬编码私钥  
❌ 在生产环境使用测试私钥  
❌ 与他人共享私钥  
❌ 在不安全的网络中传输私钥  

## 🔍 **如何验证配置正确**

### **检查账户余额**
```bash
curl -X POST -H "Content-Type: application/json" \
  --data '{"jsonrpc":"2.0","method":"eth_getBalance","params":["0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266","latest"],"id":1}' \
  http://localhost:8545
```

### **检查后端连接**
```bash
# 启动后端后访问
curl http://localhost:8080/api/v1/blockchain/status
```

## 💡 **MetaMask集成说明**

### **前端钱包连接**
- 前端使用MetaMask进行**用户交互**
- 用户通过MetaMask签名和发送交易
- 这与后端的私钥配置是**独立的**

### **后端私钥用途**
- 用于**系统级操作**（如同步、管理）
- 不涉及用户资金
- 主要用于读取链上数据和系统维护

## 🎭 **角色分离**

| 组件 | 私钥来源 | 用途 |
|------|----------|------|
| 前端 | MetaMask | 用户交易签名 |
| 后端 | 配置文件 | 系统操作、数据同步 |
| 智能合约 | 无私钥 | 执行业务逻辑 |

## 🔄 **开发流程建议**

1. **开发阶段**: 使用Hardhat测试私钥
2. **测试阶段**: 使用测试网络专用私钥
3. **生产阶段**: 使用硬件钱包或密钥管理服务

记住：**安全第一，便利第二！** 🛡️
