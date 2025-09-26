package main

import (
	"fmt"
	"nft-market/internal/config"
	"nft-market/internal/database"
	"nft-market/internal/models"
	"nft-market/internal/services"

	"github.com/joho/godotenv"
)

func main() {
	// 加载环境变量
	if err := godotenv.Load(); err != nil {
		fmt.Println("未找到.env文件，使用默认配置")
	}

	// 加载配置
	cfg := config.Load()

	// 初始化数据库
	db, err := database.Init(cfg.DatabaseURL)
	if err != nil {
		panic("数据库连接失败: " + err.Error())
	}

	// 初始化区块链服务（模拟）
	blockchainService := &services.BlockchainService{}

	// 初始化订单服务
	orderService := services.NewOrderService(db, blockchainService)

	// 测试创建订单请求
	req := &models.CreateOrderRequest{
		CollectionAddress: "0x1234567890123456789012345678901234567890",
		TokenID:           "1",
		OrderType:         models.OrderTypeListing,
		Price:             0.1,
		QuantityRemaining: 1,
		Size:              1,
		CurrencyAddress:   "0x0000000000000000000000000000000000000000",
	}

	maker := "0x9876543210987654321098765432109876543210"

	fmt.Println("创建测试订单...")
	order, err := orderService.CreateOrder(req, maker)
	if err != nil {
		fmt.Printf("创建订单失败: %v\n", err)
		return
	}

	fmt.Printf("订单创建成功! ID: %d, OrderID: %s\n", order.ID, order.OrderID)

	// 验证Item是否创建
	var item models.Item
	err = db.Where("collection_address = ? AND token_id = ?", req.CollectionAddress, req.TokenID).First(&item).Error
	if err != nil {
		fmt.Printf("Item记录未找到: %v\n", err)
		return
	}

	fmt.Printf("Item记录创建成功! ID: %d, Name: %s, Owner: %s\n", item.ID, item.Name, *item.Owner)
	fmt.Printf("ListPrice: %v, ListTime: %v\n", item.ListPrice, item.ListTime)

	fmt.Println("测试完成！订单和Item都创建成功。")
}
