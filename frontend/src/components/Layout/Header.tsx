import React from 'react';
import {
  AppBar,
  Toolbar,
  Typography,
  Button,
  Box,
  IconButton,
  Menu,
  MenuItem,
} from '@mui/material';
import { Link, useNavigate } from 'react-router-dom';
import { AccountCircle, Wallet } from '@mui/icons-material';
import { useWeb3 } from '../../contexts/Web3Context';

const Header: React.FC = () => {
  const navigate = useNavigate();
  const { account, isConnected, connectWallet, connectTestWallet, disconnectWallet, isTestMode } = useWeb3();
  const [anchorEl, setAnchorEl] = React.useState<null | HTMLElement>(null);

  const handleMenuOpen = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorEl(event.currentTarget);
  };

  const handleMenuClose = () => {
    setAnchorEl(null);
  };

  const handleProfileClick = () => {
    if (account) {
      navigate(`/profile/${account}`);
    }
    handleMenuClose();
  };

  const handleDisconnect = () => {
    disconnectWallet();
    handleMenuClose();
  };

  const formatAddress = (address: string) => {
    return `${address.slice(0, 6)}...${address.slice(-4)}`;
  };

  return (
    <AppBar position="static">
      <Toolbar>
        <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
          <Link to="/" style={{ textDecoration: 'none', color: 'inherit' }}>
            NFT市场
          </Link>
        </Typography>
        
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
          <Button color="inherit" component={Link} to="/">
            首页
          </Button>
          <Button color="inherit" component={Link} to="/marketplace">
            市场
          </Button>
          <Button color="inherit" component={Link} to="/collections">
            集合
          </Button>
          <Button color="inherit" component={Link} to="/activities">
            活动
          </Button>
          <Button color="inherit" component={Link} to="/create-order">
            创建订单
          </Button>
          <Button 
            color="inherit" 
            component={Link} 
            to="/test"
            sx={{ 
              color: '#ff9800',
              '&:hover': { 
                backgroundColor: 'rgba(255, 152, 0, 0.1)' 
              }
            }}
          >
            API测试
          </Button>
          
          {isConnected ? (
            <>
              <IconButton
                size="large"
                edge="end"
                aria-label="account"
                aria-controls="account-menu"
                aria-haspopup="true"
                onClick={handleMenuOpen}
                color="inherit"
              >
                <AccountCircle />
              </IconButton>
              <Typography variant="body2" sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                {isTestMode && <span style={{ color: '#ff9800', fontSize: '12px' }}>测试</span>}
                {formatAddress(account!)}
              </Typography>
              <Menu
                id="account-menu"
                anchorEl={anchorEl}
                anchorOrigin={{
                  vertical: 'bottom',
                  horizontal: 'right',
                }}
                keepMounted
                transformOrigin={{
                  vertical: 'top',
                  horizontal: 'right',
                }}
                open={Boolean(anchorEl)}
                onClose={handleMenuClose}
              >
                <MenuItem onClick={handleProfileClick}>我的资料</MenuItem>
                <MenuItem onClick={handleDisconnect}>断开连接</MenuItem>
              </Menu>
            </>
          ) : (
            <>
              <Button
                color="inherit"
                startIcon={<Wallet />}
                onClick={connectWallet}
                variant="outlined"
              >
                连接钱包
              </Button>
              <Button
                color="inherit"
                onClick={connectTestWallet}
                variant="outlined"
                sx={{ 
                  borderColor: '#ff9800', 
                  color: '#ff9800',
                  '&:hover': { 
                    borderColor: '#ff9800', 
                    backgroundColor: 'rgba(255, 152, 0, 0.1)' 
                  }
                }}
              >
                测试模式
              </Button>
            </>
          )}
        </Box>
      </Toolbar>
    </AppBar>
  );
};

export default Header;
