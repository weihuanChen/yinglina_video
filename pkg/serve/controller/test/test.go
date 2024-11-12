package test

import (
	"github.com/gofiber/fiber/v2"
	"time"
)

func Ping(c *fiber.Ctx) error {
	return c.SendString("Pong!\n")
}

func Hello(c *fiber.Ctx) error {
	return c.SendString("Hello, Fiber 🎉!\n")
}

func LongReq(c *fiber.Ctx) error {
	time.Sleep(20 * time.Second) // 模拟长时间处理的请求
	return c.SendString("耗时请求处理完成 !\n")
}
