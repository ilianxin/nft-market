# 🎨 NFT市场应用

基于React前端和Go后端的完整NFT交易平台，支持基于智能合约的订单簿交易系统和实时购买功能。

## ✨ 功能特性

### 🛒 核心交易功能
- **📋 订单管理**: 创建、查看、取消订单
- **💰 实时购买**: 一键购买NFT，支持钱包集成
- **🔄 订单同步**: 链上链下数据实时同步
- **📊 订单簿**: 完整的买卖订单展示
- **👤 用户管理**: 个人订单和NFT管理

### ⛓️ 区块链集成功能
- **🏪 NFT市场合约**: 完整的智能合约系统
- **💸 限价交易**: 支持限价买入/卖出订单
- **⚡ 即时成交**: 实时链上交易执行
- **🔐 安全交易**: 智能合约保障资金安全
- **📈 事件监听**: 自动监听和处理链上事件
- **🔍 链上查询**: 支持查询所有历史订单

### 🎯 用户体验功能
- **🔗 钱包连接**: 支持MetaMask等Web3钱包
- **🧪 测试模式**: 无需钱包的测试功能
- **📱 响应式设计**: 支持桌面和移动端
- **🎨 现代UI**: Material-UI设计语言
- **📊 实时数据**: 订单状态实时更新

### 🛠️ 开发者功能
- **🧪 API测试页面**: 完整的API测试工具
- **📝 结构化日志**: 前后端详细日志系统
- **🔄 数据同步**: 智能的链上链下数据同步
- **🗄️ 数据管理**: 完整的数据库管理工具

### 💻 技术架构
- **前端**: React + TypeScript + Material-UI + Ethers.js
- **后端**: Go + Gin + GORM + MySQL + 结构化日志
- **区块链**: Ethereum + Solidity + Hardhat
- **数据库**: MySQL + 事务管理
- **部署**: Docker + 自动化脚本

## 项目结构

```
NFTMarket/
├── 🎨 frontend/                    # React前端应用
│   ├── src/
│   │   ├── components/         # UI组件
│   │   │   └── Layout/        # 布局组件 (Header, Footer)
│   │   ├── pages/              # 页面组件
│   │   │   ├── HomePage.tsx   # 首页
│   │   │   ├── MarketplacePage.tsx # 市场页面
│   │   │   ├── NFTDetailPage.tsx # NFT详情页
│   │   │   ├── CreateOrderPage.tsx # 创建订单页
│   │   │   ├── UserProfilePage.tsx # 用户资料页
│   │   │   ├── CollectionsPage.tsx # 集合页面
│   │   │   ├── ActivitiesPage.tsx # 活动页面
│   │   │   └── TestPage.tsx   # 🧪 API测试页面
│   │   ├── contexts/           # React Context
│   │   │   └── Web3Context.tsx # Web3钱包连接
│   │   ├── services/           # API服务
│   │   │   └── api.ts         # 完整API服务
│   │   ├── utils/             # 工具函数
│   │   │   └── logger.ts      # 前端日志系统
│   │   └── App.tsx            # 主应用组件
│   ├── public/                # 静态资源
│   └── package.json           # 前端依赖
├── ⚙️ backend/                    # Go后端服务
│   ├── internal/
│   │   ├── api/               # API路由和处理器
│   │   │   ├── routes.go      # 路由定义
│   │   │   └── handlers/      # API处理器
│   │   │       ├── order.go   # 订单处理器 (含购买功能)
│   │   │       ├── nft.go     # NFT处理器
│   │   │       ├── collection.go # 集合处理器
│   │   │       ├── item.go    # 物品处理器
│   │   │       ├── activity.go # 活动处理器
│   │   │       └── blockchain.go # 区块链管理处理器
│   │   ├── models/            # 数据模型
│   │   │   └── models.go      # 完整数据模型定义
│   │   ├── services/          # 业务逻辑服务
│   │   │   ├── order.go       # 订单服务 (含购买逻辑)
│   │   │   ├── nft.go         # NFT服务
│   │   │   ├── blockchain.go  # 基础区块链服务
│   │   │   └── blockchain_enhanced.go # 增强区块链服务
│   │   ├── blockchain/        # 区块链交互
│   │   │   ├── contract.go    # 智能合约交互
│   │   │   └── event_listener.go # 事件监听器
│   │   ├── database/          # 数据库配置
│   │   ├── logger/            # 📝 结构化日志系统
│   │   └── config/            # 配置管理
│   ├── scripts/               # 工具脚本
│   │   ├── reset_database.go  # 数据库重置
│   │   └── setup_test_data.go # 测试数据生成
│   ├── logs/                  # 日志文件目录
│   ├── main.go               # 后端入口文件
│   └── go.mod                # Go依赖
├── 📜 contracts/                 # 智能合约
│   ├── contracts/
│   │   └── NFTMarketplace.sol # 主要市场合约
│   ├── scripts/              # 部署脚本
│   │   └── deploy.js         # 合约部署脚本
│   ├── hardhat.config.js     # Hardhat配置
│   └── package.json          # 合约依赖
├── 🐳 docker-compose.yml         # Docker服务配置
├── 📋 setup.sh / setup.bat      # 安装脚本
├── 🚀 start-dev.sh / start-dev.bat # 开发环境启动脚本
├── ⚡ start-local-ethereum.bat   # 本地以太坊启动脚本
├── 🗄️ start-mysql.sh / start-mysql.bat # MySQL启动脚本
├── 📚 BLOCKCHAIN_INTEGRATION_GUIDE.md # 区块链集成指南
├── 📖 FRONTEND_TESTING_GUIDE.md # 前端测试指南
├── 🔐 PRIVATE_KEY_SECURITY_GUIDE.md # 私钥安全指南
├── 📋 LOCAL_ETHEREUM_GUIDE.md   # 本地以太坊指南
└── 📄 README.md                 # 项目文档
```

