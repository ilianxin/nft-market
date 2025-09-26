package main

import (
	"log"
	"nft-market/internal/api"
	"nft-market/internal/config"
	"nft-market/internal/database"
	"nft-market/internal/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// 加载环境变量
	if err := godotenv.Load(); err != nil {
		log.Println("未找到.env文件，使用默认配置")
	}

	// 加载配置
	cfg := config.Load()

	// 初始化数据库
	db, err := database.Init(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}

	// 初始化区块链服务
	blockchainService, err := services.NewBlockchainService(cfg.EthereumRPC, cfg.ContractAddress, cfg.PrivateKey)
	if err != nil {
		log.Fatal("区块链服务初始化失败:", err)
	}

	// 初始化服务层
	orderService := services.NewOrderService(db, blockchainService)
	nftService := services.NewNFTService(db, blockchainService)
	collectionService := services.NewCollectionService(db)
	itemService := services.NewItemService(db)
	activityService := services.NewActivityService(db)

	// 设置Gin模式
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建路由器
	router := gin.Default()

	// 配置CORS
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:3000", "http://localhost:3001"}
	corsConfig.AllowCredentials = true
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	router.Use(cors.New(corsConfig))

	// 设置API路由
	api.SetupRoutes(router, orderService, nftService, collectionService, itemService, activityService)

	// 启动服务器
	log.Printf("服务器启动在端口 %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal("服务器启动失败:", err)
	}
}
