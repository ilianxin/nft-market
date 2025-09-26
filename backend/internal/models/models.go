package models

import (
	"time"

	"gorm.io/gorm"
)

// ChainID 链类型枚举
type ChainID int8

const (
	ChainIDEthereum ChainID = 1
)

// OrderStatus 订单状态枚举
type OrderStatus int8

const (
	OrderStatusActive    OrderStatus = 0 // 活跃
	OrderStatusFilled    OrderStatus = 1 // 已成交
	OrderStatusCancelled OrderStatus = 2 // 已取消
	OrderStatusExpired   OrderStatus = 3 // 已过期
)

// OrderType 订单类型枚举
type OrderType int8

const (
	OrderTypeListing       OrderType = 1 // 上架
	OrderTypeOffer         OrderType = 2 // 出价
	OrderTypeCollectionBid OrderType = 3 // 集合出价
	OrderTypeItemBid       OrderType = 4 // 物品出价
)

// ActivityType 活动类型枚举
type ActivityType int8

const (
	ActivityTypeBuy           ActivityType = 1  // 购买
	ActivityTypeMint          ActivityType = 2  // 铸造
	ActivityTypeList          ActivityType = 3  // 上架
	ActivityTypeCancelListing ActivityType = 4  // 取消上架
	ActivityTypeCancelOffer   ActivityType = 5  // 取消出价
	ActivityTypeMakeOffer     ActivityType = 6  // 出价
	ActivityTypeSell          ActivityType = 7  // 出售
	ActivityTypeTransfer      ActivityType = 8  // 转移
	ActivityTypeCollectionBid ActivityType = 9  // 集合出价
	ActivityTypeItemBid       ActivityType = 10 // 物品出价
)

// Collection 集合模型
type Collection struct {
	ID          uint64         `json:"id" gorm:"primaryKey;autoIncrement;comment:主键"`
	ChainID     ChainID        `json:"chain_id" gorm:"type:tinyint;default:1;not null;comment:链类型(1:以太坊)"`
	Symbol      string         `json:"symbol" gorm:"type:varchar(128);not null;comment:项目标识"`
	Name        string         `json:"name" gorm:"type:varchar(128);not null;comment:项目名称"`
	Creator     string         `json:"creator" gorm:"type:varchar(42);not null;comment:创建者"`
	Address     string         `json:"address" gorm:"type:varchar(42);not null;uniqueIndex:index_unique_address;comment:链上合约地址"`
	OwnerAmount int64          `json:"owner_amount" gorm:"type:bigint;default:0;not null;comment:拥有item人数"`
	ItemAmount  int64          `json:"item_amount" gorm:"type:bigint;default:0;not null;comment:该项目NFT的发行总量"`
	FloorPrice  *float64       `json:"floor_price" gorm:"type:decimal(30);comment:整个collection中item的最低的listing价格"`
	SalePrice   *float64       `json:"sale_price" gorm:"type:decimal(30);comment:整个collection中bid的最高的价格"`
	Description *string        `json:"description" gorm:"type:varchar(2048);comment:项目描述"`
	Website     *string        `json:"website" gorm:"type:varchar(512);comment:项目官网地址"`
	VolumeTotal *float64       `json:"volume_total" gorm:"type:decimal(30);comment:总交易量"`
	ImageURI    *string        `json:"image_uri" gorm:"type:varchar(512);comment:项目封面图的链接"`
	CreateTime  *int64         `json:"create_time" gorm:"type:bigint;comment:创建时间"`
	UpdateTime  *int64         `json:"update_time" gorm:"type:bigint;comment:更新时间"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// Item 物品模型
type Item struct {
	ID                uint64         `json:"id" gorm:"primaryKey;autoIncrement;comment:主键"`
	ChainID           ChainID        `json:"chain_id" gorm:"type:tinyint;default:1;not null;comment:链类型"`
	TokenID           string         `json:"token_id" gorm:"type:varchar(128);not null;comment:token_id"`
	Name              string         `json:"name" gorm:"type:varchar(128);not null;comment:nft名称"`
	Owner             *string        `json:"owner" gorm:"type:varchar(42);comment:拥有者"`
	CollectionAddress *string        `json:"collection_address" gorm:"type:varchar(42);comment:合约地址"`
	Creator           string         `json:"creator" gorm:"type:varchar(42);not null;comment:创建者"`
	Supply            int64          `json:"supply" gorm:"type:bigint;not null;comment:item供应量"`
	ListPrice         *float64       `json:"list_price" gorm:"type:decimal(30);comment:上架价格"`
	ListTime          *int64         `json:"list_time" gorm:"type:bigint;comment:上架时间"`
	SalePrice         *float64       `json:"sale_price" gorm:"type:decimal(30);comment:上一次成交价格"`
	CreateTime        *int64         `json:"create_time" gorm:"type:bigint;comment:创建时间"`
	UpdateTime        *int64         `json:"update_time" gorm:"type:bigint;comment:更新时间"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `json:"-" gorm:"index"`
}

