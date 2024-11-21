package redis

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"strconv"
	"time"
	"yunosphere.com/yun-fiber-scaffold/configs"
	"yunosphere.com/yun-fiber-scaffold/internal/global"
)

func New() {
	// 默认连接池
	pool := defaultPool(configs.Cfg)
	if pool == nil {
		return
	}
	global.RedisPool = pool
	global.SysLog.Infof("Redis 连接池创建成功 !")
}

func defaultPool(config configs.AppConfig) *redis.Pool {
	db, _ := strconv.Atoi(config.Redis.Db)
	return &redis.Pool{
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial(
				"tcp",
				fmt.Sprintf("%s:%s", config.Redis.Host, config.Redis.Port),
				redis.DialReadTimeout(time.Second),
				redis.DialWriteTimeout(time.Second*2),
				redis.DialConnectTimeout(10*time.Second),
				redis.DialDatabase(db),
				redis.DialPassword(config.Redis.Psw),
			)
			if err != nil {
				global.SysLog.Errorf("Redis 连接失败: %v", err)
				return nil, err // 返回 nil 而非无效的连接
			}
			return conn, nil
		},
		MaxIdle:     50,
		MaxActive:   100, // 最大连接数
		IdleTimeout: 60 * time.Second,
	}
}
