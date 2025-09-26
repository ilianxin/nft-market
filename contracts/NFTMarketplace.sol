// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "@openzeppelin/contracts/token/ERC721/IERC721.sol";
import "@openzeppelin/contracts/security/ReentrancyGuard.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

/**
 * @title NFT市场订单簿合约
 * @dev 实现基于订单簿模型的NFT交易系统
 */
contract NFTMarketplace is ReentrancyGuard, Ownable {
    
    // 订单状态枚举
    enum OrderStatus {
        Active,     // 活跃
        Filled,     // 已成交
        Cancelled,  // 已取消
        Expired     // 已过期
    }
    
    // 订单类型枚举
    enum OrderType {
        LimitSell,   // 限价卖单
        LimitBuy,    // 限价买单
        MarketSell,  // 市价卖单
        MarketBuy    // 市价买单
    }
    
    // 订单结构体
    struct Order {
        uint256 orderId;           // 订单ID
        address maker;             // 订单创建者
        address nftContract;       // NFT合约地址
        uint256 tokenId;           // NFT Token ID
        uint256 price;             // 价格 (wei)
        uint256 amount;            // 数量 (对于ERC721通常为1)
        uint256 timestamp;         // 创建时间戳
        uint256 expiration;        // 过期时间戳
        OrderType orderType;       // 订单类型
        OrderStatus status;        // 订单状态
        bytes signature;           // 签名
    }
    
    // 存储所有订单
    mapping(uint256 => Order) public orders;
    
    // 用户订单映射
    mapping(address => uint256[]) public userOrders;
    
    // NFT合约的订单映射
    mapping(address => mapping(uint256 => uint256[])) public nftOrders;
    
    // 订单计数器
    uint256 public orderCounter;
    
    // 平台手续费率 (基点, 250 = 2.5%)
    uint256 public platformFeeRate = 250;
    
    // 平台手续费收取地址
    address public feeRecipient;
    
    // 事件定义
    event OrderCreated(
        uint256 indexed orderId,
        address indexed maker,
        address indexed nftContract,
        uint256 tokenId,
        uint256 price,
        OrderType orderType
    );
    
    event OrderFilled(
        uint256 indexed orderId,
        address indexed buyer,
        address indexed seller,
        uint256 price
    );
    
    event OrderCancelled(uint256 indexed orderId, address indexed maker);
    
    event OrderExpired(uint256 indexed orderId);
    
    constructor(address _feeRecipient) {
        feeRecipient = _feeRecipient;
    }
    
    /**
     * @dev 创建限价卖单
     * @param _nftContract NFT合约地址
     * @param _tokenId NFT Token ID
     * @param _price 价格
     * @param _expiration 过期时间戳
     */
    function createLimitSellOrder(
        address _nftContract,
        uint256 _tokenId,
        uint256 _price,
        uint256 _expiration
    ) external nonReentrant {
        require(_price > 0, "Price must be greater than 0");
        require(_expiration > block.timestamp, "Expiration must be in the future");
        require(IERC721(_nftContract).ownerOf(_tokenId) == msg.sender, "Not owner of NFT");
        require(IERC721(_nftContract).isApprovedForAll(msg.sender, address(this)), "Contract not approved");
        
        uint256 orderId = ++orderCounter;
        
        orders[orderId] = Order({
            orderId: orderId,
            maker: msg.sender,
            nftContract: _nftContract,
            tokenId: _tokenId,
            price: _price,
            amount: 1,
            timestamp: block.timestamp,
            expiration: _expiration,
            orderType: OrderType.LimitSell,
            status: OrderStatus.Active,
            signature: ""
        });
        
        userOrders[msg.sender].push(orderId);
        nftOrders[_nftContract][_tokenId].push(orderId);
        
        emit OrderCreated(orderId, msg.sender, _nftContract, _tokenId, _price, OrderType.LimitSell);
    }
    
    /**
     * @dev 创建限价买单
     * @param _nftContract NFT合约地址
     * @param _tokenId NFT Token ID
     * @param _expiration 过期时间戳
     */
    function createLimitBuyOrder(
        address _nftContract,
        uint256 _tokenId,
        uint256 _expiration
    ) external payable nonReentrant {
        require(msg.value > 0, "Price must be greater than 0");
        require(_expiration > block.timestamp, "Expiration must be in the future");
        
        uint256 orderId = ++orderCounter;
        
        orders[orderId] = Order({
            orderId: orderId,
            maker: msg.sender,
            nftContract: _nftContract,
            tokenId: _tokenId,
            price: msg.value,
            amount: 1,
            timestamp: block.timestamp,
            expiration: _expiration,
            orderType: OrderType.LimitBuy,
            status: OrderStatus.Active,
            signature: ""
        });
        
        userOrders[msg.sender].push(orderId);
        nftOrders[_nftContract][_tokenId].push(orderId);
        
        emit OrderCreated(orderId, msg.sender, _nftContract, _tokenId, msg.value, OrderType.LimitBuy);
    }
    
    /**
     * @dev 创建市价卖单 (立即与最佳买单匹配)
     * @param _nftContract NFT合约地址
     * @param _tokenId NFT Token ID
     */
    function createMarketSellOrder(
        address _nftContract,
        uint256 _tokenId
    ) external nonReentrant {
        require(IERC721(_nftContract).ownerOf(_tokenId) == msg.sender, "Not owner of NFT");
        require(IERC721(_nftContract).isApprovedForAll(msg.sender, address(this)), "Contract not approved");
        
        // 查找最佳买单
        uint256 bestBuyOrderId = findBestBuyOrder(_nftContract, _tokenId);
        require(bestBuyOrderId != 0, "No matching buy order found");
        
        // 执行交易
        _executeTrade(bestBuyOrderId, msg.sender);
    }
    
    /**
     * @dev 创建市价买单 (立即与最佳卖单匹配)
     * @param _nftContract NFT合约地址
     * @param _tokenId NFT Token ID
     */
    function createMarketBuyOrder(
        address _nftContract,
        uint256 _tokenId
    ) external payable nonReentrant {
        require(msg.value > 0, "Must send ETH");
        
        // 查找最佳卖单
        uint256 bestSellOrderId = findBestSellOrder(_nftContract, _tokenId);
        require(bestSellOrderId != 0, "No matching sell order found");
        
        Order storage order = orders[bestSellOrderId];
        require(msg.value >= order.price, "Insufficient payment");
        
        // 执行交易
        _executeTrade(bestSellOrderId, msg.sender);
        
        // 退还多余的ETH
        if (msg.value > order.price) {
            payable(msg.sender).transfer(msg.value - order.price);
        }
    }
    
    /**
     * @dev 取消订单
     * @param _orderId 订单ID
     */
    function cancelOrder(uint256 _orderId) external nonReentrant {
        Order storage order = orders[_orderId];
        require(order.maker == msg.sender, "Not order maker");
        require(order.status == OrderStatus.Active, "Order not active");
        
        order.status = OrderStatus.Cancelled;
        
        // 如果是买单，退还ETH
        if (order.orderType == OrderType.LimitBuy) {
            payable(order.maker).transfer(order.price);
        }
        
        emit OrderCancelled(_orderId, msg.sender);
    }
    
    /**
     * @dev 编辑订单 (取消原订单并创建新订单)
     * @param _orderId 原订单ID
     * @param _newPrice 新价格
     * @param _newExpiration 新过期时间
     */
    function editOrder(
        uint256 _orderId,
        uint256 _newPrice,
        uint256 _newExpiration
    ) external payable nonReentrant {
        Order storage oldOrder = orders[_orderId];
        require(oldOrder.maker == msg.sender, "Not order maker");
        require(oldOrder.status == OrderStatus.Active, "Order not active");
        require(_newExpiration > block.timestamp, "Expiration must be in the future");
        
        // 取消原订单
        oldOrder.status = OrderStatus.Cancelled;
        
        // 创建新订单
        if (oldOrder.orderType == OrderType.LimitSell) {
            require(_newPrice > 0, "Price must be greater than 0");
            createLimitSellOrder(oldOrder.nftContract, oldOrder.tokenId, _newPrice, _newExpiration);
        } else if (oldOrder.orderType == OrderType.LimitBuy) {
            require(msg.value > 0, "Must send ETH for buy order");
            // 退还原订单的ETH
            payable(msg.sender).transfer(oldOrder.price);
            createLimitBuyOrder(oldOrder.nftContract, oldOrder.tokenId, _newExpiration);
        }
        
        emit OrderCancelled(_orderId, msg.sender);
    }
    
    /**
     * @dev 查找最佳买单
     * @param _nftContract NFT合约地址
     * @param _tokenId NFT Token ID
     * @return 最佳买单ID
     */
    function findBestBuyOrder(address _nftContract, uint256 _tokenId) public view returns (uint256) {
        uint256[] memory orderIds = nftOrders[_nftContract][_tokenId];
        uint256 bestOrderId = 0;
        uint256 bestPrice = 0;
        
        for (uint256 i = 0; i < orderIds.length; i++) {
            Order memory order = orders[orderIds[i]];
            if (order.status == OrderStatus.Active && 
                order.orderType == OrderType.LimitBuy &&
                order.expiration > block.timestamp &&
                order.price > bestPrice) {
                bestPrice = order.price;
                bestOrderId = orderIds[i];
            }
        }
        
        return bestOrderId;
    }
    
    /**
     * @dev 查找最佳卖单
     * @param _nftContract NFT合约地址
     * @param _tokenId NFT Token ID
     * @return 最佳卖单ID
     */
    function findBestSellOrder(address _nftContract, uint256 _tokenId) public view returns (uint256) {
        uint256[] memory orderIds = nftOrders[_nftContract][_tokenId];
        uint256 bestOrderId = 0;
        uint256 bestPrice = type(uint256).max;
        
        for (uint256 i = 0; i < orderIds.length; i++) {
            Order memory order = orders[orderIds[i]];
            if (order.status == OrderStatus.Active && 
                order.orderType == OrderType.LimitSell &&
                order.expiration > block.timestamp &&
                order.price < bestPrice) {
                bestPrice = order.price;
                bestOrderId = orderIds[i];
            }
        }
        
        return bestOrderId == type(uint256).max ? 0 : bestOrderId;
    }
    
    /**
     * @dev 执行交易
     * @param _orderId 订单ID
     * @param _taker 接受者地址
     */
    function _executeTrade(uint256 _orderId, address _taker) internal {
        Order storage order = orders[_orderId];
        require(order.status == OrderStatus.Active, "Order not active");
        require(order.expiration > block.timestamp, "Order expired");
        
        order.status = OrderStatus.Filled;
        
        // 计算平台手续费
        uint256 platformFee = (order.price * platformFeeRate) / 10000;
        uint256 sellerAmount = order.price - platformFee;
        
        if (order.orderType == OrderType.LimitSell) {
            // 限价卖单：转移NFT给买家，转移ETH给卖家
            IERC721(order.nftContract).safeTransferFrom(order.maker, _taker, order.tokenId);
            payable(order.maker).transfer(sellerAmount);
            payable(feeRecipient).transfer(platformFee);
            
            emit OrderFilled(_orderId, _taker, order.maker, order.price);
        } else if (order.orderType == OrderType.LimitBuy) {
            // 限价买单：转移NFT给买家，转移ETH给卖家
            IERC721(order.nftContract).safeTransferFrom(_taker, order.maker, order.tokenId);
            payable(_taker).transfer(sellerAmount);
            payable(feeRecipient).transfer(platformFee);
            
            emit OrderFilled(_orderId, order.maker, _taker, order.price);
        }
    }
    
    /**
     * @dev 获取用户订单
     * @param _user 用户地址
     * @return 订单ID数组
     */
    function getUserOrders(address _user) external view returns (uint256[] memory) {
        return userOrders[_user];
    }
    
    /**
     * @dev 获取NFT的所有订单
     * @param _nftContract NFT合约地址
     * @param _tokenId NFT Token ID
     * @return 订单ID数组
     */
    function getNFTOrders(address _nftContract, uint256 _tokenId) external view returns (uint256[] memory) {
        return nftOrders[_nftContract][_tokenId];
    }
    
    /**
     * @dev 批量获取订单信息
     * @param _orderIds 订单ID数组
     * @return 订单数组
     */
    function getOrdersBatch(uint256[] memory _orderIds) external view returns (Order[] memory) {
        Order[] memory result = new Order[](_orderIds.length);
        for (uint256 i = 0; i < _orderIds.length; i++) {
            result[i] = orders[_orderIds[i]];
        }
        return result;
    }
    
    /**
     * @dev 标记过期订单
     * @param _orderId 订单ID
     */
    function markOrderExpired(uint256 _orderId) external {
        Order storage order = orders[_orderId];
        require(order.expiration <= block.timestamp, "Order not expired");
        require(order.status == OrderStatus.Active, "Order not active");
        
        order.status = OrderStatus.Expired;
        
        // 如果是买单，退还ETH
        if (order.orderType == OrderType.LimitBuy) {
            payable(order.maker).transfer(order.price);
        }
        
        emit OrderExpired(_orderId);
    }
    
    /**
     * @dev 设置平台手续费率 (仅管理员)
     * @param _feeRate 手续费率 (基点)
     */
    function setPlatformFeeRate(uint256 _feeRate) external onlyOwner {
        require(_feeRate <= 1000, "Fee rate too high"); // 最大10%
        platformFeeRate = _feeRate;
    }
    
    /**
     * @dev 设置手续费收取地址 (仅管理员)
     * @param _feeRecipient 新的手续费收取地址
     */
    function setFeeRecipient(address _feeRecipient) external onlyOwner {
        require(_feeRecipient != address(0), "Invalid address");
        feeRecipient = _feeRecipient;
    }
}
