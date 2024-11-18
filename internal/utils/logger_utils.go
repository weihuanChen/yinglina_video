package utils

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"yunosphere.com/yun-fiber-scaffold/internal/global"
)

const (
	bizLog = "BizLog"
)

// SetBizLogger 设置 BizLogger 实例
func SetBizLogger(c *fiber.Ctx, logger *logrus.Entry) {
	c.Locals(bizLog, logger)
}

// BizLogger 获取 BizLogger 实例
func BizLogger(c *fiber.Ctx) *logrus.Entry {
	// 如果 context 是 nil，返回全局的系统日志对象，避免空指针
	if c == nil {
		return logrus.NewEntry(global.SysLog)
	}

	if logger, ok := c.Locals(bizLog).(*logrus.Entry); ok {
		return logger
	}
	return global.SysLog.WithFields(logrus.Fields{}) // 返回默认日志器
}
