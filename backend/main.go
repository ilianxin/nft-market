package main

import (
	"nft-market/internal/api"
	"nft-market/internal/config"
	"nft-market/internal/database"
	"nft-market/internal/logger"
	"nft-market/internal/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// 加载环境变量
	if err := godotenv.Load(); err != nil {
		// 使用临时日志，因为logger还没初始化
		println("未找到.env文件，使用默认配置")
	}

	// 初始化日志系统
	logConfig := logger.DefaultConfig()
	if err := logger.Init(logConfig); err != nil {
		panic("日志系统初始化失败: " + err.Error())
	}

	// 加载配置
	cfg := config.Load()
	logger.Info("配置加载完成", map[string]interface{}{
		"port": cfg.Port,
		"env":  cfg.Environment,
	})

	// 初始化数据库
	db, err := database.Init(cfg.DatabaseURL)
	if err != nil {
		logger.Error("数据库连接失败", err)
		panic(err)
	}
	logger.Info("数据库连接成功")

	// 初始化增强的区块链服务
	blockchainService, err := services.NewEnhancedBlockchainService(cfg.EthereumRPC, cfg.ContractAddress, cfg.PrivateKey, db)
	if err != nil {
		logger.Error("增强区块链服务初始化失败", err)
		panic(err)
	}
	logger.Info("增强区块链服务初始化成功", map[string]interface{}{
		"rpc_url":          cfg.EthereumRPC,
		"contract_address": cfg.ContractAddress,
	})

	// 初始化服务层
	orderService := services.NewOrderService(db, blockchainService)
	nftService := services.NewNFTService(db, blockchainService)
	collectionService := services.NewCollectionService(db)
	itemService := services.NewItemService(db)
	activityService := services.NewActivityService(db)
	logger.Info("所有服务初始化完成")

	// 设置Gin模式
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建路由器
	router := gin.New()

	// 使用自定义日志中间件
	router.Use(logger.GinLogger())
	router.Use(logger.GinRecovery())

	// 配置CORS
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:3000", "http://localhost:3001"}
	corsConfig.AllowCredentials = true
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization", "X-User-Address"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	router.Use(cors.New(corsConfig))

	// 设置API路由
	api.SetupRoutes(router, orderService, nftService, collectionService, itemService, activityService, blockchainService)

	// 启动服务器
	logger.Info("服务器启动", map[string]interface{}{
		"port": cfg.Port,
		"mode": gin.Mode(),
	})
	if err := router.Run(":" + cfg.Port); err != nil {
		logger.Error("服务器启动失败", err)
		panic(err)
	}
}
