import React, { useState } from 'react';
import {
  Box,
  Typography,
  Paper,
  Tabs,
  Tab,
  Grid,
  Card,
  CardContent,
  Button,
  Chip,
  CircularProgress,
  Alert,
} from '@mui/material';
import { useParams } from 'react-router-dom';
import { useQuery, useMutation, useQueryClient } from 'react-query';
import { useWeb3 } from '../contexts/Web3Context';
import { apiService } from '../services/api';

interface TabPanelProps {
  children?: React.ReactNode;
  index: number;
  value: number;
}

function TabPanel(props: TabPanelProps) {
  const { children, value, index, ...other } = props;

  return (
    <div
      role="tabpanel"
      hidden={value !== index}
      id={`profile-tabpanel-${index}`}
      aria-labelledby={`profile-tab-${index}`}
      {...other}
    >
      {value === index && <Box sx={{ pt: 3 }}>{children}</Box>}
    </div>
  );
}

const UserProfilePage: React.FC = () => {
  const { address } = useParams<{ address: string }>();
  const { account } = useWeb3();
  const [tabValue, setTabValue] = useState(0);
  const queryClient = useQueryClient();
  
  const isOwnProfile = account?.toLowerCase() === address?.toLowerCase();

  const handleTabChange = (event: React.SyntheticEvent, newValue: number) => {
    setTabValue(newValue);
  };

  // 获取用户订单
  const { data: ordersData, isLoading: ordersLoading } = useQuery(
    ['userOrders', address],
    () => apiService.getUserOrders(address!),
    {
      enabled: !!address,
    }
  );

  // 获取用户NFT
  const { data: nftsData, isLoading: nftsLoading } = useQuery(
    ['userNFTs', address],
    () => apiService.getUserNFTs(address!),
    {
      enabled: !!address,
    }
  );

  // 取消订单mutation
  const cancelOrderMutation = useMutation(
    (orderId: number) => apiService.cancelOrder(orderId),
    {
      onSuccess: () => {
        queryClient.invalidateQueries(['userOrders', address]);
        alert('订单取消成功！');
      },
      onError: (error: any) => {
        alert(`取消订单失败: ${error.response?.data?.message || error.message}`);
      },
    }
  );

  const formatAddress = (addr: string) => {
    return `${addr.slice(0, 6)}...${addr.slice(-4)}`;
  };

  const formatPrice = (price: string) => {
    try {
      const ethValue = parseFloat(price) / Math.pow(10, 18);
      return ethValue.toFixed(4) + ' ETH';
    } catch {
      return price;
    }
  };

  const getOrderStatusColor = (status: string) => {
    switch (status) {
      case 'active': return 'success';
      case 'filled': return 'info';
      case 'cancelled': return 'error';
      case 'expired': return 'warning';
      default: return 'default';
    }
  };

  const getOrderStatusText = (status: string) => {
    switch (status) {
      case 'active': return '活跃';
      case 'filled': return '已成交';
      case 'cancelled': return '已取消';
      case 'expired': return '已过期';
      default: return status;
    }
  };

  const getOrderTypeText = (orderType: string) => {
    switch (orderType) {
      case 'limit_sell': return '限价卖单';
      case 'limit_buy': return '限价买单';
      case 'market_sell': return '市价卖单';
      case 'market_buy': return '市价买单';
      default: return orderType;
    }
  };

  const handleCancelOrder = (orderId: number) => {
    if (window.confirm('确定要取消这个订单吗？')) {
      cancelOrderMutation.mutate(orderId);
    }
  };

  if (!address) {
    return (
      <Alert severity="error">
        无效的用户地址
      </Alert>
    );
  }

  return (
    <Box>
      <Paper sx={{ p: 3, mb: 3 }}>
        <Typography variant="h4" gutterBottom>
          {isOwnProfile ? '我的资料' : '用户资料'}
        </Typography>
        <Typography variant="h6" color="text.secondary">
          地址: {address}
        </Typography>
        <Typography variant="body2" color="text.secondary">
          简化地址: {formatAddress(address)}
        </Typography>
      </Paper>

      <Box sx={{ borderBottom: 1, borderColor: 'divider' }}>
        <Tabs value={tabValue} onChange={handleTabChange}>
          <Tab label="我的订单" />
          <Tab label="我的NFT" />
          <Tab label="交易历史" />
        </Tabs>
      </Box>

      <TabPanel value={tabValue} index={0}>
        {/* 订单列表 */}
        {ordersLoading ? (
          <Box className="loading-spinner">
            <CircularProgress />
          </Box>
        ) : (
          <Grid container spacing={3}>
            {ordersData?.data?.orders?.map((orderResponse: any) => {
              const order = orderResponse.order;
              return (
                <Grid item xs={12} md={6} lg={4} key={order.id}>
                  <Card>
                    <CardContent>
                      <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 2 }}>
                        <Chip 
                          label={getOrderTypeText(order.order_type)}
                          color="primary"
                          size="small"
                        />
                        <Chip
                          label={getOrderStatusText(order.status)}
                          color={getOrderStatusColor(order.status) as any}
                          size="small"
                        />
                      </Box>
                      
                      <Typography variant="h6" gutterBottom>
                        {formatPrice(order.price)}
                      </Typography>
                      
                      <Typography variant="body2" color="text.secondary" gutterBottom>
                        NFT: #{order.token_id}
                      </Typography>
                      
                      <Typography variant="body2" color="text.secondary" gutterBottom>
                        合约: {formatAddress(order.nft_contract)}
                      </Typography>
                      
                      <Typography variant="body2" color="text.secondary" gutterBottom>
                        创建时间: {new Date(order.created_at).toLocaleDateString()}
                      </Typography>
                      
                      <Typography variant="body2" color="text.secondary" gutterBottom>
                        过期时间: {new Date(order.expiration).toLocaleDateString()}
                      </Typography>

                      {isOwnProfile && order.status === 'active' && (
                        <Button
                          variant="outlined"
                          color="error"
                          size="small"
                          fullWidth
                          sx={{ mt: 2 }}
                          onClick={() => handleCancelOrder(order.id)}
                          disabled={cancelOrderMutation.isLoading}
                        >
                          取消订单
                        </Button>
                      )}
                    </CardContent>
                  </Card>
                </Grid>
              );
            })}

            {(!ordersData?.data?.orders || ordersData.data.orders.length === 0) && (
              <Grid item xs={12}>
                <Box sx={{ textAlign: 'center', py: 8 }}>
                  <Typography variant="h6" color="text.secondary">
                    暂无订单记录
                  </Typography>
                </Box>
              </Grid>
            )}
          </Grid>
        )}
      </TabPanel>

      <TabPanel value={tabValue} index={1}>
        {/* NFT列表 */}
        {nftsLoading ? (
          <Box className="loading-spinner">
            <CircularProgress />
          </Box>
        ) : (
          <Grid container spacing={3}>
            {nftsData?.data?.nfts?.map((nft: any) => (
              <Grid item xs={12} sm={6} md={4} lg={3} key={`${nft.contract}-${nft.token_id}`}>
                <Card>
                  <Box
                    sx={{
                      height: 200,
                      bgcolor: 'grey.200',
                      display: 'flex',
                      alignItems: 'center',
                      justifyContent: 'center',
                    }}
                  >
                    {nft.image ? (
                      <img
                        src={nft.image}
                        alt={nft.name || `NFT #${nft.token_id}`}
                        style={{ width: '100%', height: '100%', objectFit: 'cover' }}
                      />
                    ) : (
                      <Typography variant="h6" color="text.secondary">
                        NFT #{nft.token_id}
                      </Typography>
                    )}
                  </Box>
                  <CardContent>
                    <Typography variant="h6" gutterBottom>
                      {nft.name || `NFT #${nft.token_id}`}
                    </Typography>
                    <Typography variant="body2" color="text.secondary" gutterBottom>
                      {nft.description}
                    </Typography>
                    <Typography variant="body2" color="text.secondary">
                      合约: {formatAddress(nft.contract)}
                    </Typography>
                  </CardContent>
                </Card>
              </Grid>
            ))}

            {(!nftsData?.data?.nfts || nftsData.data.nfts.length === 0) && (
              <Grid item xs={12}>
                <Box sx={{ textAlign: 'center', py: 8 }}>
                  <Typography variant="h6" color="text.secondary">
                    暂无NFT资产
                  </Typography>
                </Box>
              </Grid>
            )}
          </Grid>
        )}
      </TabPanel>

      <TabPanel value={tabValue} index={2}>
        {/* 交易历史 */}
        <Box sx={{ textAlign: 'center', py: 8 }}>
          <Typography variant="h6" color="text.secondary">
            交易历史功能正在开发中...
          </Typography>
        </Box>
      </TabPanel>
    </Box>
  );
};

export default UserProfilePage;
