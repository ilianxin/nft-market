const { ethers } = require("hardhat");

async function main() {
  console.log("🎨 创建MetaMask可显示的NFT...\n");

  try {
    // 获取合约地址
    const contractAddresses = require('../contract-addresses.json');
    console.log("📋 使用合约地址:", contractAddresses.TestNFT);

    // 连接到TestNFT合约
    const TestNFT = await ethers.getContractFactory("TestNFT");
    const testNFT = TestNFT.attach(contractAddresses.TestNFT);

    // 获取部署者账户
    const [deployer] = await ethers.getSigners();
    console.log("👤 部署者账户:", deployer.address);

    // 使用内嵌的base64图片数据，这样MetaMask可以直接显示
    const nftData = [
      {
        tokenId: 7,
        name: "蓝色方块 NFT",
        description: "一个简单的蓝色方块，专为MetaMask显示设计。",
        // 蓝色方块的SVG，转换为base64
        image: "data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMjAwIiBoZWlnaHQ9IjIwMCIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KICA8cmVjdCB3aWR0aD0iMjAwIiBoZWlnaHQ9IjIwMCIgZmlsbD0iIzMzNzNkYyIgcng9IjIwIi8+CiAgPHRleHQgeD0iMTAwIiB5PSI4MCIgZm9udC1mYW1pbHk9IkFyaWFsLCBzYW5zLXNlcmlmIiBmb250LXNpemU9IjE4IiBmaWxsPSJ3aGl0ZSIgdGV4dC1hbmNob3I9Im1pZGRsZSI+8J+UpTwvdGV4dD4KICA8dGV4dCB4PSIxMDAiIHk9IjEyMCIgZm9udC1mYW1pbHk9IkFyaWFsLCBzYW5zLXNlcmlmIiBmb250LXNpemU9IjE2IiBmaWxsPSJ3aGl0ZSIgdGV4dC1hbmNob3I9Im1pZGRsZSI+VGVzdCBORlQ8L3RleHQ+CiAgPHRleHQgeD0iMTAwIiB5PSIxNDAiIGZvbnQtZmFtaWx5PSJBcmlhbCwgc2Fucy1zZXJpZiIgZm9udC1zaXplPSIxNiIgZmlsbD0id2hpdGUiIHRleHQtYW5jaG9yPSJtaWRkbGUiPiM3PC90ZXh0Pgo8L3N2Zz4=",
        attributes: [
          { "trait_type": "Color", "value": "Blue" },
          { "trait_type": "Shape", "value": "Square" },
          { "trait_type": "Network", "value": "Hardhat Local" }
        ]
      },
      {
        tokenId: 8,
        name: "绿色圆形 NFT",
        description: "一个优雅的绿色圆形，带有渐变效果。",
        // 绿色圆形的SVG
        image: "data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMjAwIiBoZWlnaHQ9IjIwMCIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KICA8ZGVmcz4KICAgIDxsaW5lYXJHcmFkaWVudCBpZD0iZ3JhZGllbnQiIHgxPSIwJSIgeTE9IjAlIiB4Mj0iMTAwJSIgeTI9IjEwMCUiPgogICAgICA8c3RvcCBvZmZzZXQ9IjAlIiBzdHlsZT0ic3RvcC1jb2xvcjojMTBiOTgxO3N0b3Atb3BhY2l0eToxIiAvPgogICAgICA8c3RvcCBvZmZzZXQ9IjEwMCUiIHN0eWxlPSJzdG9wLWNvbG9yOiMwNTk2Njlfc3RvcC1vcGFjaXR5OjEiIC8+CiAgICA8L2xpbmVhckdyYWRpZW50PgogIDwvZGVmcz4KICA8Y2lyY2xlIGN4PSIxMDAiIGN5PSIxMDAiIHI9IjkwIiBmaWxsPSJ1cmwoI2dyYWRpZW50KSIvPgogIDx0ZXh0IHg9IjEwMCIgeT0iODAiIGZvbnQtZmFtaWx5PSJBcmlhbCwgc2Fucy1zZXJpZiIgZm9udC1zaXplPSIxOCIgZmlsbD0id2hpdGUiIHRleHQtYW5jaG9yPSJtaWRkbGUiPvCfkZI8L3RleHQ+CiAgPHRleHQgeD0iMTAwIiB5PSIxMjAiIGZvbnQtZmFtaWx5PSJBcmlhbCwgc2Fucy1zZXJpZiIgZm9udC1zaXplPSIxNiIgZmlsbD0id2hpdGUiIHRleHQtYW5jaG9yPSJtaWRkbGUiPlRlc3QgTkZUPC90ZXh0PgogIDx0ZXh0IHg9IjEwMCIgeT0iMTQwIiBmb250LWZhbWlseT0iQXJpYWwsIHNhbnMtc2VyaWYiIGZvbnQtc2l6ZT0iMTYiIGZpbGw9IndoaXRlIiB0ZXh0LWFuY2hvcj0ibWlkZGxlIj4jODwvdGV4dD4KPC9zdmc+",
        attributes: [
          { "trait_type": "Color", "value": "Green" },
          { "trait_type": "Shape", "value": "Circle" },
          { "trait_type": "Network", "value": "Hardhat Local" }
        ]
      },
      {
        tokenId: 9,
        name: "紫色星形 NFT",
        description: "一个神秘的紫色星形，具有特殊的光芒效果。",
        // 紫色星形的SVG
        image: "data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMjAwIiBoZWlnaHQ9IjIwMCIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KICA8ZGVmcz4KICAgIDxyYWRpYWxHcmFkaWVudCBpZD0ic3RhckdyYWRpZW50IiBjeD0iNTAlIiBjeT0iNTAlIiByPSI1MCUiPgogICAgICA8c3RvcCBvZmZzZXQ9IjAlIiBzdHlsZT0ic3RvcC1jb2xvcjojYTg1NWY3O3N0b3Atb3BhY2l0eToxIiAvPgogICAgICA8c3RvcCBvZmZzZXQ9IjEwMCUiIHN0eWxlPSJzdG9wLWNvbG9yOiM3YzNhZWQ7c3RvcC1vcGFjaXR5OjEiIC8+CiAgICA8L3JhZGlhbEdyYWRpZW50PgogIDwvZGVmcz4KICA8cG9seWdvbiBwb2ludHM9IjEwMCwxMCAxMjAsMzAgMTUwLDMwIDEzMCw1MCAxNDAsODAgMTAwLDYwIDYwLDgwIDcwLDUwIDUwLDMwIDgwLDMwIiBmaWxsPSJ1cmwoI3N0YXJHcmFkaWVudCkiIHN0cm9rZT0iI2ZmZiIgc3Ryb2tlLXdpZHRoPSIyIi8+CiAgPHRleHQgeD0iMTAwIiB5PSIxMjAiIGZvbnQtZmFtaWx5PSJBcmlhbCwgc2Fucy1zZXJpZiIgZm9udC1zaXplPSIxOCIgZmlsbD0iI2ZmZiIgdGV4dC1hbmNob3I9Im1pZGRsZSI+4q2Q77iPPC90ZXh0PgogIDx0ZXh0IHg9IjEwMCIgeT0iMTUwIiBmb250LWZhbWlseT0iQXJpYWwsIHNhbnMtc2VyaWYiIGZvbnQtc2l6ZT0iMTYiIGZpbGw9IiNmZmYiIHRleHQtYW5jaG9yPSJtaWRkbGUiPlRlc3QgTkZUICM5PC90ZXh0Pgo8L3N2Zz4=",
        attributes: [
          { "trait_type": "Color", "value": "Purple" },
          { "trait_type": "Shape", "value": "Star" },
          { "trait_type": "Rarity", "value": "Legendary" },
          { "trait_type": "Network", "value": "Hardhat Local" }
        ]
      }
    ];

    console.log("🔄 铸造新的可显示NFT...\n");
    
    for (const nft of nftData) {
      // 创建内嵌的JSON元数据
      const metadata = {
        name: nft.name,
        description: nft.description,
        image: nft.image,
        attributes: nft.attributes
      };
      
      // 将元数据转换为base64编码的data URI
      const metadataJson = JSON.stringify(metadata);
      const metadataBase64 = Buffer.from(metadataJson).toString('base64');
      const tokenURI = `data:application/json;base64,${metadataBase64}`;
      
      console.log(`🎨 正在铸造 ${nft.name}...`);
      console.log(`   - Token ID: ${nft.tokenId}`);
      console.log(`   - 描述: ${nft.description}`);
      
      try {
        const tx = await testNFT.mint(deployer.address, tokenURI);
        await tx.wait();
        
        console.log(`✅ NFT铸造成功!`);
        console.log(`   - 拥有者: ${deployer.address}`);
        console.log(`   - 元数据已内嵌到TokenURI中`);
        console.log();
      } catch (error) {
        console.log(`❌ 铸造失败: ${error.message}`);
        console.log();
      }
    }

    // 获取最终状态
    const totalSupply = await testNFT.getCurrentTokenId();
    console.log(`🎉 铸造完成! 总NFT数量: ${totalSupply}`);
    
    console.log("\n🦊 MetaMask导入指南:");
    console.log("现在你可以在MetaMask中导入这些NFT：");
    console.log(`合约地址: ${contractAddresses.TestNFT}`);
    console.log("可用的Token IDs: 1, 2, 3, 4, 5, 6, 7, 8, 9");
    console.log("\n推荐导入新创建的NFT (更容易显示):");
    console.log("- Token ID 7: 蓝色方块 NFT");
    console.log("- Token ID 8: 绿色圆形 NFT");  
    console.log("- Token ID 9: 紫色星形 NFT");

    console.log("\n📱 导入步骤:");
    console.log("1. 打开MetaMask钱包");
    console.log("2. 确保连接到 'Hardhat Local' 网络");
    console.log("3. 确保使用账户: 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266");
    console.log("4. 点击 'NFTs' 标签页");
    console.log("5. 点击 '导入NFT'");
    console.log("6. 输入合约地址和Token ID (建议从7开始)");
    console.log("7. 点击 '添加'");
    console.log("\n如果NFT仍不显示，请尝试:");
    console.log("- 刷新MetaMask (锁定并重新解锁)");
    console.log("- 切换到其他网络再切换回来");
    console.log("- 重启浏览器");

  } catch (error) {
    console.error("❌ NFT创建失败:", error.message);
    process.exit(1);
  }
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
