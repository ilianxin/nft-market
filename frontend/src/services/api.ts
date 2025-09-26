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

  getNFTOrders: async (collectionAddress: string, tokenId: string, page: number = 1, pageSize: number = 20) => {
    const response = await api.get(`/orders/nft/${collectionAddress}/${tokenId}`, {
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

  // 物品相关API (原NFT API)
  getItems: async (page: number = 1, pageSize: number = 20) => {
    const response = await api.get('/items', {
      params: { page, page_size: pageSize }
    });
    return response.data;
  },

  getItemById: async (id: number) => {
    const response = await api.get(`/items/id/${id}`);
    return response.data;
  },

  getItemByTokenId: async (collectionAddress: string, tokenId: string) => {
    const response = await api.get(`/items/token/${collectionAddress}/${tokenId}`);
    return response.data;
  },

  getUserItems: async (address: string, page: number = 1, pageSize: number = 20) => {
    const response = await api.get(`/items/owner/${address}`, {
      params: { page, page_size: pageSize }
    });
    return response.data;
  },

  getItemsByCollection: async (collectionAddress: string, page: number = 1, pageSize: number = 20) => {
    const response = await api.get(`/items/collection/${collectionAddress}`, {
      params: { page, page_size: pageSize }
    });
    return response.data;
  },

  searchItems: async (keyword: string, page: number = 1, pageSize: number = 20) => {
    const response = await api.get('/items/search', {
      params: { q: keyword, page, page_size: pageSize }
    });
    return response.data;
  },

  createOrUpdateItem: async (itemData: any) => {
    const response = await api.post('/items', itemData);
    return response.data;
  },

  updateItemOwner: async (collectionAddress: string, tokenId: string, owner: string) => {
    const response = await api.put(`/items/token/${collectionAddress}/${tokenId}/owner`, { owner });
    return response.data;
  },

  updateItemPrice: async (collectionAddress: string, tokenId: string, price: number) => {
    const response = await api.put(`/items/token/${collectionAddress}/${tokenId}/price`, { price });
    return response.data;
  },

  // 集合相关API
  getCollections: async (page: number = 1, pageSize: number = 20) => {
    const response = await api.get('/collections', {
      params: { page, page_size: pageSize }
    });
    return response.data;
  },

  getCollectionById: async (id: number) => {
    const response = await api.get(`/collections/${id}`);
    return response.data;
  },

  getCollectionByAddress: async (address: string) => {
    const response = await api.get(`/collections/address/${address}`);
    return response.data;
  },

  createCollection: async (collectionData: any) => {
    const response = await api.post('/collections', collectionData);
    return response.data;
  },

  updateCollection: async (id: number, collectionData: any) => {
    const response = await api.put(`/collections/${id}`, collectionData);
    return response.data;
  },

  // 活动相关API
  getActivities: async (page: number = 1, pageSize: number = 20) => {
    const response = await api.get('/activities', {
      params: { page, page_size: pageSize }
    });
    return response.data;
  },

  getActivityById: async (id: number) => {
    const response = await api.get(`/activities/${id}`);
    return response.data;
  },

  getActivitiesByUser: async (address: string, page: number = 1, pageSize: number = 20) => {
    const response = await api.get(`/activities/user/${address}`, {
      params: { page, page_size: pageSize }
    });
    return response.data;
  },

  getActivitiesByItem: async (collectionAddress: string, tokenId: string, page: number = 1, pageSize: number = 20) => {
    const response = await api.get(`/activities/item/${collectionAddress}/${tokenId}`, {
      params: { page, page_size: pageSize }
    });
    return response.data;
  },

  createActivity: async (activityData: any) => {
    const response = await api.post('/activities', activityData);
    return response.data;
  },

  // 市场数据API
  getMarketStats: async () => {
    const response = await api.get('/market/stats');
    return response.data;
  },

  // 健康检查
  healthCheck: async () => {
    const response = await api.get('/health');
    return response.data;
  },
};

export default api;
