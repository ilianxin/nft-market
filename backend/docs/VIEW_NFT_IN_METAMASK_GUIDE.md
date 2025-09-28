# 🦊 在MetaMask中查看NFT完整指南

## 🎯 问题解决状态

✅ **已修复的问题**:
- 智能合约添加了`executeOrder`方法
- 部署了真实的ERC721 NFT合约
- 铸造了3个测试NFT
- 更新了合约地址配置

## 📍 当前合约地址

```json
{
  "TestNFT": "0x2279B7A0a67DB372996a5FaB50D91eAA73d2eBe6",
  "NFTMarketplace": "0x8A791620dd6260079BF849Dc5567aDC3F2FdC318",
  "deployer": "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
}
```

## 🔍 在MetaMask中查看NFT的步骤

### 1. **确保连接到正确的网络**

1. 打开MetaMask
2. 确认左上角显示**"Hardhat Local"**网络
3. 如果不是，请切换到本地网络

### 2. **导入NFT合约到MetaMask**

#### 方法A：自动导入（推荐）
1. **打开MetaMask**
2. **点击"NFTs"标签页**
3. **点击"导入NFT"**
4. **填入以下信息**：
   ```
   合约地址: 0x2279B7A0a67DB372996a5FaB50D91eAA73d2eBe6
   Token ID: 1 (或 2, 3)
   ```
5. **点击"添加"**

#### 方法B：手动添加
1. **点击"导入代币"**
2. **选择"自定义代币"**
3. **填入合约地址**: `0x2279B7A0a67DB372996a5FaB50D91eAA73d2eBe6`
4. **代币符号**: `TNC`
5. **小数位**: `0`

### 3. **验证NFT拥有权**

当前测试NFT的拥有者是部署账户：
```
地址: 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266
私钥: ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80
```

**确保在MetaMask中使用这个账户才能看到NFT！**

### 4. **购买NFT后查看**

购买NFT后，新的拥有者需要：

1. **在MetaMask中导入NFT合约**（如上步骤）
2. **使用购买者的账户**查看
3. **NFT会自动显示**在"NFTs"标签页

## 🧪 测试完整流程

### 步骤1：准备两个账户

**账户A (卖家)**:
```
地址: 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266
私钥: ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80
```

**账户B (买家)**:
```  
地址: 0x70997970C51812dc3A010C7d01b50e0d17dc79C8
私钥: 59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d
```

### 步骤2：验证NFT存在

1. **在MetaMask中切换到账户A**
2. **导入NFT合约**
3. **应该能看到3个NFT** (Token ID: 1, 2, 3)

### 步骤3：创建销售订单

1. **使用账户A**在前端创建订单
2. **选择要出售的NFT** (例如Token ID: 1)
3. **设置价格** (例如0.1 ETH)
4. **提交订单**

### 步骤4：购买NFT

1. **在MetaMask中切换到账户B**
2. **在前端购买订单**
3. **确认MetaMask交易**
4. **等待交易完成**

### 步骤5：验证转移

1. **在账户B中导入NFT合约**
2. **应该能看到购买的NFT**
3. **在账户A中NFT应该消失**

## 🚨 常见问题解决

### ❌ 问题1：MetaMask中看不到NFT

**可能原因**:
- 没有导入NFT合约
- 使用了错误的账户
- NFT还没有被铸造

**解决方案**:
1. 确认使用正确的账户地址
2. 手动导入NFT合约地址
3. 检查合约地址是否正确

### ❌ 问题2：购买后NFT没有转移

**可能原因**:
- 区块链交易失败
- 合约执行出错
- 账户权限问题

**解决方案**:
1. 检查MetaMask交易历史
2. 查看后端日志文件
3. 重新部署合约

### ❌ 问题3：交易失败

**可能原因**:
- Gas费用不足
- 合约方法调用失败
- 网络连接问题

**解决方案**:
1. 增加Gas限制
2. 检查Hardhat网络是否运行
3. 重新连接MetaMask

## 🔧 调试工具

### 1. **查看合约状态**

访问API测试页面：`http://localhost:3000/test`

### 2. **检查区块链状态**

```bash
# 在contracts目录下
npx hardhat console --network localhost

# 检查NFT合约
const TestNFT = await ethers.getContractFactory("TestNFT");
const testNFT = TestNFT.attach("0x2279B7A0a67DB372996a5FaB50D91eAA73d2eBe6");
console.log(await testNFT.getCurrentTokenId());
console.log(await testNFT.ownerOf(1));
```

### 3. **查看后端日志**

```bash
# 查看日志文件
tail -f backend/logs/app.log
```

## 🎉 成功标志

✅ **NFT正确显示的标志**:
- MetaMask"NFTs"标签页显示NFT
- 显示正确的合约地址和Token ID
- 显示NFT名称和图片（如果有元数据）
- 交易历史中有转移记录

## 📞 需要帮助？

如果仍然无法看到NFT，请提供：

1. **使用的账户地址**
2. **MetaMask截图**
3. **后端日志错误**
4. **交易哈希**（如果有）

这样我可以帮您进一步诊断问题！

