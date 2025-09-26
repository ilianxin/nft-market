package blockchain

import (
	"context"
	"fmt"
	"math/big"
	"nft-market/internal/logger"
	"nft-market/internal/models"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// EventListener 事件监听器
type EventListener struct {
	client          *ethclient.Client
	contractAddress common.Address
	contract        *NFTMarketplaceContract
	db              *gorm.DB
	stopChan        chan bool
	isRunning       bool
}

// EventHandler 事件处理接口
type EventHandler interface {
	HandleOrderCreated(event *OrderCreatedEvent, tx *types.Transaction, receipt *types.Receipt) error
	HandleOrderFilled(event *OrderFilledEvent, tx *types.Transaction, receipt *types.Receipt) error
	HandleOrderCancelled(event *OrderCancelledEvent, tx *types.Transaction, receipt *types.Receipt) error
}

// NewEventListener 创建新的事件监听器
func NewEventListener(client *ethclient.Client, contractAddress string, contract *NFTMarketplaceContract, db *gorm.DB) *EventListener {
	return &EventListener{
		client:          client,
		contractAddress: common.HexToAddress(contractAddress),
		contract:        contract,
		db:              db,
		stopChan:        make(chan bool),
		isRunning:       false,
	}
}

// Start 开始监听事件
func (el *EventListener) Start() error {
	if el.isRunning {
		return fmt.Errorf("事件监听器已在运行")
	}

	el.isRunning = true
	logger.Info("开始监听合约事件", logrus.Fields{
		"contract_address": el.contractAddress.Hex(),
	})

	go el.listenForEvents()
	return nil
}

// Stop 停止监听事件
func (el *EventListener) Stop() {
	if !el.isRunning {
		return
	}

	logger.Info("停止事件监听器")
	el.stopChan <- true
	el.isRunning = false
}

// listenForEvents 监听合约事件
func (el *EventListener) listenForEvents() {
	// 创建事件查询
	query := ethereum.FilterQuery{
		Addresses: []common.Address{el.contractAddress},
	}

	// 订阅日志
	logs := make(chan types.Log)
	sub, err := el.client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		logger.Error("订阅合约事件失败", err)
		return
	}
	defer sub.Unsubscribe()

	// 获取当前最新区块号
	latestBlock, err := el.client.BlockNumber(context.Background())
	if err != nil {
		logger.Error("获取最新区块号失败", err)
		return
	}

	// 同步历史事件
	go el.syncHistoricalEvents(latestBlock)

	logger.Info("事件监听器启动成功", logrus.Fields{
		"latest_block": latestBlock,
	})

	for {
		select {
		case err := <-sub.Err():
			logger.Error("事件订阅错误", err)
			return

		case vLog := <-logs:
			if err := el.processLog(vLog); err != nil {
				logger.Error("处理事件日志失败", err, logrus.Fields{
					"tx_hash":    vLog.TxHash.Hex(),
					"block_hash": vLog.BlockHash.Hex(),
					"log_index":  vLog.Index,
				})
			}

		case <-el.stopChan:
			logger.Info("收到停止信号，退出事件监听")
			return
		}
	}
}

// syncHistoricalEvents 同步历史事件
func (el *EventListener) syncHistoricalEvents(toBlock uint64) {
	logger.Info("开始同步历史事件", logrus.Fields{
		"to_block": toBlock,
	})

	// 获取数据库中最后处理的区块号
	var lastSyncedBlock uint64 = 0

	// 这里可以从数据库中获取上次同步的区块号
	// 暂时从创世块开始同步
	fromBlock := lastSyncedBlock

	// 分批同步，每次同步1000个区块
	batchSize := uint64(1000)

	for from := fromBlock; from <= toBlock; from += batchSize {
		to := from + batchSize - 1
		if to > toBlock {
			to = toBlock
		}

		if err := el.syncEventsBatch(from, to); err != nil {
			logger.Error("同步历史事件批次失败", err, logrus.Fields{
				"from_block": from,
				"to_block":   to,
			})
			// 继续处理下一批
			continue
		}

		logger.Debug("同步历史事件批次完成", logrus.Fields{
			"from_block": from,
			"to_block":   to,
		})

		// 短暂休息，避免过度请求
		time.Sleep(100 * time.Millisecond)
	}

	logger.Info("历史事件同步完成", logrus.Fields{
		"from_block": fromBlock,
		"to_block":   toBlock,
	})
}

