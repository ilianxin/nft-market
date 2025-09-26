package services

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// BlockchainService 区块链服务
type BlockchainService struct {
	client          *ethclient.Client
	contractAddress common.Address
	contractABI     abi.ABI
	privateKey      *ecdsa.PrivateKey
	publicKey       common.Address
}

// NewBlockchainService 创建新的区块链服务
func NewBlockchainService(rpcURL, contractAddress, privateKeyHex string) (*BlockchainService, error) {
	// 连接到以太坊节点
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("连接以太坊节点失败: %v", err)
	}

	// 解析合约地址
	contractAddr := common.HexToAddress(contractAddress)

	// 解析私钥
	var privateKey *ecdsa.PrivateKey
	var publicKey common.Address
	if privateKeyHex != "" {
		// 移除0x前缀如果存在
		if strings.HasPrefix(privateKeyHex, "0x") || strings.HasPrefix(privateKeyHex, "0X") {
			privateKeyHex = privateKeyHex[2:]
		}
		privateKey, err = crypto.HexToECDSA(privateKeyHex)
		if err != nil {
			return nil, fmt.Errorf("解析私钥失败: %v", err)
		}
		publicKey = crypto.PubkeyToAddress(privateKey.PublicKey)
	}

	// 加载合约ABI
	contractABI, err := abi.JSON(strings.NewReader(NFTMarketplaceABI))
	if err != nil {
		return nil, fmt.Errorf("解析合约ABI失败: %v", err)
	}

	return &BlockchainService{
		client:          client,
		contractAddress: contractAddr,
		contractABI:     contractABI,
		privateKey:      privateKey,
		publicKey:       publicKey,
	}, nil
}

// GetOrderFromChain 从链上获取订单信息
func (bs *BlockchainService) GetOrderFromChain(orderID *big.Int) (map[string]interface{}, error) {
	// 创建调用数据
	input, err := bs.contractABI.Pack("orders", orderID)
	if err != nil {
		return nil, err
	}

	// 调用合约
	msg := ethereum.CallMsg{
		To:   &bs.contractAddress,
		Data: input,
	}

	output, err := bs.client.CallContract(context.Background(), msg, nil)
	if err != nil {
		return nil, err
	}

	// 解析返回数据
	result, err := bs.contractABI.Unpack("orders", output)
	if err != nil {
		return nil, err
	}

	// 转换为map
	orderData := make(map[string]interface{})
	if len(result) > 0 {
		orderStruct := result[0].(struct {
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
		})

		orderData["orderId"] = orderStruct.OrderId
		orderData["maker"] = orderStruct.Maker.Hex()
		orderData["nftContract"] = orderStruct.NftContract.Hex()
		orderData["tokenId"] = orderStruct.TokenId
		orderData["price"] = orderStruct.Price
		orderData["amount"] = orderStruct.Amount
		orderData["timestamp"] = orderStruct.Timestamp
		orderData["expiration"] = orderStruct.Expiration
		orderData["orderType"] = orderStruct.OrderType
		orderData["status"] = orderStruct.Status
		orderData["signature"] = orderStruct.Signature
	}

	return orderData, nil
}

// GetUserOrders 获取用户订单
func (bs *BlockchainService) GetUserOrders(userAddress common.Address) ([]*big.Int, error) {
	input, err := bs.contractABI.Pack("getUserOrders", userAddress)
	if err != nil {
		return nil, err
	}

	msg := ethereum.CallMsg{
		To:   &bs.contractAddress,
		Data: input,
	}

	output, err := bs.client.CallContract(context.Background(), msg, nil)
	if err != nil {
		return nil, err
	}

	result, err := bs.contractABI.Unpack("getUserOrders", output)
	if err != nil {
		return nil, err
	}

	if len(result) > 0 {
		orderIds := result[0].([]*big.Int)
		return orderIds, nil
	}

	return []*big.Int{}, nil
}

// GetNFTOrders 获取NFT的所有订单
func (bs *BlockchainService) GetNFTOrders(nftContract common.Address, tokenID *big.Int) ([]*big.Int, error) {
	input, err := bs.contractABI.Pack("getNFTOrders", nftContract, tokenID)
	if err != nil {
		return nil, err
	}

	msg := ethereum.CallMsg{
		To:   &bs.contractAddress,
		Data: input,
	}

	output, err := bs.client.CallContract(context.Background(), msg, nil)
	if err != nil {
		return nil, err
	}

	result, err := bs.contractABI.Unpack("getNFTOrders", output)
	if err != nil {
		return nil, err
	}

	if len(result) > 0 {
		orderIds := result[0].([]*big.Int)
		return orderIds, nil
	}

	return []*big.Int{}, nil
}

// CreateLimitSellOrder 创建限价卖单
func (bs *BlockchainService) CreateLimitSellOrder(nftContract common.Address, tokenID, price, expiration *big.Int) (*types.Transaction, error) {
	if bs.privateKey == nil {
		return nil, fmt.Errorf("私钥未设置")
	}

	auth, err := bind.NewKeyedTransactorWithChainID(bs.privateKey, nil)
	if err != nil {
		return nil, err
	}

	input, err := bs.contractABI.Pack("createLimitSellOrder", nftContract, tokenID, price, expiration)
	if err != nil {
		return nil, err
	}

	gasPrice, err := bs.client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}

	gasLimit, err := bs.client.EstimateGas(context.Background(), ethereum.CallMsg{
		From: auth.From,
		To:   &bs.contractAddress,
		Data: input,
	})
	if err != nil {
		return nil, err
	}

	nonce, err := bs.client.PendingNonceAt(context.Background(), auth.From)
	if err != nil {
		return nil, err
	}

	tx := types.NewTransaction(nonce, bs.contractAddress, big.NewInt(0), gasLimit, gasPrice, input)

	chainID, err := bs.client.NetworkID(context.Background())
	if err != nil {
		return nil, err
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), bs.privateKey)
	if err != nil {
		return nil, err
	}

	err = bs.client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return nil, err
	}

	return signedTx, nil
}

