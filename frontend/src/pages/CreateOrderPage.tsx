import React, { useState } from 'react';
import {
  Box,
  Typography,
  Paper,
  TextField,
  Button,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Alert,
  CircularProgress,
  Grid,
} from '@mui/material';
import { useForm, Controller } from 'react-hook-form';
import { useWeb3 } from '../contexts/Web3Context';
import { apiService } from '../services/api';

interface OrderFormData {
  collectionAddress: string;
  tokenId: string;
  price: string;
  orderType: number;
  expirationDays: number;
}

const CreateOrderPage: React.FC = () => {
  const { account, isConnected } = useWeb3();
  const [loading, setLoading] = useState(false);
  const [success, setSuccess] = useState(false);
  const [error, setError] = useState<string>('');

  const {
    control,
    handleSubmit,
    formState: { errors },
    reset,
    watch
  } = useForm<OrderFormData>({
    defaultValues: {
      collectionAddress: '',
      tokenId: '',
      price: '',
      orderType: 1, // 1: 上架
      expirationDays: 7,
    },
  });

  const orderType = watch('orderType');

  const onSubmit = async (data: OrderFormData) => {
    console.log('开始创建订单，表单数据:', data);
    console.log('钱包连接状态:', { isConnected, account });
    
    if (!isConnected || !account) {
      setError('请先连接钱包');
      console.error('钱包未连接');
      return;
    }

    setLoading(true);
    setError('');
    setSuccess(false);

    try {
      // 计算过期时间
      const expirationDate = new Date();
      expirationDate.setDate(expirationDate.getDate() + data.expirationDays);

      const orderData = {
        collection_address: data.collectionAddress,
        token_id: data.tokenId,
        price: parseFloat(data.price) || 0, // 转换为数字
        order_type: data.orderType,
        expire_time: Math.floor(expirationDate.getTime() / 1000), // 转换为Unix时间戳
        quantity_remaining: 1,
        size: 1,
        currency_address: '0x0000000000000000000000000000000000000000', // ETH地址
      };

      // 保存用户地址到本地存储
      localStorage.setItem('userAddress', account);
      console.log('准备发送订单数据:', orderData);

      const response = await apiService.createOrder(orderData);
      console.log('订单创建成功:', response);
      
      setSuccess(true);
      reset();
    } catch (err: any) {
      console.error('创建订单失败:', err);
      setError(err.response?.data?.message || '创建订单失败，请重试');
    } finally {
      setLoading(false);
    }
  };

  if (!isConnected) {
    return (
      <Box>
        <Typography variant="h4" component="h1" gutterBottom>
          创建订单
        </Typography>
        <Alert severity="warning" sx={{ mt: 2 }}>
          请先连接您的钱包以创建订单
        </Alert>
      </Box>
    );
  }

  return (
    <Box>
      <Typography variant="h4" component="h1" gutterBottom>
        创建订单
      </Typography>
      <Typography variant="body1" color="text.secondary" sx={{ mb: 3 }}>
        在NFT市场上创建买入或卖出订单
      </Typography>

      <Grid container spacing={4}>
        <Grid item xs={12} md={8}>
          <Paper className="order-form">
            <Typography variant="h6" gutterBottom>
              订单信息
            </Typography>

            {error && (
              <Alert severity="error" sx={{ mb: 2 }}>
                {error}
              </Alert>
            )}

            {success && (
              <Alert severity="success" sx={{ mb: 2 }}>
                订单创建成功！您可以在用户资料页查看订单状态。
              </Alert>
            )}

            <form onSubmit={handleSubmit(onSubmit)}>
              <Grid container spacing={3}>
                <Grid item xs={12}>
                  <Controller
                    name="orderType"
                    control={control}
                    rules={{ required: '请选择订单类型' }}
                    render={({ field }) => (
                      <FormControl fullWidth error={!!errors.orderType}>
                        <InputLabel>订单类型</InputLabel>
                        <Select {...field} label="订单类型">
                          <MenuItem value={1}>上架</MenuItem>
                          <MenuItem value={2}>出价</MenuItem>
                          <MenuItem value={3}>集合出价</MenuItem>
                          <MenuItem value={4}>物品出价</MenuItem>
                        </Select>
                      </FormControl>
                    )}
                  />
                </Grid>

                <Grid item xs={12} sm={6}>
                  <Controller
                    name="collectionAddress"
                    control={control}
                    rules={{ 
                      required: '集合地址不能为空',
                      pattern: {
                        value: /^0x[a-fA-F0-9]{40}$/,
                        message: '请输入有效的以太坊地址'
                      }
                    }}
                    render={({ field }) => (
                      <TextField
                        {...field}
                        fullWidth
                        label="集合地址"
                        placeholder="0x..."
                        error={!!errors.collectionAddress}
                        helperText={errors.collectionAddress?.message}
                      />
                    )}
                  />
                </Grid>

                <Grid item xs={12} sm={6}>
                  <Controller
                    name="tokenId"
                    control={control}
                    rules={{ required: 'Token ID不能为空' }}
                    render={({ field }) => (
                      <TextField
                        {...field}
                        fullWidth
                        label="Token ID"
                        placeholder="1"
                        error={!!errors.tokenId}
                        helperText={errors.tokenId?.message}
                      />
                    )}
                  />
                </Grid>

                {(orderType === 1 || orderType === 2) && (
                  <Grid item xs={12} sm={6}>
                    <Controller
                      name="price"
                      control={control}
                      rules={{ 
                        required: '价格不能为空',
                        pattern: {
                          value: /^\d*\.?\d+$/,
                          message: '请输入有效的价格'
                        }
                      }}
                      render={({ field }) => (
                        <TextField
                          {...field}
                          fullWidth
                          label="价格 (ETH)"
                          placeholder="0.1"
                          type="number"
                          inputProps={{ step: '0.001', min: '0' }}
                          error={!!errors.price}
                          helperText={errors.price?.message}
                        />
                      )}
                    />
                  </Grid>
                )}

                {(orderType === 1 || orderType === 2) && (
                  <Grid item xs={12} sm={6}>
                    <Controller
                      name="expirationDays"
                      control={control}
                      rules={{ 
                        required: '过期天数不能为空',
                        min: { value: 1, message: '至少1天' },
                        max: { value: 365, message: '最多365天' }
                      }}
                      render={({ field }) => (
                        <TextField
                          {...field}
                          fullWidth
                          label="过期天数"
                          type="number"
                          inputProps={{ min: 1, max: 365 }}
                          error={!!errors.expirationDays}
                          helperText={errors.expirationDays?.message}
                        />
                      )}
                    />
                  </Grid>
                )}

                <Grid item xs={12}>
                  <Button
                    type="submit"
                    variant="contained"
                    size="large"
                    fullWidth
                    disabled={loading}
                    startIcon={loading ? <CircularProgress size={20} /> : null}
                  >
                    {loading ? '创建中...' : '创建订单'}
                  </Button>
                </Grid>
              </Grid>
            </form>
          </Paper>
        </Grid>

        <Grid item xs={12} md={4}>
          <Paper sx={{ p: 3 }}>
            <Typography variant="h6" gutterBottom>
              订单类型说明
            </Typography>
            
            <Box sx={{ mb: 2 }}>
              <Typography variant="subtitle2" gutterBottom>
                上架
              </Typography>
              <Typography variant="body2" color="text.secondary">
                将您的NFT上架到市场，设定期望售价等待买家购买
              </Typography>
            </Box>

            <Box sx={{ mb: 2 }}>
              <Typography variant="subtitle2" gutterBottom>
                出价
              </Typography>
              <Typography variant="body2" color="text.secondary">
                对特定NFT出价，等待卖家接受您的报价
              </Typography>
            </Box>

            <Box sx={{ mb: 2 }}>
              <Typography variant="subtitle2" gutterBottom>
                集合出价
              </Typography>
              <Typography variant="body2" color="text.secondary">
                对整个集合出价，购买集合中的任意NFT
              </Typography>
            </Box>

            <Box>
              <Typography variant="subtitle2" gutterBottom>
                物品出价
              </Typography>
              <Typography variant="body2" color="text.secondary">
                对特定物品出价，等待卖家接受
              </Typography>
            </Box>
          </Paper>

          <Paper sx={{ p: 3, mt: 2 }}>
            <Typography variant="h6" gutterBottom>
              注意事项
            </Typography>
            <Typography variant="body2" color="text.secondary" paragraph>
              • 创建上架订单前请确保已授权集合合约
            </Typography>
            <Typography variant="body2" color="text.secondary" paragraph>
              • 创建出价订单时ETH将被锁定直到订单成交或取消
            </Typography>
            <Typography variant="body2" color="text.secondary">
              • 所有交易都通过智能合约执行，确保安全透明
            </Typography>
          </Paper>
        </Grid>
      </Grid>
    </Box>
  );
};

export default CreateOrderPage;
