package blog

import (
	"github.com/blog_backend/entity"
	"github.com/blog_backend/entity/attachment"
)

//blog文档
type BlogEntity struct {
	entity.BaseEntity

	BlogTypeId uint64 `json:"blog_type_id"`

	YuqueFormat string `json:"yuque_format"`
	YuqueHtml   string `json:"yuque_html"`
	//文章标题
	Title string `json:"title"`
	//文章摘要
	Abstract string `json:"abstract"`
	//文章内容
	Content string `json:"content"`

	CoverPlanId uint64 `json:"cover_plan_id"`

	//附件信息
	AttachmentInfo *attachment.AttachmentEntity `json:"attachment_info"`

	BlogTypeObject *BlogTypeEntity `json:"blog_type_object"`
}

//blog_type文档
type BlogTypeEntity struct {
	entity.BaseEntity
	Title string `json:"title"`
}
