package global

import (
	"github.com/gomodule/redigo/redis"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"os"
)

// 数据库
var (
	DB        *gorm.DB    // 全局的 db 对象
	RedisPool *redis.Pool // 全局的 redis 连接池对象
)

// 日志
var (
	SysLog  *log.Logger // 全局系统级日志对象
	BizLog  *log.Entry  // 作用于每次请求的 log entry
	LogFile *os.File    // 存储日志的文件
)
