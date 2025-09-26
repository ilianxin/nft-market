package blockchain

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"nft-market/internal/logger"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
)

// NFTMarketplaceContract NFT市场合约交互类
type NFTMarketplaceContract struct {
	client          *ethclient.Client
	contractAddress common.Address
	contractABI     abi.ABI
	privateKey      *ecdsa.PrivateKey
	publicKey       *ecdsa.PublicKey
	fromAddress     common.Address
}

// ContractOrder 合约中的订单结构
type ContractOrder struct {
	OrderId     *big.Int
	Maker       common.Address
	NftContract common.Address
	TokenId     *big.Int
	Price       *big.Int
	Amount      *big.Int
	Timestamp   *big.Int
	Expiration  *big.Int
	OrderType   uint8
	Status      uint8
	Signature   []byte
}

// OrderCreatedEvent 订单创建事件
type OrderCreatedEvent struct {
	OrderId     *big.Int
	Maker       common.Address
	NftContract common.Address
	TokenId     *big.Int
	Price       *big.Int
	OrderType   uint8
}

// OrderFilledEvent 订单成交事件
type OrderFilledEvent struct {
	OrderId *big.Int
	Buyer   common.Address
	Seller  common.Address
	Price   *big.Int
}

// OrderCancelledEvent 订单取消事件
type OrderCancelledEvent struct {
	OrderId *big.Int
	Maker   common.Address
}

// NFTMarketplace合约ABI（简化版）
const NFTMarketplaceABI = `[
	{
		"inputs": [{"name": "_feeRecipient", "type": "address"}],
		"stateMutability": "nonpayable",
		"type": "constructor"
	},
	{
		"inputs": [
			{"name": "_nftContract", "type": "address"},
			{"name": "_tokenId", "type": "uint256"},
			{"name": "_price", "type": "uint256"},
			{"name": "_expiration", "type": "uint256"}
		],
		"name": "createLimitSellOrder",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [
			{"name": "_nftContract", "type": "address"},
			{"name": "_tokenId", "type": "uint256"},
			{"name": "_expiration", "type": "uint256"}
		],
		"name": "createLimitBuyOrder",
		"outputs": [],
		"stateMutability": "payable",
		"type": "function"
	},
	{
		"inputs": [{"name": "_orderId", "type": "uint256"}],
		"name": "cancelOrder",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [{"name": "_orderId", "type": "uint256"}],
		"name": "executeOrder",
		"outputs": [],
		"stateMutability": "payable",
		"type": "function"
	},
	{
		"inputs": [{"name": "", "type": "uint256"}],
		"name": "orders",
		"outputs": [
			{"name": "orderId", "type": "uint256"},
			{"name": "maker", "type": "address"},
			{"name": "nftContract", "type": "address"},
			{"name": "tokenId", "type": "uint256"},
			{"name": "price", "type": "uint256"},
			{"name": "amount", "type": "uint256"},
			{"name": "timestamp", "type": "uint256"},
			{"name": "expiration", "type": "uint256"},
			{"name": "orderType", "type": "uint8"},
			{"name": "status", "type": "uint8"},
			{"name": "signature", "type": "bytes"}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "orderCounter",
		"outputs": [{"name": "", "type": "uint256"}],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"anonymous": false,
		"inputs": [
			{"indexed": true, "name": "orderId", "type": "uint256"},
			{"indexed": true, "name": "maker", "type": "address"},
			{"indexed": true, "name": "nftContract", "type": "address"},
			{"name": "tokenId", "type": "uint256"},
			{"name": "price", "type": "uint256"},
			{"name": "orderType", "type": "uint8"}
		],
		"name": "OrderCreated",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{"indexed": true, "name": "orderId", "type": "uint256"},
			{"indexed": true, "name": "buyer", "type": "address"},
			{"indexed": true, "name": "seller", "type": "address"},
			{"name": "price", "type": "uint256"}
		],
		"name": "OrderFilled",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{"indexed": true, "name": "orderId", "type": "uint256"},
			{"indexed": true, "name": "maker", "type": "address"}
		],
		"name": "OrderCancelled",
		"type": "event"
	}
]`

