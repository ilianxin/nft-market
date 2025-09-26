import React, { useState } from 'react';
import {
  Box,
  Typography,
  Grid,
  Card,
  CardContent,
  Button,
  Tabs,
  Tab,
  Chip,
  CircularProgress,
  Pagination,
} from '@mui/material';
import { Link } from 'react-router-dom';
import { useQuery } from 'react-query';
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
      id={`marketplace-tabpanel-${index}`}
      aria-labelledby={`marketplace-tab-${index}`}
      {...other}
    >
      {value === index && <Box sx={{ pt: 3 }}>{children}</Box>}
    </div>
  );
}

const MarketplacePage: React.FC = () => {
  const [tabValue, setTabValue] = useState(0);
  const [page, setPage] = useState(1);
  const pageSize = 12;

  const handleTabChange = (event: React.SyntheticEvent, newValue: number) => {
    setTabValue(newValue);
    setPage(1);
  };

  const handlePageChange = (event: React.ChangeEvent<unknown>, value: number) => {
    setPage(value);
  };

  // 获取活跃订单
  const { data: ordersData, isLoading: ordersLoading } = useQuery(
    ['orders', tabValue, page],
    () => {
      const orderType = tabValue === 1 ? '1' : tabValue === 2 ? '2' : '';
      return apiService.getOrders(page, pageSize, orderType);
    },
    {
      keepPreviousData: true,
    }
  );

  const getOrderTypeChip = (orderType: number) => {
    const config = {
      1: { label: '上架', color: 'error' as const },
      2: { label: '出价', color: 'success' as const },
      3: { label: '集合出价', color: 'warning' as const },
      4: { label: '物品出价', color: 'info' as const },
    };
    
    return config[orderType as keyof typeof config] || { label: `类型${orderType}`, color: 'default' as const };
  };

  const formatPrice = (price: number) => {
    try {
      return price.toFixed(4) + ' ETH';
    } catch {
      return price.toString();
    }
  };

  const formatAddress = (address: string) => {
    return `${address.slice(0, 6)}...${address.slice(-4)}`;
  };

  return (
    <Box>
      <Typography variant="h4" component="h1" gutterBottom>
        NFT市场
      </Typography>
      <Typography variant="body1" color="text.secondary" sx={{ mb: 3 }}>
        浏览和交易所有活跃的NFT订单
      </Typography>

      <Box sx={{ borderBottom: 1, borderColor: 'divider' }}>
        <Tabs value={tabValue} onChange={handleTabChange}>
          <Tab label="全部订单" />
          <Tab label="上架订单" />
          <Tab label="出价订单" />
        </Tabs>
      </Box>

      <TabPanel value={tabValue} index={0}>
        {ordersLoading ? (
          <Box className="loading-spinner">
            <CircularProgress />
          </Box>
        ) : (
          <>
            <Grid container spacing={3}>
              {ordersData?.data?.orders?.map((order: any) => {
                const chipConfig = getOrderTypeChip(order.order_type);
                
                return (
                  <Grid item xs={12} sm={6} md={4} lg={3} key={order.id}>
                    <Card className="nft-card">
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
                          NFT #{order.token_id}
                        </Typography>
                      </Box>
                      <CardContent>
                        <Box sx={{ mb: 1 }}>
                          <Chip 
                            label={chipConfig.label} 
                            color={chipConfig.color}
                            size="small"
                          />
                        </Box>
                        <Typography variant="h6" gutterBottom>
                          {formatPrice(order.price)}
                        </Typography>
                        <Typography variant="body2" color="text.secondary" gutterBottom>
                          集合: {formatAddress(order.collection_address || '')}
                        </Typography>
                        <Typography variant="body2" color="text.secondary" gutterBottom>
                          创建者: {formatAddress(order.maker || '')}
                        </Typography>
                        <Typography variant="body2" color="text.secondary" gutterBottom>
                          过期时间: {order.expire_time ? new Date(order.expire_time * 1000).toLocaleDateString() : '无'}
                        </Typography>
                        <Button
                          fullWidth
                          variant="contained"
                          component={Link}
                          to={`/nft/${order.collection_address}/${order.token_id}`}
                          sx={{ mt: 2 }}
                        >
                          查看详情
                        </Button>
                      </CardContent>
                    </Card>
                  </Grid>
                );
              })}
            </Grid>

            {ordersData?.data?.total_pages > 1 && (
              <Box sx={{ display: 'flex', justifyContent: 'center', mt: 4 }}>
                <Pagination
                  count={ordersData.data.total_pages}
                  page={page}
                  onChange={handlePageChange}
                  color="primary"
                />
              </Box>
            )}

            {(!ordersData?.data?.orders || ordersData.data.orders.length === 0) && (
              <Box sx={{ textAlign: 'center', py: 8 }}>
                <Typography variant="h6" color="text.secondary">
                  暂无活跃订单
                </Typography>
                <Button
                  variant="contained"
                  component={Link}
                  to="/create-order"
                  sx={{ mt: 2 }}
                >
                  创建第一个订单
                </Button>
              </Box>
            )}
          </>
        )}
      </TabPanel>

      <TabPanel value={tabValue} index={1}>
        {/* 上架订单内容 */}
        {ordersLoading ? (
          <Box className="loading-spinner">
            <CircularProgress />
          </Box>
        ) : (
          <>
            <Grid container spacing={3}>
              {ordersData?.data?.orders?.map((order: any) => {
                return (
                  <Grid item xs={12} sm={6} md={4} lg={3} key={order.id}>
                    <Card className="nft-card">
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
                          NFT #{order.token_id}
                        </Typography>
                      </Box>
                      <CardContent>
                        <Box sx={{ mb: 1 }}>
                          <Chip label="上架" color="error" size="small" />
                        </Box>
                        <Typography variant="h6" gutterBottom>
                          {formatPrice(order.price)}
                        </Typography>
                        <Typography variant="body2" color="text.secondary" gutterBottom>
                          集合: {formatAddress(order.collection_address || '')}
                        </Typography>
                        <Typography variant="body2" color="text.secondary" gutterBottom>
                          卖家: {formatAddress(order.maker || '')}
                        </Typography>
                        <Button
                          fullWidth
                          variant="contained"
                          color="success"
                          component={Link}
                          to={`/nft/${order.collection_address}/${order.token_id}`}
                          sx={{ mt: 2 }}
                        >
                          立即购买
                        </Button>
                      </CardContent>
                    </Card>
                  </Grid>
                );
              })}
            </Grid>
          </>
        )}
      </TabPanel>

      <TabPanel value={tabValue} index={2}>
        {/* 出价订单内容 */}
        {ordersLoading ? (
          <Box className="loading-spinner">
            <CircularProgress />
          </Box>
        ) : (
          <>
            <Grid container spacing={3}>
              {ordersData?.data?.orders?.map((order: any) => {
                return (
                  <Grid item xs={12} sm={6} md={4} lg={3} key={order.id}>
                    <Card className="nft-card">
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
                          NFT #{order.token_id}
                        </Typography>
                      </Box>
                      <CardContent>
                        <Box sx={{ mb: 1 }}>
                          <Chip label="出价" color="success" size="small" />
                        </Box>
                        <Typography variant="h6" gutterBottom>
                          {formatPrice(order.price)}
                        </Typography>
                        <Typography variant="body2" color="text.secondary" gutterBottom>
                          集合: {formatAddress(order.collection_address || '')}
                        </Typography>
                        <Typography variant="body2" color="text.secondary" gutterBottom>
                          买家: {formatAddress(order.maker || '')}
                        </Typography>
                        <Button
                          fullWidth
                          variant="contained"
                          color="error"
                          component={Link}
                          to={`/nft/${order.collection_address}/${order.token_id}`}
                          sx={{ mt: 2 }}
                        >
                          立即出售
                        </Button>
                      </CardContent>
                    </Card>
                  </Grid>
                );
              })}
            </Grid>
          </>
        )}
      </TabPanel>
    </Box>
  );
};

export default MarketplacePage;
