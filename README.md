# NFT市场应用

基于React前端和Go后端的NFT交易平台，支持基于智能合约的订单簿交易系统。

## 功能特性

### 核心交易功能
- **限价订单**: 支持限价买入/卖出订单
- **市价订单**: 支持市价买入/卖出订单
- **订单管理**: 支持编辑和取消订单
- **订单查询**: 支持查询所有订单（包括过期订单）

### 智能合约功能
- **订单簿DEX**: 基于订单簿模型的去中心化交易
- **限价卖单**: 设定价格等待买家成交
- **限价买单**: 预锁定ETH等待卖家成交
- **市价卖单**: 立即与最优买单匹配
- **市价买单**: 立即与最优卖单匹配
- **订单编辑**: 取消原订单并创建新订单
- **订单取消**: 随时取消未成交订单
- **链上查询**: 支持查询所有历史订单（包括过期订单）

### 技术架构
- **前端**: React + TypeScript + Material-UI + Web3.js
- **后端**: Go + Gin + GORM + MySQL
- **区块链**: Ethereum智能合约 (Solidity)
- **数据库**: MySQL

## 项目结构

```
NFTMarket/
├── frontend/                    # React前端应用
│   ├── src/
│   │   ├── components/         # UI组件
│   │   ├── pages/              # 页面组件
│   │   ├── contexts/           # React Context
│   │   ├── services/           # API服务
│   │   └── App.tsx            # 主应用组件
│   ├── public/                # 静态资源
│   └── package.json           # 前端依赖
├── backend/                    # Go后端服务
│   ├── internal/
│   │   ├── api/               # API路由和处理器
│   │   ├── models/            # 数据模型
│   │   ├── services/          # 业务逻辑服务
│   │   ├── database/          # 数据库配置
│   │   └── config/            # 配置管理
│   ├── main.go               # 后端入口文件
│   └── go.mod                # Go依赖
├── contracts/                 # 智能合约
│   ├── NFTMarketplace.sol    # 主要市场合约
│   ├── scripts/              # 部署脚本
│   ├── hardhat.config.js     # Hardhat配置
│   └── package.json          # 合约依赖
├── setup.sh / setup.bat      # 安装脚本
├── start-dev.sh / start-dev.bat # 开发环境启动脚本
└── README.md                 # 项目文档
```

## 快速开始

### 环境要求
- Node.js 16+
- Go 1.19+
- MySQL 8.0+
- MetaMask钱包

### 一键安装（推荐）

#### Linux/macOS:
```bash
./setup.sh
```

#### Windows:
```cmd
setup.bat
```

### 一键启动开发环境

#### Linux/macOS:
```bash
./start-dev.sh
```

#### Windows:
```cmd
start-dev.bat
```

### 手动安装和运行

1. **启动MySQL数据库**：
```bash
# Linux/macOS
./start-mysql.sh

# Windows
start-mysql.bat

# 或使用Docker Compose
docker-compose up -d mysql
```

2. **安装智能合约依赖并编译**：
```bash
cd contracts
npm install
npx hardhat compile
```

3. **启动本地区块链网络**：
```bash
npx hardhat node
```

4. **部署智能合约**：
```bash
npx hardhat run scripts/deploy.js --network localhost
```

5. **安装后端依赖并启动**：
```bash
cd backend
go mod tidy
# 编辑 .env 文件配置数据库连接
go run main.go
```

6. **安装前端依赖并启动**：
```bash
cd frontend
npm install
# 编辑 .env 文件配置API地址
npm start
```

## 配置说明

### 后端配置 (backend/.env)
```env
PORT=8080
DATABASE_URL=username:password@tcp(localhost:3306)/nft_market?charset=utf8mb4&parseTime=True&loc=Local
ETHEREUM_RPC=http://localhost:8545
CONTRACT_ADDRESS=0x...
PRIVATE_KEY=your_private_key_here
JWT_SECRET=your_jwt_secret_key
```

### 前端配置 (frontend/.env)
```env
REACT_APP_API_URL=http://localhost:8080/api/v1
REACT_APP_CONTRACT_ADDRESS=0x...
REACT_APP_NETWORK_ID=1337
```

## API文档

### 订单相关接口
- `POST /api/v1/orders` - 创建订单
- `GET /api/v1/orders` - 获取订单列表
- `GET /api/v1/orders/:id` - 获取单个订单
- `PUT /api/v1/orders/:id/cancel` - 取消订单
- `GET /api/v1/orders/user/:address` - 获取用户订单
- `GET /api/v1/orders/nft/:contract/:tokenId` - 获取NFT订单
- `POST /api/v1/orders/sync/:orderid` - 从链上同步订单

### NFT相关接口
- `GET /api/v1/nfts` - 获取NFT列表
- `GET /api/v1/nfts/:contract/:tokenId` - 获取单个NFT
- `GET /api/v1/nfts/user/:address` - 获取用户NFT
- `GET /api/v1/nfts/contract/:contract` - 获取合约NFT
- `GET /api/v1/nfts/search` - 搜索NFT
- `POST /api/v1/nfts` - 创建或更新NFT

## 使用指南

### 连接钱包
1. 安装MetaMask浏览器扩展
2. 在应用中点击"连接钱包"按钮
3. 确保连接到正确的网络（本地开发为localhost:8545）

### 创建订单
1. 导航到"创建订单"页面
2. 选择订单类型（限价卖单/买单，市价卖单/买单）
3. 填写NFT合约地址和Token ID
4. 设置价格和过期时间
5. 确认并提交交易

### 浏览市场
1. 访问"市场"页面查看所有活跃订单
2. 可按订单类型筛选（全部订单/限价卖单/限价买单）
3. 点击订单卡片查看详细信息
4. 在NFT详情页面可以查看完整的订单簿

### 管理订单
1. 在"我的资料"页面查看个人订单
2. 可以取消活跃状态的订单
3. 查看订单历史和状态变化

## 开发说明

### 智能合约架构
- `NFTMarketplace.sol`: 主要的市场合约，实现所有核心功能
- 支持ERC721标准的NFT
- 使用OpenZeppelin库确保安全性
- 实现了完整的订单簿模型

### 后端架构
- 使用Gin框架提供REST API
- GORM作为ORM层操作MySQL数据库
- 集成以太坊客户端与智能合约交互
- 提供用户认证和权限管理

### 前端架构
- React + TypeScript提供类型安全
- Material-UI提供现代化UI组件
- React Router处理页面路由
- React Query管理服务器状态
- Web3.js/Ethers.js与区块链交互

## 故障排除

### 常见问题
1. **合约部署失败**: 确保本地区块链网络正在运行
2. **钱包连接失败**: 检查MetaMask是否已安装并连接到正确网络
3. **API请求失败**: 确保后端服务正在运行且MySQL数据库连接正常
4. **交易失败**: 检查账户余额和Gas费用设置
5. **数据库连接失败**: 检查MySQL服务是否启动，使用 `start-mysql.sh` 或 `start-mysql.bat` 启动

### 日志查看
- 前端日志: 浏览器开发者工具控制台
- 后端日志: 终端输出
- 区块链日志: Hardhat网络终端输出

## 贡献指南

1. Fork 项目仓库
2. 创建功能分支
3. 提交更改
4. 推送到分支
5. 创建Pull Request

## 许可证

MIT License - 详见 LICENSE 文件

## 技术支持

如有问题，请创建Issue或联系开发团队。
