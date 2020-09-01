package model

const (
	//附件功能模块
	ATTACHMENT_BLOG_Module = iota + 1
)

//数据库表
type AttachmentModel struct {
	BaseModel
	Module int64  `gorm:"cloumn:module"`
	Path   string `gorm:"cloumn:path"`
}

func (AttachmentModel) TableName() string {
	return "attachment"
}