## 🚀 快速开始

### 📋 环境要求
- **Node.js** 16+ (推荐使用 LTS 版本)
- **Go** 1.19+ 
- **MySQL** 8.0+
- **Git** (用于克隆项目)
- **MetaMask** 钱包 (用于Web3交互)

### ⚡ 一键启动（推荐）

#### 🐧 Linux/macOS:
```bash
# 克隆项目
git clone <repository-url>
cd NFTMarket

# 一键安装依赖
./setup.sh

# 一键启动开发环境
./start-dev.sh
```

#### 🪟 Windows:
```cmd
# 克隆项目
git clone <repository-url>
cd NFTMarket

# 一键安装依赖
setup.bat

# 一键启动开发环境
start-dev.bat
```

### 🎯 快速测试（无需钱包）

1. **启动服务后访问测试页面**：
   ```
   http://localhost:3000/test
   ```

2. **使用测试模式**：
   - 点击页面右上角的 **"测试模式"** 按钮
   - 无需连接MetaMask即可体验所有功能

3. **API测试**：
   - 在测试页面点击 **"设置测试用户"**
   - 逐个测试各项API功能
   - 查看详细的测试结果和日志

### 🛠️ 手动安装和运行

#### 1. **🗄️ 启动MySQL数据库**：
```bash
# Linux/macOS
./start-mysql.sh

# Windows
start-mysql.bat

# 或使用Docker Compose
docker-compose up -d mysql
```

#### 2. **⚡ 启动本地以太坊网络**：
```bash
# Windows (一键启动Hardhat网络和部署合约)
start-local-ethereum.bat

# 或手动启动
cd contracts
npm install
npx hardhat node
# 新终端窗口部署合约
npx hardhat run scripts/deploy.js --network localhost
```

#### 3. **🔧 设置测试数据**：
```bash
cd backend/scripts
go run setup_test_data.go
```

#### 4. **⚙️ 启动后端服务**：
```bash
cd backend
go mod tidy

# 设置环境变量（或创建.env文件）
export DATABASE_URL="root:123456@tcp(localhost:3306)/nft_market?charset=utf8mb4&parseTime=True&loc=Local"
export CONTRACT_ADDRESS="0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512"
export PRIVATE_KEY="0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"

go run main.go
```

#### 5. **🎨 启动前端应用**：
```bash
cd frontend
npm install
npm start
```

