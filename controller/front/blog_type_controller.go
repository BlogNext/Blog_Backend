package front

import (
	"github.com/blog_backend/exception"
	"github.com/blog_backend/help"
	"github.com/blog_backend/service/blog"
)

type BlogTypeController struct {
	BaseController
}

/**
获取博客类型列表
*/
func (c *BlogTypeController) GetList() {
	//必填字段
	type searchRequest struct {
		PerPage int `form:"per_page"`
		Page    int `form:"page"`
	}

	var search_request searchRequest

	err := c.Ctx.ShouldBind(&search_request)

	if err != nil {
		help.Gin200ErrorResponse(c.Ctx, exception.VALIDATE_ERR, err.Error(), nil)
		return
	}

	//参数默认值
	if search_request.PerPage <= 0 {
		search_request.PerPage = 10
	}
	if search_request.Page <= 0 {
		search_request.Page = 1
	}
	
	service := new(blog.BlogTypeRtService)
	result := service.GetList(search_request.PerPage, search_request.Page)

	help.Gin200SuccessResponse(c.Ctx, "成功", result)
	return
}
