package oauth

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	"github.com/blog_backend/controller"
	"github.com/blog_backend/controller/oauth"
	validate_request "github.com/blog_backend/controller/oauth/validate-request"
)

func RegisterOauthRouter(router *gin.Engine) {
	//注册验证器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		//自定义验证器，中文信息
		zh_ch := zh.New()
		uni := ut.New(zh_ch)
		trans, _ := uni.GetTranslator("zh")
		_ = zh_translations.RegisterDefaultTranslations(v, trans)

		//注册验证
		_ = v.RegisterValidation("my_phone", validate_request.MyPhone)
		_ = v.RegisterValidation("redirect_url", validate_request.RedirectUrl)
		_ = v.RegisterValidation("my_email", validate_request.MyEmail)
	}

	//后台路由
	oauth_router := router.Group("/oauth")
	{
		oauth_controller := controller.NewController(new(oauth.OauthController))
		oauth_router.Any("/:action", oauth_controller)
	}
}
