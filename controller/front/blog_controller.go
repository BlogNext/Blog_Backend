package front

import (
	"github.com/blog_backend/exception"
	"github.com/blog_backend/help"
	"github.com/blog_backend/service/blog"
)

type BlogController struct {
	BaseController
}

func (c *BlogController) Detail() {
	//必填字段
	type searchRequest struct {
		ID uint `form:"id" binding:"required"`
	}

	var search_request searchRequest

	err := c.Ctx.ShouldBind(&search_request)

	if err != nil {
		help.Gin200ErrorResponse(c.Ctx, exception.VALIDATE_ERR, err.Error(), nil)
		return
	}

	service := new(blog.BlogRtService)
	result := service.Detail(search_request.ID)

	help.Gin200SuccessResponse(c.Ctx, "成功", result)
}

/**
获取博客列表
*/
func (c *BlogController) GetList() {
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

	//过滤参数
	filter := make(map[string]string, 1)
	filter["blog_type_id"] = c.Ctx.DefaultQuery("blog_type_id", "") //分类id过滤

	service := new(blog.BlogRtService)
	result := service.GetList(filter, search_request.PerPage, search_request.Page)

	help.Gin200SuccessResponse(c.Ctx, "成功", result)
}

//搜素
func (c *BlogController) SearchBlog() {

	//非必填字段
	var search_level string
	search_level = c.Ctx.DefaultQuery("search_level", blog.MYSQL_SEARCH_LEVEL)

	//必填字段
	type searchRequest struct {
		//搜索维度
		Keyword string `form:"keyword"`
		PerPage int    `form:"per_page"`
		Page    int    `form:"page"`
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

	b_s := new(blog.BlogRtService)

	result := b_s.SearchBlog(search_level, search_request.Keyword, search_request.PerPage, search_request.Page)

	help.Gin200SuccessResponse(c.Ctx, "请求成功过", result)

	return
}