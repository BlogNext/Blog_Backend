package router

import (
	"github.com/blog_backend/router/backend"
	"github.com/blog_backend/router/front"
	"github.com/blog_backend/router/gateway"
	"github.com/gin-gonic/gin"
)

type MyRouter struct {
	routerList []RegisterRouter
}

func (mr *MyRouter) registerRouter() {
	mr.routerList = make([]RegisterRouter, 3)
	//路由注册
	mr.routerList[0] = gateway.RegisterGateWayRouter
	mr.routerList[1] = front.RegisterFrontRouter
	mr.routerList[2] = backend.RegisterBackendRouter
}

func (mr *MyRouter) RunRouter() *gin.Engine {

	//默认启动方式，包含 Logger、Recovery 中间件
	router := gin.Default()
	//加载html文件路径
	//router.LoadHTMLGlob("templates/**/*")

	//注册验证器

	//注册路由
	mr.registerRouter()

	for _, registerRouter := range mr.routerList {
		registerRouter(router)
	}

	return router
}

//注册路由的函数
type RegisterRouter func(router *gin.Engine)
