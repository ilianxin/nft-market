package main

import (
	"fmt"
	"nft-market/internal/config"
	"nft-market/internal/database"
	"nft-market/internal/logger"
	"nft-market/internal/models"
	"time"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	// 加载环境变量
	if err := godotenv.Load(); err != nil {
		fmt.Println("未找到.env文件，使用默认配置")
	}

	// 初始化日志系统
	logConfig := logger.DefaultConfig()
	if err := logger.Init(logConfig); err != nil {
		panic("日志系统初始化失败: " + err.Error())
	}

	// 加载配置
	cfg := config.Load()

	// 初始化数据库连接
	db, err := database.Init(cfg.DatabaseURL)
	if err != nil {
		panic("数据库连接失败: " + err.Error())
	}

	logger.Info("开始设置测试数据...")

	// 创建测试集合
	collection := &models.Collection{
		ChainID:     models.ChainIDEthereum,
		Address:     "0x1234567890123456789012345678901234567890",
		Name:        "测试NFT集合",
		Symbol:      "TEST",
		Description: "用于测试的NFT集合",
		Creator:     "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266",
		CreateTime:  timePtr(time.Now().Unix()),
		UpdateTime:  timePtr(time.Now().Unix()),
	}

	if err := db.Create(collection).Error; err != nil {
		logger.Warn("集合创建失败（可能已存在）", logrus.Fields{"error": err.Error()})
	} else {
		logger.Info("测试集合创建成功", logrus.Fields{"collection_id": collection.ID})
	}

	// 创建测试物品
	now := time.Now().Unix()
	items := []*models.Item{
		{
			ChainID:           models.ChainIDEthereum,
			TokenID:           "1001",
			Name:              "测试NFT #1001",
			Description:       "第一个测试NFT",
			Owner:             strPtr("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"),
			CollectionAddress: strPtr("0x1234567890123456789012345678901234567890"),
			Creator:           "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266",
			Supply:            1,
			CreateTime:        &now,
			UpdateTime:        &now,
		},
		{
			ChainID:           models.ChainIDEthereum,
			TokenID:           "1002",
			Name:              "测试NFT #1002",
			Description:       "第二个测试NFT",
			Owner:             strPtr("0x70997970C51812dc3A010C7d01b50e0d17dc79C8"),
			CollectionAddress: strPtr("0x1234567890123456789012345678901234567890"),
			Creator:           "0x70997970C51812dc3A010C7d01b50e0d17dc79C8",
			Supply:            1,
			CreateTime:        &now,
			UpdateTime:        &now,
		},
	}

	for _, item := range items {
		if err := db.Create(item).Error; err != nil {
			logger.Warn("物品创建失败（可能已存在）", logrus.Fields{
				"token_id": item.TokenID,
				"error":    err.Error(),
			})
		} else {
			logger.Info("测试物品创建成功", logrus.Fields{
				"item_id":  item.ID,
				"token_id": item.TokenID,
			})
		}
	}

	// 创建测试订单
	orders := []*models.Order{
		{
			OrderID:           fmt.Sprintf("0x%x", time.Now().UnixNano()),
			OrderType:         models.OrderTypeListing,
			OrderStatus:       models.OrderStatusActive,
			CollectionAddress: "0x1234567890123456789012345678901234567890",
			TokenID:           "1001",
			Price:             0.5,
			Maker:             "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266",
			QuantityRemaining: 1,
			Size:              1,
			CurrencyAddress:   "0x0000000000000000000000000000000000000000",
			EventTime:         &now,
			CreateTime:        &now,
			UpdateTime:        &now,
		},
		{
			OrderID:           fmt.Sprintf("0x%x", time.Now().UnixNano()+1),
			OrderType:         models.OrderTypeListing,
			OrderStatus:       models.OrderStatusActive,
			CollectionAddress: "0x1234567890123456789012345678901234567890",
			TokenID:           "1002",
			Price:             0.3,
			Maker:             "0x70997970C51812dc3A010C7d01b50e0d17dc79C8",
			QuantityRemaining: 1,
			Size:              1,
			CurrencyAddress:   "0x0000000000000000000000000000000000000000",
			EventTime:         &now,
			CreateTime:        &now,
			UpdateTime:        &now,
		},
	}

	for _, order := range orders {
		if err := db.Create(order).Error; err != nil {
			logger.Warn("订单创建失败（可能已存在）", logrus.Fields{
				"order_id": order.OrderID,
				"error":    err.Error(),
			})
		} else {
			logger.Info("测试订单创建成功", logrus.Fields{
				"db_id":    order.ID,
				"order_id": order.OrderID,
				"maker":    order.Maker,
				"price":    order.Price,
			})
		}
	}

	logger.Info("测试数据设置完成！")
	fmt.Println("测试数据设置完成！现在可以进行购买测试了。")
}

func strPtr(s string) *string {
	return &s
}

func timePtr(t int64) *int64 {
	return &t
}