// CreateLimitBuyOrder 创建限价买单
func (bs *BlockchainService) CreateLimitBuyOrder(nftContract common.Address, tokenID, expiration *big.Int, value *big.Int) (*types.Transaction, error) {
	if bs.privateKey == nil {
		return nil, fmt.Errorf("私钥未设置")
	}

	auth, err := bind.NewKeyedTransactorWithChainID(bs.privateKey, nil)
	if err != nil {
		return nil, err
	}

	input, err := bs.contractABI.Pack("createLimitBuyOrder", nftContract, tokenID, expiration)
	if err != nil {
		return nil, err
	}

	gasPrice, err := bs.client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}

	gasLimit, err := bs.client.EstimateGas(context.Background(), ethereum.CallMsg{
		From:  auth.From,
		To:    &bs.contractAddress,
		Value: value,
		Data:  input,
	})
	if err != nil {
		return nil, err
	}

	nonce, err := bs.client.PendingNonceAt(context.Background(), auth.From)
	if err != nil {
		return nil, err
	}

	tx := types.NewTransaction(nonce, bs.contractAddress, value, gasLimit, gasPrice, input)

	chainID, err := bs.client.NetworkID(context.Background())
	if err != nil {
		return nil, err
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), bs.privateKey)
	if err != nil {
		return nil, err
	}

	err = bs.client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return nil, err
	}

	return signedTx, nil
}

// CancelOrder 取消订单
func (bs *BlockchainService) CancelOrder(orderID *big.Int) (*types.Transaction, error) {
	if bs.privateKey == nil {
		return nil, fmt.Errorf("私钥未设置")
	}

	auth, err := bind.NewKeyedTransactorWithChainID(bs.privateKey, nil)
	if err != nil {
		return nil, err
	}

	input, err := bs.contractABI.Pack("cancelOrder", orderID)
	if err != nil {
		return nil, err
	}

	gasPrice, err := bs.client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}

	gasLimit, err := bs.client.EstimateGas(context.Background(), ethereum.CallMsg{
		From: auth.From,
		To:   &bs.contractAddress,
		Data: input,
	})
	if err != nil {
		return nil, err
	}

	nonce, err := bs.client.PendingNonceAt(context.Background(), auth.From)
	if err != nil {
		return nil, err
	}

	tx := types.NewTransaction(nonce, bs.contractAddress, big.NewInt(0), gasLimit, gasPrice, input)

	chainID, err := bs.client.NetworkID(context.Background())
	if err != nil {
		return nil, err
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), bs.privateKey)
	if err != nil {
		return nil, err
	}

	err = bs.client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return nil, err
	}

	return signedTx, nil
}

// NFTMarketplaceABI 合约ABI（简化版）
const NFTMarketplaceABI = `[
	{
		"inputs": [
			{
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			}
		],
		"name": "orders",
		"outputs": [
			{
				"internalType": "uint256",
				"name": "orderId",
				"type": "uint256"
			},
			{
				"internalType": "address",
				"name": "maker",
				"type": "address"
			},
			{
				"internalType": "address",
				"name": "nftContract",
				"type": "address"
			},
			{
				"internalType": "uint256",
				"name": "tokenId",
				"type": "uint256"
			},
			{
				"internalType": "uint256",
				"name": "price",
				"type": "uint256"
			},
			{
				"internalType": "uint256",
				"name": "amount",
				"type": "uint256"
			},
			{
				"internalType": "uint256",
				"name": "timestamp",
				"type": "uint256"
			},
			{
				"internalType": "uint256",
				"name": "expiration",
				"type": "uint256"
			},
			{
				"internalType": "enum NFTMarketplace.OrderType",
				"name": "orderType",
				"type": "uint8"
			},
			{
				"internalType": "enum NFTMarketplace.OrderStatus",
				"name": "status",
				"type": "uint8"
			},
			{
				"internalType": "bytes",
				"name": "signature",
				"type": "bytes"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "address",
				"name": "_user",
				"type": "address"
			}
		],
		"name": "getUserOrders",
		"outputs": [
			{
				"internalType": "uint256[]",
				"name": "",
				"type": "uint256[]"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "address",
				"name": "_nftContract",
				"type": "address"
			},
			{
				"internalType": "uint256",
				"name": "_tokenId",
				"type": "uint256"
			}
		],
		"name": "getNFTOrders",
		"outputs": [
			{
				"internalType": "uint256[]",
				"name": "",
				"type": "uint256[]"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "address",
				"name": "_nftContract",
				"type": "address"
			},
			{
				"internalType": "uint256",
				"name": "_tokenId",
				"type": "uint256"
			},
			{
				"internalType": "uint256",
				"name": "_price",
				"type": "uint256"
			},
			{
				"internalType": "uint256",
				"name": "_expiration",
				"type": "uint256"
			}
		],
		"name": "createLimitSellOrder",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "address",
				"name": "_nftContract",
				"type": "address"
			},
			{
				"internalType": "uint256",
				"name": "_tokenId",
				"type": "uint256"
			},
			{
				"internalType": "uint256",
				"name": "_expiration",
				"type": "uint256"
			}
		],
		"name": "createLimitBuyOrder",
		"outputs": [],
		"stateMutability": "payable",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "uint256",
				"name": "_orderId",
				"type": "uint256"
			}
		],
		"name": "cancelOrder",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	}
]`
