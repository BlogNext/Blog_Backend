package backend

import (
	"github.com/blog_backend/controller"
	"github.com/blog_backend/exception"
	"log"
)

type BackendController struct {
	controller.BaseController
}

//做一些鉴权操作等
func (c *BackendController) Prepare() exception.MyException {
	return nil
}

func (c *BackendController) Finish() {

	log.Println(c.UniqFullPath)
	////可以做一些释放资源的操作
	//log.Println("路径",c.Ctx.FullPath())
	//
	////测试反射的出来的是当前struct是BackendController,不能向上转
	//controller_type := reflect.TypeOf(c)
	//log.Println("当前执行的控制器是",controller_type)
	//log.Println("kind",controller_type.Kind().String())
	//log.Println("name",controller_type.String())
	//log.Println("pkgName",controller_type.PkgPath())
	//log.Println("Elem",controller_type.Elem())
	//log.Println("当前执行的方法:", c.Action)
}
