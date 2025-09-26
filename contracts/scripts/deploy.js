const { ethers } = require("hardhat");
const hre = require("hardhat");

async function main() {
  console.log("开始部署NFT市场合约...");

  // 获取部署账户
  const [deployer] = await ethers.getSigners();
  console.log("部署账户:", deployer.address);
  console.log("账户余额:", ethers.formatEther(await ethers.provider.getBalance(deployer.address)));

  // 设置手续费收取地址（使用部署者地址作为示例）
  const feeRecipient = deployer.address;

  // 部署NFTMarketplace合约
  const NFTMarketplace = await ethers.getContractFactory("NFTMarketplace");
  const marketplace = await NFTMarketplace.deploy(feeRecipient);

  await marketplace.waitForDeployment();

  console.log("NFTMarketplace合约部署到:", await marketplace.getAddress());
  console.log("手续费收取地址:", feeRecipient);

  // 验证合约部署
  console.log("验证合约部署状态...");
  const platformFeeRate = await marketplace.platformFeeRate();
  console.log("平台手续费率:", platformFeeRate.toString(), "基点");

  console.log("部署完成！");
  
  // 保存合约地址到文件
  const fs = require('fs');
  const contractAddresses = {
    NFTMarketplace: await marketplace.getAddress(),
    deployer: deployer.address,
    network: hre.network.name
  };
  
  fs.writeFileSync(
    './contract-addresses.json',
    JSON.stringify(contractAddresses, null, 2)
  );
  
  console.log("合约地址已保存到 contract-addresses.json");
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
