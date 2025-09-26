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

  // 获取物品信息
  const { data: itemData, isLoading: itemLoading } = useQuery(
    ['item', contract, tokenId],
    () => apiService.getItemByTokenId(contract!, tokenId!),
    {
      enabled: !!contract && !!tokenId,
    }
  );

  // 获取物品的所有订单
  const { data: ordersData, isLoading: ordersLoading } = useQuery(
    ['itemOrders', contract, tokenId],
    () => apiService.getNFTOrders(contract!, tokenId!),
    {
      enabled: !!contract && !!tokenId,
    }
  );

  const formatAddress = (address: string) => {
    return `${address.slice(0, 6)}...${address.slice(-4)}`;
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

  // 过滤卖单和买单
  const sellOrders = ordersData?.data?.orders?.filter((order: any) => 
    order.order_type === 1 // 上架
  ) || [];
  
  const buyOrders = ordersData?.data?.orders?.filter((order: any) => 
    order.order_type === 2 // 出价
  ) || [];

  // 获取最佳价格
  const bestSellPrice = sellOrders
    .filter((order: any) => order.order_status === 0) // 活跃状态
    .sort((a: any, b: any) => a.price - b.price)[0];
    
  const bestBuyPrice = buyOrders
    .filter((order: any) => order.order_status === 0) // 活跃状态
    .sort((a: any, b: any) => b.price - a.price)[0];

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
              {itemLoading ? (
                <CircularProgress />
              ) : (
                <Typography variant="h4" color="text.secondary">
                  Item #{tokenId}
                </Typography>
              )}
            </Box>
            <CardContent>
              <Typography variant="h5" gutterBottom>
                {itemData?.data?.item?.name || `Item #${tokenId}`}
              </Typography>
              <Typography variant="body1" color="text.secondary" paragraph>
                暂无描述
              </Typography>
              <Typography variant="body2" color="text.secondary" gutterBottom>
                集合地址: {formatAddress(contract)}
              </Typography>
              <Typography variant="body2" color="text.secondary" gutterBottom>
                Token ID: {tokenId}
              </Typography>
              {itemData?.data?.item?.owner && (
                <Typography variant="body2" color="text.secondary">
                  当前所有者: {formatAddress(itemData.data.item.owner)}
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
                    {bestSellPrice ? formatPrice(bestSellPrice.price) : '暂无'}
                  </Typography>
                </Box>
              </Grid>
              <Grid item xs={6}>
                <Box sx={{ textAlign: 'center', p: 2, bgcolor: 'success.light', borderRadius: 1 }}>
                  <Typography variant="body2" color="white">
                    最高出价
                  </Typography>
                  <Typography variant="h6" color="white">
                    {bestBuyPrice ? formatPrice(bestBuyPrice.price) : '暂无'}
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
                  立即购买 {formatPrice(bestSellPrice.price)}
                </Button>
              )}
              {bestBuyPrice && (
                <Button
                  variant="contained"
                  color="error"
                  fullWidth
                  size="large"
                >
                  立即出售 {formatPrice(bestBuyPrice.price)}
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
                        .map((order: any) => {
                          return (
                            <TableRow key={order.id}>
                              <TableCell>
                                <Typography variant="body2" fontWeight="bold">
                                  {formatPrice(order.price)}
                                </Typography>
                              </TableCell>
                              <TableCell>
                                <Chip
                                  label={getOrderStatusText(order.order_status)}
                                  color={getOrderStatusColor(order.order_status) as any}
                                  size="small"
                                />
                              </TableCell>
                              <TableCell>
                                <Typography variant="body2">
                                  {formatAddress(order.maker || '')}
                                </Typography>
                              </TableCell>
                              <TableCell>
                                {order.order_status === 0 && (
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
                {ordersData?.data?.orders?.map((order: any) => {
                  return (
                    <TableRow key={order.id}>
                      <TableCell>
                        <Chip
                          label={order.order_type === 1 ? '上架' : '出价'}
                          color={order.order_type === 1 ? 'error' : 'success'}
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
                          label={getOrderStatusText(order.order_status)}
                          color={getOrderStatusColor(order.order_status) as any}
                          size="small"
                        />
                      </TableCell>
                      <TableCell>{formatAddress(order.maker || '')}</TableCell>
                      <TableCell>
                        {order.create_time ? new Date(order.create_time * 1000).toLocaleDateString() : '未知'}
                      </TableCell>
                      <TableCell>
                        {order.expire_time ? new Date(order.expire_time * 1000).toLocaleDateString() : '无'}
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
