package attachment

import "github.com/blog_backend/entity"

//附件的doc
type AttachmentEntity struct {
	entity.BaseEntity
	//静态资源全路径
	FullUrl string `json:"full_url"`
	//静态资源功能性标识
	Module int64 `json:"module"`
	//相对路径
	Path string `json:"path"`
	//相对路径，和path是一样的作用
	Url string `json:"url"`
	//文件类型
	FileType int64 `json:"file_type"`
}
