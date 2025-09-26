package services

import (
	"fmt"
	"math/big"
	"nft-market/internal/models"
	"time"

	"gorm.io/gorm"
)

// OrderService 订单服务
type OrderService struct {
	db                *gorm.DB
	blockchainService *BlockchainService
}

// NewOrderService 创建新的订单服务
func NewOrderService(db *gorm.DB, blockchainService *BlockchainService) *OrderService {
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
		Maker:             &maker,
		QuantityRemaining: req.QuantityRemaining,
		Size:              req.Size,
		CurrencyAddress:   req.CurrencyAddress,
		EventTime:         &now,
		ExpireTime:        req.ExpireTime,
		CreateTime:        &now,
		UpdateTime:        &now,
	}

	// 保存到数据库
	if err := os.db.Create(order).Error; err != nil {
		return nil, fmt.Errorf("保存订单到数据库失败: %v", err)
	}

	// 如果需要，可以在这里调用区块链创建订单
	// 这里简化处理，实际项目中可能需要监听区块链事件

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
	if order.Maker == nil || *order.Maker != userAddress {
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
		Maker:             &maker,
		CollectionAddress: &collectionAddress,
		TokenID:           &tokenID,
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

// validateCreateOrderRequest 验证创建订单请求
func (os *OrderService) validateCreateOrderRequest(req *models.CreateOrderRequest) error {
	if req.CollectionAddress == nil || *req.CollectionAddress == "" {
		return fmt.Errorf("集合地址不能为空")
	}

	if req.TokenID == nil || *req.TokenID == "" {
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