// NewNFTMarketplaceContract 创建新的NFT市场合约实例
func NewNFTMarketplaceContract(rpcURL, contractAddr, privateKeyHex string) (*NFTMarketplaceContract, error) {
	// 连接到以太坊客户端
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("连接以太坊客户端失败: %v", err)
	}

	// 解析合约地址
	contractAddress := common.HexToAddress(contractAddr)

	// 解析合约ABI
	contractABI, err := abi.JSON(strings.NewReader(NFTMarketplaceABI))
	if err != nil {
		return nil, fmt.Errorf("解析合约ABI失败: %v", err)
	}

	// 解析私钥
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, fmt.Errorf("解析私钥失败: %v", err)
	}

	// 获取公钥和地址
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("获取公钥失败")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	logger.Info("NFT市场合约初始化成功", logrus.Fields{
		"contract_address": contractAddr,
		"from_address":     fromAddress.Hex(),
		"rpc_url":          rpcURL,
	})

	return &NFTMarketplaceContract{
		client:          client,
		contractAddress: contractAddress,
		contractABI:     contractABI,
		privateKey:      privateKey,
		publicKey:       publicKeyECDSA,
		fromAddress:     fromAddress,
	}, nil
}

// CreateLimitSellOrder 在链上创建限价卖单
func (c *NFTMarketplaceContract) CreateLimitSellOrder(nftContract, tokenID string, price float64, expiration int64) (*types.Transaction, error) {
	logger.Info("开始创建链上限价卖单", logrus.Fields{
		"nft_contract": nftContract,
		"token_id":     tokenID,
		"price":        price,
		"expiration":   expiration,
	})

	// 准备交易参数
	nftAddr := common.HexToAddress(nftContract)
	tokenIDBig := new(big.Int)
	tokenIDBig.SetString(tokenID, 10)

	// 将价格从ETH转换为Wei
	priceWei := new(big.Int)
	priceWei.SetString(fmt.Sprintf("%.0f", price*1e18), 10)

	expirationBig := big.NewInt(expiration)

	// 获取交易选项
	auth, err := c.getTransactOpts()
	if err != nil {
		return nil, fmt.Errorf("获取交易选项失败: %v", err)
	}

	// 调用合约方法
	tx, err := c.callContract(auth, "createLimitSellOrder", nftAddr, tokenIDBig, priceWei, expirationBig)
	if err != nil {
		logger.Error("创建链上限价卖单失败", err, logrus.Fields{
			"nft_contract": nftContract,
			"token_id":     tokenID,
		})
		return nil, err
	}

	logger.Info("链上限价卖单创建成功", logrus.Fields{
		"tx_hash":      tx.Hash().Hex(),
		"nft_contract": nftContract,
		"token_id":     tokenID,
	})

	return tx, nil
}

// CreateLimitBuyOrder 在链上创建限价买单
func (c *NFTMarketplaceContract) CreateLimitBuyOrder(nftContract, tokenID string, price float64, expiration int64) (*types.Transaction, error) {
	logger.Info("开始创建链上限价买单", logrus.Fields{
		"nft_contract": nftContract,
		"token_id":     tokenID,
		"price":        price,
		"expiration":   expiration,
	})

	// 准备交易参数
	nftAddr := common.HexToAddress(nftContract)
	tokenIDBig := new(big.Int)
	tokenIDBig.SetString(tokenID, 10)

	// 将价格从ETH转换为Wei
	priceWei := new(big.Int)
	priceWei.SetString(fmt.Sprintf("%.0f", price*1e18), 10)

	expirationBig := big.NewInt(expiration)

	// 获取交易选项（包含ETH值）
	auth, err := c.getTransactOpts()
	if err != nil {
		return nil, fmt.Errorf("获取交易选项失败: %v", err)
	}
	auth.Value = priceWei // 设置发送的ETH数量

	// 调用合约方法
	tx, err := c.callContract(auth, "createLimitBuyOrder", nftAddr, tokenIDBig, expirationBig)
	if err != nil {
		logger.Error("创建链上限价买单失败", err, logrus.Fields{
			"nft_contract": nftContract,
			"token_id":     tokenID,
		})
		return nil, err
	}

	logger.Info("链上限价买单创建成功", logrus.Fields{
		"tx_hash":      tx.Hash().Hex(),
		"nft_contract": nftContract,
		"token_id":     tokenID,
	})

	return tx, nil
}

