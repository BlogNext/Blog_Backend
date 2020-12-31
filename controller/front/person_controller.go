package front

import (
	"github.com/blog_backend/exception"
	"github.com/blog_backend/help"
	"github.com/blog_backend/service/blog"
)

type PersonController struct {
	BaseController
}

// @私人博客列表，这里的私人只的是登录的用户
// @Description 私人博客列表，这里的私人只的是登录的用户
// @Tags 前台-登录-私人博客
// @Security ApiKeyAuth
// @Accept  application/x-www-form-urlencoded
// @Produce  json
// @Param   per_page     query    int     true    "一页多少条"
// @Param   page     query    int     true        "第几页"
// @Param   blog_type_id     query    int     false        "博客分类"
// @Success 200 {object} interface{}	"json格式"
// @Router /front/person/blog_list [get]
func (p *PersonController) BlogList() {
	//必填字段
	type searchRequest struct {
		PerPage int `form:"per_page"`
		Page    int `form:"page"`
	}

	var sr searchRequest

	err := p.Ctx.ShouldBind(&sr)

	if err != nil {
		help.Gin200ErrorResponse(p.Ctx, exception.VALIDATE_ERR, err.Error(), nil)
		return
	}

	//参数默认值
	if sr.PerPage <= 0 {
		sr.PerPage = 10
	}
	if sr.Page <= 0 {
		sr.Page = 1
	}

	service := new(blog.BlogRtService)
	result := service.GetListByPerson(sr.PerPage, sr.Page)

	help.Gin200SuccessResponse(p.Ctx, "成功", result)
}
