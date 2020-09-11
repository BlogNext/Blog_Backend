package model

const (
	//附件功能模块
	ATTACHMENT_BLOG_Module = iota + 1
)

const (
	//图片
	ATTACHMENT_FILE_TYPE_IMAGE = iota + 1
	//视频
	ATTACHMENT_FILE_TYPE_VIDEO
)

//数据库表
type AttachmentModel struct {
	BaseModel
	Module   int64  `gorm:"cloumn:module"`
	Path     string `gorm:"cloumn:path"`
	FileType int64  `gorm:"cloumn:file_type"`
}

func (AttachmentModel) TableName() string {
	return "attachment"
}

//验证功能模块常量
func (a *AttachmentModel) CheckValidModule(module int64) bool {
	switch module {
	case ATTACHMENT_BLOG_Module:
		return true
	}
	return false
}

//验证文件类型常量
func (a *AttachmentModel) CheckValidFileType(file_type int64) bool {
	switch file_type {
	case ATTACHMENT_FILE_TYPE_IMAGE:
		return true
	case ATTACHMENT_FILE_TYPE_VIDEO:
		return true
	}
	return false
}
