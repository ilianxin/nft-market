package services

import (
	"fmt"
	"nft-market/internal/models"

	"gorm.io/gorm"
)

// NFTService NFT服务
type NFTService struct {
	db               *gorm.DB
	blockchainService *BlockchainService
}

// NewNFTService 创建新的NFT服务
func NewNFTService(db *gorm.DB, blockchainService *BlockchainService) *NFTService {
	return &NFTService{
		db:               db,
		blockchainService: blockchainService,
	}
}

// GetNFTByContractAndTokenID 根据合约地址和Token ID获取NFT
func (ns *NFTService) GetNFTByContractAndTokenID(contract, tokenID string) (*models.NFT, error) {
	var nft models.NFT
	err := ns.db.Where("contract = ? AND token_id = ?", contract, tokenID).First(&nft).Error
	if err != nil {
		return nil, err
	}
	return &nft, nil
}

// CreateOrUpdateNFT 创建或更新NFT
func (ns *NFTService) CreateOrUpdateNFT(nft *models.NFT) error {
	var existingNFT models.NFT
	err := ns.db.Where("contract = ? AND token_id = ?", nft.Contract, nft.TokenID).First(&existingNFT).Error
	
	if err == gorm.ErrRecordNotFound {
		// 创建新NFT
		return ns.db.Create(nft).Error
	} else if err != nil {
		return err
	} else {
		// 更新现有NFT
		nft.ID = existingNFT.ID
		return ns.db.Save(nft).Error
	}
}

// GetUserNFTs 获取用户的NFT列表
func (ns *NFTService) GetUserNFTs(owner string, page, pageSize int) ([]models.NFT, int64, error) {
	var nfts []models.NFT
	var total int64

	query := ns.db.Model(&models.NFT{}).Where("owner = ?", owner)

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&nfts).Error; err != nil {
		return nil, 0, err
	}

	return nfts, total, nil
}

// GetNFTsByContract 根据合约地址获取NFT列表
func (ns *NFTService) GetNFTsByContract(contract string, page, pageSize int) ([]models.NFT, int64, error) {
	var nfts []models.NFT
	var total int64

	query := ns.db.Model(&models.NFT{}).Where("contract = ?", contract)

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&nfts).Error; err != nil {
		return nil, 0, err
	}

	return nfts, total, nil
}

// SearchNFTs 搜索NFT
func (ns *NFTService) SearchNFTs(keyword string, page, pageSize int) ([]models.NFT, int64, error) {
	var nfts []models.NFT
	var total int64

	searchPattern := "%" + keyword + "%"
	query := ns.db.Model(&models.NFT{}).Where("name ILIKE ? OR description ILIKE ?", searchPattern, searchPattern)

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&nfts).Error; err != nil {
		return nil, 0, err
	}

	return nfts, total, nil
}

// UpdateNFTOwner 更新NFT所有者
func (ns *NFTService) UpdateNFTOwner(contract, tokenID, newOwner string) error {
	result := ns.db.Model(&models.NFT{}).
		Where("contract = ? AND token_id = ?", contract, tokenID).
		Update("owner", newOwner)
	
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("NFT不存在")
	}

	return nil
}
