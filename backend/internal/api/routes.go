package api

import (
	"nft-market/internal/api/handlers"
	"nft-market/internal/services"

	"github.com/gin-gonic/gin"
)

// SetupRoutes 设置API路由
func SetupRoutes(router *gin.Engine, orderService *services.OrderService, nftService *services.NFTService) {
	// 创建处理器
	orderHandler := handlers.NewOrderHandler(orderService)
	nftHandler := handlers.NewNFTHandler(nftService)

	// API版本组
	v1 := router.Group("/api/v1")
	{
		// 健康检查
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status":  "ok",
				"message": "NFT市场API服务正常运行",
			})
		})

		// 订单相关路由
		orders := v1.Group("/orders")
		{
			orders.POST("", orderHandler.CreateOrder)           // 创建订单
			orders.GET("", orderHandler.GetOrders)              // 获取订单列表
			orders.GET("/:id", orderHandler.GetOrderByID)       // 获取单个订单
			orders.PUT("/:id/cancel", orderHandler.CancelOrder) // 取消订单
			orders.GET("/user/:address", orderHandler.GetUserOrders) // 获取用户订单
			orders.GET("/nft/:contract/:tokenId", orderHandler.GetNFTOrders) // 获取NFT订单
			orders.POST("/sync/:orderid", orderHandler.SyncOrderFromChain) // 从链上同步订单
		}

		// NFT相关路由
		nfts := v1.Group("/nfts")
		{
			nfts.GET("", nftHandler.GetNFTs)                           // 获取NFT列表
			nfts.GET("/:contract/:tokenId", nftHandler.GetNFTByID)     // 获取单个NFT
			nfts.GET("/user/:address", nftHandler.GetUserNFTs)         // 获取用户NFT
			nfts.GET("/contract/:contract", nftHandler.GetNFTsByContract) // 获取合约NFT
			nfts.GET("/search", nftHandler.SearchNFTs)                 // 搜索NFT
			nfts.POST("", nftHandler.CreateOrUpdateNFT)                // 创建或更新NFT
		}

		// 市场数据路由
		market := v1.Group("/market")
		{
			market.GET("/stats", func(c *gin.Context) {
				c.JSON(200, gin.H{
					"message": "市场统计数据",
					"data":    nil,
				})
			})
			market.GET("/collections", func(c *gin.Context) {
				c.JSON(200, gin.H{
					"message": "热门合集",
					"data":    nil,
				})
			})
		}
	}
}
