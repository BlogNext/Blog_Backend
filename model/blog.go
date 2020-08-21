package model

type Blog struct {
	BaseModel
	BlogTypeId int64  `gorm:"cloumn:blog_type_id"`
	Title      string `gorm:"cloumn:title"`
	Abstract   string `gorm:"cloumn:abstract"`
	Content    string `gorm:"cloumn:content"`
}

func (Blog) TableName() string {
	return "blog"
}
