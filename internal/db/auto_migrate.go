package db

import (
	"log"
	"yunosphere.com/yun-fiber-scaffold/internal/global"
	"yunosphere.com/yun-fiber-scaffold/internal/model"
)

// AutoMigrate 执行数据库自动迁移
func AutoMigrate() {
	// 要确保数据库连接已成功初始化
	if global.DB == nil {
		log.Fatal("数据库初始化失败，无法执行自动迁移...")
	}

	// 自动迁移模型
	err := global.DB.AutoMigrate(
		model.GetAllModels()...,
	)
	if err != nil {
		log.Fatalf("数据库自动迁移失败: %v", err)
	}

	log.Println("数据库自动迁移成功")
}
