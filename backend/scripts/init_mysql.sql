-- NFT市场应用MySQL数据库初始化脚本

-- 创建数据库
CREATE DATABASE IF NOT EXISTS nft_market 
CHARACTER SET utf8mb4 
COLLATE utf8mb4_unicode_ci;

-- 使用数据库
USE nft_market;

-- 创建用户（可选，用于生产环境）
-- CREATE USER IF NOT EXISTS 'nft_market_user'@'localhost' IDENTIFIED BY 'your_secure_password';
-- GRANT ALL PRIVILEGES ON nft_market.* TO 'nft_market_user'@'localhost';
-- FLUSH PRIVILEGES;

-- 注意：实际的表结构将由GORM的AutoMigrate功能自动创建
-- 以下是表结构的参考（不需要手动执行，仅供参考）

/*
-- 订单表
CREATE TABLE IF NOT EXISTS orders (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    order_id BIGINT UNSIGNED NOT NULL UNIQUE,
    maker VARCHAR(42) NOT NULL,
    nft_contract VARCHAR(42) NOT NULL,
    token_id VARCHAR(255) NOT NULL,
    price VARCHAR(255) NOT NULL,
    amount BIGINT UNSIGNED NOT NULL DEFAULT 1,
    order_type ENUM('limit_sell', 'limit_buy', 'market_sell', 'market_buy') NOT NULL,
    status ENUM('active', 'filled', 'cancelled', 'expired') NOT NULL DEFAULT 'active',
    expiration DATETIME NOT NULL,
    signature TEXT,
    transaction_hash VARCHAR(66),
    block_number BIGINT UNSIGNED,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL,
    
    INDEX idx_maker (maker),
    INDEX idx_nft_contract (nft_contract),
    INDEX idx_token_id (token_id),
    INDEX idx_order_type (order_type),
    INDEX idx_status (status),
    INDEX idx_expiration (expiration),
    INDEX idx_transaction_hash (transaction_hash),
    INDEX idx_block_number (block_number),
    INDEX idx_deleted_at (deleted_at)
);

-- NFT表
CREATE TABLE IF NOT EXISTS nfts (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    contract VARCHAR(42) NOT NULL,
    token_id VARCHAR(255) NOT NULL,
    owner VARCHAR(42) NOT NULL,
    name VARCHAR(255),
    description TEXT,
    image TEXT,
    attributes JSON,
    metadata_uri TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL,
    
    UNIQUE KEY unique_nft (contract, token_id),
    INDEX idx_contract (contract),
    INDEX idx_token_id (token_id),
    INDEX idx_owner (owner),
    INDEX idx_deleted_at (deleted_at)
);

-- 用户表
CREATE TABLE IF NOT EXISTS users (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    address VARCHAR(42) NOT NULL UNIQUE,
    username VARCHAR(255) UNIQUE,
    email VARCHAR(255) UNIQUE,
    avatar TEXT,
    bio TEXT,
    verified BOOLEAN NOT NULL DEFAULT FALSE,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL,
    
    INDEX idx_address (address),
    INDEX idx_username (username),
    INDEX idx_email (email),
    INDEX idx_deleted_at (deleted_at)
);

-- 交易记录表
CREATE TABLE IF NOT EXISTS transactions (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    order_id BIGINT UNSIGNED NOT NULL,
    buyer VARCHAR(42) NOT NULL,
    seller VARCHAR(42) NOT NULL,
    nft_contract VARCHAR(42) NOT NULL,
    token_id VARCHAR(255) NOT NULL,
    price VARCHAR(255) NOT NULL,
    platform_fee VARCHAR(255) NOT NULL,
    transaction_hash VARCHAR(66) NOT NULL UNIQUE,
    block_number BIGINT UNSIGNED NOT NULL,
    block_hash VARCHAR(66) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL,
    
    INDEX idx_order_id (order_id),
    INDEX idx_buyer (buyer),
    INDEX idx_seller (seller),
    INDEX idx_nft_contract (nft_contract),
    INDEX idx_token_id (token_id),
    INDEX idx_transaction_hash (transaction_hash),
    INDEX idx_block_number (block_number),
    INDEX idx_deleted_at (deleted_at)
);
*/

-- 设置时区
SET time_zone = '+00:00';

-- 显示创建结果
SELECT 'NFT市场数据库初始化完成！' as message;
SHOW TABLES;
