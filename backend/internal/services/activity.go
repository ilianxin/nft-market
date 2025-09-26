package services

import (
	"fmt"
	"nft-market/internal/models"
	"time"

	"gorm.io/gorm"
)

// ActivityService 活动服务
type ActivityService struct {
	db *gorm.DB
}

// NewActivityService 创建活动服务
func NewActivityService(db *gorm.DB) *ActivityService {
	return &ActivityService{db: db}
}

// CreateActivity 创建活动
func (s *ActivityService) CreateActivity(req *models.CreateActivityRequest) (*models.Activity, error) {
	now := time.Now().Unix()
	activity := &models.Activity{
		ActivityType:      req.ActivityType,
		Maker:             req.Maker,
		Taker:             req.Taker,
		CollectionAddress: req.CollectionAddress,
		TokenID:           req.TokenID,
		Price:             req.Price,
		BlockNumber:       req.BlockNumber,
		TxHash:            req.TxHash,
		EventTime:         req.EventTime,
		CreateTime:        &now,
		UpdateTime:        &now,
	}

	if err := s.db.Create(activity).Error; err != nil {
		return nil, err
	}

	return activity, nil
}

// GetActivityByID 根据ID获取活动
func (s *ActivityService) GetActivityByID(id uint64) (*models.Activity, error) {
	var activity models.Activity
	if err := s.db.First(&activity, id).Error; err != nil {
		return nil, err
	}
	return &activity, nil
}

// ListActivities 获取活动列表
func (s *ActivityService) ListActivities(page, pageSize int, activityType *models.ActivityType, collectionAddress, tokenID, maker, taker string) (*models.ActivityListResponse, error) {
	var activities []models.Activity
	var total int64

	query := s.db.Model(&models.Activity{})

	// 活动类型过滤
	if activityType != nil {
		query = query.Where("activity_type = ?", *activityType)
	}

	// 集合地址过滤
	if collectionAddress != "" {
		query = query.Where("collection_address = ?", collectionAddress)
	}

	// 代币ID过滤
	if tokenID != "" {
		query = query.Where("token_id = ?", tokenID)
	}

	// 创建者过滤
	if maker != "" {
		query = query.Where("maker = ?", maker)
	}

	// 接受者过滤
	if taker != "" {
		query = query.Where("taker = ?", taker)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&activities).Error; err != nil {
		return nil, err
	}

	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	return &models.ActivityListResponse{
		Activities: activities,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// GetActivitiesByCollection 获取集合的活动
func (s *ActivityService) GetActivitiesByCollection(collectionAddress string, limit int) ([]models.Activity, error) {
	var activities []models.Activity
	query := s.db.Where("collection_address = ?", collectionAddress).Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&activities).Error; err != nil {
		return nil, err
	}
	return activities, nil
}

// GetActivitiesByItem 获取物品的活动
func (s *ActivityService) GetActivitiesByItem(collectionAddress, tokenID string, limit int) ([]models.Activity, error) {
	var activities []models.Activity
	query := s.db.Where("collection_address = ? AND token_id = ?", collectionAddress, tokenID).Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&activities).Error; err != nil {
		return nil, err
	}
	return activities, nil
}

// GetActivitiesByUser 获取用户的活动
func (s *ActivityService) GetActivitiesByUser(userAddress string, limit int) ([]models.Activity, error) {
	var activities []models.Activity
	query := s.db.Where("maker = ? OR taker = ?", userAddress, userAddress).Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&activities).Error; err != nil {
		return nil, err
	}
	return activities, nil
}

// GetRecentActivities 获取最近的活动
func (s *ActivityService) GetRecentActivities(limit int) ([]models.Activity, error) {
	var activities []models.Activity
	query := s.db.Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&activities).Error; err != nil {
		return nil, err
	}
	return activities, nil
}

// UpdateActivity 更新活动
func (s *ActivityService) UpdateActivity(id uint64, updates map[string]interface{}) error {
	now := time.Now().Unix()
	updates["update_time"] = now

	return s.db.Model(&models.Activity{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteActivity 删除活动
func (s *ActivityService) DeleteActivity(id uint64) error {
	return s.db.Delete(&models.Activity{}, id).Error
}

// GetActivityStats 获取活动统计
func (s *ActivityService) GetActivityStats(collectionAddress, tokenID string) (map[string]int64, error) {
	stats := make(map[string]int64)

	// 按活动类型统计
	var results []struct {
		ActivityType models.ActivityType `json:"activity_type"`
		Count        int64               `json:"count"`
	}

	query := s.db.Model(&models.Activity{})
	if collectionAddress != "" {
		query = query.Where("collection_address = ?", collectionAddress)
	}
	if tokenID != "" {
		query = query.Where("token_id = ?", tokenID)
	}

	if err := query.Select("activity_type, COUNT(*) as count").Group("activity_type").Find(&results).Error; err != nil {
		return nil, err
	}

	for _, result := range results {
		stats[fmt.Sprintf("%d", result.ActivityType)] = result.Count
	}

	return stats, nil
}
