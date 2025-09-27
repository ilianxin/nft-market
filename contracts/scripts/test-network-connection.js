const { ethers } = require("hardhat");

async function testNetworkConnection() {
  console.log("🔍 测试本地网络连接...\n");

  try {
    // 测试1: 基本网络连接
    console.log("📡 测试1: 检查网络连接");
    const provider = new ethers.JsonRpcProvider("http://localhost:8545");
    const network = await provider.getNetwork();
    console.log("✅ 网络连接成功");
    console.log(`   - 网络名称: ${network.name}`);
    console.log(`   - 链ID: ${network.chainId}`);
    console.log(`   - 区块号: ${await provider.getBlockNumber()}`);

    // 测试2: 账户连接
    console.log("\n👤 测试2: 检查测试账户");
    const testAddress = "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266";
    const balance = await provider.getBalance(testAddress);
    console.log("✅ 测试账户状态正常");
    console.log(`   - 地址: ${testAddress}`);
    console.log(`   - 余额: ${ethers.formatEther(balance)} ETH`);

    // 测试3: 合约连接
    console.log("\n📋 测试3: 检查合约连接");
    const contractAddresses = require('../contract-addresses.json');
    console.log("✅ 合约地址配置:");
    console.log(`   - TestNFT: ${contractAddresses.TestNFT}`);
    console.log(`   - NFTMarketplace: ${contractAddresses.NFTMarketplace}`);

    // 测试4: NFT合约状态
    console.log("\n🎨 测试4: 检查NFT合约");
    const TestNFT = await ethers.getContractFactory("TestNFT");
    const testNFT = TestNFT.attach(contractAddresses.TestNFT);
    
    const name = await testNFT.name();
    const symbol = await testNFT.symbol();
    const totalSupply = await testNFT.getCurrentTokenId();
    
    console.log("✅ NFT合约连接成功");
    console.log(`   - 合约名称: ${name}`);
    console.log(`   - 合约符号: ${symbol}`);
    console.log(`   - NFT总数: ${totalSupply}`);

    // MetaMask配置建议
    console.log("\n🦊 MetaMask配置信息:");
    console.log("   请在MetaMask中添加以下网络:");
    console.log(`   - 网络名称: Hardhat Local`);
    console.log(`   - RPC URL: http://localhost:8545`);
    console.log(`   - 链ID: ${network.chainId}`);
    console.log(`   - 货币符号: ETH`);
    console.log("\n   导入测试账户:");
    console.log(`   - 私钥: ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80`);
    console.log(`   - 地址: ${testAddress}`);

    console.log("\n✅ 所有测试通过！网络配置正确。");
    console.log("如果MetaMask仍无法连接，请检查:");
    console.log("1. 网络配置参数是否完全正确");
    console.log("2. 浏览器是否允许访问localhost");
    console.log("3. 是否有防火墙阻止连接");

  } catch (error) {
    console.error("❌ 网络连接测试失败:");
    console.error("   错误信息:", error.message);
    console.error("\n🔧 可能的解决方案:");
    console.error("1. 确保Hardhat网络正在运行: npx hardhat node");
    console.error("2. 检查端口8545是否被占用");
    console.error("3. 重新部署合约: npx hardhat run scripts/deploy.js --network localhost");
  }
}

// 运行测试
testNetworkConnection();
