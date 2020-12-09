package front

import (
	"github.com/blog_backend/exception"
	"github.com/blog_backend/help"
	"github.com/blog_backend/service/blog"
)

type BlogController struct {
	BaseController
}

// @博客详情
// @Description 博客详情
// @Tags 前台-博客
// @Accept  application/x-www-form-urlencoded
// @Produce  json
// @Param   id     query    uint     true        "博客id"
// @Success 200 {object} interface{}	"json格式"
// @Router /front/blog/detail [get]
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

	//获取详情
	result := service.Detail(search_request.ID)
	//浏览量加1
	service.IncBrowse(search_request.ID)

	help.Gin200SuccessResponse(c.Ctx, "成功", result)

	return
}

// @获取博客列表
// @Description 获取博客列表
// @Tags 前台-博客
// @Accept  application/x-www-form-urlencoded
// @Produce  json
// @Param   per_page     query    int     true    "一页多少条"
// @Param   page     query    int     true        "第几页"
// @Param   blog_type_id     query    int     false        "博客分类"
// @Success 200 {object} interface{}	"json格式"
// @Router /front/blog/get_list [get]
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

// @搜素博客
// @Description 搜素博客
// @Tags 前台-博客
// @Accept  application/x-www-form-urlencoded
// @Produce  json
// @Param   per_page     query    int     true    "一页多少条"
// @Param   page     query    int     true        "第几页"
// @Param   search_level     query    string     false        "搜索等级，默认mysql搜索"
// @Param   keyword     query    string     false        "搜索关键字"
// @Success 200 {object} interface{}	"json格式"
// @Router /front/blog/search_blog [get]
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

// @按排序维度获取排序博客
// @Description 按排序维度获取排序博客
// @Tags 前台-博客
// @Accept  application/x-www-form-urlencoded
// @Produce  json
// @Param   per_page     query    int     false    "一页多少条，默认值5"
// @Param   sort_dimension     query    string     false        "排序维度，默认值browse_total"
// @Success 200 {object} interface{}	"json格式"
// @Router /front/blog/get_list_by_sort [get]
func (c *BlogController) GetListBySort() {
	//必填字段
	type sortRequest struct {
		//搜索维度
		SortDimension string `form:"sort_dimension"`
		PerPage       int    `form:"per_page"`
	}

	var sort_request sortRequest

	err := c.Ctx.ShouldBind(&sort_request)

	if err != nil {
		help.Gin200ErrorResponse(c.Ctx, exception.VALIDATE_ERR, err.Error(), nil)
		return
	}

	if sort_request.PerPage <= 0 {
		sort_request.PerPage = 5
	}

	if sort_request.SortDimension == "" {
		sort_request.SortDimension = "browse_total"
	}

	b_s := new(blog.BlogRtService)
	result := b_s.GetListBySort(sort_request.SortDimension, sort_request.PerPage)

	help.Gin200SuccessResponse(c.Ctx, "请求成功过", result)

	return

}
