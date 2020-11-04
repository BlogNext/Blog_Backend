package attachment

import (
	"github.com/blog_backend/model"
	"testing"
)

func TestDownloadBlogImage(t *testing.T){
	service := new(AttachmentRtService)
	service.DownloadBlogImage("https://cdn.nlark.com/yuque/0/2020/png/345835/1604299322469-aefd2283-9524-4dfa-82b0-2dc6f30ff6e9.png",model.ATTACHMENT_BLOG_Module,model.ATTACHMENT_FILE_TYPE_IMAGE)
}