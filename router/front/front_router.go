package front

import (
	"github.com/blog_backend/controller"
	"github.com/blog_backend/controller/front"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

func RegisterFrontRouter(router *gin.Engine) {
	//注册验证器d
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		//自定义验证器，中文信息
		zh_ch := zh.New()
		uni := ut.New(zh_ch)
		trans, _ := uni.GetTranslator("zh")
		_ = zh_translations.RegisterDefaultTranslations(v, trans)
	}

	//前端公共的路由
	front_router := router.Group("/front")
	{
		//博客类型路由
		blog_type_router := front_router.Group("/blog_type")
		{
			blog_type_controller := controller.NewController(new(front.BlogTypeController))
			blog_type_router.Any("/:action", blog_type_controller)
		}

		//博客路由
		blog_router := front_router.Group("/blog")
		{
			blog_controller := controller.NewController(new(front.BlogController))
			blog_router.Any("/:action", blog_controller)
		}

		//登录
		login_router := front_router.Group("/login")
		{
			login_controller := controller.NewController(new(front.LoginController))
			login_router.Any("/:action", login_controller)
		}

	}

}
