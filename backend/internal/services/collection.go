package services

import (
	"nft-market/internal/models"
	"time"

	"gorm.io/gorm"
)

// CollectionService 集合服务
type CollectionService struct {
	db *gorm.DB
}

// NewCollectionService 创建集合服务
func NewCollectionService(db *gorm.DB) *CollectionService {
	return &CollectionService{db: db}
}

// CreateCollection 创建集合
func (s *CollectionService) CreateCollection(req *models.CreateCollectionRequest) (*models.Collection, error) {
	now := time.Now().Unix()
	collection := &models.Collection{
		ChainID:     models.ChainIDEthereum,
		Symbol:      req.Symbol,
		Name:        req.Name,
		Creator:     req.Creator,
		Address:     req.Address,
		Description: req.Description,
		Website:     req.Website,
		ImageURI:    req.ImageURI,
		CreateTime:  &now,
		UpdateTime:  &now,
	}

	if err := s.db.Create(collection).Error; err != nil {
		return nil, err
	}

	return collection, nil
}

// GetCollectionByID 根据ID获取集合
func (s *CollectionService) GetCollectionByID(id uint64) (*models.Collection, error) {
	var collection models.Collection
	if err := s.db.First(&collection, id).Error; err != nil {
		return nil, err
	}
	return &collection, nil
}

// GetCollectionByAddress 根据地址获取集合
func (s *CollectionService) GetCollectionByAddress(address string) (*models.Collection, error) {
	var collection models.Collection
	if err := s.db.Where("address = ?", address).First(&collection).Error; err != nil {
		return nil, err
	}
	return &collection, nil
}

// ListCollections 获取集合列表
func (s *CollectionService) ListCollections(page, pageSize int, keyword string) (*models.CollectionListResponse, error) {
	var collections []models.Collection
	var total int64

	query := s.db.Model(&models.Collection{})

	// 关键词搜索
	if keyword != "" {
		query = query.Where("name LIKE ? OR symbol LIKE ? OR description LIKE ?", 
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&collections).Error; err != nil {
		return nil, err
	}

	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	return &models.CollectionListResponse{
		Collections: collections,
		Total:       total,
		Page:        page,
		PageSize:    pageSize,
		TotalPages:  totalPages,
	}, nil
}

// UpdateCollection 更新集合
func (s *CollectionService) UpdateCollection(id uint64, updates map[string]interface{}) error {
	now := time.Now().Unix()
	updates["update_time"] = now

	return s.db.Model(&models.Collection{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteCollection 删除集合
func (s *CollectionService) DeleteCollection(id uint64) error {
	return s.db.Delete(&models.Collection{}, id).Error
}

// UpdateCollectionStats 更新集合统计信息
func (s *CollectionService) UpdateCollectionStats(address string, ownerAmount, itemAmount int64, floorPrice, salePrice *float64, volumeTotal *float64) error {
	updates := map[string]interface{}{
		"owner_amount": ownerAmount,
		"item_amount":  itemAmount,
		"update_time":  time.Now().Unix(),
	}

	if floorPrice != nil {
		updates["floor_price"] = floorPrice
	}
	if salePrice != nil {
		updates["sale_price"] = salePrice
	}
	if volumeTotal != nil {
		updates["volume_total"] = volumeTotal
	}

	return s.db.Model(&models.Collection{}).Where("address = ?", address).Updates(updates).Error
}
