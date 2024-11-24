package middleware

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"runtime"
	"yunosphere.com/yun-fiber-scaffold/internal/global"
)

func Recover() fiber.Handler {
	return func(c *fiber.Ctx) error {
		defer func() {
			if err := recover(); err != nil {
				// 初始化一个较大的缓冲区，动态调整直到捕获所有堆栈信息
				stackSize := 4096 // 初始缓冲区大小
				var buf []byte
				for {
					buf = make([]byte, stackSize)
					n := runtime.Stack(buf, false)
					if n < stackSize {
						// 如果捕获到的栈轨迹少于缓冲区大小，说明已经捕获完整
						buf = buf[:n]
						break
					}
					// 如果缓冲区不足，则将大小加倍
					stackSize *= 2
				}

				// 将完整的栈轨迹信息记录到日志
				global.SysLog.WithFields(map[string]interface{}{
					"stack_trace": string(buf), // 确保完整栈轨迹被记录
				}).Errorf("发生运行时异常: %v", err)

				// 返回统一的错误响应
				_ = c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"code":    fiber.StatusInternalServerError,
					"message": fmt.Sprintf("服务器内部错误: %v", err),
				})
			}
		}()

		// 继续处理请求
		return c.Next()
	}
}
