package handlers

import (
	"net/http"
	"nft-market/internal/logger"
	"nft-market/internal/models"
	"nft-market/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// OrderHandler 订单处理器
type OrderHandler struct {
	orderService *services.OrderService
}

// NewOrderHandler 创建新的订单处理器
func NewOrderHandler(orderService *services.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

// CreateOrder 创建订单
func (oh *OrderHandler) CreateOrder(c *gin.Context) {
	var req models.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "invalid_request",
			Message: "请求参数无效: " + err.Error(),
			Code:    400,
		})
		return
	}

	// 从请求头或认证中获取用户地址（这里简化处理）
	userAddress := c.GetHeader("X-User-Address")
	if userAddress == "" {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error:   "unauthorized",
			Message: "用户未认证",
			Code:    401,
		})
		return
	}

	order, err := oh.orderService.CreateOrder(&req, userAddress)
	if err != nil {
		logger.Error("创建订单失败", err, logrus.Fields{
			"user_address":       userAddress,
			"collection_address": req.CollectionAddress,
			"token_id":           req.TokenID,
			"order_type":         req.OrderType,
			"price":              req.Price,
		})
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "create_order_failed",
			Message: "创建订单失败: " + err.Error(),
			Code:    500,
		})
		return
	}

	logger.Info("订单创建成功", logrus.Fields{
		"order_id":           order.ID,
		"user_address":       userAddress,
		"collection_address": req.CollectionAddress,
		"token_id":           req.TokenID,
		"order_type":         req.OrderType,
		"price":              req.Price,
	})
	c.JSON(http.StatusCreated, gin.H{
		"message": "订单创建成功",
		"data":    order,
	})
}

// GetOrderByID 获取单个订单
func (oh *OrderHandler) GetOrderByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "invalid_id",
			Message: "无效的订单ID",
			Code:    400,
		})
		return
	}

	order, err := oh.orderService.GetOrderByID(uint(id))
	if err != nil {
		logger.Warn("订单不存在", logrus.Fields{
			"order_id": id,
			"error":    err.Error(),
		})
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   "order_not_found",
			Message: "订单不存在",
			Code:    404,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "获取订单成功",
		"data":    order,
	})
}

// GetOrders 获取订单列表
func (oh *OrderHandler) GetOrders(c *gin.Context) {
	// 解析查询参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	orderType := c.Query("order_type")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	orders, err := oh.orderService.GetActiveOrders(page, pageSize, orderType)
	if err != nil {
		logger.Error("获取订单列表失败", err, logrus.Fields{
			"page":       page,
			"page_size":  pageSize,
			"order_type": orderType,
		})
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "get_orders_failed",
			Message: "获取订单列表失败: " + err.Error(),
			Code:    500,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "获取订单列表成功",
		"data":    orders,
	})
}

// GetUserOrders 获取用户订单
func (oh *OrderHandler) GetUserOrders(c *gin.Context) {
	userAddress := c.Param("address")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	status := c.Query("status")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	orders, err := oh.orderService.GetUserOrders(userAddress, page, pageSize, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "get_user_orders_failed",
			Message: "获取用户订单失败: " + err.Error(),
			Code:    500,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "获取用户订单成功",
		"data":    orders,
	})
}

// GetNFTOrders 获取NFT订单
func (oh *OrderHandler) GetNFTOrders(c *gin.Context) {
	collectionAddress := c.Param("collection_address")
	tokenId := c.Param("token_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	orders, err := oh.orderService.GetNFTOrders(collectionAddress, tokenId, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "get_nft_orders_failed",
			Message: "获取NFT订单失败: " + err.Error(),
			Code:    500,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "获取NFT订单成功",
		"data":    orders,
	})
}

// CancelOrder 取消订单
func (oh *OrderHandler) CancelOrder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "invalid_id",
			Message: "无效的订单ID",
			Code:    400,
		})
		return
	}

	// 从请求头获取用户地址
	userAddress := c.GetHeader("X-User-Address")
	if userAddress == "" {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error:   "unauthorized",
			Message: "用户未认证",
			Code:    401,
		})
		return
	}

	err = oh.orderService.CancelOrder(id, userAddress)
	if err != nil {
		logger.Error("取消订单失败", err, logrus.Fields{
			"order_id":     id,
			"user_address": userAddress,
		})
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "cancel_order_failed",
			Message: "取消订单失败: " + err.Error(),
			Code:    500,
		})
		return
	}

	logger.Info("订单取消成功", logrus.Fields{
		"order_id":     id,
		"user_address": userAddress,
	})
	c.JSON(http.StatusOK, gin.H{
		"message": "订单取消成功",
	})
}

// SyncOrderFromChain 从链上同步订单
func (oh *OrderHandler) SyncOrderFromChain(c *gin.Context) {
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

	order, err := oh.orderService.SyncOrderFromChain(orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "sync_order_failed",
			Message: "同步订单失败: " + err.Error(),
			Code:    500,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "订单同步成功",
		"data":    order,
	})
}
