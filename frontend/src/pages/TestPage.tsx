import React, { useState } from 'react';
import {
  Container,
  Typography,
  Box,
  Button,
  TextField,
  Card,
  CardContent,
  Alert,
  Divider,
  Grid,
  Paper,
  CircularProgress,
  Chip
} from '@mui/material';
import { apiService } from '../services/api';
import logger from '../utils/logger';

interface TestResult {
  success: boolean;
  data?: any;
  error?: string;
  timestamp: string;
}

const TestPage: React.FC = () => {
  const [loading, setLoading] = useState(false);
  const [results, setResults] = useState<{ [key: string]: TestResult }>({});
  
  // 测试用的表单数据
  const [testOrderData, setTestOrderData] = useState({
    collection_address: '0x1234567890123456789012345678901234567890',
    token_id: '123',
    order_type: 1, // 1 = listing
    price: 0.5,
    expire_time: Math.floor(Date.now() / 1000) + (7 * 24 * 60 * 60), // 7天后
    quantity_remaining: 1,
    size: 1,
    currency_address: '0x0000000000000000000000000000000000000000'
  });

  const addResult = (testName: string, result: TestResult) => {
    setResults(prev => ({
      ...prev,
      [testName]: result
    }));
  };

  const runTest = async (testName: string, testFunction: () => Promise<any>) => {
    setLoading(true);
    try {
      logger.info(`开始测试: ${testName}`);
      const data = await testFunction();
      addResult(testName, {
        success: true,
        data,
        timestamp: new Date().toLocaleString()
      });
      logger.info(`测试成功: ${testName}`, data);
    } catch (error: any) {
      addResult(testName, {
        success: false,
        error: error.message || String(error),
        timestamp: new Date().toLocaleString()
      });
      logger.error(`测试失败: ${testName}`, error);
    } finally {
      setLoading(false);
    }
  };

  // 设置测试用户地址
  const setupTestUser = () => {
    const testAddress = '0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266';
    localStorage.setItem('userAddress', testAddress);
    addResult('设置测试用户', {
      success: true,
      data: { address: testAddress },
      timestamp: new Date().toLocaleString()
    });
  };

  const testAPIs = [
    {
      name: '获取订单列表',
      test: () => apiService.getOrders()
    },
    {
      name: '创建订单',
      test: () => apiService.createOrder(testOrderData)
    },
    {
      name: '购买订单',
      test: async () => {
        // 先获取订单列表，找到一个可购买的订单
        const ordersResponse = await apiService.getOrders();
        const orders = ordersResponse?.data?.orders || [];
        const buyableOrder = orders.find((order: any) => 
          order.order_status === 0 && 
          order.order_type === 1 && 
          order.maker !== '0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266'
        );
        
        if (!buyableOrder) {
          // 如果没有可购买的订单，创建一个测试订单
          await apiService.createOrder({
            ...testOrderData,
            order_type: 1, // listing
            price: 0.1
          });
          throw new Error('没有找到可购买的订单，已创建测试订单，请重新测试购买功能');
        }
        
        return apiService.purchaseOrder(buyableOrder.id, buyableOrder.price);
      }
    },
    {
      name: '获取物品列表',
      test: () => apiService.getItems()
    },
    {
      name: '获取集合列表',
      test: () => apiService.getCollections()
    },
    {
      name: '获取活动列表',
      test: () => apiService.getActivities()
    },
    {
      name: '获取区块链状态',
      test: () => fetch('http://localhost:8080/api/v1/blockchain/status').then(r => r.json())
    }
  ];

  const clearResults = () => {
    setResults({});
  };

  const renderResult = (testName: string, result: TestResult) => (
    <Card key={testName} sx={{ mb: 2 }}>
      <CardContent>
        <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
          <Typography variant="h6" sx={{ flexGrow: 1 }}>
            {testName}
          </Typography>
          <Chip 
            label={result.success ? '成功' : '失败'} 
            color={result.success ? 'success' : 'error'}
            size="small"
          />
        </Box>
        <Typography variant="caption" color="text.secondary" sx={{ display: 'block', mb: 1 }}>
          {result.timestamp}
        </Typography>
        {result.success ? (
          <Box>
            <Typography variant="body2" color="success.main" sx={{ mb: 1 }}>
              ✅ 测试通过
            </Typography>
            <Paper sx={{ p: 2, backgroundColor: '#f5f5f5' }}>
              <pre style={{ margin: 0, fontSize: '12px', overflow: 'auto' }}>
                {JSON.stringify(result.data, null, 2)}
              </pre>
            </Paper>
          </Box>
        ) : (
          <Alert severity="error">
            <Typography variant="body2">
              {result.error}
            </Typography>
          </Alert>
        )}
      </CardContent>
    </Card>
  );

  return (
    <Box sx={{ height: '100vh', display: 'flex', flexDirection: 'column' }}>
      <Container maxWidth="lg" sx={{ py: 2, flex: 1, display: 'flex', flexDirection: 'column' }}>
        <Typography variant="h4" gutterBottom>
          🧪 API测试页面
        </Typography>
        
        <Alert severity="info" sx={{ mb: 2 }}>
          这个页面可以在不连接MetaMask的情况下测试所有API功能。
          使用Hardhat测试账户进行模拟操作。
        </Alert>

        <Grid container spacing={2} sx={{ flex: 1, minHeight: 0 }}>
          <Grid item xs={12} md={4}>
            <Card sx={{ height: 'fit-content' }}>
              <CardContent>
              <Typography variant="h6" gutterBottom>
                🛠️ 测试控制
              </Typography>
              
              <Box sx={{ mb: 2 }}>
                <Button 
                  variant="contained" 
                  onClick={setupTestUser}
                  fullWidth
                  sx={{ mb: 1 }}
                >
                  设置测试用户
                </Button>
                <Typography variant="caption" color="text.secondary">
                  设置Hardhat测试账户地址
                </Typography>
              </Box>

              <Divider sx={{ my: 2 }} />

              <Typography variant="subtitle2" gutterBottom>
                快速测试
              </Typography>
              {testAPIs.map((api) => (
                <Button
                  key={api.name}
                  variant="outlined"
                  onClick={() => runTest(api.name, api.test)}
                  disabled={loading}
                  fullWidth
                  sx={{ mb: 1 }}
                >
                  {loading ? <CircularProgress size={20} /> : api.name}
                </Button>
              ))}

              <Divider sx={{ my: 2 }} />

              <Button 
                variant="text" 
                color="secondary"
                onClick={clearResults}
                fullWidth
              >
                清除结果
              </Button>
            </CardContent>
          </Card>

            <Card sx={{ mt: 2 }}>
              <CardContent>
                <Typography variant="h6" gutterBottom>
                  📝 测试订单数据
                </Typography>
                <TextField
                  label="合约地址"
                  value={testOrderData.collection_address}
                  onChange={(e) => setTestOrderData(prev => ({ ...prev, collection_address: e.target.value }))}
                  fullWidth
                  margin="dense"
                  size="small"
                />
                <TextField
                  label="Token ID"
                  value={testOrderData.token_id}
                  onChange={(e) => setTestOrderData(prev => ({ ...prev, token_id: e.target.value }))}
                  fullWidth
                  margin="dense"
                  size="small"
                />
                <TextField
                  label="价格 (ETH)"
                  type="number"
                  value={testOrderData.price}
                  onChange={(e) => setTestOrderData(prev => ({ ...prev, price: parseFloat(e.target.value) }))}
                  fullWidth
                  margin="dense"
                  size="small"
                />
                <Box sx={{ mt: 2, p: 1, bgcolor: 'info.light', borderRadius: 1 }}>
                  <Typography variant="caption" color="info.dark">
                    💡 购买测试提示：购买功能会自动寻找可购买的订单，如果没有找到会先创建一个测试订单。
                  </Typography>
                </Box>
              </CardContent>
            </Card>
          </Grid>

          <Grid item xs={12} md={8} sx={{ display: 'flex', flexDirection: 'column', minHeight: 0 }}>
            <Typography variant="h6" gutterBottom>
              📊 测试结果
            </Typography>
          
            {Object.keys(results).length === 0 ? (
              <Paper sx={{ p: 4, textAlign: 'center' }}>
                <Typography variant="body1" color="text.secondary">
                  还没有测试结果。点击左侧按钮开始测试。
                </Typography>
              </Paper>
            ) : (
              <Box sx={{ 
                flex: 1, 
                overflowY: 'auto', 
                pr: 1,
                '&::-webkit-scrollbar': {
                  width: '8px',
                },
                '&::-webkit-scrollbar-track': {
                  backgroundColor: '#f1f1f1',
                  borderRadius: '4px',
                },
                '&::-webkit-scrollbar-thumb': {
                  backgroundColor: '#c1c1c1',
                  borderRadius: '4px',
                },
                '&::-webkit-scrollbar-thumb:hover': {
                  backgroundColor: '#a8a8a8',
                },
              }}>
                {Object.entries(results).map(([testName, result]) => 
                  renderResult(testName, result)
                )}
              </Box>
            )}
          </Grid>
        </Grid>

        <Box sx={{ mt: 2, flexShrink: 0 }}>
          <Alert severity="warning">
            <Typography variant="body2">
              <strong>注意：</strong> 确保后端服务已启动 (localhost:8080) 并且本地以太坊网络正在运行 (localhost:8545)
            </Typography>
          </Alert>
        </Box>
      </Container>
    </Box>
  );
};

export default TestPage;
