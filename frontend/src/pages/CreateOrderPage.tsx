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
  nftContract: string;
  tokenId: string;
  price: string;
  orderType: string;
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
      nftContract: '',
      tokenId: '',
      price: '',
      orderType: 'limit_sell',
      expirationDays: 7,
    },
  });

  const orderType = watch('orderType');

  const onSubmit = async (data: OrderFormData) => {
    if (!isConnected || !account) {
      setError('请先连接钱包');
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
        nft_contract: data.nftContract,
        token_id: data.tokenId,
        price: data.price || '0', // 买单价格通过发送ETH设置
        order_type: data.orderType,
        expiration: expirationDate.toISOString(),
        signature: '', // 实际应用中需要用户签名
      };

      // 保存用户地址到本地存储
      localStorage.setItem('userAddress', account);

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
                          <MenuItem value="limit_sell">限价卖单</MenuItem>
                          <MenuItem value="limit_buy">限价买单</MenuItem>
                          <MenuItem value="market_sell">市价卖单</MenuItem>
                          <MenuItem value="market_buy">市价买单</MenuItem>
                        </Select>
                      </FormControl>
                    )}
                  />
                </Grid>

                <Grid item xs={12} sm={6}>
                  <Controller
                    name="nftContract"
                    control={control}
                    rules={{ 
                      required: 'NFT合约地址不能为空',
                      pattern: {
                        value: /^0x[a-fA-F0-9]{40}$/,
                        message: '请输入有效的以太坊地址'
                      }
                    }}
                    render={({ field }) => (
                      <TextField
                        {...field}
                        fullWidth
                        label="NFT合约地址"
                        placeholder="0x..."
                        error={!!errors.nftContract}
                        helperText={errors.nftContract?.message}
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

                {(orderType === 'limit_sell' || orderType === 'market_sell') && (
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

                {(orderType === 'limit_buy') && (
                  <Grid item xs={12} sm={6}>
                    <Controller
                      name="price"
                      control={control}
                      rules={{ 
                        required: '出价不能为空',
                        pattern: {
                          value: /^\d*\.?\d+$/,
                          message: '请输入有效的出价'
                        }
                      }}
                      render={({ field }) => (
                        <TextField
                          {...field}
                          fullWidth
                          label="出价 (ETH)"
                          placeholder="0.1"
                          type="number"
                          inputProps={{ step: '0.001', min: '0' }}
                          error={!!errors.price}
                          helperText={errors.price?.message || '创建买单时需要预先锁定这些ETH'}
                        />
                      )}
                    />
                  </Grid>
                )}

                {(orderType === 'limit_sell' || orderType === 'limit_buy') && (
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
                限价卖单
              </Typography>
              <Typography variant="body2" color="text.secondary">
                设定期望售价，等待买家以该价格购买您的NFT
              </Typography>
            </Box>

            <Box sx={{ mb: 2 }}>
              <Typography variant="subtitle2" gutterBottom>
                限价买单
              </Typography>
              <Typography variant="body2" color="text.secondary">
                设定出价并锁定ETH，等待卖家接受您的报价
              </Typography>
            </Box>

            <Box sx={{ mb: 2 }}>
              <Typography variant="subtitle2" gutterBottom>
                市价卖单
              </Typography>
              <Typography variant="body2" color="text.secondary">
                立即以当前最高买价出售您的NFT
              </Typography>
            </Box>

            <Box>
              <Typography variant="subtitle2" gutterBottom>
                市价买单
              </Typography>
              <Typography variant="body2" color="text.secondary">
                立即以当前最低售价购买指定的NFT
              </Typography>
            </Box>
          </Paper>

          <Paper sx={{ p: 3, mt: 2 }}>
            <Typography variant="h6" gutterBottom>
              注意事项
            </Typography>
            <Typography variant="body2" color="text.secondary" paragraph>
              • 创建卖单前请确保已授权NFT合约
            </Typography>
            <Typography variant="body2" color="text.secondary" paragraph>
              • 创建买单时ETH将被锁定直到订单成交或取消
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
