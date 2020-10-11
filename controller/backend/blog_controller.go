package backend

import (
	"github.com/blog_backend/exception"
	"github.com/blog_backend/help"
	"github.com/blog_backend/service/backend"
	"github.com/blog_backend/service/front"
)

type BlogController struct {
	BackendController
}

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

	b_s := new(backend.BlogService)
	result := b_s.GetList(search_request.PerPage, search_request.Page)
	help.Gin200SuccessResponse(c.Ctx, "成功", result)
	return
}

func (c *BlogController) AddBlog() {

	//定义获取方法参数
	type addRequest struct {
		BlogTypeId  int64  `form:"blog_type_id" binding:"required"`
		Title       string `form:"title"  binding:"required"`
		Abstract    string `form:"abstract"  binding:"required"`
		Content     string `form:"content"  binding:"required"`
		CoverPlanId int64  `form:"cover_plan_id"`
	}

	var add_request addRequest

	err := c.Ctx.ShouldBind(&add_request)
	if err != nil {
		help.Gin200ErrorResponse(c.Ctx, exception.VALIDATE_ERR, err.Error(), nil)
		return
	}

	//调用服务
	b_s := new(backend.BlogService)
	b_s.AddBlog(add_request.BlogTypeId, add_request.CoverPlanId, add_request.Title, add_request.Abstract, add_request.Content)

	//网络响应
	help.Gin200SuccessResponse(c.Ctx, "添加成功", nil)

	return
}

//更新
func (c *BlogController) UpdateBlog() {
	//定义获取方法参数
	type updateRequest struct {
		ID          int64  `form:"id" binding:"required"`
		BlogTypeId  int64  `form:"blog_type_id" binding:"required"`
		Title       string `form:"title"  binding:"required"`
		Abstract    string `form:"abstract"  binding:"required"`
		Content     string `form:"content"  binding:"required"`
		CoverPlanId int64  `form:"cover_plan_id"`
	}

	var update_request updateRequest

	err := c.Ctx.ShouldBind(&update_request)
	if err != nil {
		help.Gin200ErrorResponse(c.Ctx, exception.VALIDATE_ERR, err.Error(), nil)
		return
	}

	//调用服务
	b_s := new(backend.BlogService)
	b_s.UpdateBlog(update_request.ID, update_request.BlogTypeId, update_request.CoverPlanId,
		update_request.Title, update_request.Abstract, update_request.Content)

	//网络响应
	help.Gin200SuccessResponse(c.Ctx, "更新成功", nil)

	return
}

//导入数据到es中
func (c *BlogController) ImportData() {

	type importRequest struct {
		Password string `form:"password" binding:"required"`
	}

	var import_request importRequest

	err := c.Ctx.ShouldBind(&import_request)
	if err != nil {
		help.Gin200ErrorResponse(c.Ctx, exception.VALIDATE_ERR, "不要乱动这个方法，这个方法不对外提供的，请联系ly", nil)
		return
	}

	b_s := new(backend.BlogService)
	b_s.ImportDataToEs()

	help.Gin200SuccessResponse(c.Ctx, "导入完毕", nil)
	return
}

//搜素
func (c *BlogController) SearchBlog() {

	//非必填字段
	var search_level string
	search_level = c.Ctx.DefaultQuery("search_level", front.ES_SEARCH_LEVEL)

	//必填字段
	type searchRequest struct {
		//搜索维度
		Keyword string `form:"keyword" binding:"required"`
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

	b_s := new(front.BlogService)

	result := b_s.SearchBlog(search_level, search_request.Keyword, search_request.PerPage, search_request.Page)

	help.Gin200SuccessResponse(c.Ctx, "请求成功过", result)

	return
}
