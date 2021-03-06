package backend

import (
	"github.com/blog_backend/help"
	"github.com/blog_backend/service/attachment"
)

type AttachmentController struct {
	BaseController
}

func (a *AttachmentController) UploadBlog() {

	attachmentService := new(attachment.AttachmentBkService)
	uploadResule := attachmentService.UploadBlog(a.Ctx)

	help.Gin200SuccessResponse(a.Ctx, "上传成功", uploadResule)

	return
}
