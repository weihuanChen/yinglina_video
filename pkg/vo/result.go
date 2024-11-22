package vo

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"time"
	bizerr "yunosphere.com/yun-fiber-scaffold/internal/error"
)

type Result struct {
	Code      int         `json:"code"`
	Msg       string      `json:"msg"`
	Data      interface{} `json:"data"`
	RequestId interface{} `json:"requestId"`
	TimeStamp int64       `json:"timeStamp"`
}

func Success(data interface{}, c *fiber.Ctx) Result {
	return Result{
		Code:      0, // 0 表示成功
		Msg:       "success",
		Data:      data,
		RequestId: c.Locals(requestid.ConfigDefault.ContextKey),
		TimeStamp: time.Now().UnixMilli(), // 使用毫秒级时间戳
	}
}

func Fail(err *bizerr.Err, data interface{}, c *fiber.Ctx) Result {
	return Result{
		Code:      err.Code,
		Msg:       err.Msg,
		Data:      data,
		RequestId: c.Locals(requestid.ConfigDefault.ContextKey),
		TimeStamp: time.Now().UnixMilli(),
	}
}
