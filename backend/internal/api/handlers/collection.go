package handlers

import (
	"net/http"
	"nft-market/internal/models"
	"nft-market/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CollectionHandler 集合处理器
type CollectionHandler struct {
	collectionService *services.CollectionService
}

// NewCollectionHandler 创建集合处理器
func NewCollectionHandler(collectionService *services.CollectionService) *CollectionHandler {
	return &CollectionHandler{
		collectionService: collectionService,
	}
}

// CreateCollection 创建集合
func (h *CollectionHandler) CreateCollection(c *gin.Context) {
	var req models.CreateCollectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	collection, err := h.collectionService.CreateCollection(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to create collection",
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusCreated, models.CollectionResponse{
		Collection: collection,
	})
}

// GetCollection 获取集合详情
func (h *CollectionHandler) GetCollection(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid collection ID",
			Message: "Collection ID must be a valid number",
			Code:    http.StatusBadRequest,
		})
		return
	}

	collection, err := h.collectionService.GetCollectionByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   "Collection not found",
			Message: err.Error(),
			Code:    http.StatusNotFound,
		})
		return
	}

	c.JSON(http.StatusOK, models.CollectionResponse{
		Collection: collection,
	})
}

// GetCollectionByAddress 根据地址获取集合
func (h *CollectionHandler) GetCollectionByAddress(c *gin.Context) {
	address := c.Param("address")
	if address == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid address",
			Message: "Address parameter is required",
			Code:    http.StatusBadRequest,
		})
		return
	}

	collection, err := h.collectionService.GetCollectionByAddress(address)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   "Collection not found",
			Message: err.Error(),
			Code:    http.StatusNotFound,
		})
		return
	}

	c.JSON(http.StatusOK, models.CollectionResponse{
		Collection: collection,
	})
}

// ListCollections 获取集合列表
func (h *CollectionHandler) ListCollections(c *gin.Context) {
	// 获取查询参数
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "20")
	keyword := c.Query("keyword")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	response, err := h.collectionService.ListCollections(page, pageSize, keyword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to get collections",
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// UpdateCollection 更新集合
func (h *CollectionHandler) UpdateCollection(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid collection ID",
			Message: "Collection ID must be a valid number",
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

	err = h.collectionService.UpdateCollection(id, updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to update collection",
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Collection updated successfully"})
}

// DeleteCollection 删除集合
func (h *CollectionHandler) DeleteCollection(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid collection ID",
			Message: "Collection ID must be a valid number",
			Code:    http.StatusBadRequest,
		})
		return
	}

	err = h.collectionService.DeleteCollection(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to delete collection",
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Collection deleted successfully"})
}
