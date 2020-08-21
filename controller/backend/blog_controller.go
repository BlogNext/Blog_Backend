package backend

import (
	"github.com/blog_backend/controller"
	"github.com/blog_backend/exception"
	"github.com/blog_backend/help"
	"github.com/blog_backend/service/backend"
)

type BlogController struct {
	controller.BaseController
}

func (c *BlogController) GetList() {
	b_s := new(backend.BlogService)
	result := b_s.GetList()
	help.Gin200SuccessResponse(c.Ctx, "成功", result)
	return
}

func (c *BlogController) AddBlog() {

	//定义获取方法参数
	type addRequest struct {
		BlogTypeId int64  `form:"blog_type_id" binding:"required"`
		Title      string `form:"title"  binding:"required"`
		Abstract   string `form:"abstract"  binding:"required"`
		Content    string `form:"content"  binding:"required"`
	}

	var add_request addRequest

	err := c.Ctx.ShouldBind(&add_request)
	if err != nil {
		help.Gin200ErrorResponse(c.Ctx, exception.VALIDATE_ERR, err.Error(), nil)
		return
	}

	//调用服务
	b_s := new(backend.BlogService)
	b_s.AddBlog(add_request.BlogTypeId, add_request.Title, add_request.Abstract, add_request.Content)

	//网络响应
	help.Gin200SuccessResponse(c.Ctx, "添加成功", nil)

	return
}

//更新
func (c *BlogController) UpdateBlog() {
	//定义获取方法参数
	type updateRequest struct {
		ID         int64  `form:"id" binding:"required"`
		BlogTypeId int64  `form:"blog_type_id" binding:"required"`
		Title      string `form:"title"  binding:"required"`
		Abstract   string `form:"abstract"  binding:"required"`
		Content    string `form:"content"  binding:"required"`
	}

	var update_request updateRequest

	err := c.Ctx.ShouldBind(&update_request)
	if err != nil {
		help.Gin200ErrorResponse(c.Ctx, exception.VALIDATE_ERR, err.Error(), nil)
		return
	}

	//调用服务
	b_s := new(backend.BlogService)
	b_s.UpdateBlog(update_request.ID, update_request.BlogTypeId, update_request.Title, update_request.Abstract, update_request.Content)

	//网络响应
	help.Gin200SuccessResponse(c.Ctx, "更新成功", nil)

	return
}
