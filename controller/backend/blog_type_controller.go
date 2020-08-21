package backend

import (
	"github.com/blog_backend/controller"
	"log"
)

type BlogTypeController struct {
	controller.BaseController
}

//列表接口
func (c *BlogTypeController) Get() {
	log.Println("asdfasdf")
}
