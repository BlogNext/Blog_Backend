package model

type BlogTypeModel struct {
	BaseModel
	Title string `gorm:"cloumn:title"`
}

func (BlogTypeModel) TableName() string {
	return "blog_type"
}
