package gateway

import (
	"github.com/blog_backend/controller"
	"github.com/blog_backend/controller/gateway"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

func RegisterGateWayRouter(router *gin.Engine) {
	//注册验证器d
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		//自定义验证器，中文信息
		zhCh := zh.New()
		uni := ut.New(zhCh)
		trans, _ := uni.GetTranslator("zh")
		_ = zh_translations.RegisterDefaultTranslations(v, trans)
	}

	//gateway公共的路由
	gatewayRouter := router.Group("/gateway")
	{
		//语雀路由
		yuqueRouter := gatewayRouter.Group("/yuque")
		{
			yuqueController := controller.NewController(new(gateway.YuqueController))
			yuqueRouter.Any("/:action", yuqueController)
		}

	}

}