// CancelOrder 在链上取消订单
func (c *NFTMarketplaceContract) CancelOrder(orderID uint64) (*types.Transaction, error) {
	logger.Info("开始取消链上订单", logrus.Fields{
		"order_id": orderID,
	})

	orderIDBig := big.NewInt(int64(orderID))

	// 获取交易选项
	auth, err := c.getTransactOpts()
	if err != nil {
		return nil, fmt.Errorf("获取交易选项失败: %v", err)
	}

	// 调用合约方法
	tx, err := c.callContract(auth, "cancelOrder", orderIDBig)
	if err != nil {
		logger.Error("取消链上订单失败", err, logrus.Fields{
			"order_id": orderID,
		})
		return nil, err
	}

	logger.Info("链上订单取消成功", logrus.Fields{
		"tx_hash":  tx.Hash().Hex(),
		"order_id": orderID,
	})

	return tx, nil
}

// ExecuteOrder 在链上执行订单
func (c *NFTMarketplaceContract) ExecuteOrder(orderID uint64, price float64) (*types.Transaction, error) {
	logger.Info("开始执行链上订单", logrus.Fields{
		"order_id": orderID,
		"price":    price,
	})

	orderIDBig := big.NewInt(int64(orderID))

	// 获取交易选项
	auth, err := c.getTransactOpts()
	if err != nil {
		return nil, fmt.Errorf("获取交易选项失败: %v", err)
	}

	// 如果是买单执行，需要发送ETH
	if price > 0 {
		priceWei := new(big.Int)
		priceWei.SetString(fmt.Sprintf("%.0f", price*1e18), 10)
		auth.Value = priceWei
	}

	// 调用合约方法
	tx, err := c.callContract(auth, "executeOrder", orderIDBig)
	if err != nil {
		logger.Error("执行链上订单失败", err, logrus.Fields{
			"order_id": orderID,
		})
		return nil, err
	}

	logger.Info("链上订单执行成功", logrus.Fields{
		"tx_hash":  tx.Hash().Hex(),
		"order_id": orderID,
	})

	return tx, nil
}

// GetOrder 从链上获取订单信息
func (c *NFTMarketplaceContract) GetOrder(orderID uint64) (*ContractOrder, error) {
	orderIDBig := big.NewInt(int64(orderID))

	// 调用合约的只读方法
	result, err := c.callContractRead("orders", orderIDBig)
	if err != nil {
		return nil, fmt.Errorf("获取链上订单失败: %v", err)
	}

	// 解析返回结果
	if len(result) < 11 {
		return nil, fmt.Errorf("链上订单数据格式错误")
	}

	order := &ContractOrder{
		OrderId:     result[0].(*big.Int),
		Maker:       result[1].(common.Address),
		NftContract: result[2].(common.Address),
		TokenId:     result[3].(*big.Int),
		Price:       result[4].(*big.Int),
		Amount:      result[5].(*big.Int),
		Timestamp:   result[6].(*big.Int),
		Expiration:  result[7].(*big.Int),
		OrderType:   result[8].(uint8),
		Status:      result[9].(uint8),
		Signature:   result[10].([]byte),
	}

	return order, nil
}

