package handlers

import (
	"net/http"
	"nft-market/internal/models"
	"nft-market/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ActivityHandler 活动处理器
type ActivityHandler struct {
	activityService *services.ActivityService
}

// NewActivityHandler 创建活动处理器
func NewActivityHandler(activityService *services.ActivityService) *ActivityHandler {
	return &ActivityHandler{
		activityService: activityService,
	}
}

// CreateActivity 创建活动
func (h *ActivityHandler) CreateActivity(c *gin.Context) {
	var req models.CreateActivityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	activity, err := h.activityService.CreateActivity(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to create activity",
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusCreated, models.ActivityResponse{
		Activity: activity,
	})
}

// GetActivity 获取活动详情
func (h *ActivityHandler) GetActivity(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid activity ID",
			Message: "Activity ID must be a valid number",
			Code:    http.StatusBadRequest,
		})
		return
	}

	activity, err := h.activityService.GetActivityByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   "Activity not found",
			Message: err.Error(),
			Code:    http.StatusNotFound,
		})
		return
	}

	c.JSON(http.StatusOK, models.ActivityResponse{
		Activity: activity,
	})
}

// ListActivities 获取活动列表
func (h *ActivityHandler) ListActivities(c *gin.Context) {
	// 获取查询参数
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "20")
	activityTypeStr := c.Query("activity_type")
	collectionAddress := c.Query("collection_address")
	tokenID := c.Query("token_id")
	maker := c.Query("maker")
	taker := c.Query("taker")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var activityType *models.ActivityType
	if activityTypeStr != "" {
		if at, err := strconv.Atoi(activityTypeStr); err == nil {
			atType := models.ActivityType(at)
			activityType = &atType
		}
	}

	response, err := h.activityService.ListActivities(page, pageSize, activityType, collectionAddress, tokenID, maker, taker)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to get activities",
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetActivitiesByCollection 获取集合的活动
func (h *ActivityHandler) GetActivitiesByCollection(c *gin.Context) {
	collectionAddress := c.Param("collection_address")
	if collectionAddress == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid collection address",
			Message: "Collection address is required",
			Code:    http.StatusBadRequest,
		})
		return
	}

	limitStr := c.DefaultQuery("limit", "50")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 1000 {
		limit = 50
	}

	activities, err := h.activityService.GetActivitiesByCollection(collectionAddress, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to get activities",
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"activities": activities})
}

// GetActivitiesByItem 获取物品的活动
func (h *ActivityHandler) GetActivitiesByItem(c *gin.Context) {
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

	limitStr := c.DefaultQuery("limit", "50")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 1000 {
		limit = 50
	}

	activities, err := h.activityService.GetActivitiesByItem(collectionAddress, tokenID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to get activities",
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"activities": activities})
}

// GetActivitiesByUser 获取用户的活动
func (h *ActivityHandler) GetActivitiesByUser(c *gin.Context) {
	userAddress := c.Param("user_address")
	if userAddress == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid user address",
			Message: "User address is required",
			Code:    http.StatusBadRequest,
		})
		return
	}

	limitStr := c.DefaultQuery("limit", "50")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 1000 {
		limit = 50
	}

	activities, err := h.activityService.GetActivitiesByUser(userAddress, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to get activities",
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"activities": activities})
}

// GetRecentActivities 获取最近的活动
func (h *ActivityHandler) GetRecentActivities(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "50")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 1000 {
		limit = 50
	}

	activities, err := h.activityService.GetRecentActivities(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to get activities",
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"activities": activities})
}

// UpdateActivity 更新活动
func (h *ActivityHandler) UpdateActivity(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid activity ID",
			Message: "Activity ID must be a valid number",
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

	err = h.activityService.UpdateActivity(id, updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to update activity",
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Activity updated successfully"})
}

// DeleteActivity 删除活动
func (h *ActivityHandler) DeleteActivity(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid activity ID",
			Message: "Activity ID must be a valid number",
			Code:    http.StatusBadRequest,
		})
		return
	}

	err = h.activityService.DeleteActivity(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to delete activity",
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Activity deleted successfully"})
}

// GetActivityStats 获取活动统计
func (h *ActivityHandler) GetActivityStats(c *gin.Context) {
	collectionAddress := c.Query("collection_address")
	tokenID := c.Query("token_id")

	stats, err := h.activityService.GetActivityStats(collectionAddress, tokenID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to get activity stats",
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"stats": stats})
}
