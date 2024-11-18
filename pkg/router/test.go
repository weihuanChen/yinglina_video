package router

import (
	"github.com/gofiber/fiber/v2"
	"yunosphere.com/yun-fiber-scaffold/pkg/serve/controller/test"
)

// testRouter 测试路由
func testRouter(r ...fiber.Router) {
	// api v1 group
	apiV1 := r[0]
	testGroupV1 := apiV1.Group("/test")
	testGroupV1.Get("/ping", test.Ping)
	testGroupV1.Get("/hello", test.Hello)
	testGroupV1.Get("/testLogger", test.TestLogger)

	// api v2 group
	apiV2 := r[1]
	testGroupV2 := apiV2.Group("/test")
	testGroupV2.Get("/long", test.LongReq)
}
