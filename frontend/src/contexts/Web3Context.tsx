import React, { createContext, useContext, useState, useEffect, useCallback, ReactNode } from 'react';
import { ethers } from 'ethers';

interface Web3ContextType {
  account: string | null;
  provider: ethers.BrowserProvider | null;
  signer: ethers.JsonRpcSigner | null;
  isConnected: boolean;
  connectWallet: () => Promise<void>;
  connectTestWallet: () => void;
  disconnectWallet: () => void;
  switchNetwork: (chainId: string) => Promise<void>;
  chainId: string | null;
  isTestMode: boolean;
}

const Web3Context = createContext<Web3ContextType | undefined>(undefined);

interface Web3ProviderProps {
  children: ReactNode;
}

export const Web3Provider: React.FC<Web3ProviderProps> = ({ children }) => {
  const [account, setAccount] = useState<string | null>(null);
  const [provider, setProvider] = useState<ethers.BrowserProvider | null>(null);
  const [signer, setSigner] = useState<ethers.JsonRpcSigner | null>(null);
  const [chainId, setChainId] = useState<string | null>(null);
  const [isTestMode, setIsTestMode] = useState<boolean>(false);

  const isConnected = account !== null;

  const disconnectWallet = useCallback(() => {
    setAccount(null);
    setProvider(null);
    setSigner(null);
    setChainId(null);
    setIsTestMode(false);
  }, []);

  const handleAccountsChanged = useCallback((accounts: string[]) => {
    if (accounts.length === 0) {
      disconnectWallet();
    } else {
      setAccount(accounts[0]);
    }
  }, [disconnectWallet]);

  const handleChainChanged = useCallback((chainId: string) => {
    setChainId(parseInt(chainId, 16).toString());
    // 重新加载页面以确保状态同步
    window.location.reload();
  }, []);

  // 检查是否已连接钱包
  useEffect(() => {
    checkConnection();
    
    // 监听账户和网络变化
    if (window.ethereum) {
      window.ethereum.on('accountsChanged', handleAccountsChanged);
      window.ethereum.on('chainChanged', handleChainChanged);
    }

    return () => {
      if (window.ethereum) {
        window.ethereum.removeListener('accountsChanged', handleAccountsChanged);
        window.ethereum.removeListener('chainChanged', handleChainChanged);
      }
    };
  }, [handleAccountsChanged, handleChainChanged]);

  const checkConnection = async () => {
    if (window.ethereum) {
      try {
        const provider = new ethers.BrowserProvider(window.ethereum);
        const accounts = await provider.listAccounts();
        
        if (accounts.length > 0) {
          const signer = await provider.getSigner();
          const network = await provider.getNetwork();
          
          setAccount(accounts[0].address);
          setProvider(provider);
          setSigner(signer);
          setChainId(network.chainId.toString());
        }
      } catch (error) {
        console.error('检查连接状态失败:', error);
      }
    }
  };

  const connectWallet = async () => {
    if (!window.ethereum) {
      alert('请安装MetaMask钱包！');
      return;
    }

    try {
      // 请求连接钱包
      await window.ethereum.request({ method: 'eth_requestAccounts' });
      
      const provider = new ethers.BrowserProvider(window.ethereum);
      const signer = await provider.getSigner();
      const address = await signer.getAddress();
      const network = await provider.getNetwork();

      setAccount(address);
      setProvider(provider);
      setSigner(signer);
      setChainId(network.chainId.toString());
    } catch (error) {
      console.error('连接钱包失败:', error);
      alert('连接钱包失败，请重试。');
    }
  };

  const switchNetwork = async (targetChainId: string) => {
    if (!window.ethereum) return;

    try {
      await window.ethereum.request({
        method: 'wallet_switchEthereumChain',
        params: [{ chainId: `0x${parseInt(targetChainId).toString(16)}` }],
      });
    } catch (error: any) {
      if (error.code === 4902) {
        // 网络不存在，需要添加
        console.log('需要添加网络');
      } else {
        console.error('切换网络失败:', error);
      }
    }
  };

  // 测试模式连接（使用Hardhat测试账户）
  const connectTestWallet = () => {
    const testAccounts = [
      '0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266', // Account #0
      '0x70997970C51812dc3A010C7d01b50e0d17dc79C8', // Account #1
      '0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC', // Account #2
    ];
    
    const selectedAccount = testAccounts[0]; // 使用第一个测试账户
    
    setAccount(selectedAccount);
    setChainId('31337'); // Hardhat链ID
    setIsTestMode(true);
    
    // 保存到localStorage供API使用
    localStorage.setItem('userAddress', selectedAccount);
    
    console.log('测试模式已启用，使用账户:', selectedAccount);
  };

  const value: Web3ContextType = {
    account,
    provider,
    signer,
    isConnected,
    connectWallet,
    connectTestWallet,
    disconnectWallet,
    switchNetwork,
    chainId,
    isTestMode,
  };

  return <Web3Context.Provider value={value}>{children}</Web3Context.Provider>;
};

export const useWeb3 = (): Web3ContextType => {
  const context = useContext(Web3Context);
  if (context === undefined) {
    throw new Error('useWeb3 must be used within a Web3Provider');
  }
  return context;
};

// 声明window.ethereum类型
declare global {
  interface Window {
    ethereum?: any;
  }
}
