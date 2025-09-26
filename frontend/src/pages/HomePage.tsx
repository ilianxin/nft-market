import React from 'react';
import {
  Box,
  Typography,
  Grid,
  Card,
  CardContent,
  Button,
  Container,
} from '@mui/material';
import { Link } from 'react-router-dom';
import { TrendingUp, ShowChart, AccountBalanceWallet } from '@mui/icons-material';

const HomePage: React.FC = () => {
  return (
    <Box>
      {/* Hero Section */}
      <Box className="hero-section">
        <Container>
          <Typography variant="h2" component="h1" gutterBottom>
            欢迎来到NFT市场
          </Typography>
          <Typography variant="h5" component="h2" gutterBottom>
            基于区块链的订单簿模型数字藏品交易平台
          </Typography>
          <Box sx={{ mt: 4, display: 'flex', gap: 2, justifyContent: 'center' }}>
            <Button
              variant="contained"
              size="large"
              component={Link}
              to="/marketplace"
              sx={{ bgcolor: 'white', color: 'primary.main' }}
            >
              浏览市场
            </Button>
            <Button
              variant="outlined"
              size="large"
              component={Link}
              to="/create-order"
              sx={{ color: 'white', borderColor: 'white' }}
            >
              创建订单
            </Button>
          </Box>
        </Container>
      </Box>

      {/* Features Section */}
      <Container sx={{ py: 8 }}>
        <Typography variant="h3" component="h2" textAlign="center" gutterBottom>
          平台特色
        </Typography>
        <Typography variant="body1" textAlign="center" sx={{ mb: 6 }}>
          基于智能合约的安全、透明、去中心化交易体验
        </Typography>
        
        <Grid container spacing={4}>
          <Grid item xs={12} md={4}>
            <Card className="stats-card">
              <CardContent>
                <TrendingUp sx={{ fontSize: 48, color: 'primary.main', mb: 2 }} />
                <Typography variant="h5" component="h3" gutterBottom>
                  订单簿交易
                </Typography>
                <Typography variant="body1" color="text.secondary">
                  支持限价和市价订单，提供专业的交易体验，精确控制买卖价格
                </Typography>
              </CardContent>
            </Card>
          </Grid>
          
          <Grid item xs={12} md={4}>
            <Card className="stats-card">
              <CardContent>
                <ShowChart sx={{ fontSize: 48, color: 'primary.main', mb: 2 }} />
                <Typography variant="h5" component="h3" gutterBottom>
                  智能合约
                </Typography>
                <Typography variant="body1" color="text.secondary">
                  基于以太坊智能合约，确保交易安全、透明，无需信任第三方
                </Typography>
              </CardContent>
            </Card>
          </Grid>
          
          <Grid item xs={12} md={4}>
            <Card className="stats-card">
              <CardContent>
                <AccountBalanceWallet sx={{ fontSize: 48, color: 'primary.main', mb: 2 }} />
                <Typography variant="h5" component="h3" gutterBottom>
                  钱包集成
                </Typography>
                <Typography variant="body1" color="text.secondary">
                  支持MetaMask等主流钱包，轻松管理您的数字资产和NFT
                </Typography>
              </CardContent>
            </Card>
          </Grid>
        </Grid>
      </Container>

      {/* Trading Features Section */}
      <Box sx={{ bgcolor: 'grey.100', py: 8 }}>
        <Container>
          <Typography variant="h3" component="h2" textAlign="center" gutterBottom>
            交易功能
          </Typography>
          
          <Grid container spacing={4} sx={{ mt: 4 }}>
            <Grid item xs={12} md={6}>
              <Box sx={{ p: 3 }}>
                <Typography variant="h5" gutterBottom>
                  限价订单
                </Typography>
                <Typography variant="body1" paragraph>
                  • 限价买单：设定理想价格等待卖家成交
                </Typography>
                <Typography variant="body1" paragraph>
                  • 限价卖单：设定期望价格等待买家成交
                </Typography>
                <Typography variant="body1">
                  精确控制交易价格，适合专业交易者
                </Typography>
              </Box>
            </Grid>
            
            <Grid item xs={12} md={6}>
              <Box sx={{ p: 3 }}>
                <Typography variant="h5" gutterBottom>
                  市价订单
                </Typography>
                <Typography variant="body1" paragraph>
                  • 市价买单：立即以最优价格购买NFT
                </Typography>
                <Typography variant="body1" paragraph>
                  • 市价卖单：立即以最优价格出售NFT
                </Typography>
                <Typography variant="body1">
                  快速成交，适合急需交易的用户
                </Typography>
              </Box>
            </Grid>
            
            <Grid item xs={12} md={6}>
              <Box sx={{ p: 3 }}>
                <Typography variant="h5" gutterBottom>
                  订单管理
                </Typography>
                <Typography variant="body1" paragraph>
                  • 随时取消未成交的订单
                </Typography>
                <Typography variant="body1" paragraph>
                  • 编辑订单价格和过期时间
                </Typography>
                <Typography variant="body1">
                  灵活管理您的交易策略
                </Typography>
              </Box>
            </Grid>
            
            <Grid item xs={12} md={6}>
              <Box sx={{ p: 3 }}>
                <Typography variant="h5" gutterBottom>
                  历史查询
                </Typography>
                <Typography variant="body1" paragraph>
                  • 查看所有历史订单记录
                </Typography>
                <Typography variant="body1" paragraph>
                  • 包括已过期和已取消的订单
                </Typography>
                <Typography variant="body1">
                  完整的交易历史追踪
                </Typography>
              </Box>
            </Grid>
          </Grid>
        </Container>
      </Box>
    </Box>
  );
};

export default HomePage;
