package models

import (
	"time"

	"gorm.io/gorm"
)

// OrderStatus 订单状态枚举
type OrderStatus string

const (
	OrderStatusActive    OrderStatus = "active"
	OrderStatusFilled    OrderStatus = "filled"
	OrderStatusCancelled OrderStatus = "cancelled"
	OrderStatusExpired   OrderStatus = "expired"
)

// OrderType 订单类型枚举
type OrderType string

const (
	OrderTypeLimitSell  OrderType = "limit_sell"
	OrderTypeLimitBuy   OrderType = "limit_buy"
	OrderTypeMarketSell OrderType = "market_sell"
	OrderTypeMarketBuy  OrderType = "market_buy"
)

// Order 订单模型
type Order struct {
	ID              uint          `json:"id" gorm:"primaryKey"`
	OrderID         uint64        `json:"order_id" gorm:"uniqueIndex;not null"`
	Maker           string        `json:"maker" gorm:"not null;index"`
	NFTContract     string        `json:"nft_contract" gorm:"not null;index"`
	TokenID         string        `json:"token_id" gorm:"not null;index"`
	Price           string        `json:"price" gorm:"not null"`
	Amount          uint64        `json:"amount" gorm:"not null;default:1"`
	OrderType       OrderType     `json:"order_type" gorm:"not null;index"`
	Status          OrderStatus   `json:"status" gorm:"not null;index;default:'active'"`
	Expiration      time.Time     `json:"expiration" gorm:"not null;index"`
	Signature       string        `json:"signature"`
	TransactionHash string        `json:"transaction_hash" gorm:"index"`
	BlockNumber     uint64        `json:"block_number" gorm:"index"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
}

// NFT NFT模型
type NFT struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Contract    string         `json:"contract" gorm:"not null;index"`
	TokenID     string         `json:"token_id" gorm:"not null;index"`
	Owner       string         `json:"owner" gorm:"not null;index"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Image       string         `json:"image"`
	Attributes  string         `json:"attributes" gorm:"type:json"`
	MetadataURI string         `json:"metadata_uri"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// User 用户模型
type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Address   string         `json:"address" gorm:"uniqueIndex;not null"`
	Username  string         `json:"username" gorm:"uniqueIndex"`
	Email     string         `json:"email" gorm:"uniqueIndex"`
	Avatar    string         `json:"avatar"`
	Bio       string         `json:"bio"`
	Verified  bool           `json:"verified" gorm:"default:false"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// Transaction 交易记录模型
type Transaction struct {
	ID              uint           `json:"id" gorm:"primaryKey"`
	OrderID         uint64         `json:"order_id" gorm:"not null;index"`
	Buyer           string         `json:"buyer" gorm:"not null;index"`
	Seller          string         `json:"seller" gorm:"not null;index"`
	NFTContract     string         `json:"nft_contract" gorm:"not null;index"`
	TokenID         string         `json:"token_id" gorm:"not null;index"`
	Price           string         `json:"price" gorm:"not null"`
	PlatformFee     string         `json:"platform_fee" gorm:"not null"`
	TransactionHash string         `json:"transaction_hash" gorm:"uniqueIndex;not null"`
	BlockNumber     uint64         `json:"block_number" gorm:"not null;index"`
	BlockHash       string         `json:"block_hash" gorm:"not null"`
	Status          string         `json:"status" gorm:"not null;default:'pending'"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
}

// CreateOrderRequest 创建订单请求
type CreateOrderRequest struct {
	NFTContract string    `json:"nft_contract" binding:"required"`
	TokenID     string    `json:"token_id" binding:"required"`
	Price       string    `json:"price"`
	OrderType   OrderType `json:"order_type" binding:"required"`
	Expiration  time.Time `json:"expiration" binding:"required"`
	Signature   string    `json:"signature"`
}

// OrderResponse 订单响应
type OrderResponse struct {
	Order *Order `json:"order"`
	NFT   *NFT   `json:"nft,omitempty"`
}

// OrderListResponse 订单列表响应
type OrderListResponse struct {
	Orders     []OrderResponse `json:"orders"`
	Total      int64           `json:"total"`
	Page       int             `json:"page"`
	PageSize   int             `json:"page_size"`
	TotalPages int             `json:"total_pages"`
}

// ErrorResponse 错误响应
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}
