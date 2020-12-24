package front

import (
	"github.com/blog_backend/exception"
	"github.com/blog_backend/help"
	"github.com/blog_backend/service/blog"
)

type PersonController struct {
	BaseController
}

func (p *PersonController) Blog() {
	//必填字段
	type searchRequest struct {
		PerPage int `form:"per_page"`
		Page    int `form:"page"`
	}

	var search_request searchRequest

	err := p.Ctx.ShouldBind(&search_request)

	if err != nil {
		help.Gin200ErrorResponse(p.Ctx, exception.VALIDATE_ERR, err.Error(), nil)
		return
	}

	//参数默认值
	if search_request.PerPage <= 0 {
		search_request.PerPage = 10
	}
	if search_request.Page <= 0 {
		search_request.Page = 1
	}

	//过滤参数
	filter := make(map[string]string, 1)
	filter["blog_type_id"] = p.Ctx.DefaultQuery("blog_type_id", "") //分类id过滤

	service := new(blog.BlogRtService)
	result := service.GetList(filter, search_request.PerPage, search_request.Page)

	help.Gin200SuccessResponse(p.Ctx, "成功", result)
}
