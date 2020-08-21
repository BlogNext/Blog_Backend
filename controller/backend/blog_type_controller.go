package backend

import (
	"github.com/blog_backend/controller"
	"github.com/blog_backend/help"
	"github.com/blog_backend/service/backend"
	"log"
)

type BlogTypeController struct {
	controller.BaseController
}

//列表接口
func (c *BlogTypeController) Get() {
	log.Println("asdfasdf")
}

//添加一个类型
func (c *BlogTypeController) Post() {
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
func (c *BlogTypeController) Put() {

}
