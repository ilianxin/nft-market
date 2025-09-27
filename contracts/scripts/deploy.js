const { ethers } = require("hardhat");
const hre = require("hardhat");

async function main() {
  console.log("开始部署NFT市场和测试NFT合约...");

  // 获取部署账户
  const [deployer] = await ethers.getSigners();
  console.log("部署账户:", deployer.address);
  console.log("账户余额:", ethers.formatEther(await ethers.provider.getBalance(deployer.address)));

  // 设置手续费收取地址（使用部署者地址作为示例）
  const feeRecipient = deployer.address;

  // 1. 部署测试NFT合约
  console.log("\n1. 部署测试NFT合约...");
  const TestNFT = await ethers.getContractFactory("TestNFT");
  const testNFT = await TestNFT.deploy();
  await testNFT.waitForDeployment();
  console.log("TestNFT合约部署到:", await testNFT.getAddress());

  // 2. 部署NFTMarketplace合约
  console.log("\n2. 部署NFT市场合约...");
  const NFTMarketplace = await ethers.getContractFactory("NFTMarketplace");
  const marketplace = await NFTMarketplace.deploy(feeRecipient);
  await marketplace.waitForDeployment();
  console.log("NFTMarketplace合约部署到:", await marketplace.getAddress());
  console.log("手续费收取地址:", feeRecipient);

  // 3. 验证合约部署状态
  console.log("\n3. 验证合约部署状态...");
  const platformFeeRate = await marketplace.platformFeeRate();
  console.log("平台手续费率:", platformFeeRate.toString(), "基点");
  
  const nftName = await testNFT.name();
  const nftSymbol = await testNFT.symbol();
  console.log("NFT合约名称:", nftName);
  console.log("NFT合约符号:", nftSymbol);

  // 4. 铸造一些测试NFT
  console.log("\n4. 铸造测试NFT...");
  const testTokenURIs = [
    "https://ipfs.io/ipfs/QmYourHash1/metadata.json",
    "https://ipfs.io/ipfs/QmYourHash2/metadata.json", 
    "https://ipfs.io/ipfs/QmYourHash3/metadata.json"
  ];
  
  for (let i = 0; i < testTokenURIs.length; i++) {
    const tx = await testNFT.mint(deployer.address, testTokenURIs[i]);
    await tx.wait();
    console.log(`铸造NFT #${i + 1} 成功`);
  }
  
  const totalSupply = await testNFT.getCurrentTokenId();
  console.log("总共铸造了", totalSupply.toString(), "个NFT");

  console.log("\n部署完成！");
  
  // 5. 保存合约地址到文件
  const fs = require('fs');
  const contractAddresses = {
    TestNFT: await testNFT.getAddress(),
    NFTMarketplace: await marketplace.getAddress(),
    deployer: deployer.address,
    network: hre.network.name,
    totalNFTs: totalSupply.toString()
  };
  
  fs.writeFileSync(
    './contract-addresses.json',
    JSON.stringify(contractAddresses, null, 2)
  );
  
  console.log("合约地址已保存到 contract-addresses.json");
  
  console.log("\n📋 部署摘要:");
  console.log("- TestNFT:", await testNFT.getAddress());
  console.log("- NFTMarketplace:", await marketplace.getAddress());
  console.log("- 测试NFT数量:", totalSupply.toString());
  console.log("- 网络:", hre.network.name);
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
