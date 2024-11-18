package logger

import (
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"log"
	"os"
	"path"
	"sync"
	"time"
	"yunosphere.com/yun-fiber-scaffold/configs"
	"yunosphere.com/yun-fiber-scaffold/internal/global"
	"yunosphere.com/yun-fiber-scaffold/internal/utils"

	"github.com/gofiber/fiber/v2"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

var once sync.Once // 确保日志只初始化一次

func initLogger() {
	// 确保日志只初始化一次
	once.Do(func() {
		logFilePath := configs.Cfg.Logger.LogFilePath
		logFileName := configs.Cfg.Logger.LogFileName
		logTimestampFmt := configs.Cfg.Logger.LogTimestampFmt
		logLevel := configs.Cfg.Logger.LogLevel
		logMaxAge := configs.Cfg.Logger.LogMaxAge
		logRotationTime := configs.Cfg.Logger.LogRotationTime

		// 创建日志目录
		if err := os.MkdirAll(logFilePath, 0755); err != nil {
			log.Fatalf("创建日志目录失败: %v", err)
		}

		// 初始化 logrus
		logger := logrus.New()

		// 设置日志格式
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: logTimestampFmt,
		})

		// 设置日志级别
		logLevelParsed, err := logrus.ParseLevel(logLevel)
		if err != nil {
			log.Fatalf("日志级别解析失败: %v", err)
		}
		logger.SetLevel(logLevelParsed)

		// 配置日志轮转
		fileName := path.Join(logFilePath, logFileName)
		maxAge := time.Duration(logMaxAge) * time.Hour
		rotationTime := time.Duration(logRotationTime) * time.Hour

		writer, err := rotatelogs.New(
			path.Join(logFilePath, "%Y%m%d.log"),
			rotatelogs.WithLinkName(fileName),
			rotatelogs.WithMaxAge(maxAge),
			rotatelogs.WithRotationTime(rotationTime),
		)
		if err != nil {
			log.Fatalf("设置日志轮转失败: %v", err)
		}

		// 配置日志级别与轮转日志的映射
		writerMap := lfshook.WriterMap{
			logrus.InfoLevel:  writer,
			logrus.FatalLevel: writer,
			logrus.DebugLevel: writer,
			logrus.WarnLevel:  writer,
			logrus.ErrorLevel: writer,
			logrus.PanicLevel: writer,
		}

		// 添加 Hook，使用 JSONFormatter
		hook := lfshook.NewHook(writerMap, &logrus.JSONFormatter{
			TimestampFormat: logTimestampFmt,
		})

		// 将 logger 的默认输出设置为标准输出，文件输出由 hook 处理
		logger.SetOutput(os.Stdout)
		logger.AddHook(hook)

		global.SysLog = logger
	})
}

// New 返回 Fiber 中间件
func New() fiber.Handler {
	initLogger() // 确保日志已初始化

	return func(c *fiber.Ctx) error {
		// 从 Fiber 上下文获取请求 ID
		reqId := c.Locals(requestid.ConfigDefault.ContextKey)

		// 创建带有请求相关字段的日志实例
		bizLog := global.SysLog.WithFields(logrus.Fields{
			"requestId": reqId,
			"requestIp": c.IP(), // Fiber 获取客户端IP的方法
		})

		// 将日志实例存储到 Fiber 上下文中
		utils.SetBizLogger(c, bizLog)

		// 继续处理请求
		return c.Next()
	}
}
