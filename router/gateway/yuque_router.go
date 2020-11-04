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

func RegisterYuqueRouter(router *gin.Engine) {
	//注册验证器d
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		//自定义验证器，中文信息
		zh_ch := zh.New()
		uni := ut.New(zh_ch)
		trans, _ := uni.GetTranslator("zh")
		_ = zh_translations.RegisterDefaultTranslations(v, trans)
	}

	//后端公共的路由
	gateway_router := router.Group("/gateway")
	{
		//语雀路由
		yuque_router := gateway_router.Group("/yuque")
		{
			yuque_controller := controller.NewController(new(gateway.YuqueController))
			yuque_router.Any("/:action", yuque_controller)
		}

	}

}
