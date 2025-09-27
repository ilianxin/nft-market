const { ethers } = require("hardhat");

async function main() {
  console.log("🔍 检查NFT合约状态...");

  try {
    // 获取合约地址
    const contractAddresses = require('../contract-addresses.json');
    console.log("📋 合约地址信息:", contractAddresses);

    // 连接到TestNFT合约
    const TestNFT = await ethers.getContractFactory("TestNFT");
    const testNFT = TestNFT.attach(contractAddresses.TestNFT);

    console.log("\n🎯 TestNFT合约状态:");
    console.log("- 合约地址:", contractAddresses.TestNFT);

    // 检查合约基本信息
    const name = await testNFT.name();
    const symbol = await testNFT.symbol();
    const totalSupply = await testNFT.getCurrentTokenId();

    console.log("- 合约名称:", name);
    console.log("- 合约符号:", symbol);
    console.log("- 当前Token ID计数:", totalSupply.toString());

    // 检查每个NFT的状态
    console.log("\n🖼️ NFT详细信息:");
    for (let i = 1; i <= parseInt(totalSupply.toString()); i++) {
      try {
        const owner = await testNFT.ownerOf(i);
        const tokenURI = await testNFT.tokenURI(i);
        const creator = await testNFT.getCreator(i);
        
        console.log(`\nNFT #${i}:`);
        console.log("  - 拥有者:", owner);
        console.log("  - 创建者:", creator);
        console.log("  - TokenURI:", tokenURI);
        
        // 检查是否是部署者账户
        if (owner.toLowerCase() === contractAddresses.deployer.toLowerCase()) {
          console.log("  - ✅ 由部署者账户拥有");
        } else {
          console.log("  - ⚠️ 拥有者不是部署者账户");
        }
      } catch (error) {
        console.log(`NFT #${i}: ❌ 不存在或查询失败 -`, error.message);
      }
    }

    // 检查部署者账户余额
    const [deployer] = await ethers.getSigners();
    const balance = await ethers.provider.getBalance(deployer.address);
    console.log("\n💰 部署者账户信息:");
    console.log("- 地址:", deployer.address);
    console.log("- ETH余额:", ethers.formatEther(balance));

    console.log("\n✅ NFT合约状态检查完成！");

  } catch (error) {
    console.error("❌ 检查失败:", error.message);
    process.exit(1);
  }
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