// GetOrderCounter 获取订单计数器
func (c *NFTMarketplaceContract) GetOrderCounter() (*big.Int, error) {
	result, err := c.callContractRead("orderCounter")
	if err != nil {
		return nil, fmt.Errorf("获取订单计数器失败: %v", err)
	}

	if len(result) == 0 {
		return big.NewInt(0), nil
	}

	return result[0].(*big.Int), nil
}

// getTransactOpts 获取交易选项
func (c *NFTMarketplaceContract) getTransactOpts() (*bind.TransactOpts, error) {
	nonce, err := c.client.PendingNonceAt(context.Background(), c.fromAddress)
	if err != nil {
		return nil, fmt.Errorf("获取nonce失败: %v", err)
	}

	gasPrice, err := c.client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, fmt.Errorf("获取gas价格失败: %v", err)
	}

	chainID, err := c.client.NetworkID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("获取链ID失败: %v", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(c.privateKey, chainID)
	if err != nil {
		return nil, fmt.Errorf("创建交易签名器失败: %v", err)
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(300000)
	auth.GasPrice = gasPrice

	return auth, nil
}

// callContract 调用合约方法
func (c *NFTMarketplaceContract) callContract(auth *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	// 编码方法调用
	data, err := c.contractABI.Pack(method, params...)
	if err != nil {
		return nil, fmt.Errorf("编码合约调用失败: %v", err)
	}

	// 创建交易
	tx := types.NewTransaction(
		auth.Nonce.Uint64(),
		c.contractAddress,
		auth.Value,
		auth.GasLimit,
		auth.GasPrice,
		data,
	)

	// 签名交易
	chainID, err := c.client.NetworkID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("获取链ID失败: %v", err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), c.privateKey)
	if err != nil {
		return nil, fmt.Errorf("签名交易失败: %v", err)
	}

	// 发送交易
	err = c.client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return nil, fmt.Errorf("发送交易失败: %v", err)
	}

	return signedTx, nil
}

// callContractRead 调用合约只读方法
func (c *NFTMarketplaceContract) callContractRead(method string, params ...interface{}) ([]interface{}, error) {
	// 编码方法调用
	data, err := c.contractABI.Pack(method, params...)
	if err != nil {
		return nil, fmt.Errorf("编码合约调用失败: %v", err)
	}

	// 创建调用消息
	msg := ethereum.CallMsg{
		To:   &c.contractAddress,
		Data: data,
	}

	// 执行调用
	result, err := c.client.CallContract(context.Background(), msg, nil)
	if err != nil {
		return nil, fmt.Errorf("调用合约失败: %v", err)
	}

	// 解码结果
	outputs, err := c.contractABI.Unpack(method, result)
	if err != nil {
		return nil, fmt.Errorf("解码合约返回值失败: %v", err)
	}

	return outputs, nil
}

// WaitForTransaction 等待交易确认
func (c *NFTMarketplaceContract) WaitForTransaction(tx *types.Transaction, timeout time.Duration) (*types.Receipt, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	logger.Info("等待交易确认", logrus.Fields{
		"tx_hash": tx.Hash().Hex(),
		"timeout": timeout.String(),
	})

	receipt, err := bind.WaitMined(ctx, c.client, tx)
	if err != nil {
		logger.Error("等待交易确认失败", err, logrus.Fields{
			"tx_hash": tx.Hash().Hex(),
		})
		return nil, err
	}

	logger.Info("交易确认成功", logrus.Fields{
		"tx_hash":      tx.Hash().Hex(),
		"block_number": receipt.BlockNumber.String(),
		"gas_used":     receipt.GasUsed,
		"status":       receipt.Status,
	})

	return receipt, nil
}

// Close 关闭客户端连接
func (c *NFTMarketplaceContract) Close() {
	if c.client != nil {
		c.client.Close()
	}
}