// syncEventsBatch 同步事件批次
func (el *EventListener) syncEventsBatch(fromBlock, toBlock uint64) error {
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(int64(fromBlock)),
		ToBlock:   big.NewInt(int64(toBlock)),
		Addresses: []common.Address{el.contractAddress},
	}

	logs, err := el.client.FilterLogs(context.Background(), query)
	if err != nil {
		return fmt.Errorf("获取历史日志失败: %v", err)
	}

	for _, vLog := range logs {
		if err := el.processLog(vLog); err != nil {
			logger.Error("处理历史事件失败", err, logrus.Fields{
				"tx_hash":    vLog.TxHash.Hex(),
				"block_hash": vLog.BlockHash.Hex(),
			})
		}
	}

	return nil
}

// processLog 处理日志事件
func (el *EventListener) processLog(vLog types.Log) error {
	// 获取交易和收据信息
	tx, _, err := el.client.TransactionByHash(context.Background(), vLog.TxHash)
	if err != nil {
		return fmt.Errorf("获取交易信息失败: %v", err)
	}

	receipt, err := el.client.TransactionReceipt(context.Background(), vLog.TxHash)
	if err != nil {
		return fmt.Errorf("获取交易收据失败: %v", err)
	}

	// 根据事件主题处理不同类型的事件
	if len(vLog.Topics) == 0 {
		return fmt.Errorf("事件主题为空")
	}

	eventSignature := vLog.Topics[0]

	// OrderCreated事件签名: keccak256("OrderCreated(uint256,address,address,uint256,uint256,uint8)")
	orderCreatedSig := common.HexToHash("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef") // 实际需要计算

	// OrderFilled事件签名: keccak256("OrderFilled(uint256,address,address,uint256)")
	orderFilledSig := common.HexToHash("0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890") // 实际需要计算

	// OrderCancelled事件签名: keccak256("OrderCancelled(uint256,address)")
	orderCancelledSig := common.HexToHash("0xfedcba0987654321fedcba0987654321fedcba0987654321fedcba0987654321") // 实际需要计算

	switch eventSignature {
	case orderCreatedSig:
		return el.handleOrderCreatedEvent(vLog, tx, receipt)
	case orderFilledSig:
		return el.handleOrderFilledEvent(vLog, tx, receipt)
	case orderCancelledSig:
		return el.handleOrderCancelledEvent(vLog, tx, receipt)
	default:
		logger.Debug("未知事件类型", logrus.Fields{
			"event_signature": eventSignature.Hex(),
			"tx_hash":         vLog.TxHash.Hex(),
		})
		return nil
	}
}

// handleOrderCreatedEvent 处理订单创建事件
func (el *EventListener) handleOrderCreatedEvent(vLog types.Log, tx *types.Transaction, receipt *types.Receipt) error {
	// 解析事件数据
	event, err := el.parseOrderCreatedEvent(vLog)
	if err != nil {
		return fmt.Errorf("解析OrderCreated事件失败: %v", err)
	}

	logger.Info("处理OrderCreated事件", logrus.Fields{
		"order_id":     event.OrderId.String(),
		"maker":        event.Maker.Hex(),
		"nft_contract": event.NftContract.Hex(),
		"token_id":     event.TokenId.String(),
		"price":        event.Price.String(),
		"order_type":   event.OrderType,
		"tx_hash":      tx.Hash().Hex(),
	})

	// 检查订单是否已存在
	var existingOrder models.Order
	result := el.db.Where("order_id = ?", fmt.Sprintf("0x%x", event.OrderId)).First(&existingOrder)
	if result.Error == nil {
		logger.Debug("订单已存在，跳过处理", logrus.Fields{
			"order_id": event.OrderId.String(),
		})
		return nil
	}

	// 创建新的订单记录
	now := time.Now().Unix()
	order := &models.Order{
		OrderID:           fmt.Sprintf("0x%x", event.OrderId),
		OrderType:         models.OrderType(event.OrderType + 1), // 合约从0开始，模型从1开始
		OrderStatus:       models.OrderStatusActive,
		CollectionAddress: event.NftContract.Hex(),
		TokenID:           event.TokenId.String(),
		Price:             float64(event.Price.Int64()) / 1e18, // 转换为ETH
		Maker:             event.Maker.Hex(),
		EventTime:         &now,
		CreateTime:        &now,
		UpdateTime:        &now,
		QuantityRemaining: 1,
		Size:              1,
		CurrencyAddress:   "0x0000000000000000000000000000000000000000",
	}

	// 保存到数据库
	if err := el.db.Create(order).Error; err != nil {
		return fmt.Errorf("保存订单到数据库失败: %v", err)
	}

	// 创建对应的Item记录
	if err := el.createOrUpdateItemFromEvent(event, now); err != nil {
		logger.Error("从事件创建Item记录失败", err)
		// 不返回错误，因为订单已经创建成功
	}

	// 创建活动记录
	if err := el.createActivityFromOrderEvent(event, tx, receipt, now); err != nil {
		logger.Error("创建活动记录失败", err)
	}

	logger.Info("OrderCreated事件处理完成", logrus.Fields{
		"order_id": event.OrderId.String(),
		"db_id":    order.ID,
	})

	return nil
}

