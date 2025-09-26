package services

import (
	"fmt"
	"math/big"
	"nft-market/internal/logger"
	"nft-market/internal/models"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// OrderService 订单服务
type OrderService struct {
	db                *gorm.DB
	blockchainService *EnhancedBlockchainService
}

// NewOrderService 创建新的订单服务
func NewOrderService(db *gorm.DB, blockchainService *EnhancedBlockchainService) *OrderService {
	return &OrderService{
		db:                db,
		blockchainService: blockchainService,
	}
}

// CreateOrder 创建订单
func (os *OrderService) CreateOrder(req *models.CreateOrderRequest, maker string) (*models.Order, error) {
	// 验证输入
	if err := os.validateCreateOrderRequest(req); err != nil {
		return nil, err
	}

	// 生成订单ID（这里简化处理，实际项目中应该从区块链获取）
	orderID := fmt.Sprintf("0x%x", time.Now().UnixNano())
	now := time.Now().Unix()

	// 创建订单模型
	order := &models.Order{
		OrderID:           orderID,
		OrderType:         req.OrderType,
		OrderStatus:       models.OrderStatusActive,
		CollectionAddress: req.CollectionAddress,
		TokenID:           req.TokenID,
		Price:             req.Price,
		Maker:             maker,
		QuantityRemaining: req.QuantityRemaining,
		Size:              req.Size,
		CurrencyAddress:   req.CurrencyAddress,
		EventTime:         &now,
		ExpireTime:        req.ExpireTime,
		CreateTime:        &now,
		UpdateTime:        &now,
	}

	// 开始数据库事务
	tx := os.db.Begin()
	if tx.Error != nil {
		return nil, fmt.Errorf("开始事务失败: %v", tx.Error)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 保存订单到数据库
	if err := tx.Create(order).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("保存订单到数据库失败: %v", err)
	}

	// 创建或更新Item记录
	if err := os.createOrUpdateItem(tx, req, maker, now); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("创建Item记录失败: %v", err)
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("提交事务失败: %v", err)
	}

	// 在区块链上创建订单
	if os.blockchainService != nil {
		logger.Info("开始在区块链上创建订单", logrus.Fields{
			"order_id": order.OrderID,
		})

		// 异步创建链上订单，避免阻塞用户操作
		go func() {
			tx, err := os.blockchainService.CreateOrderOnChain(order)
			if err != nil {
				logger.Error("链上创建订单失败", err, logrus.Fields{
					"order_id": order.OrderID,
				})
				return
			}

			// 等待交易确认
			receipt, err := os.blockchainService.WaitForTransactionConfirmation(tx, 5*time.Minute)
			if err != nil {
				logger.Error("等待交易确认失败", err, logrus.Fields{
					"order_id": order.OrderID,
					"tx_hash":  tx.Hash().Hex(),
				})
				return
			}

			logger.Info("链上订单创建并确认成功", logrus.Fields{
				"order_id":     order.OrderID,
				"tx_hash":      tx.Hash().Hex(),
				"block_number": receipt.BlockNumber.String(),
			})
		}()
	}

	return order, nil
}

