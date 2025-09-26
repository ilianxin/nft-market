# 🦊 MetaMask连接本地合约指南

本指南将帮助您将MetaMask钱包连接到本地Hardhat网络，并与我们的NFT市场合约进行交互。

## 📋 前提条件

- ✅ 已安装MetaMask浏览器扩展
- ✅ 本地Hardhat网络正在运行
- ✅ NFT市场合约已部署

## 🚀 步骤1：启动本地区块链

### Windows:
```cmd
# 使用一键启动脚本
start-local-ethereum.bat

# 或手动启动
cd contracts
npx hardhat node
```

### Linux/macOS:
```bash
# 手动启动
cd contracts
npx hardhat node
```

**✅ 成功标志**: 终端显示类似以下内容：
```
Started HTTP and WebSocket JSON-RPC server at http://127.0.0.1:8545/
```

## ⚙️ 步骤2：配置MetaMask网络

### 1. 打开MetaMask并添加网络

1. **点击网络下拉菜单**（通常显示"以太坊主网"）
2. **选择"添加网络"**
3. **点击"手动添加网络"**

### 2. 填入网络配置

```
网络名称: Hardhat Local
新RPC URL: http://localhost:8545
链ID: 31337
货币符号: ETH
区块浏览器URL: (留空，或填入 http://localhost:8545)
```

### 3. 保存并切换网络

- 点击**"保存"**
- MetaMask会自动切换到新网络
- 确认左上角显示**"Hardhat Local"**

## 💰 步骤3：导入测试账户

### Hardhat提供的测试账户

Hardhat启动时会创建20个测试账户，每个都有10000 ETH。以下是前几个：

```
账户 #0:
地址: 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266
私钥: 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80

账户 #1:
地址: 0x70997970C51812dc3A010C7d01b50e0d17dc79C8
私钥: 0x59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d

账户 #2:
地址: 0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC
私钥: 0x5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a
```

### 导入账户到MetaMask

1. **点击MetaMask右上角的账户图标**
2. **选择"导入账户"**
3. **选择"私钥"**
4. **粘贴私钥**（例如：`0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80`）
5. **点击"导入"**

**✅ 成功标志**: 账户显示 **10000 ETH** 余额

## 📜 步骤4：验证合约连接

### 当前部署的合约信息

```json
{
  "NFTMarketplace": "0x5FC8d32690cc91D4c39d9d3abcBD16989F875707",
  "deployer": "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266",
  "network": "localhost"
}
```

### 在应用中测试连接

1. **启动前端应用**: `http://localhost:3000`
2. **点击"连接钱包"**按钮
3. **选择MetaMask**并授权连接
4. **确认连接成功**: 右上角显示您的地址

## 🧪 步骤5：测试合约交互

### 方法1：使用前端应用

1. **访问市场页面**: 浏览可购买的NFT
2. **创建订单**: 在"创建订单"页面上架NFT
3. **购买NFT**: 点击"立即购买"按钮
4. **查看交易**: 在MetaMask中确认交易

### 方法2：使用测试页面

1. **访问**: `http://localhost:3000/test`
2. **点击"设置测试用户"**
3. **运行各项API测试**
4. **查看区块链状态**

## 🔧 常见问题解决

### ❌ 问题1：MetaMask显示"无法连接"

**解决方案**:
- 确保Hardhat网络正在运行 (`npx hardhat node`)
- 检查RPC URL是否正确: `http://localhost:8545`
- 重启MetaMask扩展

### ❌ 问题2：余额显示为0

**解决方案**:
- 确认已切换到正确的网络 (Hardhat Local)
- 重新导入私钥
- 刷新MetaMask

### ❌ 问题3：交易失败

**解决方案**:
- 检查Gas费用设置
- 确认合约地址正确
- 查看Hardhat终端的错误信息

### ❌ 问题4：合约交互失败

**解决方案**:
- 确认合约已正确部署
- 检查合约地址是否匹配
- 重新部署合约: `npx hardhat run scripts/deploy.js --network localhost`

## 🎯 验证连接成功

### 检查清单

- [ ] MetaMask显示"Hardhat Local"网络
- [ ] 账户余额显示10000 ETH
- [ ] 前端应用可以连接钱包
- [ ] 可以查看和创建订单
- [ ] 可以进行NFT交易
- [ ] 交易在MetaMask中正确显示

### 成功标志

1. **网络连接**: MetaMask左上角显示"Hardhat Local"
2. **账户余额**: 显示10000 ETH（或接近这个数值）
3. **合约交互**: 可以成功创建和购买订单
4. **交易历史**: MetaMask显示交易记录

## 🚨 安全提醒

⚠️ **重要**: 这些私钥仅用于本地开发测试！

- ❌ **绝对不要**在主网或测试网使用这些私钥
- ❌ **绝对不要**向这些地址发送真实的ETH
- ✅ **仅在本地开发环境**使用
- ✅ **测试完成后**可以删除这些账户

## 🎉 下一步

连接成功后，您可以：

1. **体验完整功能**: 创建、购买、管理NFT
2. **开发新功能**: 基于现有合约扩展功能
3. **学习区块链**: 观察交易和合约交互
4. **测试调试**: 使用测试页面进行API测试

---

**🎊 恭喜！您已成功连接MetaMask到本地NFT市场！**

如有任何问题，请查看终端日志或访问测试页面进行调试。
