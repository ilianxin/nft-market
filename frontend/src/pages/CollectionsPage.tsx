import React, { useState } from 'react';
import {
  Box,
  Typography,
  Grid,
  Card,
  CardContent,
  CardMedia,
  Button,
  Chip,
  CircularProgress,
  Pagination,
  TextField,
  InputAdornment,
} from '@mui/material';
import { Search as SearchIcon } from '@mui/icons-material';
import { Link } from 'react-router-dom';
import { useQuery } from 'react-query';
import { apiService } from '../services/api';

const CollectionsPage: React.FC = () => {
  const [page, setPage] = useState(1);
  const [searchKeyword, setSearchKeyword] = useState('');
  const pageSize = 12;

  const handlePageChange = (event: React.ChangeEvent<unknown>, value: number) => {
    setPage(value);
  };

  const handleSearchChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setSearchKeyword(event.target.value);
    setPage(1); // 重置到第一页
  };

  // 获取集合列表
  const { data: collectionsData, isLoading: collectionsLoading } = useQuery(
    ['collections', page, searchKeyword],
    () => apiService.getCollections(page, pageSize),
    {
      keepPreviousData: true,
    }
  );

  const formatAddress = (address: string) => {
    return `${address.slice(0, 6)}...${address.slice(-4)}`;
  };

  const formatPrice = (price: number | null) => {
    if (!price) return '暂无价格';
    return price.toFixed(4) + ' ETH';
  };

  return (
    <Box>
      <Typography variant="h4" component="h1" gutterBottom>
        NFT集合
      </Typography>
      <Typography variant="body1" color="text.secondary" sx={{ mb: 3 }}>
        浏览所有NFT集合，发现您喜欢的项目
      </Typography>

      {/* 搜索框 */}
      <Box sx={{ mb: 3 }}>
        <TextField
          fullWidth
          placeholder="搜索集合名称或描述..."
          value={searchKeyword}
          onChange={handleSearchChange}
          InputProps={{
            startAdornment: (
              <InputAdornment position="start">
                <SearchIcon />
              </InputAdornment>
            ),
          }}
          sx={{ maxWidth: 600 }}
        />
      </Box>

      {collectionsLoading ? (
        <Box className="loading-spinner">
          <CircularProgress />
        </Box>
      ) : (
        <>
          <Grid container spacing={3}>
            {collectionsData?.data?.collections?.map((collection: any) => (
              <Grid item xs={12} sm={6} md={4} lg={3} key={collection.id}>
                <Card className="collection-card">
                  <CardMedia
                    component="div"
                    sx={{
                      height: 200,
                      bgcolor: 'grey.200',
                      display: 'flex',
                      alignItems: 'center',
                      justifyContent: 'center',
                    }}
                  >
                    {collection.image_uri ? (
                      <img
                        src={collection.image_uri}
                        alt={collection.name}
                        style={{ width: '100%', height: '100%', objectFit: 'cover' }}
                      />
                    ) : (
                      <Typography variant="h6" color="text.secondary">
                        {collection.name}
                      </Typography>
                    )}
                  </CardMedia>
                  <CardContent>
                    <Typography variant="h6" gutterBottom>
                      {collection.name}
                    </Typography>
                    <Typography variant="body2" color="text.secondary" gutterBottom>
                      {collection.symbol}
                    </Typography>
                    <Typography variant="body2" color="text.secondary" gutterBottom>
                      创建者: {formatAddress(collection.creator)}
                    </Typography>
                    <Typography variant="body2" color="text.secondary" gutterBottom>
                      地址: {formatAddress(collection.address)}
                    </Typography>
                    
                    <Box sx={{ display: 'flex', gap: 1, mb: 2 }}>
                      <Chip 
                        label={`${collection.item_amount} 个物品`} 
                        size="small" 
                        color="primary" 
                      />
                      <Chip 
                        label={`${collection.owner_amount} 个拥有者`} 
                        size="small" 
                        color="secondary" 
                      />
                    </Box>

                    <Box sx={{ mb: 2 }}>
                      <Typography variant="body2" color="text.secondary">
                        地板价: {formatPrice(collection.floor_price)}
                      </Typography>
                      <Typography variant="body2" color="text.secondary">
                        最高出价: {formatPrice(collection.sale_price)}
                      </Typography>
                      {collection.volume_total && (
                        <Typography variant="body2" color="text.secondary">
                          总交易量: {formatPrice(collection.volume_total)}
                        </Typography>
                      )}
                    </Box>

                    <Button
                      fullWidth
                      variant="contained"
                      component={Link}
                      to={`/collection/${collection.address}`}
                    >
                      查看集合
                    </Button>
                  </CardContent>
                </Card>
              </Grid>
            ))}
          </Grid>

          {collectionsData?.data?.total_pages > 1 && (
            <Box sx={{ display: 'flex', justifyContent: 'center', mt: 4 }}>
              <Pagination
                count={collectionsData.data.total_pages}
                page={page}
                onChange={handlePageChange}
                color="primary"
              />
            </Box>
          )}

          {(!collectionsData?.data?.collections || collectionsData.data.collections.length === 0) && (
            <Box sx={{ textAlign: 'center', py: 8 }}>
              <Typography variant="h6" color="text.secondary">
                暂无集合数据
              </Typography>
            </Box>
          )}
        </>
      )}
    </Box>
  );
};

export default CollectionsPage;
