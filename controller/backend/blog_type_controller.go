package backend

import (
	"github.com/blog_backend/exception"
	"github.com/blog_backend/help"
	"github.com/blog_backend/service/backend"
)

type BlogTypeController struct {
	BackendController
}

//列表接口
func (c *BlogTypeController) GetList() {

	//必填字段
	type searchRequest struct {
		PerPage int `form:"per_page" binding:"required"`
		Page    int `form:"page" binding:"required"`
	}

	var search_request searchRequest

	err := c.Ctx.ShouldBind(&search_request)

	if err != nil {
		help.Gin200ErrorResponse(c.Ctx, exception.VALIDATE_ERR, err.Error(), nil)
		return
	}

	b_t_s := new(backend.BlogTypeService)
	result := b_t_s.List(search_request.PerPage, search_request.Page)
	help.Gin200SuccessResponse(c.Ctx, "成功", result)

	return
}

//添加一个类型
func (c *BlogTypeController) AddType() {
	//获取参数
	title := c.Ctx.PostForm("title")

	//调用服务
	b_t_s := new(backend.BlogTypeService)
	b_t_s.Add(title)

	//网络响应
	help.Gin200SuccessResponse(c.Ctx, "添加成功", nil)

	return
}

//修改一个类型
func (c *BlogTypeController) UpdateType() {

	//定义获取方法参数
	type updateRequest struct {
		Id    int64  `form:"id" binding:"required"`
		Title string `form:"title"  binding:"required"`
	}
	var update_request updateRequest

	err := c.Ctx.ShouldBind(&update_request)
	if err != nil {
		help.Gin200ErrorResponse(c.Ctx, exception.VALIDATE_ERR, err.Error(), nil)
		return
	}

	//调用服务
	b_t_s := new(backend.BlogTypeService)
	b_t_s.Update(update_request.Id, update_request.Title)

	//网络响应
	help.Gin200SuccessResponse(c.Ctx, "修改成功", nil)

	return
}
