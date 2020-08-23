package common

import (
	"fmt"
	"github.com/blog_backend/controller"
	"github.com/blog_backend/exception"
	"github.com/blog_backend/help"
	"github.com/blog_backend/model"
	"github.com/blog_backend/service/common"
	"strconv"
)

type AttachmentController struct {
	controller.BaseController
}

func (a *AttachmentController) UploadBlog() {
	form, _ := a.Ctx.MultipartForm()

	var upload_resule interface{}

	if modules, ok := form.Value["modules"]; ok {

		attachment_service := new(common.AttachmentService)

		for _, module := range modules {

			switch int_module, _ := strconv.ParseInt(module, 10, 64); int_module {
			case model.ATTACHMENT_BLOG_Module:
				//博客图片
				upload_resule = attachment_service.UploadBlog(a.Ctx)
				break;
			default:
				panic(exception.NewException(exception.VALIDATE_ERR, fmt.Sprintf("没有支持的模块%d", int_module)))
			}
		}
	}

	help.Gin200SuccessResponse(a.Ctx,"上传成功",upload_resule)
	return
}
