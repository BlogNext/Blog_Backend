package backend

type BlogController struct {
	BaseController
}

////测试添加一个es文档
//func (c *BlogController) AddEs() {
//
//	type importRequest struct {
//		ID uint `form:"id" binding:"required"`
//	}
//
//	var import_request importRequest
//
//	err := c.Ctx.ShouldBind(&import_request)
//	if err != nil {
//		help.Gin200ErrorResponse(c.Ctx, exception.VALIDATE_ERR, "不要乱动这个方法，这个方法不对外提供的，请联系ly", nil)
//		return
//	}
//
//	db := mysql.GetDefaultDBConnect()
//	blog_model := new(model.BlogModel)
//	db.First(blog_model, import_request.ID)
//	blog_list_entity := blog.ChangeToBlogEntity(blog_model)
//	b_s := new(blog.BlogEsBkService)
//	b_s.AddDoc(blog_list_entity)
//
//	help.Gin200SuccessResponse(c.Ctx, "添加完毕", nil)
//	return
//}
//
////测试获取一个文档
//func (c *BlogController) GetDetailEs() {
//	type importRequest struct {
//		ID uint `form:"id" binding:"required"`
//	}
//
//	var import_request importRequest
//
//	err := c.Ctx.ShouldBind(&import_request)
//	if err != nil {
//		help.Gin200ErrorResponse(c.Ctx, exception.VALIDATE_ERR, "不要乱动这个方法，这个方法不对外提供的，请联系ly", nil)
//		return
//	}
//
//	blog_entity := new(blog2.BlogEntity)
//	b_s := new(blog.BlogEsBkService)
//	b_s.GetDocByMysqlId(import_request.ID, blog_entity)
//
//	help.Gin200SuccessResponse(c.Ctx, "获取es文档", blog_entity)
//	return
//}
//
//func (c *BlogController) CreateIndex() {
//	type importRequest struct {
//		Password string `form:"password" binding:"required"`
//	}
//
//	var import_request importRequest
//
//	err := c.Ctx.ShouldBind(&import_request)
//	if err != nil {
//		help.Gin200ErrorResponse(c.Ctx, exception.VALIDATE_ERR, "不要乱动这个方法，这个方法不对外提供的，请联系ly", nil)
//		return
//	}
//
//	if strings.Compare(import_request.Password, "ly123") != 0 {
//		help.Gin200ErrorResponse(c.Ctx, exception.VALIDATE_ERR, "密码不对", nil)
//		return
//	}
//
//	b_s := new(blog.BlogEsBkService)
//	b_s.CreateIndex()
//
//	help.Gin200SuccessResponse(c.Ctx, "创建完毕", nil)
//	return
//}
