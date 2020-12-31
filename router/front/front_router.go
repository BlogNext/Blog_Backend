package front

import (
	"github.com/blog_backend/controller"
	"github.com/blog_backend/controller/front"
	middleware_front "github.com/blog_backend/middleware/front"
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
		zhCh := zh.New()
		uni := ut.New(zhCh)
		trans, _ := uni.GetTranslator("zh")
		_ = zh_translations.RegisterDefaultTranslations(v, trans)
	}

	//前端公共的路由
	frontRouter := router.Group("/front")
	{
		//博客类型路由
		blogTypeRouter := frontRouter.Group("/blog_type")
		{
			blogTypeController := controller.NewController(new(front.BlogTypeController))
			blogTypeRouter.Any("/:action", blogTypeController)
		}

		//博客路由
		blogRouter := frontRouter.Group("/blog")
		{
			blogController := controller.NewController(new(front.BlogController))
			blogRouter.Any("/:action", blogController)
		}

		//登录
		loginRouter := frontRouter.Group("/login")
		{
			loginController := controller.NewController(new(front.LoginController))
			loginRouter.Any("/:action", loginController)
		}

		//私人空间路由
		personRouter := frontRouter.Group("/person")
		{
			personRouter.Use(middleware_front.LoginMiddleware())
			personController := controller.NewController(new(front.PersonController))
			personRouter.Any("/:action", personController)
		}

	}

}
