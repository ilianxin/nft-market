import React, { useState } from 'react';
import {
  Box,
  Typography,
  Grid,
  Card,
  CardContent,
  Chip,
  CircularProgress,
  Pagination,
  TextField,
  InputAdornment,
  Tabs,
  Tab,
} from '@mui/material';
import { Search as SearchIcon } from '@mui/icons-material';
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
      id={`activities-tabpanel-${index}`}
      aria-labelledby={`activities-tab-${index}`}
      {...other}
    >
      {value === index && <Box sx={{ pt: 3 }}>{children}</Box>}
    </div>
  );
}

const ActivitiesPage: React.FC = () => {
  const [tabValue, setTabValue] = useState(0);
  const [page, setPage] = useState(1);
  const [searchKeyword, setSearchKeyword] = useState('');
  const pageSize = 20;

  const handleTabChange = (event: React.SyntheticEvent, newValue: number) => {
    setTabValue(newValue);
    setPage(1);
  };

  const handlePageChange = (event: React.ChangeEvent<unknown>, value: number) => {
    setPage(value);
  };

  const handleSearchChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setSearchKeyword(event.target.value);
    setPage(1);
  };

  // 获取活动列表
  const { data: activitiesData, isLoading: activitiesLoading } = useQuery(
    ['activities', page, searchKeyword],
    () => apiService.getActivities(page, pageSize),
    {
      keepPreviousData: true,
    }
  );

  const formatAddress = (address: string) => {
    return `${address.slice(0, 6)}...${address.slice(-4)}`;
  };

  const formatPrice = (price: number) => {
    return price.toFixed(4) + ' ETH';
  };

  const getActivityTypeText = (activityType: number) => {
    switch (activityType) {
      case 1: return '购买';
      case 2: return '铸造';
      case 3: return '上架';
      case 4: return '取消上架';
      case 5: return '取消出价';
      case 6: return '出价';
      case 7: return '出售';
      case 8: return '转移';
      case 9: return '集合出价';
      case 10: return '物品出价';
      default: return `活动${activityType}`;
    }
  };

  const getActivityTypeColor = (activityType: number) => {
    switch (activityType) {
      case 1: return 'success'; // 购买
      case 2: return 'info';    // 铸造
      case 3: return 'warning'; // 上架
      case 4: return 'error';   // 取消上架
      case 5: return 'error';   // 取消出价
      case 6: return 'primary'; // 出价
      case 7: return 'success'; // 出售
      case 8: return 'default'; // 转移
      case 9: return 'primary'; // 集合出价
      case 10: return 'primary'; // 物品出价
      default: return 'default';
    }
  };

  const formatTimestamp = (timestamp: number | null) => {
    if (!timestamp) return '未知时间';
    return new Date(timestamp * 1000).toLocaleString();
  };

  return (
    <Box>
      <Typography variant="h4" component="h1" gutterBottom>
        活动记录
      </Typography>
      <Typography variant="body1" color="text.secondary" sx={{ mb: 3 }}>
        查看所有NFT交易和活动记录
      </Typography>

      {/* 搜索框 */}
      <Box sx={{ mb: 3 }}>
        <TextField
          fullWidth
          placeholder="搜索活动记录..."
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

      <Box sx={{ borderBottom: 1, borderColor: 'divider' }}>
        <Tabs value={tabValue} onChange={handleTabChange}>
          <Tab label="全部活动" />
          <Tab label="交易活动" />
          <Tab label="上架活动" />
          <Tab label="出价活动" />
        </Tabs>
      </Box>

      <TabPanel value={tabValue} index={0}>
        {activitiesLoading ? (
          <Box className="loading-spinner">
            <CircularProgress />
          </Box>
        ) : (
          <>
            <Grid container spacing={2}>
              {activitiesData?.data?.activities?.map((activity: any) => (
                <Grid item xs={12} key={activity.id}>
                  <Card>
                    <CardContent>
                      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
                        <Chip
                          label={getActivityTypeText(activity.activity_type)}
                          color={getActivityTypeColor(activity.activity_type) as any}
                          size="small"
                        />
                        <Typography variant="body2" color="text.secondary">
                          {formatTimestamp(activity.event_time)}
                        </Typography>
                      </Box>

                      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                        <Box>
                          <Typography variant="body1" gutterBottom>
                            {getActivityTypeText(activity.activity_type)} - Token ID #{activity.token_id}
                          </Typography>
                          <Typography variant="body2" color="text.secondary">
                            集合: {formatAddress(activity.collection_address || '')}
                          </Typography>
                          {activity.maker && (
                            <Typography variant="body2" color="text.secondary">
                              发起方: {formatAddress(activity.maker)}
                            </Typography>
                          )}
                          {activity.taker && (
                            <Typography variant="body2" color="text.secondary">
                              接受方: {formatAddress(activity.taker)}
                            </Typography>
                          )}
                        </Box>
                        <Box sx={{ textAlign: 'right' }}>
                          <Typography variant="h6" color="primary">
                            {formatPrice(activity.price)}
                          </Typography>
                          {activity.tx_hash && (
                            <Typography variant="body2" color="text.secondary">
                              交易: {formatAddress(activity.tx_hash)}
                            </Typography>
                          )}
                        </Box>
                      </Box>
                    </CardContent>
                  </Card>
                </Grid>
              ))}
            </Grid>

            {activitiesData?.data?.total_pages > 1 && (
              <Box sx={{ display: 'flex', justifyContent: 'center', mt: 4 }}>
                <Pagination
                  count={activitiesData.data.total_pages}
                  page={page}
                  onChange={handlePageChange}
                  color="primary"
                />
              </Box>
            )}

            {(!activitiesData?.data?.activities || activitiesData.data.activities.length === 0) && (
              <Box sx={{ textAlign: 'center', py: 8 }}>
                <Typography variant="h6" color="text.secondary">
                  暂无活动记录
                </Typography>
              </Box>
            )}
          </>
        )}
      </TabPanel>

      <TabPanel value={tabValue} index={1}>
        {/* 交易活动 - 购买、出售、转移 */}
        {activitiesLoading ? (
          <Box className="loading-spinner">
            <CircularProgress />
          </Box>
        ) : (
          <>
            <Grid container spacing={2}>
              {activitiesData?.data?.activities?.filter((activity: any) => 
                [1, 7, 8].includes(activity.activity_type)
              ).map((activity: any) => (
                <Grid item xs={12} key={activity.id}>
                  <Card>
                    <CardContent>
                      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
                        <Chip
                          label={getActivityTypeText(activity.activity_type)}
                          color={getActivityTypeColor(activity.activity_type) as any}
                          size="small"
                        />
                        <Typography variant="body2" color="text.secondary">
                          {formatTimestamp(activity.event_time)}
                        </Typography>
                      </Box>

                      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                        <Box>
                          <Typography variant="body1" gutterBottom>
                            {getActivityTypeText(activity.activity_type)} - Token ID #{activity.token_id}
                          </Typography>
                          <Typography variant="body2" color="text.secondary">
                            集合: {formatAddress(activity.collection_address || '')}
                          </Typography>
                          {activity.maker && (
                            <Typography variant="body2" color="text.secondary">
                              发起方: {formatAddress(activity.maker)}
                            </Typography>
                          )}
                          {activity.taker && (
                            <Typography variant="body2" color="text.secondary">
                              接受方: {formatAddress(activity.taker)}
                            </Typography>
                          )}
                        </Box>
                        <Box sx={{ textAlign: 'right' }}>
                          <Typography variant="h6" color="primary">
                            {formatPrice(activity.price)}
                          </Typography>
                          {activity.tx_hash && (
                            <Typography variant="body2" color="text.secondary">
                              交易: {formatAddress(activity.tx_hash)}
                            </Typography>
                          )}
                        </Box>
                      </Box>
                    </CardContent>
                  </Card>
                </Grid>
              ))}
            </Grid>
          </>
        )}
      </TabPanel>

      <TabPanel value={tabValue} index={2}>
        {/* 上架活动 */}
        {activitiesLoading ? (
          <Box className="loading-spinner">
            <CircularProgress />
          </Box>
        ) : (
          <>
            <Grid container spacing={2}>
              {activitiesData?.data?.activities?.filter((activity: any) => 
                [3, 4].includes(activity.activity_type)
              ).map((activity: any) => (
                <Grid item xs={12} key={activity.id}>
                  <Card>
                    <CardContent>
                      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
                        <Chip
                          label={getActivityTypeText(activity.activity_type)}
                          color={getActivityTypeColor(activity.activity_type) as any}
                          size="small"
                        />
                        <Typography variant="body2" color="text.secondary">
                          {formatTimestamp(activity.event_time)}
                        </Typography>
                      </Box>

                      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                        <Box>
                          <Typography variant="body1" gutterBottom>
                            {getActivityTypeText(activity.activity_type)} - Token ID #{activity.token_id}
                          </Typography>
                          <Typography variant="body2" color="text.secondary">
                            集合: {formatAddress(activity.collection_address || '')}
                          </Typography>
                          {activity.maker && (
                            <Typography variant="body2" color="text.secondary">
                              操作者: {formatAddress(activity.maker)}
                            </Typography>
                          )}
                        </Box>
                        <Box sx={{ textAlign: 'right' }}>
                          <Typography variant="h6" color="primary">
                            {formatPrice(activity.price)}
                          </Typography>
                          {activity.tx_hash && (
                            <Typography variant="body2" color="text.secondary">
                              交易: {formatAddress(activity.tx_hash)}
                            </Typography>
                          )}
                        </Box>
                      </Box>
                    </CardContent>
                  </Card>
                </Grid>
              ))}
            </Grid>
          </>
        )}
      </TabPanel>

      <TabPanel value={tabValue} index={3}>
        {/* 出价活动 */}
        {activitiesLoading ? (
          <Box className="loading-spinner">
            <CircularProgress />
          </Box>
        ) : (
          <>
            <Grid container spacing={2}>
              {activitiesData?.data?.activities?.filter((activity: any) => 
                [5, 6, 9, 10].includes(activity.activity_type)
              ).map((activity: any) => (
                <Grid item xs={12} key={activity.id}>
                  <Card>
                    <CardContent>
                      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
                        <Chip
                          label={getActivityTypeText(activity.activity_type)}
                          color={getActivityTypeColor(activity.activity_type) as any}
                          size="small"
                        />
                        <Typography variant="body2" color="text.secondary">
                          {formatTimestamp(activity.event_time)}
                        </Typography>
                      </Box>

                      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                        <Box>
                          <Typography variant="body1" gutterBottom>
                            {getActivityTypeText(activity.activity_type)} - Token ID #{activity.token_id}
                          </Typography>
                          <Typography variant="body2" color="text.secondary">
                            集合: {formatAddress(activity.collection_address || '')}
                          </Typography>
                          {activity.maker && (
                            <Typography variant="body2" color="text.secondary">
                              出价者: {formatAddress(activity.maker)}
                            </Typography>
                          )}
                        </Box>
                        <Box sx={{ textAlign: 'right' }}>
                          <Typography variant="h6" color="primary">
                            {formatPrice(activity.price)}
                          </Typography>
                          {activity.tx_hash && (
                            <Typography variant="body2" color="text.secondary">
                              交易: {formatAddress(activity.tx_hash)}
                            </Typography>
                          )}
                        </Box>
                      </Box>
                    </CardContent>
                  </Card>
                </Grid>
              ))}
            </Grid>
          </>
        )}
      </TabPanel>
    </Box>
  );
};

export default ActivitiesPage;
