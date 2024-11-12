package router

import (
	"github.com/gofiber/fiber/v2"
)

// YunRouter 初始化路由
func YunRouter(app *fiber.App) {
	// 创建一个带统一前缀的 API 路由组 /api/v1
	api1 := app.Group("/api/v1")
	// 创建一个带统一前缀的 API 路由组 /api/v2
	api2 := app.Group("/api/v2")

	// 添加 test 相关的多版本 api
	testRouter(api1, api2)
}
