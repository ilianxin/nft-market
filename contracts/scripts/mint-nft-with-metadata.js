const { ethers } = require("hardhat");
const fs = require('fs');
const path = require('path');

async function main() {
  console.log("🎨 创建带有完整元数据的测试NFT...\n");

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

    // 创建本地元数据文件
    const metadataDir = path.join(__dirname, '../metadata');
    if (!fs.existsSync(metadataDir)) {
      fs.mkdirSync(metadataDir, { recursive: true });
    }

    // NFT元数据模板
    const nftMetadata = [
      {
        name: "Local Test NFT #1",
        description: "这是一个在本地Hardhat网络上创建的测试NFT，用于验证MetaMask显示功能。",
        image: "data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMjAwIiBoZWlnaHQ9IjIwMCIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KICA8cmVjdCB3aWR0aD0iMjAwIiBoZWlnaHQ9IjIwMCIgZmlsbD0iIzMzNzNkYyIvPgogIDx0ZXh0IHg9IjUwJSIgeT0iNTAlIiBmb250LWZhbWlseT0iQXJpYWwsIHNhbnMtc2VyaWYiIGZvbnQtc2l6ZT0iMjQiIGZpbGw9IndoaXRlIiB0ZXh0LWFuY2hvcj0ibWlkZGxlIiBkeT0iLjNlbSI+TkZUICMxPC90ZXh0Pgo8L3N2Zz4=",
        attributes: [
          {
            "trait_type": "Type",
            "value": "Test NFT"
          },
          {
            "trait_type": "Network",
            "value": "Hardhat Local"
          },
          {
            "trait_type": "Rarity",
            "value": "Common"
          }
        ]
      },
      {
        name: "Local Test NFT #2",
        description: "这是第二个测试NFT，具有不同的颜色和属性。",
        image: "data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMjAwIiBoZWlnaHQ9IjIwMCIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KICA8cmVjdCB3aWR0aD0iMjAwIiBoZWlnaHQ9IjIwMCIgZmlsbD0iIzEwYjk4MSIvPgogIDx0ZXh0IHg9IjUwJSIgeT0iNTAlIiBmb250LWZhbWlseT0iQXJpYWwsIHNhbnMtc2VyaWYiIGZvbnQtc2l6ZT0iMjQiIGZpbGw9IndoaXRlIiB0ZXh0LWFuY2hvcj0ibWlkZGxlIiBkeT0iLjNlbSI+TkZUICMyPC90ZXh0Pgo8L3N2Zz4=",
        attributes: [
          {
            "trait_type": "Type",
            "value": "Test NFT"
          },
          {
            "trait_type": "Network",
            "value": "Hardhat Local"
          },
          {
            "trait_type": "Rarity",
            "value": "Rare"
          }
        ]
      },
      {
        name: "Local Test NFT #3",
        description: "这是第三个测试NFT，具有特殊的紫色主题。",
        image: "data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMjAwIiBoZWlnaHQ9IjIwMCIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KICA8cmVjdCB3aWR0aD0iMjAwIiBoZWlnaHQ9IjIwMCIgZmlsbD0iIzk5MzNmZiIvPgogIDx0ZXh0IHg9IjUwJSIgeT0iNTAlIiBmb250LWZhbWlseT0iQXJpYWwsIHNhbnMtc2VyaWYiIGZvbnQtc2l6ZT0iMjQiIGZpbGw9IndoaXRlIiB0ZXh0LWFuY2hvcj0ibWlkZGxlIiBkeT0iLjNlbSI+TkZUICMzPC90ZXh0Pgo8L3N2Zz4=",
        attributes: [
          {
            "trait_type": "Type",
            "value": "Test NFT"
          },
          {
            "trait_type": "Network",
            "value": "Hardhat Local"
          },
          {
            "trait_type": "Rarity",
            "value": "Epic"
          }
        ]
      }
    ];

    // 保存元数据文件并铸造新的NFT
    console.log("🔄 铸造新的NFT...");
    
    for (let i = 0; i < nftMetadata.length; i++) {
      const metadata = nftMetadata[i];
      const metadataPath = path.join(metadataDir, `${i + 4}.json`); // 从ID 4开始，避免与现有NFT冲突
      
      // 保存元数据文件
      fs.writeFileSync(metadataPath, JSON.stringify(metadata, null, 2));
      console.log(`📄 元数据文件已保存: ${metadataPath}`);
      
      // 使用本地文件路径作为tokenURI
      const tokenURI = `file:///${metadataPath.replace(/\\/g, '/')}`;
      
      // 铸造NFT
      console.log(`🎨 正在铸造 ${metadata.name}...`);
      const tx = await testNFT.mint(deployer.address, tokenURI);
      await tx.wait();
      
      const newTokenId = await testNFT.getCurrentTokenId();
      console.log(`✅ NFT铸造成功! Token ID: ${newTokenId}`);
      console.log(`   - 名称: ${metadata.name}`);
      console.log(`   - 拥有者: ${deployer.address}`);
      console.log(`   - TokenURI: ${tokenURI}`);
      console.log();
    }

    // 获取最终状态
    const totalSupply = await testNFT.getCurrentTokenId();
    console.log(`🎉 铸造完成! 总NFT数量: ${totalSupply}`);
    
    // 显示所有NFT的拥有情况
    console.log("\n📊 所有NFT拥有情况:");
    for (let i = 1; i <= parseInt(totalSupply.toString()); i++) {
      try {
        const owner = await testNFT.ownerOf(i);
        const tokenURI = await testNFT.tokenURI(i);
        console.log(`NFT #${i}:`);
        console.log(`  - 拥有者: ${owner}`);
        console.log(`  - TokenURI: ${tokenURI.substring(0, 80)}...`);
      } catch (error) {
        console.log(`NFT #${i}: 不存在`);
      }
    }

    console.log("\n🦊 MetaMask配置信息:");
    console.log("请在MetaMask中手动导入以下NFT:");
    console.log(`合约地址: ${contractAddresses.TestNFT}`);
    console.log("Token IDs: 1, 2, 3, 4, 5, 6");
    console.log("\n步骤:");
    console.log("1. 打开MetaMask");
    console.log("2. 点击 'NFTs' 标签页");
    console.log("3. 点击 '导入NFT'");
    console.log("4. 输入合约地址和Token ID");
    console.log("5. 点击 '添加'");

  } catch (error) {
    console.error("❌ NFT铸造失败:", error.message);
    process.exit(1);
  }
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
