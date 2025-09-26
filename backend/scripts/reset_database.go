package main

import (
	"fmt"
	"nft-market/internal/config"
	"nft-market/internal/database"
	"nft-market/internal/models"

	"github.com/joho/godotenv"
)

func main() {
	// 加载环境变量
	if err := godotenv.Load(); err != nil {
		fmt.Println("未找到.env文件，使用默认配置")
	}

	// 加载配置
	cfg := config.Load()

	// 初始化数据库连接
	db, err := database.InitWithoutMigration(cfg.DatabaseURL)
	if err != nil {
		panic("数据库连接失败: " + err.Error())
	}

	fmt.Println("开始重置数据库表...")

	// 删除现有表（如果存在）
	tables := []string{
		"orders", "activities", "items", "collections", "users",
	}

	for _, table := range tables {
		if err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", table)).Error; err != nil {
			fmt.Printf("删除表 %s 失败: %v\n", table, err)
		} else {
			fmt.Printf("删除表 %s 成功\n", table)
		}
	}

	fmt.Println("开始重新创建表...")

	// 重新创建表
	err = db.AutoMigrate(
		&models.Collection{},
		&models.Item{},
		&models.Order{},
		&models.Activity{},
		&models.User{},
	)
	if err != nil {
		panic("表创建失败: " + err.Error())
	}

	fmt.Println("数据库表重置完成！")

	// 验证表结构
	var result []map[string]interface{}
	if err := db.Raw("DESCRIBE orders").Scan(&result).Error; err != nil {
		fmt.Printf("查看orders表结构失败: %v\n", err)
	} else {
		fmt.Println("orders表结构:")
		for _, row := range result {
			fmt.Printf("  %s: %s\n", row["Field"], row["Type"])
		}
	}
}
