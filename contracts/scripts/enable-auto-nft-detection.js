const { ethers } = require("hardhat");

async function main() {
  console.log("🔍 配置MetaMask自动NFT检测...\n");

  try {
    // 获取合约地址
    const contractAddresses = require('../contract-addresses.json');
    console.log("📋 NFT合约地址:", contractAddresses.TestNFT);

    // 连接到TestNFT合约
    const TestNFT = await ethers.getContractFactory("TestNFT");
    const testNFT = TestNFT.attach(contractAddresses.TestNFT);

    // 获取部署者账户
    const [deployer] = await ethers.getSigners();
    console.log("👤 账户地址:", deployer.address);

    // 检查合约是否支持必要的接口
    console.log("🔧 检查合约接口支持...");
    
    // 检查ERC721接口
    const ERC721_INTERFACE_ID = "0x80ac58cd";
    const ERC721_METADATA_INTERFACE_ID = "0x5b5e139f";
    const ERC721_ENUMERABLE_INTERFACE_ID = "0x780e9d63";

    try {
      const supportsERC721 = await testNFT.supportsInterface(ERC721_INTERFACE_ID);
      const supportsMetadata = await testNFT.supportsInterface(ERC721_METADATA_INTERFACE_ID);
      
      console.log(`✅ ERC721接口: ${supportsERC721 ? '支持' : '不支持'}`);
      console.log(`✅ ERC721Metadata接口: ${supportsMetadata ? '支持' : '不支持'}`);
      
      try {
        const supportsEnumerable = await testNFT.supportsInterface(ERC721_ENUMERABLE_INTERFACE_ID);
        console.log(`✅ ERC721Enumerable接口: ${supportsEnumerable ? '支持' : '不支持'}`);
      } catch (error) {
        console.log(`⚠️ ERC721Enumerable接口: 不支持 (这是正常的)`);
      }
    } catch (error) {
      console.log("⚠️ 无法检查接口支持，但这不影响NFT显示");
    }

    // 获取所有NFT信息
    console.log("\n📊 检查所有NFT...");
    const totalSupply = await testNFT.getCurrentTokenId();
    console.log(`总NFT数量: ${totalSupply}`);

    const nftList = [];
    for (let i = 1; i <= parseInt(totalSupply.toString()); i++) {
      try {
        const owner = await testNFT.ownerOf(i);
        const tokenURI = await testNFT.tokenURI(i);
        
        nftList.push({
          tokenId: i,
          owner: owner,
          tokenURI: tokenURI,
          isOwnedByUser: owner.toLowerCase() === deployer.address.toLowerCase()
        });
        
        console.log(`NFT #${i}:`);
        console.log(`  - 拥有者: ${owner}`);
        console.log(`  - 是否属于当前账户: ${owner.toLowerCase() === deployer.address.toLowerCase() ? '是' : '否'}`);
        console.log(`  - TokenURI长度: ${tokenURI.length} 字符`);
        
        // 检查TokenURI类型
        if (tokenURI.startsWith('data:application/json')) {
          console.log(`  - 元数据类型: 内嵌JSON (推荐)`);
        } else if (tokenURI.startsWith('data:')) {
          console.log(`  - 元数据类型: Data URI`);
        } else if (tokenURI.startsWith('http')) {
          console.log(`  - 元数据类型: HTTP URL`);
        } else if (tokenURI.startsWith('ipfs://')) {
          console.log(`  - 元数据类型: IPFS`);
        } else {
          console.log(`  - 元数据类型: 其他 (${tokenURI.substring(0, 20)}...)`);
        }
        console.log();
      } catch (error) {
        console.log(`NFT #${i}: 不存在或查询失败`);
      }
    }

    // 统计用户拥有的NFT
    const userNFTs = nftList.filter(nft => nft.isOwnedByUser);
    console.log(`🎯 当前账户拥有的NFT数量: ${userNFTs.length}`);

    // MetaMask自动检测的条件和建议
    console.log("\n🦊 MetaMask自动检测NFT的条件:");
    console.log("1. ✅ 合约必须是标准的ERC721合约");
    console.log("2. ✅ 合约必须部署在支持的网络上");
    console.log("3. ✅ 账户必须拥有NFT");
    console.log("4. ✅ NFT必须有有效的tokenURI");
    console.log("5. ⚠️ MetaMask的自动检测在本地网络上可能不稳定");

    console.log("\n🔧 启用自动检测的方法:");
    console.log("\n方法1: 在MetaMask中启用自动检测");
    console.log("1. 打开MetaMask");
    console.log("2. 点击右上角的头像 → 设置");
    console.log("3. 点击 '安全与隐私'");
    console.log("4. 开启 '自动检测NFT' (Autodetect NFTs)");
    console.log("5. 开启 '自动检测代币' (Autodetect tokens)");

    console.log("\n方法2: 使用前端应用触发检测");
    console.log("1. 在前端应用中连接MetaMask");
    console.log("2. 调用钱包相关的API");
    console.log("3. 这可能触发MetaMask重新扫描NFT");

    console.log("\n方法3: 手动刷新MetaMask");
    console.log("1. 锁定MetaMask钱包");
    console.log("2. 重新解锁");
    console.log("3. 切换网络后再切换回来");
    console.log("4. 重启浏览器");

    // 为本地测试网络的特殊建议
    console.log("\n⚠️ 本地网络特殊说明:");
    console.log("- Hardhat本地网络不在MetaMask的默认检测列表中");
    console.log("- 自动检测在本地网络上成功率较低");
    console.log("- 建议使用手动导入方式");
    console.log("- 或者部署到测试网络 (如Goerli, Sepolia)");

    // 生成批量导入命令
    console.log("\n📝 批量导入命令 (如果自动检测失败):");
    console.log("MetaMask导入信息:");
    console.log(`合约地址: ${contractAddresses.TestNFT}`);
    console.log("Token IDs:", userNFTs.map(nft => nft.tokenId).join(", "));

    // 检查是否有推荐的NFT (Token ID 7-9)
    const recommendedNFTs = userNFTs.filter(nft => nft.tokenId >= 7 && nft.tokenId <= 9);
    if (recommendedNFTs.length > 0) {
      console.log("\n⭐ 推荐优先导入 (更容易显示):");
      recommendedNFTs.forEach(nft => {
        console.log(`- Token ID ${nft.tokenId} (优化的可显示NFT)`);
      });
    }

    console.log("\n🎉 配置完成！");
    console.log("如果启用了自动检测，MetaMask应该会在几分钟内自动显示NFT。");
    console.log("如果没有自动显示，请使用手动导入方式。");

  } catch (error) {
    console.error("❌ 配置失败:", error.message);
    process.exit(1);
  }
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