#### 6. **🧪 开始测试**：
- 访问 `http://localhost:3000`
- 点击右上角 **"测试模式"** 或访问 `/test` 页面
- 开始体验完整的NFT交易功能！

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

## 📡 API文档

### 🛒 订单相关接口
- `POST /api/v1/orders` - 创建订单
- `GET /api/v1/orders` - 获取订单列表
- `GET /api/v1/orders/:id` - 获取单个订单
- `PUT /api/v1/orders/:id/cancel` - 取消订单
- `POST /api/v1/orders/:id/purchase` - **💰 购买订单** (新增)
- `GET /api/v1/orders/user/:address` - 获取用户订单
- `GET /api/v1/orders/nft/:collection_address/:token_id` - 获取NFT订单
- `POST /api/v1/orders/sync/:orderid` - 从链上同步订单

### 🎨 物品相关接口
- `GET /api/v1/items` - 获取物品列表
- `GET /api/v1/items/id/:id` - 获取单个物品
- `GET /api/v1/items/token/:collection_address/:token_id` - 根据Token ID获取物品
- `POST /api/v1/items` - 创建物品
- `PUT /api/v1/items/id/:id` - 更新物品
- `DELETE /api/v1/items/id/:id` - 删除物品
- `GET /api/v1/items/collection/:collection_address` - 获取集合下的物品
- `GET /api/v1/items/owner/:owner` - 获取用户拥有的物品

### 📚 集合相关接口
- `GET /api/v1/collections` - 获取集合列表
- `GET /api/v1/collections/:id` - 获取单个集合
- `GET /api/v1/collections/address/:address` - 根据地址获取集合
- `POST /api/v1/collections` - 创建集合
- `PUT /api/v1/collections/:id` - 更新集合
- `DELETE /api/v1/collections/:id` - 删除集合

### 📊 活动相关接口
- `GET /api/v1/activities` - 获取活动列表
- `GET /api/v1/activities/:id` - 获取单个活动
- `GET /api/v1/activities/user/:address` - 获取用户活动
- `GET /api/v1/activities/stats` - 获取活动统计

### ⛓️ 区块链管理接口
- `GET /api/v1/blockchain/status` - 获取区块链服务状态
- `GET /api/v1/blockchain/counter` - 获取订单计数器
- `GET /api/v1/blockchain/order/:orderid` - 获取链上订单信息
- `POST /api/v1/blockchain/sync/order/:orderid` - 同步单个订单
- `POST /api/v1/blockchain/sync/all` - 同步所有订单
- `POST /api/v1/blockchain/execute/:orderid` - 执行订单

## 📖 使用指南

### 🔗 连接钱包
1. **安装MetaMask**: 下载并安装MetaMask浏览器扩展
2. **连接钱包**: 点击页面右上角 **"连接钱包"** 按钮
3. **网络配置**: 确保连接到正确的网络（本地开发为 `localhost:8545`）
4. **测试模式**: 或点击 **"测试模式"** 按钮无需钱包即可体验

### 🛒 购买NFT（核心功能）
1. **浏览市场**: 访问 **"市场"** 页面查看所有可购买的NFT
2. **选择NFT**: 点击心仪的NFT卡片
3. **立即购买**: 点击 **"立即购买"** 按钮
4. **确认交易**: 确认购买信息并等待区块链处理
5. **交易完成**: 系统会自动更新NFT拥有权和订单状态

### 📋 创建订单
1. **导航页面**: 前往 **"创建订单"** 页面
2. **选择类型**: 选择订单类型（上架出售/出价购买）
3. **填写信息**: 输入NFT合约地址、Token ID、价格等
4. **设置期限**: 选择订单过期时间
5. **提交订单**: 确认信息并创建订单

### 🔍 浏览和搜索
1. **市场页面**: 查看所有活跃的买卖订单
2. **NFT详情**: 点击任意NFT查看详细信息和完整订单簿
3. **个人资料**: 在 **"我的资料"** 页面管理个人订单和NFT
4. **集合浏览**: 访问 **"集合"** 页面探索不同的NFT集合
5. **活动记录**: 在 **"活动"** 页面查看所有交易历史

