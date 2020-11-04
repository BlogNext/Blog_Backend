package backend

import (
	"github.com/blog_backend/help"
	"github.com/blog_backend/service/attachment"
)

type AttachmentController struct {
	BaseController
}

func (a *AttachmentController) UploadBlog() {

	attachment_service := new(attachment.AttachmentBkService)
	upload_resule := attachment_service.UploadBlog(a.Ctx)

	help.Gin200SuccessResponse(a.Ctx, "上传成功", upload_resule)

	return
}
