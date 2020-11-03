package test

import (
	"github.com/blog_backend/model"
	"github.com/blog_backend/service/attachment"
	"testing"
)

func TestDownloadBlogImage(t testing.T){
	service := new(attachment.AttachmentService)
	service.DownloadBlogImage("",model.ATTACHMENT_BLOG_Module,model.ATTACHMENT_FILE_TYPE_IMAGE)
}