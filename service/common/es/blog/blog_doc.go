package blog

import (
	"github.com/blog_backend/service/common/es"
	"github.com/blog_backend/service/common/es/attachment"
)

//blog文档
type BlogDoc struct {
	es.BaseDoc
	//文章标题
	Title string `json:"title"`
	//文章摘要
	Abstract string `json:"abstract"`
	//文章内容
	Content string `json:"content"`

	//附件信息
	AttachmentInfo attachment.AttachmentDoc `json:"attachment_info"`

	BlogTypeObject BlogTypeDoc `json:"blog_type_object"`
}


//blog_type文档
type BlogTypeDoc struct {
	es.BaseDoc
	Title      string `json:"title"`
}
