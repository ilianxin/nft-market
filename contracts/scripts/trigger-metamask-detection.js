// 这个脚本可以在浏览器控制台中运行，帮助触发MetaMask的NFT检测

const triggerNFTDetection = async () => {
  console.log("🔍 尝试触发MetaMask NFT自动检测...");

  if (typeof window.ethereum === 'undefined') {
    console.error("❌ MetaMask未安装或未连接");
    return;
  }

  try {
    // 1. 请求账户连接
    console.log("📱 请求MetaMask连接...");
    const accounts = await window.ethereum.request({ method: 'eth_requestAccounts' });
    console.log("✅ 已连接账户:", accounts[0]);

    // 2. 检查网络
    console.log("🌐 检查网络...");
    const chainId = await window.ethereum.request({ method: 'eth_chainId' });
    console.log("✅ 当前网络链ID:", parseInt(chainId, 16));
    
    if (parseInt(chainId, 16) !== 31337) {
      console.warn("⚠️ 当前不在Hardhat本地网络 (链ID应为31337)");
    }

    // 3. 获取账户余额 (这可能触发MetaMask刷新)
    console.log("💰 获取账户余额...");
    const balance = await window.ethereum.request({
      method: 'eth_getBalance',
      params: [accounts[0], 'latest']
    });
    const ethBalance = parseInt(balance, 16) / Math.pow(10, 18);
    console.log("✅ 账户余额:", ethBalance.toFixed(4), "ETH");

    // 4. 合约信息
    const contractAddress = "0x5FbDB2315678afecb367f032d93F642f64180aa3";
    console.log("📋 NFT合约地址:", contractAddress);

    // 5. 尝试调用合约方法 (这可能触发NFT检测)
    console.log("🔍 调用合约方法检查NFT...");
    
    // ERC721的balanceOf方法
    const balanceOfData = "0x70a08231" + accounts[0].slice(2).padStart(64, '0');
    
    try {
      const nftBalance = await window.ethereum.request({
        method: 'eth_call',
        params: [{
          to: contractAddress,
          data: balanceOfData
        }, 'latest']
      });
      
      const nftCount = parseInt(nftBalance, 16);
      console.log("✅ 账户拥有的NFT数量:", nftCount);
      
      if (nftCount > 0) {
        console.log("🎉 检测到NFT! MetaMask应该会自动显示它们。");
        
        // 6. 建议用户检查MetaMask
        console.log("\n📱 请检查MetaMask:");
        console.log("1. 切换到 'NFTs' 标签页");
        console.log("2. 如果没有自动显示，尝试:");
        console.log("   - 锁定并重新解锁MetaMask");
        console.log("   - 切换网络再切换回来");
        console.log("   - 刷新浏览器页面");
        
        console.log("\n🔧 如果仍然没有显示，请手动导入:");
        console.log("合约地址:", contractAddress);
        console.log("推荐Token IDs: 7, 8, 9 (这些有完整的图片和元数据)");
        
      } else {
        console.log("❌ 账户没有NFT，请确认:");
        console.log("1. 使用正确的账户:", accounts[0]);
        console.log("2. 连接到正确的网络 (Hardhat Local, 链ID: 31337)");
        console.log("3. NFT已经铸造到这个账户");
      }
      
    } catch (contractError) {
      console.error("❌ 合约调用失败:", contractError.message);
      console.log("这可能意味着:");
      console.log("1. 合约地址不正确");
      console.log("2. 网络配置有问题"); 
      console.log("3. 本地Hardhat网络没有运行");
    }

    // 7. 触发一些可能有助于NFT检测的事件
    console.log("🔄 触发检测事件...");
    
    // 发送一个0值交易给自己 (有时能触发余额刷新)
    try {
      const txHash = await window.ethereum.request({
        method: 'eth_sendTransaction',
        params: [{
          from: accounts[0],
          to: accounts[0],
          value: '0x0',
          gas: '0x5208', // 21000 gas
        }]
      });
      console.log("✅ 触发交易已发送:", txHash);
      console.log("这可能帮助MetaMask刷新账户状态");
    } catch (txError) {
      console.log("⚠️ 无法发送触发交易 (这是正常的):", txError.message);
    }

  } catch (error) {
    console.error("❌ 检测过程出错:", error.message);
  }
};

// 导出函数以便在浏览器中使用
if (typeof window !== 'undefined') {
  window.triggerNFTDetection = triggerNFTDetection;
  console.log("✅ 函数已加载! 在浏览器控制台中运行: triggerNFTDetection()");
} else {
  // Node.js环境
  console.log("这个脚本需要在浏览器环境中运行");
  console.log("请复制以下代码到浏览器控制台中运行:");
  console.log("\n" + triggerNFTDetection.toString() + "\n\ntriggerNFTDetection();");
}

// 如果直接运行，执行检测
if (typeof window !== 'undefined' && window.ethereum) {
  triggerNFTDetection();
}
