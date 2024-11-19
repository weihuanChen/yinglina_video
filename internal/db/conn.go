package db

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
	"yunosphere.com/yun-fiber-scaffold/configs"
	"yunosphere.com/yun-fiber-scaffold/internal/global"
)

// 初始化数据库实例
func initDB() error {
	dsn := buildDSN(configs.Cfg)
	db, err := gorm.Open(postgres.Open(dsn), gormConfig())
	if err != nil {
		handleDBError(err)
		return err
	}

	global.DB = db
	err = setupDBConnectionPool()
	if err != nil {
		return err
	}
	return nil
}

// 配置数据库连接池
func setupDBConnectionPool() error {
	SqlDB, err := global.DB.DB()
	if err != nil {
		handleDBError(err)
		return err
	}

	SqlDB.SetMaxOpenConns(100)                 // 设置数据库的最大打开连接数为100
	SqlDB.SetMaxIdleConns(10)                  // 设置数据库的最大空闲连接数为10
	SqlDB.SetConnMaxLifetime(10 * time.Second) // 设置数据库连接的最大存活时间为10秒

	return nil
}

func handleDBError(err error) {
	global.SysLog.Errorf("数据库连接失败: %v", err)
}

func gormConfig() *gorm.Config {
	return &gorm.Config{
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold:             2 * time.Second,
				Colorful:                  true,
				IgnoreRecordNotFoundError: false,
				LogLevel:                  logger.Silent,
			}),
	}
}

func buildDSN(config configs.AppConfig, tempDbName ...string) string {
	var DbName string

	// 如果传入 tempDbName 参数，则使用传入的数据库名称
	if len(tempDbName) > 0 && tempDbName[0] != "" {
		DbName = tempDbName[0]
	} else {
		DbName = config.Db.Name // 默认使用配置文件中的数据库名称
	}

	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Db.Host,
		config.Db.Port,
		config.Db.User,
		config.Db.Psw,
		DbName,
	)
}

func localDBExists(config configs.AppConfig) bool {
	dsn := buildDSN(config)
	db, err := gorm.Open(postgres.Open(dsn), gormConfig())
	if err != nil {
		// 数据库不存在，返回 false
		return false
	}
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	// 如果没有错误（err == nil），说明数据库存在，返回 true；否则返回 false
	return db.Exec("SELECT 1").Error == nil
}

func checkAndCreateDB(config configs.AppConfig) error {
	// 连接到指定 postgres 数据库
	postgresDSN := buildDSN(config, "postgres")

	// 使用 gorm 连接到 postgres 数据库
	db, err := gorm.Open(postgres.Open(postgresDSN), gormConfig())
	if err != nil {
		handleDBError(err)
		return err
	}
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	// 检查数据库是否存在
	var exists bool
	checkDBQuery := fmt.Sprintf("SELECT 1 FROM pg_database WHERE datname = '%s'", config.Db.Name)
	err = db.Raw(checkDBQuery).Scan(&exists).Error
	if err != nil {
		handleDBError(err)
		return err
	}

	// 如果数据库不存在，则创建并指定编码为 UTF8
	if !exists {
		global.SysLog.Infof("PostgreSQL 数据库 %s 不存在，正在创建...", config.Db.Name)
		createDBQuery := fmt.Sprintf("CREATE DATABASE %s WITH ENCODING 'UTF8'", config.Db.Name)
		err = db.Exec(createDBQuery).Error
		if err != nil {
			handleDBError(err)
			return err
		}
		global.SysLog.Infof("PostgreSQL 数据库 %s 创建成功", config.Db.Name)
	}

	return nil
}

func New() {
	// 检查本地是否存在目标数据库
	if localDBExists(configs.Cfg) {
		err := initDB()
		if err == nil {
			global.SysLog.Infof("PostgreSQL 数据库连接成功!")
		}
	} else {
		// 数据库不存在，尝试创建
		err := checkAndCreateDB(configs.Cfg)
		if err != nil {
			_ = fmt.Errorf("DB Create Error: %v", err)
			return
		}
		// 创建成功后初始化数据库连接
		err = initDB()
		if err == nil {
			global.SysLog.Infof("PostgreSQL 数据库连接成功!")
		}
	}
}
