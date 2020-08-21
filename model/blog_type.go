package model

type BlogType struct {
	BaseModel
	Title string `gorm:"cloumn:title"`
}

func (BlogType) TableName() string {
	return "blog_type"
}
