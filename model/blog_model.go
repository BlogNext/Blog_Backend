package model

type BlogModel struct {
	BaseModel
	BlogTypeId  int64  `gorm:"cloumn:blog_type_id"`
	CoverPlanId int64  `gorm:"cloumn:cover_plan_id"`
	Title       string `gorm:"cloumn:title"`
	Abstract    string `gorm:"cloumn:abstract"`
	Content     string `gorm:"cloumn:content"`
	DocID       string `gorm:"cloumn:doc_id"`
}

func (BlogModel) TableName() string {
	return "blog"
}