// handleOrderFilledEvent 处理订单成交事件
func (el *EventListener) handleOrderFilledEvent(vLog types.Log, tx *types.Transaction, receipt *types.Receipt) error {
	// 解析事件数据
	event, err := el.parseOrderFilledEvent(vLog)
	if err != nil {
		return fmt.Errorf("解析OrderFilled事件失败: %v", err)
	}

	logger.Info("处理OrderFilled事件", logrus.Fields{
		"order_id": event.OrderId.String(),
		"buyer":    event.Buyer.Hex(),
		"seller":   event.Seller.Hex(),
		"price":    event.Price.String(),
		"tx_hash":  tx.Hash().Hex(),
	})

	// 更新订单状态
	result := el.db.Model(&models.Order{}).
		Where("order_id = ?", fmt.Sprintf("0x%x", event.OrderId)).
		Update("order_status", models.OrderStatusFilled)

	if result.Error != nil {
		return fmt.Errorf("更新订单状态失败: %v", result.Error)
	}

	// 创建活动记录
	now := time.Now().Unix()
	sellerHex := event.Seller.Hex()
	buyerHex := event.Buyer.Hex()
	txHashHex := tx.Hash().Hex()
	activity := &models.Activity{
		ActivityType:    models.ActivityTypeBuy,
		Maker:           &sellerHex,
		Taker:           &buyerHex,
		Price:           float64(event.Price.Int64()) / 1e18,
		TxHash:          &txHashHex,
		BlockNumber:     int64(receipt.BlockNumber.Uint64()),
		EventTime:       &now,
		CreateTime:      &now,
		UpdateTime:      &now,
		CurrencyAddress: "1", // ETH
	}

	if err := el.db.Create(activity).Error; err != nil {
		logger.Error("创建成交活动记录失败", err)
	}

	logger.Info("OrderFilled事件处理完成", logrus.Fields{
		"order_id": event.OrderId.String(),
	})

	return nil
}

// handleOrderCancelledEvent 处理订单取消事件
func (el *EventListener) handleOrderCancelledEvent(vLog types.Log, tx *types.Transaction, receipt *types.Receipt) error {
	// 解析事件数据
	event, err := el.parseOrderCancelledEvent(vLog)
	if err != nil {
		return fmt.Errorf("解析OrderCancelled事件失败: %v", err)
	}

	logger.Info("处理OrderCancelled事件", logrus.Fields{
		"order_id": event.OrderId.String(),
		"maker":    event.Maker.Hex(),
		"tx_hash":  tx.Hash().Hex(),
	})

	// 更新订单状态
	result := el.db.Model(&models.Order{}).
		Where("order_id = ?", fmt.Sprintf("0x%x", event.OrderId)).
		Update("order_status", models.OrderStatusCancelled)

	if result.Error != nil {
		return fmt.Errorf("更新订单状态失败: %v", result.Error)
	}

	// 创建活动记录
	now := time.Now().Unix()
	makerHex := event.Maker.Hex()
	txHashHex := tx.Hash().Hex()
	activity := &models.Activity{
		ActivityType: models.ActivityTypeCancelListing,
		Maker:        &makerHex,
		TxHash:       &txHashHex,
		BlockNumber:  int64(receipt.BlockNumber.Uint64()),
		EventTime:    &now,
		CreateTime:   &now,
		UpdateTime:   &now,
	}

	if err := el.db.Create(activity).Error; err != nil {
		logger.Error("创建取消活动记录失败", err)
	}

	logger.Info("OrderCancelled事件处理完成", logrus.Fields{
		"order_id": event.OrderId.String(),
	})

	return nil
}

// parseOrderCreatedEvent 解析OrderCreated事件
func (el *EventListener) parseOrderCreatedEvent(vLog types.Log) (*OrderCreatedEvent, error) {
	// 这里需要根据实际的ABI来解析事件
	// 简化处理，假设事件数据在Topics和Data中
	if len(vLog.Topics) < 4 {
		return nil, fmt.Errorf("OrderCreated事件Topics数量不足")
	}

	event := &OrderCreatedEvent{
		OrderId:     new(big.Int).SetBytes(vLog.Topics[1][:]),
		Maker:       common.BytesToAddress(vLog.Topics[2][:]),
		NftContract: common.BytesToAddress(vLog.Topics[3][:]),
	}

	// 解析Data中的其他字段（TokenId, Price, OrderType）
	// 这里需要根据实际ABI进行解析
	// 简化处理
	if len(vLog.Data) >= 96 { // 3个uint256字段
		event.TokenId = new(big.Int).SetBytes(vLog.Data[0:32])
		event.Price = new(big.Int).SetBytes(vLog.Data[32:64])
		event.OrderType = uint8(vLog.Data[95]) // 最后一个字节
	}

	return event, nil
}

