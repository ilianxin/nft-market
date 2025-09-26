import axios from 'axios';

const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080/api/v1';

// 创建axios实例
const api = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// 请求拦截器
api.interceptors.request.use(
  (config) => {
    // 添加用户地址到请求头（如果存在）
    const userAddress = localStorage.getItem('userAddress');
    if (userAddress) {
      config.headers['X-User-Address'] = userAddress;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// 响应拦截器
api.interceptors.response.use(
  (response) => {
    return response;
  },
  (error) => {
    console.error('API错误:', error);
    return Promise.reject(error);
  }
);

export const apiService = {
  // 订单相关API
  createOrder: async (orderData: any) => {
    const response = await api.post('/orders', orderData);
    return response.data;
  },

  getOrders: async (page: number = 1, pageSize: number = 20, orderType: string = '') => {
    const params: any = { page, page_size: pageSize };
    if (orderType) {
      params.order_type = orderType;
    }
    const response = await api.get('/orders', { params });
    return response.data;
  },

  getOrderById: async (id: number) => {
    const response = await api.get(`/orders/${id}`);
    return response.data;
  },

  getUserOrders: async (address: string, page: number = 1, pageSize: number = 20, status: string = '') => {
    const params: any = { page, page_size: pageSize };
    if (status) {
      params.status = status;
    }
    const response = await api.get(`/orders/user/${address}`, { params });
    return response.data;
  },

  getNFTOrders: async (contract: string, tokenId: string, page: number = 1, pageSize: number = 20) => {
    const response = await api.get(`/orders/nft/${contract}/${tokenId}`, {
      params: { page, page_size: pageSize }
    });
    return response.data;
  },

  cancelOrder: async (id: number) => {
    const response = await api.put(`/orders/${id}/cancel`);
    return response.data;
  },

  syncOrderFromChain: async (orderID: number) => {
    const response = await api.post(`/orders/sync/${orderID}`);
    return response.data;
  },

  // NFT相关API
  getNFTs: async (page: number = 1, pageSize: number = 20) => {
    const response = await api.get('/nfts', {
      params: { page, page_size: pageSize }
    });
    return response.data;
  },

  getNFTById: async (contract: string, tokenId: string) => {
    const response = await api.get(`/nfts/${contract}/${tokenId}`);
    return response.data;
  },

  getUserNFTs: async (address: string, page: number = 1, pageSize: number = 20) => {
    const response = await api.get(`/nfts/user/${address}`, {
      params: { page, page_size: pageSize }
    });
    return response.data;
  },

  getNFTsByContract: async (contract: string, page: number = 1, pageSize: number = 20) => {
    const response = await api.get(`/nfts/contract/${contract}`, {
      params: { page, page_size: pageSize }
    });
    return response.data;
  },

  searchNFTs: async (keyword: string, page: number = 1, pageSize: number = 20) => {
    const response = await api.get('/nfts/search', {
      params: { q: keyword, page, page_size: pageSize }
    });
    return response.data;
  },

  createOrUpdateNFT: async (nftData: any) => {
    const response = await api.post('/nfts', nftData);
    return response.data;
  },

  // 市场数据API
  getMarketStats: async () => {
    const response = await api.get('/market/stats');
    return response.data;
  },

  getCollections: async () => {
    const response = await api.get('/market/collections');
    return response.data;
  },

  // 健康检查
  healthCheck: async () => {
    const response = await api.get('/health');
    return response.data;
  },
};

export default api;
