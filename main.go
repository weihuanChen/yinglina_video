package main

import "github.com/gofiber/fiber/v2"

func main() {
	// 初始化配置
	initConfig := fiber.Config{
		AppName: "yun-fiber-scaffold",
	}

	// 创建新的 fiber app
	app := fiber.New(initConfig)

	// 指定 / 路由的功能
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, Fiber 🎉!\n")
	})

	// 启动应用
	if err := app.Listen(":8080"); err != nil {
		panic(err)
	}
}
