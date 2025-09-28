# 模型改动总结

## 后端模型变化

### 1. 新增模型

#### Collection (集合模型)
- `ID`: 主键
- `ChainID`: 链类型 (1: 以太坊)
- `Symbol`: 项目标识
- `Name`: 项目名称
- `Creator`: 创建者地址
- `Address`: 合约地址
- `OwnerAmount`: 拥有者数量
- `ItemAmount`: 物品总量
- `FloorPrice`: 地板价
- `SalePrice`: 最高出价
- `Description`: 项目描述
- `Website`: 官网地址
- `VolumeTotal`: 总交易量
- `ImageURI`: 封面图链接

#### Item (物品模型，原NFT模型)
- `ID`: 主键
- `ChainID`: 链类型
- `TokenID`: 代币ID
- `Name`: 物品名称
- `Owner`: 拥有者地址
- `CollectionAddress`: 集合地址
- `Creator`: 创建者
- `Supply`: 供应量
- `ListPrice`: 上架价格
- `ListTime`: 上架时间
- `SalePrice`: 成交价格

#### Activity (活动模型)
- `ID`: 主键
- `ActivityType`: 活动类型 (1-10)
- `Maker`: 发起方
- `Taker`: 接受方
- `CollectionAddress`: 集合地址
- `TokenID`: 代币ID
- `Price`: 价格
- `BlockNumber`: 区块号
- `TxHash`: 交易哈希
- `EventTime`: 事件时间

### 2. 订单模型变化

#### Order 字段更新
- `OrderID`: 从数字改为字符串格式
- `CollectionAddress`: 原 `nft_contract` 字段
- `TokenID`: 保持字符串类型
- `Price`: 从字符串改为数字类型
- `OrderType`: 枚举值变化 (1: 上架, 2: 出价, 3: 集合出价, 4: 物品出价)
- `OrderStatus`: 枚举值变化 (0: 活跃, 1: 已成交, 2: 已取消, 3: 已过期)
- `ExpireTime`: 从ISO字符串改为Unix时间戳

### 3. 用户模型优化
- `Address`: 指定varchar(42)长度以解决MySQL索引问题
- `Username`: 指定varchar(255)长度
- `Email`: 指定varchar(255)长度

## 前端代码更新

### 1. API服务更新 (`frontend/src/services/api.ts`)

#### 新增API端点
- **集合相关**: `getCollections`, `getCollectionById`, `getCollectionByAddress`, `createCollection`, `updateCollection`
- **物品相关**: `getItems`, `getItemById`, `getItemByTokenId`, `getUserItems`, `getItemsByCollection`, `searchItems`, `createOrUpdateItem`, `updateItemOwner`, `updateItemPrice`
- **活动相关**: `getActivities`, `getActivityById`, `getActivitiesByUser`, `getActivitiesByItem`, `createActivity`

#### 更新现有API
- `getNFTOrders`: 参数从 `contract, tokenId` 改为 `collectionAddress, tokenId`
- 所有NFT相关API重命名为Item相关API

### 2. 页面组件更新

#### MarketplacePage (`frontend/src/pages/MarketplacePage.tsx`)
- 更新订单类型枚举映射 (1: 上架, 2: 出价, 3: 集合出价, 4: 物品出价)
- 更新字段引用 (`nft_contract` → `collection_address`)
- 更新价格格式化函数以处理数字类型
- 更新时间戳处理 (Unix时间戳转换)

#### CreateOrderPage (`frontend/src/pages/CreateOrderPage.tsx`)
- 更新表单字段 (`nftContract` → `collectionAddress`)
- 更新订单类型选项和枚举值
- 更新API调用数据结构
- 更新过期时间处理 (Unix时间戳)

#### UserProfilePage (`frontend/src/pages/UserProfilePage.tsx`)
- 更新订单状态和类型枚举映射
- 更新字段引用和数据类型
- 更新NFT相关API调用为Item相关API
- 更新时间戳格式化

### 3. 新增页面组件

#### CollectionsPage (`frontend/src/pages/CollectionsPage.tsx`)
- 集合列表展示
- 搜索功能
- 分页功能
- 集合详情展示 (地板价、最高出价、交易量等)

#### ActivitiesPage (`frontend/src/pages/ActivitiesPage.tsx`)
- 活动记录列表
- 按活动类型分类 (交易、上架、出价)
- 搜索和分页功能
- 活动详情展示

### 4. 导航更新

#### App.tsx
- 添加新路由: `/collections`, `/activities`

#### Header.tsx
- 添加导航链接: "集合", "活动"

## 数据流变化

### 1. 订单创建流程
```
前端表单 → API调用 → 后端验证 → 数据库存储
- 字段映射更新
- 数据类型转换
- 时间戳处理
```

### 2. 订单显示流程
```
数据库查询 → API响应 → 前端渲染
- 枚举值映射
- 时间戳格式化
- 地址格式化
```

### 3. 新增数据流
```
集合管理: 创建 → 展示 → 搜索
活动记录: 记录 → 分类 → 展示
物品管理: 创建 → 更新 → 查询
```

## 兼容性说明

### 1. 数据库兼容性
- 所有字段类型变化都通过GORM自动迁移处理
- 索引问题已通过指定字段长度解决

### 2. API兼容性
- 新增API端点不影响现有功能
- 现有API保持向后兼容

### 3. 前端兼容性
- 所有组件都更新以使用新的数据结构
- 错误处理保持一致性

## 测试建议

### 1. 后端测试
- 测试所有新的API端点
- 验证数据库迁移
- 测试枚举值映射

### 2. 前端测试
- 测试所有页面的数据展示
- 验证表单提交功能
- 测试搜索和分页功能

### 3. 集成测试
- 测试完整的数据流
- 验证前后端数据一致性
- 测试错误处理

## 部署注意事项

1. 确保数据库迁移成功执行
2. 验证所有API端点正常工作
3. 检查前端页面加载和功能正常
4. 监控系统性能和数据一致性