// GetOrderByID 根据ID获取订单
func (os *OrderService) GetOrderByID(id uint) (*models.Order, error) {
	var order models.Order
	if err := os.db.First(&order, id).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

// GetOrderByOrderID 根据链上订单ID获取订单
func (os *OrderService) GetOrderByOrderID(orderID string) (*models.Order, error) {
	var order models.Order
	if err := os.db.Where("order_id = ?", orderID).First(&order).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

// GetUserOrders 获取用户订单
func (os *OrderService) GetUserOrders(userAddress string, page, pageSize int, status string) (*models.OrderListResponse, error) {
	var orders []models.Order
	var total int64

	query := os.db.Model(&models.Order{}).Where("maker = ?", userAddress)

	if status != "" {
		query = query.Where("order_status = ?", status)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&orders).Error; err != nil {
		return nil, err
	}

	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	return &models.OrderListResponse{
		Orders:     orders,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// GetNFTOrders 获取NFT的所有订单
func (os *OrderService) GetNFTOrders(collectionAddress, tokenID string, page, pageSize int) (*models.OrderListResponse, error) {
	var orders []models.Order
	var total int64

	query := os.db.Model(&models.Order{}).Where("collection_address = ? AND token_id = ?", collectionAddress, tokenID)

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&orders).Error; err != nil {
		return nil, err
	}

	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	return &models.OrderListResponse{
		Orders:     orders,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// GetActiveOrders 获取活跃订单
func (os *OrderService) GetActiveOrders(page, pageSize int, orderType string) (*models.OrderListResponse, error) {
	var orders []models.Order
	var total int64

	now := time.Now().Unix()
	query := os.db.Model(&models.Order{}).Where("order_status = ? AND (expire_time IS NULL OR expire_time > ?)",
		models.OrderStatusActive, now)

	if orderType != "" {
		query = query.Where("order_type = ?", orderType)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&orders).Error; err != nil {
		return nil, err
	}

	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	return &models.OrderListResponse{
		Orders:     orders,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// CancelOrder 取消订单
func (os *OrderService) CancelOrder(orderID uint64, userAddress string) error {
	// 查找订单
	var order models.Order
	if err := os.db.First(&order, orderID).Error; err != nil {
		return fmt.Errorf("订单不存在: %v", err)
	}

	// 验证权限
	if order.Maker != userAddress {
		return fmt.Errorf("无权限取消此订单")
	}

	// 验证状态
	if order.OrderStatus != models.OrderStatusActive {
		return fmt.Errorf("订单状态不允许取消")
	}

	// 更新状态
	now := time.Now().Unix()
	order.OrderStatus = models.OrderStatusCancelled
	order.UpdateTime = &now

	if err := os.db.Save(&order).Error; err != nil {
		return fmt.Errorf("更新订单状态失败: %v", err)
	}

	// 在区块链上取消订单
	if os.blockchainService != nil {
		logger.Info("开始在区块链上取消订单", logrus.Fields{
			"order_id": orderID,
		})

		// 异步取消链上订单
		go func() {
			tx, err := os.blockchainService.CancelOrderOnChain(orderID)
			if err != nil {
				logger.Error("链上取消订单失败", err, logrus.Fields{
					"order_id": orderID,
				})
				return
			}

			// 等待交易确认
			receipt, err := os.blockchainService.WaitForTransactionConfirmation(tx, 2*time.Minute)
			if err != nil {
				logger.Error("等待取消交易确认失败", err, logrus.Fields{
					"order_id": orderID,
					"tx_hash":  tx.Hash().Hex(),
				})
				return
			}

			logger.Info("链上订单取消并确认成功", logrus.Fields{
				"order_id":     orderID,
				"tx_hash":      tx.Hash().Hex(),
				"block_number": receipt.BlockNumber.String(),
			})
		}()
	}

	return nil
}

// PurchaseOrder 购买订单
func (os *OrderService) PurchaseOrder(orderID uint64, buyerAddress string, offeredPrice float64) error {
	logger.Info("开始购买订单", logrus.Fields{
		"order_id":      orderID,
		"buyer_address": buyerAddress,
		"offered_price": offeredPrice,
	})

	// 查找订单
	var order models.Order
	if err := os.db.First(&order, orderID).Error; err != nil {
		logger.Error("订单查找失败", err, logrus.Fields{
			"order_id": orderID,
		})
		return fmt.Errorf("订单不存在: %v", err)
	}

	logger.Info("找到订单", logrus.Fields{
		"order_id":     orderID,
		"order_maker":  order.Maker,
		"order_price":  order.Price,
		"order_status": order.OrderStatus,
		"order_type":   order.OrderType,
	})

	// 验证订单状态
	if order.OrderStatus != models.OrderStatusActive {
		return fmt.Errorf("订单状态不允许购买，当前状态: %d", order.OrderStatus)
	}

	// 验证订单类型（只能购买卖单）
	if order.OrderType != models.OrderTypeListing {
		return fmt.Errorf("只能购买上架订单（listing），当前订单类型: %d", order.OrderType)
	}

	// 验证买家不是卖家
	if order.Maker == buyerAddress {
		return fmt.Errorf("不能购买自己的订单")
	}

	// 验证价格（如果提供了价格，必须匹配或更高）
	if offeredPrice > 0 && offeredPrice < order.Price {
		return fmt.Errorf("出价过低，订单价格: %.6f ETH，您的出价: %.6f ETH", order.Price, offeredPrice)
	}

	// 检查订单是否过期
	if order.ExpireTime != nil && time.Now().Unix() > *order.ExpireTime {
		return fmt.Errorf("订单已过期")
	}

	logger.Info("开始处理订单购买", logrus.Fields{
		"order_id":      orderID,
		"buyer":         buyerAddress,
		"seller":        order.Maker,
		"price":         order.Price,
		"offered_price": offeredPrice,
	})

	// 开始数据库事务
	tx := os.db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("开始事务失败: %v", tx.Error)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 更新订单状态
	now := time.Now().Unix()
	updateData := map[string]interface{}{
		"order_status": models.OrderStatusFilled,
		"taker":        buyerAddress,
		"update_time":  now,
	}

	if err := tx.Model(&order).Updates(updateData).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("更新订单状态失败: %v", err)
	}

	// 更新Item的拥有者
	if err := os.updateItemOwner(tx, order.CollectionAddress, order.TokenID, buyerAddress, now); err != nil {
		tx.Rollback()
		return fmt.Errorf("更新物品拥有者失败: %v", err)
	}

	// 创建交易活动记录
	if err := os.createPurchaseActivity(tx, &order, buyerAddress, now); err != nil {
		tx.Rollback()
		return fmt.Errorf("创建交易活动记录失败: %v", err)
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("提交事务失败: %v", err)
	}

	// 在区块链上执行订单（异步）
	if os.blockchainService != nil {
		logger.Info("开始在区块链上执行订单", logrus.Fields{
			"order_id": orderID,
		})

		// 异步执行链上订单
		go func() {
			finalPrice := order.Price
			if offeredPrice > 0 {
				finalPrice = offeredPrice
			}

			tx, err := os.blockchainService.ExecuteOrderOnChain(orderID, finalPrice)
			if err != nil {
				logger.Error("链上执行订单失败", err, logrus.Fields{
					"order_id": orderID,
				})
				return
			}

			// 等待交易确认
			receipt, err := os.blockchainService.WaitForTransactionConfirmation(tx, 5*time.Minute)
			if err != nil {
				logger.Error("等待执行交易确认失败", err, logrus.Fields{
					"order_id": orderID,
					"tx_hash":  tx.Hash().Hex(),
				})
				return
			}

			logger.Info("链上订单执行并确认成功", logrus.Fields{
				"order_id":     orderID,
				"tx_hash":      tx.Hash().Hex(),
				"block_number": receipt.BlockNumber.String(),
			})
		}()
	}

	logger.Info("订单购买成功", logrus.Fields{
		"order_id": orderID,
		"buyer":    buyerAddress,
		"seller":   order.Maker,
		"price":    order.Price,
	})

	return nil
}

// updateItemOwner 更新物品拥有者
func (os *OrderService) updateItemOwner(tx *gorm.DB, collectionAddress, tokenID, newOwner string, now int64) error {
	updateData := map[string]interface{}{
		"owner":       newOwner,
		"list_price":  nil, // 清除上架价格
		"list_time":   nil, // 清除上架时间
		"update_time": now,
	}

	result := tx.Model(&models.Item{}).
		Where("collection_address = ? AND token_id = ?", collectionAddress, tokenID).
		Updates(updateData)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		logger.Warn("未找到要更新的Item记录", logrus.Fields{
			"collection_address": collectionAddress,
			"token_id":           tokenID,
		})
	} else {
		logger.Info("Item拥有者更新成功", logrus.Fields{
			"collection_address": collectionAddress,
			"token_id":           tokenID,
			"new_owner":          newOwner,
		})
	}

	return nil
}

// createPurchaseActivity 创建购买活动记录
func (os *OrderService) createPurchaseActivity(tx *gorm.DB, order *models.Order, buyer string, now int64) error {
	activity := &models.Activity{
		ActivityType:      models.ActivityTypeBuy,
		Maker:             &order.Maker,
		Taker:             &buyer,
		CollectionAddress: &order.CollectionAddress,
		TokenID:           &order.TokenID,
		Price:             order.Price,
		BlockNumber:       0, // 链上确认后会更新
		EventTime:         &now,
		CreateTime:        &now,
		UpdateTime:        &now,
		CurrencyAddress:   "1", // ETH
	}

	if err := tx.Create(activity).Error; err != nil {
		return err
	}

	logger.Info("购买活动记录创建成功", logrus.Fields{
		"order_id": order.ID,
		"buyer":    buyer,
		"seller":   order.Maker,
		"price":    order.Price,
	})

	return nil
}

// SyncOrderFromChain 从链上同步订单信息
func (os *OrderService) SyncOrderFromChain(orderID uint64) (*models.Order, error) {
	// 从链上获取订单信息
	chainOrder, err := os.blockchainService.GetOrderFromChain(big.NewInt(int64(orderID)))
	if err != nil {
		return nil, fmt.Errorf("从链上获取订单失败: %v", err)
	}

	// 检查是否已存在
	var existingOrder models.Order
	err = os.db.Where("order_id = ?", orderID).First(&existingOrder).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// 构建订单模型
	now := time.Now().Unix()
	maker := chainOrder["maker"].(string)
	tokenID := chainOrder["tokenId"].(*big.Int).String()
	price := float64(chainOrder["price"].(*big.Int).Int64()) / 1e18 // 转换为ETH单位
	collectionAddress := chainOrder["nftContract"].(string)

	order := &models.Order{
		OrderID:           fmt.Sprintf("0x%x", orderID),
		Maker:             maker,
		CollectionAddress: collectionAddress,
		TokenID:           tokenID,
		Price:             price,
		OrderType:         models.OrderType(chainOrder["orderType"].(uint8)),
		OrderStatus:       models.OrderStatus(chainOrder["status"].(uint8)),
		EventTime:         &now,
		CreateTime:        &now,
		UpdateTime:        &now,
	}

	if err == gorm.ErrRecordNotFound {
		// 创建新订单
		if err := os.db.Create(order).Error; err != nil {
			return nil, err
		}
	} else {
		// 更新现有订单
		order.ID = existingOrder.ID
		if err := os.db.Save(order).Error; err != nil {
			return nil, err
		}
	}

	return order, nil
}

// MarkExpiredOrders 标记过期订单
func (os *OrderService) MarkExpiredOrders() error {
	now := time.Now().Unix()
	result := os.db.Model(&models.Order{}).
		Where("order_status = ? AND expire_time IS NOT NULL AND expire_time <= ?", models.OrderStatusActive, now).
		Update("order_status", models.OrderStatusExpired)

	if result.Error != nil {
		return result.Error
	}

	fmt.Printf("标记了 %d 个过期订单\n", result.RowsAffected)
	return nil
}

// createOrUpdateItem 创建或更新Item记录
func (os *OrderService) createOrUpdateItem(tx *gorm.DB, req *models.CreateOrderRequest, maker string, now int64) error {
	var item models.Item

	// 检查Item是否已存在
	err := tx.Where("collection_address = ? AND token_id = ?", req.CollectionAddress, req.TokenID).First(&item).Error

	if err == gorm.ErrRecordNotFound {
		// 创建新的Item记录
		item = models.Item{
			ChainID:           models.ChainIDEthereum,
			TokenID:           req.TokenID,
			Name:              fmt.Sprintf("NFT #%s", req.TokenID), // 默认名称
			Owner:             &maker,
			CollectionAddress: &req.CollectionAddress,
			Creator:           maker,
			Supply:            1, // 默认供应量为1
			CreateTime:        &now,
			UpdateTime:        &now,
		}

		// 如果是上架订单，设置上架价格和时间
		if req.OrderType == models.OrderTypeListing {
			item.ListPrice = &req.Price
			item.ListTime = &now
		}

		if err := tx.Create(&item).Error; err != nil {
			return fmt.Errorf("创建Item失败: %v", err)
		}

		logger.Info("创建新Item记录成功", logrus.Fields{
			"collection_address": req.CollectionAddress,
			"token_id":           req.TokenID,
			"owner":              maker,
			"order_type":         req.OrderType,
		})
	} else if err != nil {
		return fmt.Errorf("查询Item失败: %v", err)
	} else {
		// 更新现有Item记录
		updateData := map[string]interface{}{
			"update_time": now,
		}

		// 如果是上架订单，更新上架价格和时间
		if req.OrderType == models.OrderTypeListing {
			updateData["list_price"] = req.Price
			updateData["list_time"] = now
			updateData["owner"] = maker // 更新拥有者
		}

		if err := tx.Model(&item).Updates(updateData).Error; err != nil {
			return fmt.Errorf("更新Item失败: %v", err)
		}

		logger.Info("更新Item记录成功", logrus.Fields{
			"collection_address": req.CollectionAddress,
			"token_id":           req.TokenID,
			"owner":              maker,
			"order_type":         req.OrderType,
			"update_data":        updateData,
		})
	}

	return nil
}

// validateCreateOrderRequest 验证创建订单请求
func (os *OrderService) validateCreateOrderRequest(req *models.CreateOrderRequest) error {
	if req.CollectionAddress == "" {
		return fmt.Errorf("集合地址不能为空")
	}

	if req.TokenID == "" {
		return fmt.Errorf("Token ID不能为空")
	}

	if req.OrderType == models.OrderTypeListing || req.OrderType == models.OrderTypeOffer {
		if req.Price <= 0 {
			return fmt.Errorf("订单必须指定有效价格")
		}
	}

	if req.ExpireTime != nil && *req.ExpireTime <= time.Now().Unix() {
		return fmt.Errorf("过期时间必须在未来")
	}

	return nil
}
