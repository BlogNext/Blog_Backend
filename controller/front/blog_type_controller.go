package front

import (
	"github.com/blog_backend/exception"
	"github.com/blog_backend/help"
	"github.com/blog_backend/service/blog"
)

type BlogTypeController struct {
	BaseController
}

// @获取博客类型列表
// @Description 获取博客类型列表
// @Tags 前台-博客类型（知识库）
// @Accept  application/x-www-form-urlencoded
// @Produce  json
// @Param   per_page     query    int     true        "每页多少条记录"
// @Param   page     query    int     true        "第几页"
// @Success 200 {object} interface{}	"json格式"
// @Router /front/blog_type/get_list [get]
func (c *BlogTypeController) GetList() {
	//必填字段
	type searchRequest struct {
		PerPage int `form:"per_page"`
		Page    int `form:"page"`
	}

	var sr searchRequest

	err := c.Ctx.ShouldBind(&sr)

	if err != nil {
		help.Gin200ErrorResponse(c.Ctx, exception.VALIDATE_ERR, err.Error(), nil)
		return
	}

	//参数默认值
	if sr.PerPage <= 0 {
		sr.PerPage = 10
	}
	if sr.Page <= 0 {
		sr.Page = 1
	}

	service := new(blog.BlogTypeRtService)
	result := service.GetList(sr.PerPage, sr.Page)

	help.Gin200SuccessResponse(c.Ctx, "成功", result)
	return
}