// Order 订单模型
type Order struct {
	ID                uint64         `json:"id" gorm:"primaryKey;autoIncrement;comment:主键"`
	MarketplaceID     int8           `json:"marketplace_id" gorm:"type:tinyint;default:0;not null;comment:0.local"`
	OrderID           string         `json:"order_id" gorm:"type:varchar(66);not null;uniqueIndex:index_hash;comment:订单hash"`
	OrderStatus       OrderStatus    `json:"order_status" gorm:"type:tinyint;default:0;not null;comment:标记订单状态"`
	OrderType         OrderType      `json:"order_type" gorm:"type:tinyint;not null;comment:1: listing 2:offer 3:collection bid 4:item bid"`
	EventTime         *int64         `json:"event_time" gorm:"type:bigint;comment:订单时间"`
	CollectionAddress *string        `json:"collection_address" gorm:"type:varchar(42);comment:集合地址"`
	TokenID           *string        `json:"token_id" gorm:"type:varchar(128);comment:代币ID"`
	ExpireTime        *int64         `json:"expire_time" gorm:"type:bigint;comment:过期时间"`
	Price             float64        `json:"price" gorm:"type:decimal(30);default:0;not null;comment:价格"`
	Maker             *string        `json:"maker" gorm:"type:varchar(42);comment:创建者"`
	Taker             *string        `json:"taker" gorm:"type:varchar(42);comment:接受者"`
	QuantityRemaining int64          `json:"quantity_remaining" gorm:"type:bigint;default:1;not null;comment:erc721: 1, erc1155: n"`
	Size              int64          `json:"size" gorm:"type:bigint;default:1;not null;comment:数量"`
	Salt              *int64         `json:"salt" gorm:"type:bigint;default:0;comment:随机数"`
	CurrencyAddress   string         `json:"currency_address" gorm:"type:varchar(42);default:'0x0';not null;comment:货币地址"`
	CreateTime        *int64         `json:"create_time" gorm:"type:bigint;comment:创建时间"`
	UpdateTime        *int64         `json:"update_time" gorm:"type:bigint;comment:更新时间"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `json:"-" gorm:"index"`
}

// Activity 活动模型
type Activity struct {
	ID                uint64         `json:"id" gorm:"primaryKey;autoIncrement;comment:主键"`
	ActivityType      ActivityType   `json:"activity_type" gorm:"type:tinyint;not null;comment:(1:Buy,2:Mint,3:List,4:Cancel Listing,5:Cancel Offer,6.Make Offer,7.Sell,8.Transfer,9.Collection-bid,10.Item-bid)"`
	Maker             *string        `json:"maker" gorm:"type:varchar(42);comment:对于buy,sell,listing,transfer类型指的是nft流转的起始方，即卖方address。对于其他类型可以理解为发起方，如make offer谁发起的from就是谁的地址"`
	Taker             *string        `json:"taker" gorm:"type:varchar(42);comment:目标方,和maker相对"`
	MarketplaceID     int8           `json:"marketplace_id" gorm:"type:tinyint;default:0;not null;comment:市场ID"`
	CollectionAddress *string        `json:"collection_address" gorm:"type:varchar(42);comment:集合地址"`
	TokenID           *string        `json:"token_id" gorm:"type:varchar(128);comment:代币ID"`
	CurrencyAddress   string         `json:"currency_address" gorm:"type:varchar(42);default:'1';not null;comment:货币类型(1表示eth)"`
	Price             float64        `json:"price" gorm:"type:decimal(30);default:0;not null;comment:nft 价格"`
	BlockNumber       int64          `json:"block_number" gorm:"type:bigint;default:0;not null;comment:区块号"`
	TxHash            *string        `json:"tx_hash" gorm:"type:varchar(66);comment:交易事务hash"`
	EventTime         *int64         `json:"event_time" gorm:"type:bigint;comment:链上事件发生的时间"`
	CreateTime        *int64         `json:"create_time" gorm:"type:bigint;comment:创建时间"`
	UpdateTime        *int64         `json:"update_time" gorm:"type:bigint;comment:更新时间"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `json:"-" gorm:"index"`
}

