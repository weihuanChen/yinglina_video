package test

import (
	"github.com/gofiber/fiber/v2"
	"time"
	"yunosphere.com/yun-fiber-scaffold/internal/error"
	"yunosphere.com/yun-fiber-scaffold/internal/global"
	"yunosphere.com/yun-fiber-scaffold/internal/utils"
	"yunosphere.com/yun-fiber-scaffold/pkg/vo"
)

func Ping(c *fiber.Ctx) error {
	return c.SendString("Hello go!\n")
}

func Hello(c *fiber.Ctx) error {
	return c.SendString("Hello, Fiber 🎉!\n")
}

func LongReq(c *fiber.Ctx) error {
	time.Sleep(20 * time.Second) // 模拟长时间处理的请求
	return c.SendString("耗时请求处理完成 !\n")
}

func TestLogger(c *fiber.Ctx) error {
	utils.BizLogger(c).Infof("TestLogger...")
	return c.SendString("测试日志成功!")
}

func TestRedis(c *fiber.Ctx) error {
	utils.BizLogger(c).Infof("开始写入缓存...")
	// 初始化缓存连接
	conn := global.RedisPool.Get()
	defer conn.Close()
	_, err := conn.Do("SET", "TEST:", "测试 value")
	if err != nil {
		utils.BizLogger(c).Errorf("测试写入缓存失败: %v", err)
		return err
	}
	utils.BizLogger(c).Infof("写入缓存成功...")

	utils.BizLogger(c).Infof("开始读取缓存...")
	// 这里可以复用 conn 打开的连接
	articlesCache, err := conn.Do("GET", "TEST:")
	if err != nil {
		utils.BizLogger(c).Errorf("测试读取缓存失败: %v", err)
		return err
	}
	utils.BizLogger(c).Infof("读取缓存成功, key: %s , value: %s", "TEST:", articlesCache)
	return c.SendString("测试缓存功能完成!")
}

func TestSuccRes(c *fiber.Ctx) error {
	return c.JSON(vo.Success(nil, c))
}

func TestErrRes(c *fiber.Ctx) error {
	return c.JSON(vo.Fail(biz_err.New(biz_err.ServerError), nil, c))
}

func TestErrorMiddleware(c *fiber.Ctx) error {
	// 模拟业务异常
	panic("发生业务异常")
}
