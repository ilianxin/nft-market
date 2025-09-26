package services

import (
	"fmt"
	"math/big"
	"nft-market/internal/blockchain"
	"nft-market/internal/logger"
	"nft-market/internal/models"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// EnhancedBlockchainService 增强的区块链服务
type EnhancedBlockchainService struct {
	*BlockchainService // 嵌入原有服务
	contract           *blockchain.NFTMarketplaceContract
	eventListener      *blockchain.EventListener
	db                 *gorm.DB
}

// NewEnhancedBlockchainService 创建增强的区块链服务
func NewEnhancedBlockchainService(rpcURL, contractAddress, privateKey string, db *gorm.DB) (*EnhancedBlockchainService, error) {
	// 创建基础区块链服务
	baseService, err := NewBlockchainService(rpcURL, contractAddress, privateKey)
	if err != nil {
		return nil, fmt.Errorf("创建基础区块链服务失败: %v", err)
	}

	// 创建增强的合约实例
	contract, err := blockchain.NewNFTMarketplaceContract(rpcURL, contractAddress, privateKey)
	if err != nil {
		return nil, fmt.Errorf("创建增强合约实例失败: %v", err)
	}

	// 创建事件监听器
	eventListener := blockchain.NewEventListener(baseService.client, contractAddress, contract, db)

	service := &EnhancedBlockchainService{
		BlockchainService: baseService,
		contract:          contract,
		eventListener:     eventListener,
		db:                db,
	}

	// 启动事件监听
	if err := eventListener.Start(); err != nil {
		logger.Error("启动事件监听器失败", err)
		// 不返回错误，允许服务在没有事件监听的情况下运行
	} else {
		logger.Info("区块链事件监听器启动成功")
	}

	return service, nil
}

// CreateOrderOnChain 在区块链上创建订单（增强版）
func (ebs *EnhancedBlockchainService) CreateOrderOnChain(order *models.Order) (*types.Transaction, error) {
	logger.Info("开始在链上创建订单", logrus.Fields{
		"order_id":           order.OrderID,
		"collection_address": order.CollectionAddress,
		"token_id":           order.TokenID,
		"order_type":         order.OrderType,
		"price":              order.Price,
	})

	var tx *types.Transaction
	var err error

	switch order.OrderType {
	case models.OrderTypeListing:
		// 创建限价卖单
		expiration := int64(time.Now().Add(30 * 24 * time.Hour).Unix()) // 默认30天过期
		if order.ExpireTime != nil {
			expiration = *order.ExpireTime
		}

		tx, err = ebs.contract.CreateLimitSellOrder(
			order.CollectionAddress,
			order.TokenID,
			order.Price,
			expiration,
		)

	case models.OrderTypeOffer:
		// 创建限价买单
		expiration := int64(time.Now().Add(7 * 24 * time.Hour).Unix()) // 默认7天过期
		if order.ExpireTime != nil {
			expiration = *order.ExpireTime
		}

		tx, err = ebs.contract.CreateLimitBuyOrder(
			order.CollectionAddress,
			order.TokenID,
			order.Price,
			expiration,
		)

	default:
		return nil, fmt.Errorf("不支持的订单类型: %v", order.OrderType)
	}

	if err != nil {
		logger.Error("链上创建订单失败", err, logrus.Fields{
			"order_id": order.OrderID,
		})
		return nil, err
	}

	logger.Info("链上订单创建成功", logrus.Fields{
		"order_id": order.OrderID,
		"tx_hash":  tx.Hash().Hex(),
	})

	return tx, nil
}

// CancelOrderOnChain 在区块链上取消订单（增强版）
func (ebs *EnhancedBlockchainService) CancelOrderOnChain(orderID uint64) (*types.Transaction, error) {
	logger.Info("开始在链上取消订单", logrus.Fields{
		"order_id": orderID,
	})

	tx, err := ebs.contract.CancelOrder(orderID)
	if err != nil {
		logger.Error("链上取消订单失败", err, logrus.Fields{
			"order_id": orderID,
		})
		return nil, err
	}

	logger.Info("链上订单取消成功", logrus.Fields{
		"order_id": orderID,
		"tx_hash":  tx.Hash().Hex(),
	})

	return tx, nil
}

// ExecuteOrderOnChain 在区块链上执行订单
func (ebs *EnhancedBlockchainService) ExecuteOrderOnChain(orderID uint64, price float64) (*types.Transaction, error) {
	logger.Info("开始在链上执行订单", logrus.Fields{
		"order_id": orderID,
		"price":    price,
	})

	tx, err := ebs.contract.ExecuteOrder(orderID, price)
	if err != nil {
		logger.Error("链上执行订单失败", err, logrus.Fields{
			"order_id": orderID,
		})
		return nil, err
	}

	logger.Info("链上订单执行成功", logrus.Fields{
		"order_id": orderID,
		"tx_hash":  tx.Hash().Hex(),
	})

	return tx, nil
}

// GetOrderFromChainEnhanced 从链上获取订单信息（增强版）
func (ebs *EnhancedBlockchainService) GetOrderFromChainEnhanced(orderID uint64) (map[string]interface{}, error) {
	logger.Debug("从链上获取订单信息", logrus.Fields{
		"order_id": orderID,
	})

	contractOrder, err := ebs.contract.GetOrder(orderID)
	if err != nil {
		logger.Error("从链上获取订单失败", err, logrus.Fields{
			"order_id": orderID,
		})
		return nil, err
	}

	// 转换为map格式返回
	orderData := map[string]interface{}{
		"orderId":     contractOrder.OrderId,
		"maker":       contractOrder.Maker.Hex(),
		"nftContract": contractOrder.NftContract.Hex(),
		"tokenId":     contractOrder.TokenId,
		"price":       contractOrder.Price,
		"amount":      contractOrder.Amount,
		"timestamp":   contractOrder.Timestamp,
		"expiration":  contractOrder.Expiration,
		"orderType":   contractOrder.OrderType,
		"status":      contractOrder.Status,
		"signature":   contractOrder.Signature,
	}

	logger.Debug("从链上获取订单成功", logrus.Fields{
		"order_id": orderID,
		"maker":    contractOrder.Maker.Hex(),
		"status":   contractOrder.Status,
	})

	return orderData, nil
}

// WaitForTransactionConfirmation 等待交易确认
func (ebs *EnhancedBlockchainService) WaitForTransactionConfirmation(tx *types.Transaction, timeout time.Duration) (*types.Receipt, error) {
	return ebs.contract.WaitForTransaction(tx, timeout)
}

// GetOrderCounter 获取订单计数器
func (ebs *EnhancedBlockchainService) GetOrderCounter() (*big.Int, error) {
	return ebs.contract.GetOrderCounter()
}

// SyncOrderFromChain 从链上同步订单到数据库（增强版）
func (ebs *EnhancedBlockchainService) SyncOrderFromChain(orderID uint64) (*models.Order, error) {
	logger.Info("开始从链上同步订单", logrus.Fields{
		"order_id": orderID,
	})

	// 从链上获取订单数据
	chainOrder, err := ebs.GetOrderFromChainEnhanced(orderID)
	if err != nil {
		return nil, err
	}

	// 检查数据库中是否已存在
	var existingOrder models.Order
	result := ebs.db.Where("order_id = ?", fmt.Sprintf("0x%x", chainOrder["orderId"])).First(&existingOrder)

	now := time.Now().Unix()
	maker := chainOrder["maker"].(string)
	tokenID := chainOrder["tokenId"].(*big.Int).String()
	price := float64(chainOrder["price"].(*big.Int).Int64()) / 1e18 // 转换为ETH单位
	collectionAddress := chainOrder["nftContract"].(string)

	order := &models.Order{
		OrderID:           fmt.Sprintf("0x%x", chainOrder["orderId"]),
		Maker:             maker,
		CollectionAddress: collectionAddress,
		TokenID:           tokenID,
		Price:             price,
		OrderType:         models.OrderType(chainOrder["orderType"].(uint8) + 1), // 合约从0开始，模型从1开始
		OrderStatus:       models.OrderStatus(chainOrder["status"].(uint8)),
		EventTime:         &now,
		CreateTime:        &now,
		UpdateTime:        &now,
		QuantityRemaining: chainOrder["amount"].(*big.Int).Int64(),
		Size:              chainOrder["amount"].(*big.Int).Int64(),
		CurrencyAddress:   "0x0000000000000000000000000000000000000000",
	}

	if result.Error == gorm.ErrRecordNotFound {
		// 创建新订单
		if err := ebs.db.Create(order).Error; err != nil {
			return nil, fmt.Errorf("保存同步订单失败: %v", err)
		}

		logger.Info("订单同步成功（新创建）", logrus.Fields{
			"order_id": orderID,
			"db_id":    order.ID,
		})
	} else if result.Error != nil {
		return nil, fmt.Errorf("查询现有订单失败: %v", result.Error)
	} else {
		// 更新现有订单
		if err := ebs.db.Model(&existingOrder).Updates(map[string]interface{}{
			"order_status": order.OrderStatus,
			"price":        order.Price,
			"update_time":  now,
		}).Error; err != nil {
			return nil, fmt.Errorf("更新同步订单失败: %v", err)
		}

		order = &existingOrder
		logger.Info("订单同步成功（更新现有）", logrus.Fields{
			"order_id": orderID,
			"db_id":    order.ID,
		})
	}

	return order, nil
}

// SyncAllOrdersFromChain 从链上同步所有订单
func (ebs *EnhancedBlockchainService) SyncAllOrdersFromChain() error {
	logger.Info("开始同步所有链上订单")

	// 获取订单计数器
	orderCounter, err := ebs.GetOrderCounter()
	if err != nil {
		return fmt.Errorf("获取订单计数器失败: %v", err)
	}

	totalOrders := orderCounter.Uint64()
	logger.Info("开始同步订单", logrus.Fields{
		"total_orders": totalOrders,
	})

	// 同步每个订单
	for i := uint64(1); i <= totalOrders; i++ {
		if _, err := ebs.SyncOrderFromChain(i); err != nil {
			logger.Error("同步订单失败", err, logrus.Fields{
				"order_id": i,
			})
			// 继续同步其他订单
			continue
		}

		// 每同步10个订单休息一下
		if i%10 == 0 {
			time.Sleep(100 * time.Millisecond)
			logger.Debug("同步进度", logrus.Fields{
				"synced":       i,
				"total_orders": totalOrders,
			})
		}
	}

	logger.Info("所有订单同步完成", logrus.Fields{
		"total_orders": totalOrders,
	})

	return nil
}

// Close 关闭增强区块链服务
func (ebs *EnhancedBlockchainService) Close() {
	if ebs.eventListener != nil {
		ebs.eventListener.Stop()
	}
	if ebs.contract != nil {
		ebs.contract.Close()
	}
	if ebs.BlockchainService != nil && ebs.BlockchainService.client != nil {
		ebs.BlockchainService.client.Close()
	}
	logger.Info("增强区块链服务已关闭")
}
