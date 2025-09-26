package api

import (
	"nft-market/internal/api/handlers"
	"nft-market/internal/services"

	"github.com/gin-gonic/gin"
)

// SetupRoutes 设置API路由
func SetupRoutes(router *gin.Engine, orderService *services.OrderService, nftService *services.NFTService, collectionService *services.CollectionService, itemService *services.ItemService, activityService *services.ActivityService) {
	// 创建处理器
	orderHandler := handlers.NewOrderHandler(orderService)
	nftHandler := handlers.NewNFTHandler(nftService)
	collectionHandler := handlers.NewCollectionHandler(collectionService)
	itemHandler := handlers.NewItemHandler(itemService)
	activityHandler := handlers.NewActivityHandler(activityService)

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
			orders.POST("", orderHandler.CreateOrder)                                   // 创建订单
			orders.GET("", orderHandler.GetOrders)                                      // 获取订单列表
			orders.GET("/:id", orderHandler.GetOrderByID)                               // 获取单个订单
			orders.PUT("/:id/cancel", orderHandler.CancelOrder)                         // 取消订单
			orders.GET("/user/:address", orderHandler.GetUserOrders)                    // 获取用户订单
			orders.GET("/nft/:collection_address/:token_id", orderHandler.GetNFTOrders) // 获取NFT订单
			orders.POST("/sync/:orderid", orderHandler.SyncOrderFromChain)              // 从链上同步订单
		}

		// NFT相关路由
		nfts := v1.Group("/nfts")
		{
			nfts.GET("", nftHandler.GetNFTs)                              // 获取NFT列表
			nfts.GET("/:contract/:tokenId", nftHandler.GetNFTByID)        // 获取单个NFT
			nfts.GET("/user/:address", nftHandler.GetUserNFTs)            // 获取用户NFT
			nfts.GET("/contract/:contract", nftHandler.GetNFTsByContract) // 获取合约NFT
			nfts.GET("/search", nftHandler.SearchNFTs)                    // 搜索NFT
			nfts.POST("", nftHandler.CreateOrUpdateNFT)                   // 创建或更新NFT
		}

		// 集合相关路由
		collections := v1.Group("/collections")
		{
			collections.POST("", collectionHandler.CreateCollection)                       // 创建集合
			collections.GET("", collectionHandler.ListCollections)                         // 获取集合列表
			collections.GET("/:id", collectionHandler.GetCollection)                       // 获取集合详情
			collections.GET("/address/:address", collectionHandler.GetCollectionByAddress) // 根据地址获取集合
			collections.PUT("/:id", collectionHandler.UpdateCollection)                    // 更新集合
			collections.DELETE("/:id", collectionHandler.DeleteCollection)                 // 删除集合
		}

		// 物品相关路由
		items := v1.Group("/items")
		{
			items.POST("", itemHandler.CreateItem)                                               // 创建物品
			items.GET("", itemHandler.ListItems)                                                 // 获取物品列表
			items.GET("/id/:id", itemHandler.GetItem)                                            // 获取物品详情
			items.GET("/token/:collection_address/:token_id", itemHandler.GetItemByTokenID)      // 根据代币ID获取物品
			items.PUT("/id/:id", itemHandler.UpdateItem)                                         // 更新物品
			items.PUT("/token/:collection_address/:token_id/owner", itemHandler.UpdateItemOwner) // 更新物品拥有者
			items.PUT("/token/:collection_address/:token_id/price", itemHandler.UpdateItemPrice) // 更新物品价格
			items.DELETE("/id/:id", itemHandler.DeleteItem)                                      // 删除物品
			items.GET("/collection/:collection_address", itemHandler.GetItemsByCollection)       // 获取集合下的所有物品
			items.GET("/owner/:owner", itemHandler.GetItemsByOwner)                              // 获取用户拥有的所有物品
		}

		// 活动相关路由
		activities := v1.Group("/activities")
		{
			activities.POST("", activityHandler.CreateActivity)                                          // 创建活动
			activities.GET("", activityHandler.ListActivities)                                           // 获取活动列表
			activities.GET("/:id", activityHandler.GetActivity)                                          // 获取活动详情
			activities.GET("/collection/:collection_address", activityHandler.GetActivitiesByCollection) // 获取集合的活动
			activities.GET("/item/:collection_address/:token_id", activityHandler.GetActivitiesByItem)   // 获取物品的活动
			activities.GET("/user/:user_address", activityHandler.GetActivitiesByUser)                   // 获取用户的活动
			activities.GET("/recent", activityHandler.GetRecentActivities)                               // 获取最近的活动
			activities.PUT("/:id", activityHandler.UpdateActivity)                                       // 更新活动
			activities.DELETE("/:id", activityHandler.DeleteActivity)                                    // 删除活动
			activities.GET("/stats", activityHandler.GetActivityStats)                                   // 获取活动统计
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
		}
	}
}
