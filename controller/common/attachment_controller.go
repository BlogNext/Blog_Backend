package common

import (
	"github.com/blog_backend/controller"
	"github.com/blog_backend/help"
	"github.com/blog_backend/service/attachment"
)

type AttachmentController struct {
	controller.BaseController
}

func (a *AttachmentController) UploadBlog() {

	attachment_service := new(attachment.AttachmentBkService)
	upload_resule := attachment_service.UploadBlog(a.Ctx)

	help.Gin200SuccessResponse(a.Ctx, "上传成功", upload_resule)

	return
}