### 🛠️ 管理功能
1. **订单管理**: 查看、取消个人活跃订单
2. **NFT管理**: 管理拥有的NFT和上架状态
3. **交易历史**: 查看完整的买卖记录
4. **API测试**: 使用 `/test` 页面进行功能测试和调试

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

## 🔧 故障排除

### ❓ 常见问题

#### 🚀 启动问题
1. **端口被占用**: 确保 3000、8080、8545、3306 端口未被占用
2. **MySQL连接失败**: 检查MySQL服务是否启动，密码是否为 `123456`
3. **Go编译错误**: 确保Go版本 1.19+ 并运行 `go mod tidy`
4. **Node.js版本问题**: 使用Node.js 16+ LTS版本

#### ⛓️ 区块链问题  
1. **合约部署失败**: 确保Hardhat网络正在运行 (`npx hardhat node`)
2. **私钥错误**: 使用测试私钥 `0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80`
3. **网络连接失败**: 检查 `localhost:8545` 是否可访问
4. **交易失败**: 确保测试账户有足够的ETH余额

#### 🛒 购买问题
1. **404错误**: 确保数据库中有测试数据，运行 `backend/scripts/setup_test_data.go`
2. **购买失败**: 检查订单状态和用户地址是否正确
3. **钱包问题**: 使用 **"测试模式"** 绕过钱包连接

#### 🗄️ 数据库问题
1. **表不存在**: 运行 `backend/scripts/reset_database.go` 重置数据库
2. **数据为空**: 运行 `backend/scripts/setup_test_data.go` 创建测试数据
3. **连接超时**: 检查MySQL服务状态和端口 3306

### 📊 日志查看

#### 📱 前端日志
- **浏览器控制台**: F12 → Console 查看详细日志
- **网络请求**: F12 → Network 查看API调用
- **测试页面**: `/test` 页面查看API测试结果

#### ⚙️ 后端日志  
- **终端输出**: 实时查看服务器日志
- **日志文件**: `backend/logs/app.log` 查看详细日志
- **数据库日志**: 查看GORM SQL执行日志

#### ⛓️ 区块链日志
- **Hardhat终端**: 查看合约交互和交易日志
- **MetaMask**: 查看钱包交易历史

### 🧪 调试建议

1. **使用测试模式**: 点击 **"测试模式"** 按钮避免钱包问题
2. **访问测试页面**: `http://localhost:3000/test` 进行API测试
3. **检查数据库状态**: 在测试页面运行 **"检查数据库状态"**
4. **查看详细日志**: 检查前端控制台和后端日志文件
5. **重置环境**: 必要时重新运行安装和设置脚本

## 🤝 贡献指南

我们欢迎所有形式的贡献！

### 📋 如何贡献
1. **Fork** 项目仓库到您的GitHub账户
2. **Clone** 到本地：`git clone <your-fork-url>`
3. **创建分支**：`git checkout -b feature/your-feature-name`
4. **开发功能**：实现您的功能或修复
5. **测试验证**：确保所有功能正常工作
6. **提交更改**：`git commit -m "feat: add your feature"`
7. **推送分支**：`git push origin feature/your-feature-name`
8. **创建PR**：在GitHub上创建Pull Request

### 🎯 贡献方向
- 🐛 **Bug修复**: 发现并修复问题
- ✨ **新功能**: 添加有用的新特性
- 📚 **文档改进**: 完善文档和指南
- 🎨 **UI/UX优化**: 改善用户体验
- ⚡ **性能优化**: 提升应用性能
- 🧪 **测试覆盖**: 增加测试用例

### 📝 代码规范
- **前端**: 遵循TypeScript和React最佳实践
- **后端**: 遵循Go语言规范和RESTful API设计
- **智能合约**: 遵循Solidity安全最佳实践
- **提交信息**: 使用约定式提交格式

## 📄 许可证

本项目采用 **MIT License** 开源协议 - 详见 [LICENSE](LICENSE) 文件

## 🆘 技术支持


### 联系方式
- **邮箱**: ilianxin@126.com
- **项目地址**: [GitHub Repository](https://github.com/your-repo)

---

## 🎉 特别感谢

感谢所有为这个项目做出贡献的开发者！

**如果这个项目对您有帮助，请给我们一个 ⭐ Star！**
