package model

type BlogModel struct {
	BaseModel
	DocID         string `gorm:"cloumn:doc_id"` //es文档id
	CoverPlanId   int64  `gorm:"cloumn:cover_plan_id"`
	BlogTypeId    int64  `gorm:"cloumn:blog_type_id"`
	YuqueId       int64  `gorm:"cloumn:yuque_id"` //语雀文档id
	YuqueSlug     string `gorm:"cloumn:yuque_slug"`
	YuqueIdFormat string `gorm:"cloumn:yuque_format"`
	YuqueHtml     string `gorm:"cloumn:yuque_html"`
	YuqueLake     string `gorm:"cloumn:yuque_lake"`
	Title         string `gorm:"cloumn:title"`
	Abstract      string `gorm:"cloumn:abstract"`
	Content       string `gorm:"cloumn:content"`
}

func (BlogModel) TableName() string {
	return "blog"
}
