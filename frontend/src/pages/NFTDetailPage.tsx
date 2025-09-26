import React, { useState } from 'react';
import {
  Box,
  Typography,
  Grid,
  Card,
  CardContent,
  Button,
  Chip,
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  CircularProgress,
  Alert,
} from '@mui/material';
import { useParams } from 'react-router-dom';
import { useQuery } from 'react-query';
import { apiService } from '../services/api';

const NFTDetailPage: React.FC = () => {
  const { contract, tokenId } = useParams<{ contract: string; tokenId: string }>();
  const [selectedOrderType, setSelectedOrderType] = useState<'sell' | 'buy'>('sell');

  // 获取NFT信息
  const { data: nftData, isLoading: nftLoading } = useQuery(
    ['nft', contract, tokenId],
    () => apiService.getNFTById(contract!, tokenId!),
    {
      enabled: !!contract && !!tokenId,
    }
  );

  // 获取NFT的所有订单
  const { data: ordersData, isLoading: ordersLoading } = useQuery(
    ['nftOrders', contract, tokenId],
    () => apiService.getNFTOrders(contract!, tokenId!),
    {
      enabled: !!contract && !!tokenId,
    }
  );

  const formatAddress = (address: string) => {
    return `${address.slice(0, 6)}...${address.slice(-4)}`;
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

  // 过滤卖单和买单
  const sellOrders = ordersData?.data?.orders?.filter((orderResponse: any) => 
    orderResponse.order.order_type === 'limit_sell' || orderResponse.order.order_type === 'market_sell'
  ) || [];
  
  const buyOrders = ordersData?.data?.orders?.filter((orderResponse: any) => 
    orderResponse.order.order_type === 'limit_buy' || orderResponse.order.order_type === 'market_buy'
  ) || [];

  // 获取最佳价格
  const bestSellPrice = sellOrders
    .filter((orderResponse: any) => orderResponse.order.status === 'active')
    .sort((a: any, b: any) => parseFloat(a.order.price) - parseFloat(b.order.price))[0];
    
  const bestBuyPrice = buyOrders
    .filter((orderResponse: any) => orderResponse.order.status === 'active')
    .sort((a: any, b: any) => parseFloat(b.order.price) - parseFloat(a.order.price))[0];

  if (!contract || !tokenId) {
    return (
      <Alert severity="error">
        无效的NFT参数
      </Alert>
    );
  }

  return (
    <Box>
      <Grid container spacing={4}>
        {/* NFT信息 */}
        <Grid item xs={12} md={6}>
          <Card>
            <Box
              sx={{
                height: 400,
                bgcolor: 'grey.200',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
              }}
            >
              {nftLoading ? (
                <CircularProgress />
              ) : nftData?.data?.image ? (
                <img
                  src={nftData.data.image}
                  alt={nftData.data.name || `NFT #${tokenId}`}
                  style={{ width: '100%', height: '100%', objectFit: 'cover' }}
                />
              ) : (
                <Typography variant="h4" color="text.secondary">
                  NFT #{tokenId}
                </Typography>
              )}
            </Box>
            <CardContent>
              <Typography variant="h5" gutterBottom>
                {nftData?.data?.name || `NFT #${tokenId}`}
              </Typography>
              <Typography variant="body1" color="text.secondary" paragraph>
                {nftData?.data?.description || '暂无描述'}
              </Typography>
              <Typography variant="body2" color="text.secondary" gutterBottom>
                合约地址: {formatAddress(contract)}
              </Typography>
              <Typography variant="body2" color="text.secondary" gutterBottom>
                Token ID: {tokenId}
              </Typography>
              {nftData?.data?.owner && (
                <Typography variant="body2" color="text.secondary">
                  当前所有者: {formatAddress(nftData.data.owner)}
                </Typography>
              )}
            </CardContent>
          </Card>
        </Grid>

        {/* 价格信息和快速交易 */}
        <Grid item xs={12} md={6}>
          <Paper sx={{ p: 3, mb: 3 }}>
            <Typography variant="h6" gutterBottom>
              当前最佳价格
            </Typography>
            
            <Grid container spacing={2}>
              <Grid item xs={6}>
                <Box sx={{ textAlign: 'center', p: 2, bgcolor: 'error.light', borderRadius: 1 }}>
                  <Typography variant="body2" color="white">
                    最低售价
                  </Typography>
                  <Typography variant="h6" color="white">
                    {bestSellPrice ? formatPrice(bestSellPrice.order.price) : '暂无'}
                  </Typography>
                </Box>
              </Grid>
              <Grid item xs={6}>
                <Box sx={{ textAlign: 'center', p: 2, bgcolor: 'success.light', borderRadius: 1 }}>
                  <Typography variant="body2" color="white">
                    最高出价
                  </Typography>
                  <Typography variant="h6" color="white">
                    {bestBuyPrice ? formatPrice(bestBuyPrice.order.price) : '暂无'}
                  </Typography>
                </Box>
              </Grid>
            </Grid>

            <Box sx={{ mt: 3, display: 'flex', gap: 2 }}>
              {bestSellPrice && (
                <Button
                  variant="contained"
                  color="success"
                  fullWidth
                  size="large"
                >
                  立即购买 {formatPrice(bestSellPrice.order.price)}
                </Button>
              )}
              {bestBuyPrice && (
                <Button
                  variant="contained"
                  color="error"
                  fullWidth
                  size="large"
                >
                  立即出售 {formatPrice(bestBuyPrice.order.price)}
                </Button>
              )}
            </Box>
          </Paper>

          {/* 订单簿预览 */}
          <Paper sx={{ p: 3 }}>
            <Typography variant="h6" gutterBottom>
              订单簿
            </Typography>
            
            <Box sx={{ display: 'flex', gap: 1, mb: 2 }}>
              <Button
                variant={selectedOrderType === 'sell' ? 'contained' : 'outlined'}
                onClick={() => setSelectedOrderType('sell')}
                color="error"
                size="small"
              >
                卖单 ({sellOrders.length})
              </Button>
              <Button
                variant={selectedOrderType === 'buy' ? 'contained' : 'outlined'}
                onClick={() => setSelectedOrderType('buy')}
                color="success"
                size="small"
              >
                买单 ({buyOrders.length})
              </Button>
            </Box>

            <Box sx={{ maxHeight: 300, overflowY: 'auto' }}>
              {ordersLoading ? (
                <Box sx={{ textAlign: 'center', py: 2 }}>
                  <CircularProgress size={24} />
                </Box>
              ) : (
                <TableContainer>
                  <Table size="small">
                    <TableHead>
                      <TableRow>
                        <TableCell>价格</TableCell>
                        <TableCell>状态</TableCell>
                        <TableCell>创建者</TableCell>
                        <TableCell>操作</TableCell>
                      </TableRow>
                    </TableHead>
                    <TableBody>
                      {(selectedOrderType === 'sell' ? sellOrders : buyOrders)
                        .slice(0, 10)
                        .map((orderResponse: any) => {
                          const order = orderResponse.order;
                          return (
                            <TableRow key={order.id}>
                              <TableCell>
                                <Typography variant="body2" fontWeight="bold">
                                  {formatPrice(order.price)}
                                </Typography>
                              </TableCell>
                              <TableCell>
                                <Chip
                                  label={getOrderStatusText(order.status)}
                                  color={getOrderStatusColor(order.status) as any}
                                  size="small"
                                />
                              </TableCell>
                              <TableCell>
                                <Typography variant="body2">
                                  {formatAddress(order.maker)}
                                </Typography>
                              </TableCell>
                              <TableCell>
                                {order.status === 'active' && (
                                  <Button
                                    size="small"
                                    variant="outlined"
                                    color={selectedOrderType === 'sell' ? 'success' : 'error'}
                                  >
                                    {selectedOrderType === 'sell' ? '购买' : '出售'}
                                  </Button>
                                )}
                              </TableCell>
                            </TableRow>
                          );
                        })}
                    </TableBody>
                  </Table>
                </TableContainer>
              )}
            </Box>
          </Paper>
        </Grid>
      </Grid>

      {/* 完整订单历史 */}
      <Paper sx={{ mt: 4, p: 3 }}>
        <Typography variant="h6" gutterBottom>
          所有订单历史
        </Typography>
        
        {ordersLoading ? (
          <Box className="loading-spinner">
            <CircularProgress />
          </Box>
        ) : (
          <TableContainer>
            <Table>
              <TableHead>
                <TableRow>
                  <TableCell>类型</TableCell>
                  <TableCell>价格</TableCell>
                  <TableCell>状态</TableCell>
                  <TableCell>创建者</TableCell>
                  <TableCell>创建时间</TableCell>
                  <TableCell>过期时间</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {ordersData?.data?.orders?.map((orderResponse: any) => {
                  const order = orderResponse.order;
                  return (
                    <TableRow key={order.id}>
                      <TableCell>
                        <Chip
                          label={order.order_type === 'limit_sell' || order.order_type === 'market_sell' ? '卖单' : '买单'}
                          color={order.order_type === 'limit_sell' || order.order_type === 'market_sell' ? 'error' : 'success'}
                          size="small"
                        />
                      </TableCell>
                      <TableCell>
                        <Typography variant="body2" fontWeight="bold">
                          {formatPrice(order.price)}
                        </Typography>
                      </TableCell>
                      <TableCell>
                        <Chip
                          label={getOrderStatusText(order.status)}
                          color={getOrderStatusColor(order.status) as any}
                          size="small"
                        />
                      </TableCell>
                      <TableCell>{formatAddress(order.maker)}</TableCell>
                      <TableCell>
                        {new Date(order.created_at).toLocaleDateString()}
                      </TableCell>
                      <TableCell>
                        {new Date(order.expiration).toLocaleDateString()}
                      </TableCell>
                    </TableRow>
                  );
                })}
              </TableBody>
            </Table>
          </TableContainer>
        )}

        {(!ordersData?.data?.orders || ordersData.data.orders.length === 0) && (
          <Box sx={{ textAlign: 'center', py: 4 }}>
            <Typography variant="body1" color="text.secondary">
              此NFT暂无订单记录
            </Typography>
          </Box>
        )}
      </Paper>
    </Box>
  );
};

export default NFTDetailPage;
