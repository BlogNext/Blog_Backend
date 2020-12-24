package front

import (
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
// @Param   password  query   string     true        "登录密码"
// @Success 200 {object} interface{}	"json格式"
// @Router /front/login/Login_by_yuque [post]
func (u *LoginController) LoginByYuque() {

	//必填字段
	type loginRequest struct {
		Login    string `form:"login" binding:"required"`
		Password string `form:"password" binding:"required"`
	}

	var login_request loginRequest

	err := u.Ctx.ShouldBind(&login_request)

	if err != nil {
		help.Gin200ErrorResponse(u.Ctx, exception.VALIDATE_ERR, err.Error(), nil)
		return
	}

	service := new(login.LoginRtService)

	login_entity := service.LoginByYuque(login_request.Login, login_request.Password)

	help.Gin200SuccessResponse(u.Ctx, "成功", login_entity)
	return
}
