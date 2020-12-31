package backend

import (
	"github.com/blog_backend/controller"
	"github.com/blog_backend/controller/backend"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

func RegisterBackendRouter(router *gin.Engine) {
	//注册验证器d
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		//自定义验证器，中文信息
		zhCh := zh.New()
		uni := ut.New(zhCh)
		trans, _ := uni.GetTranslator("zh")
		_ = zh_translations.RegisterDefaultTranslations(v, trans)
	}

	//后端公共的路由
	gatewayRouter := router.Group("/backend")
	{
		//博客路由
		attachmentRouter := gatewayRouter.Group("/attachment")
		{
			attachmentController := controller.NewController(new(backend.AttachmentController))
			attachmentRouter.Any("/:action", attachmentController)
		}

		//博客路由
		blogRouter := gatewayRouter.Group("/blog")
		{
			blogController := controller.NewController(new(backend.BlogController))
			blogRouter.Any("/:action", blogController)
		}

	}

}
