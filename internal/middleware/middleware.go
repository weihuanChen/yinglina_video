package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"yunosphere.com/yun-fiber-scaffold/internal/logger"
)

func InitMiddleware(app *fiber.App) {
	// 初始化请求 ID
	app.Use(requestid.New())
	// 初始化日志库
	app.Use(logger.New())
	// 初始化全局错误处理中间件
	app.Use(GlobalErrorHandler())
	// 初始化全局运行时异常捕获中间件
	app.Use(Recover())
}
