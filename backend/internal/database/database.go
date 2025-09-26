package database

import (
	"nft-market/internal/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Init 初始化数据库连接
func Init(databaseURL string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(databaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	// 自动迁移数据库表结构
	err = db.AutoMigrate(
		&models.Collection{},
		&models.Item{},
		&models.Order{},
		&models.Activity{},
		&models.User{},
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}
