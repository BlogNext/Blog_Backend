package model

type BlogTypeModel struct {
	BaseModel
	YuqueName string `gorm:"cloumn:yuque_name"`
	YuqueId   int64  `gorm:"cloumn:yuque_id"`
	YuqueType string `gorm:"cloumn:yuque_type"`
}

func (BlogTypeModel) TableName() string {
	return "blog_type"
}
