import React, { createContext, useContext, useState, useEffect, useCallback, ReactNode } from 'react';
import { ethers } from 'ethers';

interface Web3ContextType {
  account: string | null;
  provider: ethers.BrowserProvider | null;
  signer: ethers.JsonRpcSigner | null;
  isConnected: boolean;
  connectWallet: () => Promise<void>;
  disconnectWallet: () => void;
  switchNetwork: (chainId: string) => Promise<void>;
  chainId: string | null;
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

  const isConnected = account !== null;

  const disconnectWallet = useCallback(() => {
    setAccount(null);
    setProvider(null);
    setSigner(null);
    setChainId(null);
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

  const value: Web3ContextType = {
    account,
    provider,
    signer,
    isConnected,
    connectWallet,
    disconnectWallet,
    switchNetwork,
    chainId,
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