// User 用户模型（保留原有结构）
type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Address   string         `json:"address" gorm:"type:varchar(42);uniqueIndex;not null"`
	Username  string         `json:"username" gorm:"type:varchar(255);uniqueIndex"`
	Email     string         `json:"email" gorm:"type:varchar(255);uniqueIndex"`
	Avatar    string         `json:"avatar" gorm:"type:varchar(512)"`
	Bio       string         `json:"bio" gorm:"type:text"`
	Verified  bool           `json:"verified" gorm:"default:false"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// 请求和响应结构体

// CreateCollectionRequest 创建集合请求
type CreateCollectionRequest struct {
	Symbol      string  `json:"symbol" binding:"required"`
	Name        string  `json:"name" binding:"required"`
	Creator     string  `json:"creator" binding:"required"`
	Address     string  `json:"address" binding:"required"`
	Description *string `json:"description"`
	Website     *string `json:"website"`
	ImageURI    *string `json:"image_uri"`
}

// CreateItemRequest 创建物品请求
type CreateItemRequest struct {
	TokenID           string  `json:"token_id" binding:"required"`
	Name              string  `json:"name" binding:"required"`
	Owner             *string `json:"owner"`
	CollectionAddress *string `json:"collection_address"`
	Creator           string  `json:"creator" binding:"required"`
	Supply            int64   `json:"supply" binding:"required"`
}

// CreateOrderRequest 创建订单请求
type CreateOrderRequest struct {
	CollectionAddress *string   `json:"collection_address"`
	TokenID           *string   `json:"token_id"`
	OrderType         OrderType `json:"order_type" binding:"required"`
	Price             float64   `json:"price" binding:"required"`
	ExpireTime        *int64    `json:"expire_time"`
	QuantityRemaining int64     `json:"quantity_remaining"`
	Size              int64     `json:"size"`
	CurrencyAddress   string    `json:"currency_address"`
}

// CreateActivityRequest 创建活动请求
type CreateActivityRequest struct {
	ActivityType      ActivityType `json:"activity_type" binding:"required"`
	Maker             *string      `json:"maker"`
	Taker             *string      `json:"taker"`
	CollectionAddress *string      `json:"collection_address"`
	TokenID           *string      `json:"token_id"`
	Price             float64      `json:"price" binding:"required"`
	BlockNumber       int64        `json:"block_number" binding:"required"`
	TxHash            *string      `json:"tx_hash"`
	EventTime         *int64       `json:"event_time"`
}

// CollectionResponse 集合响应
type CollectionResponse struct {
	Collection *Collection `json:"collection"`
}

// ItemResponse 物品响应
type ItemResponse struct {
	Item *Item `json:"item"`
}

// OrderResponse 订单响应
type OrderResponse struct {
	Order *Order `json:"order"`
}

// ActivityResponse 活动响应
type ActivityResponse struct {
	Activity *Activity `json:"activity"`
}

// CollectionListResponse 集合列表响应
type CollectionListResponse struct {
	Collections []Collection `json:"collections"`
	Total       int64        `json:"total"`
	Page        int          `json:"page"`
	PageSize    int          `json:"page_size"`
	TotalPages  int          `json:"total_pages"`
}

// ItemListResponse 物品列表响应
type ItemListResponse struct {
	Items      []Item `json:"items"`
	Total      int64  `json:"total"`
	Page       int    `json:"page"`
	PageSize   int    `json:"page_size"`
	TotalPages int    `json:"total_pages"`
}

// OrderListResponse 订单列表响应
type OrderListResponse struct {
	Orders     []Order `json:"orders"`
	Total      int64   `json:"total"`
	Page       int     `json:"page"`
	PageSize   int     `json:"page_size"`
	TotalPages int     `json:"total_pages"`
}

// ActivityListResponse 活动列表响应
type ActivityListResponse struct {
	Activities []Activity `json:"activities"`
	Total      int64      `json:"total"`
	Page       int        `json:"page"`
	PageSize   int        `json:"page_size"`
	TotalPages int        `json:"total_pages"`
}

// ErrorResponse 错误响应
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}
