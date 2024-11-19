package cmd

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"yunosphere.com/yun-fiber-scaffold/configs"
	"yunosphere.com/yun-fiber-scaffold/internal/db"
	"yunosphere.com/yun-fiber-scaffold/internal/logger"
	"yunosphere.com/yun-fiber-scaffold/pkg/router"
)

func Start() {
	// 初始化配置
	configs.InitCfg()
	// 启动应用
	runApp()
}

func runApp() {
	// 创建配置
	initConfig := fiber.Config{
		AppName: configs.Cfg.App.Name,
		// 添加一些优雅关闭相关的配置，确保资源能及时释放
		IdleTimeout:  time.Second * time.Duration(configs.Cfg.App.IdleTimeOut),
		ReadTimeout:  time.Second * time.Duration(configs.Cfg.App.ReadTimeOut),
		WriteTimeout: time.Second * time.Duration(configs.Cfg.App.WriteTimeOut),
	}
	// 初始化 app
	app := fiber.New(initConfig)

	// 初始化请求 ID
	app.Use(requestid.New())
	// 初始化日志库
	app.Use(logger.New())

	// 初始化数据库
	db.New()

	// 设置路由
	setupRoutes(app)

	// 创建一个用于通知服务启动完成的 channel
	started := make(chan struct{})

	// 在 goroutine 中启动服务
	go func() {
		fmt.Println("正在启动服务...")
		// 通知服务开始启动
		close(started)

		if err := app.Listen(fmt.Sprintf(":%s", "8080")); err != nil {
			log.Fatalf("服务启动失败: %v", err)
		}
	}()

	// 等待服务启动
	<-started
	fmt.Println("服务成功启动在端口: 8080 !")

	// 创建通道接收信号（监听系统中断信号（Ctrl+C）和终止信号）
	// 使用缓冲 channel 确保信号不会丢失
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 等待中断信号
	<-quit
	fmt.Println("收到关闭信号，正在优雅关闭服务...")

	// 创建一个带超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 执行清理操作
	if err := gracefulShutdown(ctx, app); err != nil {
		log.Printf("服务关闭过程中发生错误: %v", err)
	}

	fmt.Println("服务已优雅停止...")
}

// setupRoutes 设置路由
func setupRoutes(app *fiber.App) {
	router.YunRouter(app)
}

// gracefulShutdown 处理优雅关闭应用
func gracefulShutdown(ctx context.Context, app *fiber.App) error {
	// 创建一个 channel 来接收关闭结果
	done := make(chan error, 1)
	go func() {
		// 执行 Fiber app 的关闭
		err := app.Shutdown()
		done <- err
	}()

	// 等待关闭完成或上下文超时
	select {
	case <-ctx.Done():
		return fmt.Errorf("服务关闭超时: %v", ctx.Err())
	case err := <-done:
		if err != nil {
			return fmt.Errorf("服务关闭出错: %v", err)
		}
		return nil
	}
}