// parseOrderFilledEvent 解析OrderFilled事件
func (el *EventListener) parseOrderFilledEvent(vLog types.Log) (*OrderFilledEvent, error) {
	if len(vLog.Topics) < 4 {
		return nil, fmt.Errorf("OrderFilled事件Topics数量不足")
	}

	event := &OrderFilledEvent{
		OrderId: new(big.Int).SetBytes(vLog.Topics[1][:]),
		Buyer:   common.BytesToAddress(vLog.Topics[2][:]),
		Seller:  common.BytesToAddress(vLog.Topics[3][:]),
	}

	// 解析Price
	if len(vLog.Data) >= 32 {
		event.Price = new(big.Int).SetBytes(vLog.Data[0:32])
	}

	return event, nil
}

// parseOrderCancelledEvent 解析OrderCancelled事件
func (el *EventListener) parseOrderCancelledEvent(vLog types.Log) (*OrderCancelledEvent, error) {
	if len(vLog.Topics) < 3 {
		return nil, fmt.Errorf("OrderCancelled事件Topics数量不足")
	}

	event := &OrderCancelledEvent{
		OrderId: new(big.Int).SetBytes(vLog.Topics[1][:]),
		Maker:   common.BytesToAddress(vLog.Topics[2][:]),
	}

	return event, nil
}

// createOrUpdateItemFromEvent 从事件创建或更新Item记录
func (el *EventListener) createOrUpdateItemFromEvent(event *OrderCreatedEvent, now int64) error {
	var item models.Item

	// 检查Item是否已存在
	err := el.db.Where("collection_address = ? AND token_id = ?",
		event.NftContract.Hex(), event.TokenId.String()).First(&item).Error

	if err == gorm.ErrRecordNotFound {
		// 创建新的Item记录
		makerHex := event.Maker.Hex()
		contractHex := event.NftContract.Hex()
		item = models.Item{
			ChainID:           models.ChainIDEthereum,
			TokenID:           event.TokenId.String(),
			Name:              fmt.Sprintf("NFT #%s", event.TokenId.String()),
			Owner:             &makerHex,
			CollectionAddress: &contractHex,
			Creator:           event.Maker.Hex(),
			Supply:            1,
			CreateTime:        &now,
			UpdateTime:        &now,
		}

		// 如果是上架订单，设置价格
		if event.OrderType == 0 { // LimitSell
			price := float64(event.Price.Int64()) / 1e18
			item.ListPrice = &price
			item.ListTime = &now
		}

		return el.db.Create(&item).Error
	} else if err != nil {
		return err
	} else {
		// 更新现有Item记录
		updateData := map[string]interface{}{
			"update_time": now,
		}

		if event.OrderType == 0 { // LimitSell
			price := float64(event.Price.Int64()) / 1e18
			updateData["list_price"] = price
			updateData["list_time"] = now
		}

		return el.db.Model(&item).Updates(updateData).Error
	}
}

// createActivityFromOrderEvent 从订单事件创建活动记录
func (el *EventListener) createActivityFromOrderEvent(event *OrderCreatedEvent, tx *types.Transaction, receipt *types.Receipt, now int64) error {
	var activityType models.ActivityType

	switch event.OrderType {
	case 0: // LimitSell
		activityType = models.ActivityTypeList
	case 1: // LimitBuy
		activityType = models.ActivityTypeMakeOffer
	default:
		activityType = models.ActivityTypeList
	}

	makerHex := event.Maker.Hex()
	contractHex := event.NftContract.Hex()
	tokenIdStr := event.TokenId.String()
	txHashHex := tx.Hash().Hex()
	activity := &models.Activity{
		ActivityType:      activityType,
		Maker:             &makerHex,
		CollectionAddress: &contractHex,
		TokenID:           &tokenIdStr,
		Price:             float64(event.Price.Int64()) / 1e18,
		TxHash:            &txHashHex,
		BlockNumber:       int64(receipt.BlockNumber.Uint64()),
		EventTime:         &now,
		CreateTime:        &now,
		UpdateTime:        &now,
		CurrencyAddress:   "1", // ETH
	}

	return el.db.Create(activity).Error
}
