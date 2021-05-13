package front

import (
	"github.com/blog_backend/entity/login/front"
	"github.com/blog_backend/exception"
	"github.com/blog_backend/help"
	"github.com/blog_backend/service/login"
)

type LoginController struct {
	BaseController
}

// @语雀账号登录
// @Description 语雀账号登录
// @Tags 前台-登录
// @Accept  application/x-www-form-urlencoded
// @Produce  json
// @Param   login  query   string     true        "语雀login"
// @Param   password  query   string     true        "xiaochen123"
// @Success 200 {object} interface{}	"json格式"
// @Router /front/login/Login_by_yuque [post]
func (u *LoginController) LoginByYuque() {

	//必填字段
	type loginRequest struct {
		Login    string `form:"login" binding:"required"`
		Password string `form:"password" binding:"required"`
	}

	var lr loginRequest

	err := u.Ctx.ShouldBind(&lr)

	if err != nil {
		help.Gin200ErrorResponse(u.Ctx, exception.VALIDATE_ERR, err.Error(), nil)
		return
	}

	service := new(login.LoginRtService)

	tokenString := service.LoginByYuque(lr.Login, lr.Password)

	help.Gin200SuccessResponse(u.Ctx, "成功", tokenString)
	return
}

func (u *LoginController) IsLogin() {
	token := u.Ctx.GetHeader("x-access-token")
	service := new(login.LoginRtService)
	loginEntity := new(front.LoginEntity)
	isLogin := service.IsLogin(token, loginEntity)

	result := make(map[string]interface{})
	result["is_login"] = isLogin
	result["login_entity"] = loginEntity

	help.Gin200SuccessResponse(u.Ctx, "成功", result)
	return
}

// @oauthSSo单点登录
// @Description 基于BlogNext的用户体系和oauth协议做的预授权码单点登录
// @Tags 前台-登录
// @Accept  application/x-www-form-urlencoded
// @Produce  json
// @Param   pre_code  query   string     true        "预授权码"
// @Success 200 {object} interface{}	"json格式"
// @Router /front/login/login_blog_next_pre_code [post]
func (u *LoginController) LoginBlogNextPreCode() {

	var request front.BlogNextPreCodeRequest

	err := u.Ctx.ShouldBindJSON(&request)

	if err != nil {
		help.Gin200ErrorResponse(u.Ctx, exception.VALIDATE_ERR, err.Error(), nil)
		return
	}

	service := new(login.LoginRtService)
	response := service.BlogNextPreCode(&request)

	help.Gin200SuccessResponse(u.Ctx, "成功", response)
	return
}
