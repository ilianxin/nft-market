package handlers

import (
	"net/http"
	"nft-market/internal/logger"
	"nft-market/internal/models"
	"nft-market/internal/services"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// BlockchainHandler 区块链处理器
type BlockchainHandler struct {
	blockchainService *services.EnhancedBlockchainService
}

// NewBlockchainHandler 创建新的区块链处理器
func NewBlockchainHandler(blockchainService *services.EnhancedBlockchainService) *BlockchainHandler {
	return &BlockchainHandler{
		blockchainService: blockchainService,
	}
}

// SyncOrderFromChain 从链上同步单个订单
func (bh *BlockchainHandler) SyncOrderFromChain(c *gin.Context) {
	orderIDStr := c.Param("orderid")
	orderID, err := strconv.ParseUint(orderIDStr, 10, 64)
	if err != nil {
		logger.Warn("无效的订单ID", logrus.Fields{
			"order_id_str": orderIDStr,
			"error":        err.Error(),
		})
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "invalid_order_id",
			Message: "无效的订单ID",
			Code:    400,
		})
		return
	}

	logger.Info("开始从链上同步订单", logrus.Fields{
		"order_id": orderID,
	})

	order, err := bh.blockchainService.SyncOrderFromChain(orderID)
	if err != nil {
		logger.Error("从链上同步订单失败", err, logrus.Fields{
			"order_id": orderID,
		})
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "sync_order_failed",
			Message: "同步订单失败: " + err.Error(),
			Code:    500,
		})
		return
	}

	logger.Info("从链上同步订单成功", logrus.Fields{
		"order_id": orderID,
		"db_id":    order.ID,
	})

	c.JSON(http.StatusOK, gin.H{
		"message": "订单同步成功",
		"data":    order,
	})
}

// SyncAllOrdersFromChain 从链上同步所有订单
func (bh *BlockchainHandler) SyncAllOrdersFromChain(c *gin.Context) {
	logger.Info("开始从链上同步所有订单")

	// 异步执行同步，避免阻塞请求
	go func() {
		if err := bh.blockchainService.SyncAllOrdersFromChain(); err != nil {
			logger.Error("同步所有订单失败", err)
		}
	}()

	c.JSON(http.StatusAccepted, gin.H{
		"message": "开始同步所有订单，请查看日志了解进度",
	})
}

// GetOrderCounter 获取链上订单计数器
func (bh *BlockchainHandler) GetOrderCounter(c *gin.Context) {
	logger.Debug("获取链上订单计数器")

	counter, err := bh.blockchainService.GetOrderCounter()
	if err != nil {
		logger.Error("获取订单计数器失败", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "get_counter_failed",
			Message: "获取订单计数器失败: " + err.Error(),
			Code:    500,
		})
		return
	}

	logger.Info("获取订单计数器成功", logrus.Fields{
		"counter": counter.String(),
	})

	c.JSON(http.StatusOK, gin.H{
		"message": "获取订单计数器成功",
		"data": gin.H{
			"counter": counter.String(),
		},
	})
}

// GetChainOrderInfo 获取链上订单信息
func (bh *BlockchainHandler) GetChainOrderInfo(c *gin.Context) {
	orderIDStr := c.Param("orderid")
	orderID, err := strconv.ParseUint(orderIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "invalid_order_id",
			Message: "无效的订单ID",
			Code:    400,
		})
		return
	}

	logger.Debug("获取链上订单信息", logrus.Fields{
		"order_id": orderID,
	})

	orderInfo, err := bh.blockchainService.GetOrderFromChainEnhanced(orderID)
	if err != nil {
		logger.Error("获取链上订单信息失败", err, logrus.Fields{
			"order_id": orderID,
		})
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "get_chain_order_failed",
			Message: "获取链上订单信息失败: " + err.Error(),
			Code:    500,
		})
		return
	}

	logger.Debug("获取链上订单信息成功", logrus.Fields{
		"order_id": orderID,
	})

	c.JSON(http.StatusOK, gin.H{
		"message": "获取链上订单信息成功",
		"data":    orderInfo,
	})
}

// ExecuteOrder 执行订单
func (bh *BlockchainHandler) ExecuteOrder(c *gin.Context) {
	orderIDStr := c.Param("orderid")
	orderID, err := strconv.ParseUint(orderIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "invalid_order_id",
			Message: "无效的订单ID",
			Code:    400,
		})
		return
	}

	var req struct {
		Price float64 `json:"price" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "invalid_request",
			Message: "请求参数无效: " + err.Error(),
			Code:    400,
		})
		return
	}

	logger.Info("开始执行链上订单", logrus.Fields{
		"order_id": orderID,
		"price":    req.Price,
	})

	// 异步执行订单
	go func() {
		tx, err := bh.blockchainService.ExecuteOrderOnChain(orderID, req.Price)
		if err != nil {
			logger.Error("执行链上订单失败", err, logrus.Fields{
				"order_id": orderID,
			})
			return
		}

		// 等待交易确认
		receipt, err := bh.blockchainService.WaitForTransactionConfirmation(tx, 5*time.Minute)
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

	c.JSON(http.StatusAccepted, gin.H{
		"message": "开始执行订单，请查看日志了解进度",
	})
}

// GetBlockchainStatus 获取区块链服务状态
func (bh *BlockchainHandler) GetBlockchainStatus(c *gin.Context) {
	// 获取最新区块号
	blockNumber, err := bh.blockchainService.GetBlockNumber()
	if err != nil {
		logger.Error("获取区块号失败", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "get_block_number_failed",
			Message: "获取区块号失败: " + err.Error(),
			Code:    500,
		})
		return
	}

	// 获取订单计数器
	counter, err := bh.blockchainService.GetOrderCounter()
	if err != nil {
		logger.Error("获取订单计数器失败", err)
		counter = nil
	}

	status := gin.H{
		"latest_block": blockNumber,
		"service":      "running",
	}

	if counter != nil {
		status["order_counter"] = counter.String()
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "区块链服务状态正常",
		"data":    status,
	})
}
