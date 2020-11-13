package backend

import (
	"github.com/blog_backend/common-lib/db/mysql"
	"github.com/blog_backend/exception"
	"github.com/blog_backend/help"
	"github.com/blog_backend/model"
	"github.com/blog_backend/service/blog"
	"strings"
)

type BlogController struct {
	BaseController
}

func (c *BlogController) AddEs() {

	type importRequest struct {
		ID uint `form:"id" binding:"required"`
	}

	var import_request importRequest

	err := c.Ctx.ShouldBind(&import_request)
	if err != nil {
		help.Gin200ErrorResponse(c.Ctx, exception.VALIDATE_ERR, "不要乱动这个方法，这个方法不对外提供的，请联系ly", nil)
		return
	}

	db := mysql.GetDefaultDBConnect()
	blog_model := new(model.BlogModel)
	db.First(blog_model, import_request.ID)
	blog_list_entity := blog.ChangeToBlogEntity(blog_model)
	b_s := new(blog.BlogEsBkService)
	b_s.AddDoc(blog_list_entity)

	help.Gin200SuccessResponse(c.Ctx, "添加完毕", nil)
	return
}

func (c *BlogController) CreateIndex() {
	type importRequest struct {
		Password string `form:"password" binding:"required"`
	}

	var import_request importRequest

	err := c.Ctx.ShouldBind(&import_request)
	if err != nil {
		help.Gin200ErrorResponse(c.Ctx, exception.VALIDATE_ERR, "不要乱动这个方法，这个方法不对外提供的，请联系ly", nil)
		return
	}

	if strings.Compare(import_request.Password, "ly123") != 0 {
		help.Gin200ErrorResponse(c.Ctx, exception.VALIDATE_ERR, "密码不对", nil)
		return
	}

	b_s := new(blog.BlogEsBkService)
	b_s.CreateIndex()

	help.Gin200SuccessResponse(c.Ctx, "创建完毕", nil)
	return
}
