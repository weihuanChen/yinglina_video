package middleware

import (
	"errors"
	"fmt"
	"yunosphere.com/yun-fiber-scaffold/internal/error"
	"yunosphere.com/yun-fiber-scaffold/pkg/vo"

	"github.com/gofiber/fiber/v2"
)

func GlobalErrorHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 继续处理请求
		err := c.Next()

		// 如果有错误
		if err != nil {
			// 默认的响应码
			code := fiber.StatusInternalServerError

			// 判断是否为自定义的业务异常
			var e *biz_err.Err
			if errors.As(err, &e) {
				code = e.Code
			}

			// 打印错误信息，方便调试
			fmt.Printf("请求异常: %v\n", err)

			// 返回封装后的 Result 结构
			return c.Status(fiber.StatusInternalServerError).JSON(vo.Fail(biz_err.New(code, err.Error()), nil, c))
		}

		// 如果没有错误，继续正常流程
		return nil
	}
}
