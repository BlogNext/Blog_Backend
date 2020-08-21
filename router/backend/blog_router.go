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

func RegisterBlogRouter(router *gin.Engine) {
	//注册验证器d
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		//自定义验证器，中文信息
		zh_ch := zh.New()
		uni := ut.New(zh_ch)
		trans, _ := uni.GetTranslator("zh")
		_ = zh_translations.RegisterDefaultTranslations(v, trans)
	}

	//后台博客类型路由
	blog_router := router.Group("/blog_type")
	{
		blog_type_controller := controller.NewController(new(backend.BlogTypeController))
		blog_router.Any("/:action", blog_type_controller)
	}
}