package attachment

import "github.com/blog_backend/entity"

//附件的doc
type AttachmentDoc struct {
	*entity.BaseEntity
	//静态资源全路径
	FullUrl string `json:"full_url"`
	//静态资源功能性标识
	Module string `json:"module"`
	//相对路径
	Path string `json:"path"`
	//相对路径，和path是一样的作用
	Url string `json:"url"`
}
