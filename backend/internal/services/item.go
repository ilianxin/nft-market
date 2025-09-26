package services

import (
	"nft-market/internal/models"
	"time"

	"gorm.io/gorm"
)

// ItemService 物品服务
type ItemService struct {
	db *gorm.DB
}

// NewItemService 创建物品服务
func NewItemService(db *gorm.DB) *ItemService {
	return &ItemService{db: db}
}

// CreateItem 创建物品
func (s *ItemService) CreateItem(req *models.CreateItemRequest) (*models.Item, error) {
	now := time.Now().Unix()
	item := &models.Item{
		ChainID:           models.ChainIDEthereum,
		TokenID:           req.TokenID,
		Name:              req.Name,
		Owner:             req.Owner,
		CollectionAddress: req.CollectionAddress,
		Creator:           req.Creator,
		Supply:            req.Supply,
		CreateTime:        &now,
		UpdateTime:        &now,
	}

	if err := s.db.Create(item).Error; err != nil {
		return nil, err
	}

	return item, nil
}

// GetItemByID 根据ID获取物品
func (s *ItemService) GetItemByID(id uint64) (*models.Item, error) {
	var item models.Item
	if err := s.db.First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

// GetItemByTokenID 根据代币ID获取物品
func (s *ItemService) GetItemByTokenID(collectionAddress, tokenID string) (*models.Item, error) {
	var item models.Item
	if err := s.db.Where("collection_address = ? AND token_id = ?", collectionAddress, tokenID).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

// ListItems 获取物品列表
func (s *ItemService) ListItems(page, pageSize int, collectionAddress, owner, keyword string) (*models.ItemListResponse, error) {
	var items []models.Item
	var total int64

	query := s.db.Model(&models.Item{})

	// 集合地址过滤
	if collectionAddress != "" {
		query = query.Where("collection_address = ?", collectionAddress)
	}

	// 拥有者过滤
	if owner != "" {
		query = query.Where("owner = ?", owner)
	}

	// 关键词搜索
	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&items).Error; err != nil {
		return nil, err
	}

	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	return &models.ItemListResponse{
		Items:      items,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// UpdateItem 更新物品
func (s *ItemService) UpdateItem(id uint64, updates map[string]interface{}) error {
	now := time.Now().Unix()
	updates["update_time"] = now

	return s.db.Model(&models.Item{}).Where("id = ?", id).Updates(updates).Error
}

// UpdateItemOwner 更新物品拥有者
func (s *ItemService) UpdateItemOwner(collectionAddress, tokenID, newOwner string) error {
	now := time.Now().Unix()
	updates := map[string]interface{}{
		"owner":       newOwner,
		"update_time": now,
	}

	return s.db.Model(&models.Item{}).Where("collection_address = ? AND token_id = ?", collectionAddress, tokenID).Updates(updates).Error
}

// UpdateItemPrice 更新物品价格
func (s *ItemService) UpdateItemPrice(collectionAddress, tokenID string, listPrice, salePrice *float64) error {
	now := time.Now().Unix()
	updates := map[string]interface{}{
		"update_time": now,
	}

	if listPrice != nil {
		updates["list_price"] = listPrice
		updates["list_time"] = now
	}
	if salePrice != nil {
		updates["sale_price"] = salePrice
	}

	return s.db.Model(&models.Item{}).Where("collection_address = ? AND token_id = ?", collectionAddress, tokenID).Updates(updates).Error
}

// DeleteItem 删除物品
func (s *ItemService) DeleteItem(id uint64) error {
	return s.db.Delete(&models.Item{}, id).Error
}

// GetItemsByCollection 获取集合下的所有物品
func (s *ItemService) GetItemsByCollection(collectionAddress string) ([]models.Item, error) {
	var items []models.Item
	if err := s.db.Where("collection_address = ?", collectionAddress).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// GetItemsByOwner 获取用户拥有的所有物品
func (s *ItemService) GetItemsByOwner(owner string) ([]models.Item, error) {
	var items []models.Item
	if err := s.db.Where("owner = ?", owner).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}
