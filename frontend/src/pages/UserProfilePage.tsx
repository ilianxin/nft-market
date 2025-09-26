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

  // 获取用户物品
  const { data: itemsData, isLoading: itemsLoading } = useQuery(
    ['userItems', address],
    () => apiService.getUserItems(address!),
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

  const formatPrice = (price: number) => {
    try {
      return price.toFixed(4) + ' ETH';
    } catch {
      return price.toString();
    }
  };

  const getOrderStatusColor = (status: number) => {
    switch (status) {
      case 0: return 'success'; // 活跃
      case 1: return 'info';    // 已成交
      case 2: return 'error';   // 已取消
      case 3: return 'warning'; // 已过期
      default: return 'default';
    }
  };

  const getOrderStatusText = (status: number) => {
    switch (status) {
      case 0: return '活跃';
      case 1: return '已成交';
      case 2: return '已取消';
      case 3: return '已过期';
      default: return `状态${status}`;
    }
  };

  const getOrderTypeText = (orderType: number) => {
    switch (orderType) {
      case 1: return '上架';
      case 2: return '出价';
      case 3: return '集合出价';
      case 4: return '物品出价';
      default: return `类型${orderType}`;
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
          <Tab label="我的物品" />
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
            {ordersData?.data?.orders?.map((order: any) => {
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
                          label={getOrderStatusText(order.order_status)}
                          color={getOrderStatusColor(order.order_status) as any}
                          size="small"
                        />
                      </Box>
                      
                      <Typography variant="h6" gutterBottom>
                        {formatPrice(order.price)}
                      </Typography>
                      
                      <Typography variant="body2" color="text.secondary" gutterBottom>
                        Token ID: #{order.token_id}
                      </Typography>
                      
                      <Typography variant="body2" color="text.secondary" gutterBottom>
                        集合: {formatAddress(order.collection_address || '')}
                      </Typography>
                      
                      <Typography variant="body2" color="text.secondary" gutterBottom>
                        创建时间: {order.create_time ? new Date(order.create_time * 1000).toLocaleDateString() : '未知'}
                      </Typography>
                      
                      <Typography variant="body2" color="text.secondary" gutterBottom>
                        过期时间: {order.expire_time ? new Date(order.expire_time * 1000).toLocaleDateString() : '无'}
                      </Typography>

                      {isOwnProfile && order.order_status === 0 && (
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
        {/* 物品列表 */}
        {itemsLoading ? (
          <Box className="loading-spinner">
            <CircularProgress />
          </Box>
        ) : (
          <Grid container spacing={3}>
            {itemsData?.data?.items?.map((item: any) => (
              <Grid item xs={12} sm={6} md={4} lg={3} key={`${item.collection_address}-${item.token_id}`}>
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
                    <Typography variant="h6" color="text.secondary">
                      Item #{item.token_id}
                    </Typography>
                  </Box>
                  <CardContent>
                    <Typography variant="h6" gutterBottom>
                      {item.name || `Item #${item.token_id}`}
                    </Typography>
                    <Typography variant="body2" color="text.secondary" gutterBottom>
                      拥有者: {formatAddress(item.owner || '')}
                    </Typography>
                    <Typography variant="body2" color="text.secondary" gutterBottom>
                      集合: {formatAddress(item.collection_address || '')}
                    </Typography>
                    {item.list_price && (
                      <Typography variant="body2" color="text.secondary">
                        上架价格: {formatPrice(item.list_price)}
                      </Typography>
                    )}
                  </CardContent>
                </Card>
              </Grid>
            ))}

            {(!itemsData?.data?.items || itemsData.data.items.length === 0) && (
              <Grid item xs={12}>
                <Box sx={{ textAlign: 'center', py: 8 }}>
                  <Typography variant="h6" color="text.secondary">
                    暂无物品资产
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
