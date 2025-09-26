package handlers

import (
	"net/http"
	"nft-market/internal/models"
	"nft-market/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

// NFTHandler NFT处理器
type NFTHandler struct {
	nftService *services.NFTService
}

// NewNFTHandler 创建新的NFT处理器
func NewNFTHandler(nftService *services.NFTService) *NFTHandler {
	return &NFTHandler{
		nftService: nftService,
	}
}

// GetNFTByID 获取单个NFT
func (nh *NFTHandler) GetNFTByID(c *gin.Context) {
	contract := c.Param("contract")
	tokenId := c.Param("tokenId")

	nft, err := nh.nftService.GetNFTByContractAndTokenID(contract, tokenId)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   "nft_not_found",
			Message: "NFT不存在",
			Code:    404,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "获取NFT成功",
		"data":    nft,
	})
}

// GetNFTs 获取NFT列表
func (nh *NFTHandler) GetNFTs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// 这里可以根据需要实现获取所有NFT的逻辑
	c.JSON(http.StatusOK, gin.H{
		"message": "获取NFT列表成功",
		"data": gin.H{
			"nfts":       []models.NFT{},
			"total":      0,
			"page":       page,
			"page_size":  pageSize,
			"total_pages": 0,
		},
	})
}

// GetUserNFTs 获取用户NFT
func (nh *NFTHandler) GetUserNFTs(c *gin.Context) {
	userAddress := c.Param("address")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	nfts, total, err := nh.nftService.GetUserNFTs(userAddress, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "get_user_nfts_failed",
			Message: "获取用户NFT失败: " + err.Error(),
			Code:    500,
		})
		return
	}

	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	c.JSON(http.StatusOK, gin.H{
		"message": "获取用户NFT成功",
		"data": gin.H{
			"nfts":        nfts,
			"total":       total,
			"page":        page,
			"page_size":   pageSize,
			"total_pages": totalPages,
		},
	})
}

// GetNFTsByContract 获取合约NFT
func (nh *NFTHandler) GetNFTsByContract(c *gin.Context) {
	contract := c.Param("contract")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	nfts, total, err := nh.nftService.GetNFTsByContract(contract, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "get_contract_nfts_failed",
			Message: "获取合约NFT失败: " + err.Error(),
			Code:    500,
		})
		return
	}

	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	c.JSON(http.StatusOK, gin.H{
		"message": "获取合约NFT成功",
		"data": gin.H{
			"nfts":        nfts,
			"total":       total,
			"page":        page,
			"page_size":   pageSize,
			"total_pages": totalPages,
		},
	})
}

// SearchNFTs 搜索NFT
func (nh *NFTHandler) SearchNFTs(c *gin.Context) {
	keyword := c.Query("q")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if keyword == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "missing_keyword",
			Message: "搜索关键词不能为空",
			Code:    400,
		})
		return
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	nfts, total, err := nh.nftService.SearchNFTs(keyword, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "search_nfts_failed",
			Message: "搜索NFT失败: " + err.Error(),
			Code:    500,
		})
		return
	}

	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	c.JSON(http.StatusOK, gin.H{
		"message": "搜索NFT成功",
		"data": gin.H{
			"nfts":        nfts,
			"total":       total,
			"page":        page,
			"page_size":   pageSize,
			"total_pages": totalPages,
			"keyword":     keyword,
		},
	})
}

// CreateOrUpdateNFT 创建或更新NFT
func (nh *NFTHandler) CreateOrUpdateNFT(c *gin.Context) {
	var nft models.NFT
	if err := c.ShouldBindJSON(&nft); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "invalid_request",
			Message: "请求参数无效: " + err.Error(),
			Code:    400,
		})
		return
	}

	err := nh.nftService.CreateOrUpdateNFT(&nft)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "create_nft_failed",
			Message: "创建或更新NFT失败: " + err.Error(),
			Code:    500,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "NFT创建或更新成功",
		"data":    nft,
	})
}
