package handlers

import (
	"net/http"
	"nft-market/internal/models"
	"nft-market/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ItemHandler 物品处理器
type ItemHandler struct {
	itemService *services.ItemService
}

// NewItemHandler 创建物品处理器
func NewItemHandler(itemService *services.ItemService) *ItemHandler {
	return &ItemHandler{
		itemService: itemService,
	}
}

// CreateItem 创建物品
func (h *ItemHandler) CreateItem(c *gin.Context) {
	var req models.CreateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	item, err := h.itemService.CreateItem(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to create item",
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusCreated, models.ItemResponse{
		Item: item,
	})
}

// GetItem 获取物品详情
func (h *ItemHandler) GetItem(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid item ID",
			Message: "Item ID must be a valid number",
			Code:    http.StatusBadRequest,
		})
		return
	}

	item, err := h.itemService.GetItemByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   "Item not found",
			Message: err.Error(),
			Code:    http.StatusNotFound,
		})
		return
	}

	c.JSON(http.StatusOK, models.ItemResponse{
		Item: item,
	})
}

// GetItemByTokenID 根据代币ID获取物品
func (h *ItemHandler) GetItemByTokenID(c *gin.Context) {
	collectionAddress := c.Param("collection_address")
	tokenID := c.Param("token_id")

	if collectionAddress == "" || tokenID == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid parameters",
			Message: "Collection address and token ID are required",
			Code:    http.StatusBadRequest,
		})
		return
	}

	item, err := h.itemService.GetItemByTokenID(collectionAddress, tokenID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   "Item not found",
			Message: err.Error(),
			Code:    http.StatusNotFound,
		})
		return
	}

	c.JSON(http.StatusOK, models.ItemResponse{
		Item: item,
	})
}

// ListItems 获取物品列表
func (h *ItemHandler) ListItems(c *gin.Context) {
	// 获取查询参数
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "20")
	collectionAddress := c.Query("collection_address")
	owner := c.Query("owner")
	keyword := c.Query("keyword")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	response, err := h.itemService.ListItems(page, pageSize, collectionAddress, owner, keyword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to get items",
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// UpdateItem 更新物品
func (h *ItemHandler) UpdateItem(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid item ID",
			Message: "Item ID must be a valid number",
			Code:    http.StatusBadRequest,
		})
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	err = h.itemService.UpdateItem(id, updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to update item",
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item updated successfully"})
}

// UpdateItemOwner 更新物品拥有者
func (h *ItemHandler) UpdateItemOwner(c *gin.Context) {
	collectionAddress := c.Param("collection_address")
	tokenID := c.Param("token_id")

	if collectionAddress == "" || tokenID == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid parameters",
			Message: "Collection address and token ID are required",
			Code:    http.StatusBadRequest,
		})
		return
	}

	var req struct {
		Owner string `json:"owner" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	err := h.itemService.UpdateItemOwner(collectionAddress, tokenID, req.Owner)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to update item owner",
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item owner updated successfully"})
}

// UpdateItemPrice 更新物品价格
func (h *ItemHandler) UpdateItemPrice(c *gin.Context) {
	collectionAddress := c.Param("collection_address")
	tokenID := c.Param("token_id")

	if collectionAddress == "" || tokenID == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid parameters",
			Message: "Collection address and token ID are required",
			Code:    http.StatusBadRequest,
		})
		return
	}

	var req struct {
		ListPrice *float64 `json:"list_price"`
		SalePrice *float64 `json:"sale_price"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	err := h.itemService.UpdateItemPrice(collectionAddress, tokenID, req.ListPrice, req.SalePrice)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to update item price",
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item price updated successfully"})
}

// DeleteItem 删除物品
func (h *ItemHandler) DeleteItem(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid item ID",
			Message: "Item ID must be a valid number",
			Code:    http.StatusBadRequest,
		})
		return
	}

	err = h.itemService.DeleteItem(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to delete item",
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item deleted successfully"})
}

// GetItemsByCollection 获取集合下的所有物品
func (h *ItemHandler) GetItemsByCollection(c *gin.Context) {
	collectionAddress := c.Param("collection_address")
	if collectionAddress == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid collection address",
			Message: "Collection address is required",
			Code:    http.StatusBadRequest,
		})
		return
	}

	items, err := h.itemService.GetItemsByCollection(collectionAddress)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to get items",
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": items})
}

// GetItemsByOwner 获取用户拥有的所有物品
func (h *ItemHandler) GetItemsByOwner(c *gin.Context) {
	owner := c.Param("owner")
	if owner == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid owner address",
			Message: "Owner address is required",
			Code:    http.StatusBadRequest,
		})
		return
	}

	items, err := h.itemService.GetItemsByOwner(owner)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to get items",
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": items})
}
